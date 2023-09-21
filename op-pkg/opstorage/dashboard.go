package opstorage

import (
	"bytes"
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
	"github.com/grafana/grafana/pkg/services/dashboards"
)

type dashboardStorage struct {
	client *client.Client
}

type Dashboard struct {
	ID        int64            `json:"id"`
	UID       string           `json:"uid"`
	Slug      string           `json:"slug"`
	Title     string           `json:"title"`
	OrgID     int64            `json:"orgId"`
	GnetID    int64            `json:"gnetId"`
	Version   int              `json:"version"`
	PluginID  string           `json:"pluginId"`
	FolderID  int64            `json:"folderId"`
	IsFolder  bool             `json:"isFolder"`
	HasACL    bool             `json:"hasACL"`
	Data      *simplejson.Json `json:"data"`
	UpdatedBy int64            `json:"updatedBy"`
	CreatedBy int64            `json:"createdBy"`
	Created   time.Time        `json:"created"`
	Updated   time.Time        `json:"updated"`
}

func (dashboard *Dashboard) ToModel() *dashboards.Dashboard {
	return &dashboards.Dashboard{
		ID:        dashboard.ID,
		UID:       dashboard.UID,
		Slug:      dashboard.Slug,
		OrgID:     dashboard.OrgID,
		GnetID:    dashboard.GnetID,
		Version:   dashboard.Version,
		PluginID:  dashboard.PluginID,
		Created:   dashboard.Created,
		Updated:   dashboard.Updated,
		UpdatedBy: dashboard.UpdatedBy,
		CreatedBy: dashboard.CreatedBy,
		FolderID:  dashboard.FolderID,
		IsFolder:  dashboard.IsFolder,
		HasACL:    dashboard.HasACL,
		Title:     dashboard.Title,
		Data:      dashboard.Data,
	}
}

type SaveDashboardQuery struct {
	Dashboard    *simplejson.Json `json:"dashboard"`
	UserID       int64            `json:"userId"`
	Overwrite    bool             `json:"overwrite"`
	Message      string           `json:"message"`
	OrgID        int64            `json:"orgId"`
	RestoredFrom int              `json:"restoredFrom"`
	PluginID     string           `json:"pluginId"`
	FolderID     int64            `json:"folderId"`
	FolderUID    string           `json:"folderUid"`
	IsFolder     bool             `json:"isFolder"`
	UpdatedAt    time.Time        `json:"updatedAt"`
}

func (s *dashboardStorage) SaveDashboard(ctx context.Context, query *SaveDashboardQuery) (*Dashboard, error) {
	var (
		requestContextData = middleware.GetRequestContextData(ctx)
		userSessionData    = middleware.GetUserSessionData(ctx)
	)

	if userSessionData == "" {
		return nil, ErrEmptyUserSession
	}

	payload, err := json.Marshal(query)
	if err != nil {
		return nil, err
	}
	data, err := s.client.Post(ctx, "dashboard/saveDashboard",
		bytes.NewReader(payload),
		interceptor.WithRequestHeader("X-REQUEST-CONTEXT", requestContextData),
		interceptor.WithRequestCookie(&http.Cookie{Name: "user_session", Value: userSessionData}),
	)
	if err != nil {
		return nil, err
	}
	var resp Dashboard
	err = json.Unmarshal(data, &resp)
	if err != nil {
		return nil, err
	}
	return &resp, err
}

type GetDashboardQuery struct {
	ID       int64
	UID      string
	Title    string
	Slug     string
	FolderID *int64
	Type     string
	OrgID    int64
}

func (s *dashboardStorage) GetDashboard(ctx context.Context, query *GetDashboardQuery) (*Dashboard, error) {
	var (
		requestContextData = middleware.GetRequestContextData(ctx)
		userSessionData    = middleware.GetUserSessionData(ctx)
	)

	if userSessionData == "" {
		return nil, ErrEmptyUserSession
	}

	params := url.Values{}
	if query.ID != 0 {
		params.Set("id", strconv.FormatInt(query.ID, 10))
	}
	if query.UID != "" {
		params.Set("uid", query.UID)
	}
	if query.Title != "" {
		params.Set("title", query.Title)
	}
	if query.Slug != "" {
		params.Set("slug", query.Slug)
	}
	if query.FolderID != nil {
		params.Set("folder_id", strconv.FormatInt(*query.FolderID, 10))
	}
	if query.Type != "" {
		params.Set("type", query.Type)
	}
	params.Set("org_id", strconv.FormatInt(query.OrgID, 10))

	data, err := s.client.Get(ctx, "dashboard/getDashboard",
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
		var resp Dashboard
		err = json.Unmarshal(data, &resp)
		if err != nil {
			return nil, err
		}
		return &resp, err
	}
}

type GetDashboardRefQuery struct {
	ID int64
}

type DashboardRef struct {
	UID  string `json:"uid"`
	Slug string `json:"slug"`
}

func (s *dashboardStorage) GetDashboardRef(ctx context.Context, query *GetDashboardRefQuery) (*DashboardRef, error) {
	var (
		requestContextData = middleware.GetRequestContextData(ctx)
		userSessionData    = middleware.GetUserSessionData(ctx)
	)

	if userSessionData == "" {
		return nil, ErrEmptyUserSession
	}

	params := url.Values{}

	if query.ID != 0 {
		params.Set("id", strconv.FormatInt(query.ID, 10))
	}

	data, err := s.client.Get(ctx, "dashboard/getDashboardRef",
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
		var resp DashboardRef
		err = json.Unmarshal(data, &resp)
		if err != nil {
			return nil, err
		}
		return &resp, err
	}
}

type GetDashboardTagsQuery struct {
	OrgID int64
}

type DashboardTag struct {
	Name  string `json:"name"`
	Count int    `json:"count"`
}

func (s *dashboardStorage) GetDashboardTags(ctx context.Context, query *GetDashboardTagsQuery) ([]*DashboardTag, error) {
	var (
		requestContextData = middleware.GetRequestContextData(ctx)
		userSessionData    = middleware.GetUserSessionData(ctx)
	)

	if userSessionData == "" {
		return nil, ErrEmptyUserSession
	}

	params := url.Values{}
	params.Set("org_id", strconv.FormatInt(query.OrgID, 10))

	data, err := s.client.Get(ctx, "dashboard/getDashboardTags",
		interceptor.WithRequestQueryParams(params),
		interceptor.WithRequestHeader("X-REQUEST-CONTEXT", requestContextData),
		interceptor.WithRequestCookie(&http.Cookie{Name: "user_session", Value: userSessionData}),
	)
	if err != nil {
		return nil, err
	}
	var resp struct {
		List []*DashboardTag `json:"list"`
	}
	err = json.Unmarshal(data, &resp)
	if err != nil {
		return nil, err
	}
	return resp.List, err
}

type FindDashboardsQuery struct {
	Title         string
	OrgID         int64
	DashboardIDs  []int64
	DashboardUIDs []string
	Type          string
	FolderIDs     []int64
	Tags          []string
	Limit         int64
	Page          int64
}

func (s *dashboardStorage) FindDashboards(ctx context.Context, query *FindDashboardsQuery) ([]*Dashboard, error) {
	var (
		requestContextData = middleware.GetRequestContextData(ctx)
		userSessionData    = middleware.GetUserSessionData(ctx)
	)

	if userSessionData == "" {
		return nil, ErrEmptyUserSession
	}

	params := url.Values{}
	if query.Title != "" {
		params.Set("title", query.Title)
	}
	if query.Type != "" {
		params.Set("type", query.Type)
	}
	if query.Limit != 0 {
		params.Set("limit", strconv.FormatInt(query.Limit, 10))
	}
	if query.Page != 0 {
		params.Set("page", strconv.FormatInt(query.Page, 10))
	}
	for _, dashboardID := range query.DashboardIDs {
		params.Add("dashboard_ids[]", strconv.FormatInt(dashboardID, 10))
	}
	for _, dashboardUID := range query.DashboardUIDs {
		params.Add("dashboard_uids[]", dashboardUID)
	}
	for _, folderID := range query.FolderIDs {
		params.Add("folder_ids[]", strconv.FormatInt(folderID, 10))
	}
	for _, tag := range query.Tags {
		params.Add("tags[]", tag)
	}
	params.Set("org_id", strconv.FormatInt(query.OrgID, 10))

	data, err := s.client.Get(ctx, "dashboard/findDashboards",
		interceptor.WithRequestQueryParams(params),
		interceptor.WithRequestHeader("X-REQUEST-CONTEXT", requestContextData),
		interceptor.WithRequestCookie(&http.Cookie{Name: "user_session", Value: userSessionData}),
	)
	if err != nil {
		return nil, err
	}
	var resp struct {
		List []*Dashboard `json:"list"`
	}
	err = json.Unmarshal(data, &resp)
	if err != nil {
		return nil, err
	}
	return resp.List, err
}

type GetDashboardsQuery struct {
	DashboardIDs  []int64
	DashboardUIDs []string
	OrgID         int64
}

func (s *dashboardStorage) GetDashboards(ctx context.Context, query *GetDashboardsQuery) ([]*Dashboard, error) {
	var (
		requestContextData = middleware.GetRequestContextData(ctx)
		userSessionData    = middleware.GetUserSessionData(ctx)
	)

	if userSessionData == "" {
		return nil, ErrEmptyUserSession
	}

	params := url.Values{}
	for _, dashboardID := range query.DashboardIDs {
		params.Add("dashboard_ids[]", strconv.FormatInt(dashboardID, 10))
	}
	for _, dashboardUID := range query.DashboardUIDs {
		params.Add("dashboard_uids[]", dashboardUID)
	}
	params.Set("org_id", strconv.FormatInt(query.OrgID, 10))

	data, err := s.client.Get(ctx, "dashboard/getDashboards",
		interceptor.WithRequestQueryParams(params),
		interceptor.WithRequestHeader("X-REQUEST-CONTEXT", requestContextData),
		interceptor.WithRequestCookie(&http.Cookie{Name: "user_session", Value: userSessionData}),
	)
	if err != nil {
		return nil, err
	}
	var resp struct {
		List []*Dashboard `json:"list"`
	}
	err = json.Unmarshal(data, &resp)
	if err != nil {
		return nil, err
	}
	return resp.List, err
}

type GetDashboardsByPluginIDQuery struct {
	PluginID string
	OrgID    int64
}

func (s *dashboardStorage) GetDashboardsByPluginID(ctx context.Context, query *GetDashboardsByPluginIDQuery) ([]*Dashboard, error) {
	var (
		requestContextData = middleware.GetRequestContextData(ctx)
		userSessionData    = middleware.GetUserSessionData(ctx)
	)

	if userSessionData == "" {
		return nil, ErrEmptyUserSession
	}

	params := url.Values{}
	if query.PluginID != "" {
		params.Set("plugin_id", query.PluginID)
	}
	params.Set("org_id", strconv.FormatInt(query.OrgID, 10))

	data, err := s.client.Get(ctx, "dashboard/getDashboardsByPluginID",
		interceptor.WithRequestQueryParams(params),
		interceptor.WithRequestHeader("X-REQUEST-CONTEXT", requestContextData),
		interceptor.WithRequestCookie(&http.Cookie{Name: "user_session", Value: userSessionData}),
	)
	if err != nil {
		return nil, err
	}
	var resp struct {
		List []*Dashboard `json:"list"`
	}
	err = json.Unmarshal(data, &resp)
	if err != nil {
		return nil, err
	}
	return resp.List, err
}

type CountDashboardsInFolderQuery struct {
	FolderID int64
	OrgID    int64
}

func (s *dashboardStorage) CountDashboardsInFolder(ctx context.Context, query *CountDashboardsInFolderQuery) (int64, error) {
	var (
		requestContextData = middleware.GetRequestContextData(ctx)
		userSessionData    = middleware.GetUserSessionData(ctx)
	)

	if userSessionData == "" {
		return 0, ErrEmptyUserSession
	}

	params := url.Values{}
	params.Set("folder_id", strconv.FormatInt(query.FolderID, 10))
	params.Set("org_id", strconv.FormatInt(query.OrgID, 10))

	data, err := s.client.Get(ctx, "dashboard/countDashboardsInFolder",
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
	return resp.Count, err
}

type CountDashboardsQuery struct {
	OrgID  int64
	UserID int64
}

func (s *dashboardStorage) Count(ctx context.Context, query *CountDashboardsQuery) (int64, error) {
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

	data, err := s.client.Get(ctx, "dashboard/count",
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
	return resp.Count, err
}

type DeleteDashboardQuery struct {
	ID    int64
	OrgID int64
}

func (s *dashboardStorage) DeleteDashboard(ctx context.Context, query *DeleteDashboardQuery) error {
	var (
		requestContextData = middleware.GetRequestContextData(ctx)
		userSessionData    = middleware.GetUserSessionData(ctx)
	)

	if userSessionData == "" {
		return ErrEmptyUserSession
	}

	params := url.Values{}
	if query.ID != 0 {
		params.Set("id", strconv.FormatInt(query.ID, 10))
	}
	params.Set("org_id", strconv.FormatInt(query.OrgID, 10))

	_, err := s.client.Delete(ctx, "dashboard/deleteDashboard",
		interceptor.WithRequestQueryParams(params),
		interceptor.WithRequestHeader("X-REQUEST-CONTEXT", requestContextData),
		interceptor.WithRequestCookie(&http.Cookie{Name: "user_session", Value: userSessionData}),
	)
	if errors.Is(err, client.ErrNotFound) {
		return ErrNotFound
	}
	return err
}
