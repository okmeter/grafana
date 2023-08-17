package datasource

import (
	"context"
	"fmt"
	"time"

	op_pkg "github.com/grafana/grafana/op-pkg"
	"github.com/grafana/grafana/op-pkg/sdk/middleware"
	"github.com/grafana/grafana/op-pkg/store"

	"github.com/grafana/grafana/pkg/infra/db"
	"github.com/grafana/grafana/pkg/infra/localcache"
	"github.com/grafana/grafana/pkg/infra/log"
	"github.com/grafana/grafana/pkg/services/datasources"
	"github.com/grafana/grafana/pkg/services/user"
)

const (
	DefaultCacheTTL = 5 * time.Second
)

func ProvideCacheService(cacheService *localcache.CacheService, sqlStore db.DB) *Service {
	logger := log.New("datasources")
	return &Service{
		logger:       logger,
		cacheTTL:     DefaultCacheTTL,
		CacheService: cacheService,
		store:        op_pkg.GetDatasourceStore(logger),
	}
}

type Service struct {
	logger       log.Logger
	cacheTTL     time.Duration
	store        *store.DatasourceStore
	CacheService *localcache.CacheService
}

func (dc *Service) GetDatasource(
	ctx context.Context,
	datasourceID int64,
	user *user.SignedInUser,
	skipCache bool,
) (*datasources.DataSource, error) {
	var (
		userSessionData    = middleware.GetUserSessionData(ctx)
		requestContextData = middleware.GetRequestContextData(ctx)
		cacheKey           = idKey(requestContextData, userSessionData, datasourceID)
	)

	if !skipCache {
		dc.logger.FromContext(ctx).Debug("Querying for data source via cache", "key", cacheKey)
		if cached, found := dc.CacheService.Get(cacheKey); found {
			ds := cached.(*datasources.DataSource)
			return ds, nil
		}
	}

	dc.logger.FromContext(ctx).Debug("Querying for data source via store", "id", datasourceID, "orgId", user.OrgID)

	query := &datasources.GetDataSourceQuery{ID: datasourceID, OrgID: user.OrgID}
	ds, err := dc.store.GetDataSource(ctx, query)
	if err != nil {
		return nil, err
	}

	if ds.UID != "" {
		dc.CacheService.Set(uidKey(requestContextData, userSessionData, ds.UID), ds, time.Second*5)
	}
	dc.CacheService.Set(cacheKey, ds, dc.cacheTTL)
	return ds, nil
}

func (dc *Service) GetDatasourceByUID(
	ctx context.Context,
	datasourceUID string,
	user *user.SignedInUser,
	skipCache bool,
) (*datasources.DataSource, error) {
	if datasourceUID == "" {
		return nil, fmt.Errorf("can not get data source by uid, uid is empty")
	}
	if user.OrgID == 0 {
		return nil, fmt.Errorf("can not get data source by uid, orgId is missing")
	}

	var (
		userSessionData    = middleware.GetUserSessionData(ctx)
		requestContextData = middleware.GetRequestContextData(ctx)
		uidCacheKey        = uidKey(requestContextData, userSessionData, datasourceUID)
	)

	if !skipCache {
		dc.logger.FromContext(ctx).Debug("Querying for data source via cache", "key", uidCacheKey)
		if cached, found := dc.CacheService.Get(uidCacheKey); found {
			ds := cached.(*datasources.DataSource)
			return ds, nil
		}
	}

	dc.logger.FromContext(ctx).Debug("Querying for data source via store", "uid", datasourceUID, "orgId", user.OrgID)
	query := &datasources.GetDataSourceQuery{UID: datasourceUID, OrgID: user.OrgID}
	ds, err := dc.store.GetDataSource(ctx, query)
	if err != nil {
		return nil, err
	}

	dc.CacheService.Set(uidCacheKey, ds, dc.cacheTTL)
	dc.CacheService.Set(idKey(requestContextData, userSessionData, ds.ID), ds, dc.cacheTTL)
	return ds, nil
}

func idKey(requestContext, userSession string, id int64) string {
	return fmt.Sprintf("ds-id-%s-%s-%d", requestContext, userSession, id)
}

func uidKey(requestContext, userSession string, uid string) string {
	return fmt.Sprintf("ds-uid-%s-%s-%s", requestContext, userSession, uid)
}
