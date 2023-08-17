package store

import (
	"context"
	"errors"
	"fmt"

	"github.com/grafana/grafana/op-pkg/opstorage"
	"github.com/grafana/grafana/op-pkg/sdk/middleware"

	"github.com/grafana/grafana/pkg/infra/log"
	"github.com/grafana/grafana/pkg/infra/metrics"
	"github.com/grafana/grafana/pkg/services/datasources"
	"github.com/grafana/grafana/pkg/services/quota"
)

type DatasourceStore struct {
	logger    log.Logger
	opStorage *opstorage.Storage
}

func NewDatasourceStore(logger log.Logger, opStorage *opstorage.Storage) *DatasourceStore {
	return &DatasourceStore{logger: logger, opStorage: opStorage}
}

func (d *DatasourceStore) GetDataSource(ctx context.Context, query *datasources.GetDataSourceQuery) (*datasources.DataSource, error) {
	ctx = middleware.NewQuerierContext(ctx, "GetDataSource")
	metrics.MDBDataSourceQueryByID.Inc()

	if query.OrgID == 0 || (query.ID == 0 && len(query.Name) == 0 && len(query.UID) == 0) {
		return nil, datasources.ErrDataSourceIdentifierNotSet
	}

	datasource, err := d.opStorage.Datasource.GetDatasource(ctx, &opstorage.GetDataSourceQuery{
		ID:    query.ID,
		UID:   query.UID,
		Name:  query.Name,
		OrgID: query.OrgID,
	})
	switch {
	case errors.Is(err, opstorage.ErrNotFound):
		return nil, datasources.ErrDataSourceNotFound
	case err != nil:
		d.logger.Error("failed getting data source", "err", err, "uid", query.UID, "id", query.ID, "name", query.Name, "orgId", query.OrgID)
		return nil, err
	default:
		return datasource.ToModel(), nil
	}
}

func (d *DatasourceStore) GetDataSources(ctx context.Context, query *datasources.GetDataSourcesQuery) ([]*datasources.DataSource, error) {
	ctx = middleware.NewQuerierContext(ctx, "GetDataSources")

	list, err := d.opStorage.Datasource.GetDatasources(ctx, &opstorage.GetDatasourcesQuery{
		OrgID: query.OrgID,
		Limit: query.DataSourceLimit,
	})
	if err != nil {
		return nil, err
	}

	result := make([]*datasources.DataSource, 0, len(list))
	for _, item := range list {
		result = append(result, item.ToModel())
	}
	return result, nil
}

func (d *DatasourceStore) GetAllDataSources(ctx context.Context, query *datasources.GetAllDataSourcesQuery) ([]*datasources.DataSource, error) {
	ctx = middleware.NewQuerierContext(ctx, "GetAllDataSources")

	list, err := d.opStorage.Datasource.GetAllDatasources(ctx)
	if err != nil {
		return nil, err
	}

	result := make([]*datasources.DataSource, 0, len(list))
	for _, item := range list {
		result = append(result, item.ToModel())
	}
	return result, nil
}

func (d *DatasourceStore) GetDataSourcesByType(ctx context.Context, query *datasources.GetDataSourcesByTypeQuery) ([]*datasources.DataSource, error) {
	ctx = middleware.NewQuerierContext(ctx, "GetDataSourcesByType")

	if query.Type == "" {
		return nil, fmt.Errorf("datasource type cannot be empty")
	}

	list, err := d.opStorage.Datasource.GetDatasourcesByType(ctx, &opstorage.GetDatasourcesByTypeQuery{
		OrgID: query.OrgID,
		Type:  query.Type,
	})
	if err != nil {
		return nil, err
	}

	result := make([]*datasources.DataSource, 0, len(list))
	for _, item := range list {
		result = append(result, item.ToModel())
	}
	return result, nil
}

func (d *DatasourceStore) GetDefaultDataSource(ctx context.Context, query *datasources.GetDefaultDataSourceQuery) (*datasources.DataSource, error) {
	ctx = middleware.NewQuerierContext(ctx, "GetDefaultDataSource")

	datasource, err := d.opStorage.Datasource.GetDefaultDatasource(ctx, &opstorage.GetDefaultDataSourceQuery{
		OrgID: query.OrgID,
	})
	switch {
	case errors.Is(err, opstorage.ErrNotFound):
		return nil, datasources.ErrDataSourceNotFound
	case err != nil:
		d.logger.Error("failed getting data source", "err", err, "orgId", query.OrgID)
		return nil, err
	default:
		return datasource.ToModel(), nil
	}
}

func (d *DatasourceStore) DeleteDataSource(ctx context.Context, cmd *datasources.DeleteDataSourceCommand) error {
	return datasources.ErrDatasourceIsReadOnly
}

func (d *DatasourceStore) Count(ctx context.Context, scopeParams *quota.ScopeParameters) (*quota.Map, error) {
	ctx = middleware.NewQuerierContext(ctx, "Count")

	u := &quota.Map{}

	var (
		count int64
		tag   quota.Tag
		err   error
	)

	count, err = d.opStorage.Datasource.Count(ctx, &opstorage.CountDatasourceQuery{
		OrgID:  scopeParams.OrgID,
		UserID: scopeParams.UserID,
	})
	if err != nil {
		return nil, err
	}

	if scopeParams != nil && scopeParams.OrgID != 0 {
		tag, err = quota.NewTag(datasources.QuotaTargetSrv, datasources.QuotaTarget, quota.OrgScope)
		if err != nil {
			return u, err
		}
		u.Set(tag, count)
	} else {
		tag, err = quota.NewTag(datasources.QuotaTargetSrv, datasources.QuotaTarget, quota.GlobalScope)
		if err != nil {
			return u, err
		}
		u.Set(tag, count)
	}

	return u, nil
}

func (d *DatasourceStore) AddDataSource(ctx context.Context, cmd *datasources.AddDataSourceCommand) (*datasources.DataSource, error) {
	return nil, datasources.ErrDatasourceIsReadOnly
}

func (d *DatasourceStore) UpdateDataSource(ctx context.Context, cmd *datasources.UpdateDataSourceCommand) (*datasources.DataSource, error) {
	return nil, datasources.ErrDatasourceIsReadOnly
}
