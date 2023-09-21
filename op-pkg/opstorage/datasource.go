package opstorage

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/grafana/grafana/op-pkg/sdk/client"
	"github.com/grafana/grafana/op-pkg/sdk/client/interceptor"
	"github.com/grafana/grafana/op-pkg/sdk/middleware"

	"github.com/grafana/grafana/pkg/components/simplejson"
	"github.com/grafana/grafana/pkg/services/datasources"
)

type datasourceStorage struct {
	client *client.Client
}

type Datasource struct {
	ID              int64             `json:"id"`
	UID             string            `json:"uid"`
	OrgID           int64             `json:"orgId"`
	Version         int               `json:"version"`
	Name            string            `json:"name"`
	Type            string            `json:"type"`
	Access          string            `json:"access"`
	URL             string            `json:"url"`
	User            string            `json:"user"`
	Database        string            `json:"database"`
	BasicAuth       bool              `json:"basicAuth"`
	BasicAuthUser   string            `json:"basicAuthUser"`
	WithCredentials bool              `json:"withCredentials"`
	IsDefault       bool              `json:"isDefault"`
	JsonData        *simplejson.Json  `json:"jsonData"`
	SecureJsonData  map[string]string `json:"secureJsonData"`
	ReadOnly        bool              `json:"readOnly"`
	Created         time.Time         `json:"created"`
	Updated         time.Time         `json:"updated"`
}

func (datasource *Datasource) ToModel() *datasources.DataSource {
	secureJsonData := make(map[string][]byte)
	for k, v := range datasource.SecureJsonData {
		secureJsonData[k] = []byte(v)
	}
	return &datasources.DataSource{
		ID:              datasource.ID,
		OrgID:           datasource.OrgID,
		Version:         datasource.Version,
		Name:            datasource.Name,
		Type:            datasource.Type,
		Access:          datasources.DsAccess(datasource.Access),
		URL:             datasource.URL,
		User:            datasource.User,
		Database:        datasource.Database,
		BasicAuth:       datasource.BasicAuth,
		BasicAuthUser:   datasource.BasicAuthUser,
		WithCredentials: datasource.WithCredentials,
		IsDefault:       datasource.IsDefault,
		JsonData:        datasource.JsonData,
		SecureJsonData:  secureJsonData,
		ReadOnly:        datasource.ReadOnly,
		UID:             datasource.UID,
		Created:         datasource.Created,
		Updated:         datasource.Updated,
	}
}

type GetDataSourceQuery struct {
	ID    int64
	UID   string
	Name  string
	OrgID int64
}

func (s *datasourceStorage) GetDatasource(ctx context.Context, query *GetDataSourceQuery) (*Datasource, error) {
	var (
		requestContextData = middleware.GetRequestContextData(ctx)
		userSessionData    = middleware.GetUserSessionData(ctx)
	)

	if userSessionData == "" {
		return nil, ErrEmptyUserSession
	}

	params := url.Values{}
	if query.ID > 0 {
		params.Set("id", strconv.FormatInt(query.ID, 10))
	}
	if query.UID != "" {
		params.Set("uid", query.UID)
	}
	if query.Name != "" {
		params.Set("name", query.Name)
	}
	params.Set("org_id", strconv.FormatInt(query.OrgID, 10))

	data, err := s.client.Get(ctx, "datasource/getDatasource",
		interceptor.WithRequestQueryParams(params),
		interceptor.WithResponseCodeCustomError(http.StatusNoContent, client.ErrNotFound),
		interceptor.WithRequestHeader("X-REQUEST-CONTEXT", requestContextData),
		interceptor.WithRequestCookie(&http.Cookie{Name: "user_session", Value: userSessionData}),
	)
	switch {
	case errors.Is(err, client.ErrNotFound):
		return nil, ErrNotFound
	case err != nil:
		return nil, err
	default:
		var resp Datasource
		err = json.Unmarshal(data, &resp)
		if err != nil {
			return nil, err
		}
		return &resp, err
	}
}

type GetDefaultDataSourceQuery struct {
	OrgID int64
}

func (s *datasourceStorage) GetDefaultDatasource(ctx context.Context, query *GetDefaultDataSourceQuery) (*Datasource, error) {
	var (
		requestContextData = middleware.GetRequestContextData(ctx)
		userSessionData    = middleware.GetUserSessionData(ctx)
	)

	if userSessionData == "" {
		return nil, ErrEmptyUserSession
	}

	params := url.Values{}
	params.Set("org_id", strconv.FormatInt(query.OrgID, 10))

	data, err := s.client.Get(ctx, "datasource/getDefaultDatasource",
		interceptor.WithRequestQueryParams(params),
		interceptor.WithResponseCodeCustomError(http.StatusNoContent, client.ErrNotFound),
		interceptor.WithRequestHeader("X-REQUEST-CONTEXT", requestContextData),
		interceptor.WithRequestCookie(&http.Cookie{Name: "user_session", Value: userSessionData}),
	)
	switch {
	case errors.Is(err, client.ErrNotFound):
		return nil, ErrNotFound
	case err != nil:
		return nil, err
	default:
		var resp Datasource
		err = json.Unmarshal(data, &resp)
		if err != nil {
			return nil, err
		}
		return &resp, err
	}
}

func (s *datasourceStorage) GetAllDatasources(ctx context.Context) ([]*Datasource, error) {
	var (
		requestContextData = middleware.GetRequestContextData(ctx)
		userSessionData    = middleware.GetUserSessionData(ctx)
	)

	if userSessionData == "" {
		return nil, ErrEmptyUserSession
	}

	data, err := s.client.Get(ctx, "datasource/getAllDatasources",
		interceptor.WithRequestHeader("X-REQUEST-CONTEXT", requestContextData),
		interceptor.WithRequestCookie(&http.Cookie{Name: "user_session", Value: userSessionData}),
	)
	if err != nil {
		return nil, err
	}
	var resp struct {
		List []*Datasource `json:"list"`
	}
	err = json.Unmarshal(data, &resp)
	if err != nil {
		return nil, err
	}
	return resp.List, err
}

type GetDatasourcesQuery struct {
	OrgID int64
	Limit int
}

func (s *datasourceStorage) GetDatasources(ctx context.Context, query *GetDatasourcesQuery) ([]*Datasource, error) {
	var (
		requestContextData = middleware.GetRequestContextData(ctx)
		userSessionData    = middleware.GetUserSessionData(ctx)
	)

	if userSessionData == "" {
		return nil, ErrEmptyUserSession
	}

	params := url.Values{}
	if query.Limit > 0 {
		params.Set("limit", strconv.Itoa(query.Limit))
	}
	params.Set("org_id", strconv.FormatInt(query.OrgID, 10))

	data, err := s.client.Get(ctx, "datasource/getDatasources",
		interceptor.WithRequestQueryParams(params),
		interceptor.WithRequestHeader("X-REQUEST-CONTEXT", requestContextData),
		interceptor.WithRequestCookie(&http.Cookie{Name: "user_session", Value: userSessionData}),
	)
	if err != nil {
		return nil, err
	}
	var resp struct {
		List []*Datasource `json:"list"`
	}
	err = json.Unmarshal(data, &resp)
	if err != nil {
		return nil, err
	}
	return resp.List, err
}

type GetDatasourcesByTypeQuery struct {
	OrgID int64
	Type  string
}

func (s *datasourceStorage) GetDatasourcesByType(ctx context.Context, query *GetDatasourcesByTypeQuery) ([]*Datasource, error) {
	var (
		requestContextData = middleware.GetRequestContextData(ctx)
		userSessionData    = middleware.GetUserSessionData(ctx)
	)

	if userSessionData == "" {
		return nil, ErrEmptyUserSession
	}

	params := url.Values{}
	if query.Type != "" {
		params.Set("type", query.Type)
	}
	params.Set("org_id", strconv.FormatInt(query.OrgID, 10))

	data, err := s.client.Get(ctx, "datasource/getDatasourcesByType",
		interceptor.WithRequestQueryParams(params),
		interceptor.WithRequestHeader("X-REQUEST-CONTEXT", requestContextData),
		interceptor.WithRequestCookie(&http.Cookie{Name: "user_session", Value: userSessionData}),
	)
	if err != nil {
		return nil, err
	}
	var resp struct {
		List []*Datasource `json:"list"`
	}
	err = json.Unmarshal(data, &resp)
	if err != nil {
		return nil, err
	}
	return resp.List, err
}

type CountDatasourceQuery struct {
	OrgID  int64
	UserID int64
}

func (s *datasourceStorage) Count(ctx context.Context, query *CountDatasourceQuery) (int64, error) {
	var (
		requestContextData = middleware.GetRequestContextData(ctx)
		userSessionData    = middleware.GetUserSessionData(ctx)
	)

	if userSessionData == "" {
		return 0, ErrEmptyUserSession
	}

	params := url.Values{}
	if query.UserID != 0 {
		params.Set("user_id", strconv.FormatInt(query.UserID, 10))
	}
	params.Set("org_id", strconv.FormatInt(query.OrgID, 10))

	data, err := s.client.Get(ctx, "datasource/count",
		interceptor.WithRequestQueryParams(params),
		interceptor.WithRequestHeader("X-REQUEST-CONTEXT", requestContextData),
		interceptor.WithRequestCookie(&http.Cookie{Name: "user_session", Value: userSessionData}),
	)
	if err != nil {
		return 0, err
	}
	var resp struct {
		Count int64 `json:"count"`
	}
	err = json.Unmarshal(data, &resp)
	if err != nil {
		return 0, err
	}
	return resp.Count, nil
}
