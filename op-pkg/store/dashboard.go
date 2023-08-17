package store

import (
	"context"
	"errors"

	"github.com/grafana/grafana/op-pkg/opstorage"
	"github.com/grafana/grafana/op-pkg/sdk/middleware"

	"github.com/grafana/grafana/pkg/infra/log"
	alertmodels "github.com/grafana/grafana/pkg/services/alerting/models"
	"github.com/grafana/grafana/pkg/services/dashboards"
	"github.com/grafana/grafana/pkg/services/folder"
	"github.com/grafana/grafana/pkg/services/quota"
)

const (
	DashboardTypeFolder      = "dash-folder"
	DashboardTypeDashboard   = "dash-db"
	DashboardTypeAlertFolder = "dash-folder-alerting"
)

type DashboardStore struct {
	logger    log.Logger
	opStorage *opstorage.Storage
}

func NewDashboardStore(logger log.Logger, opStorage *opstorage.Storage) *DashboardStore {
	return &DashboardStore{logger: logger, opStorage: opStorage}
}

func (d *DashboardStore) GetDashboardACLInfoList(ctx context.Context, query *dashboards.GetDashboardACLInfoListQuery) ([]*dashboards.DashboardACLInfoDTO, error) {
	return nil, nil
}

func (d *DashboardStore) HasAdminPermissionInDashboardsOrFolders(ctx context.Context, query *folder.HasAdminPermissionInDashboardsOrFoldersQuery) (bool, error) {
	return true, nil
}

func (d *DashboardStore) HasEditPermissionInFolders(ctx context.Context, query *folder.HasEditPermissionInFoldersQuery) (bool, error) {
	return true, nil
}

func (d *DashboardStore) ValidateDashboardBeforeSave(ctx context.Context, dashboard *dashboards.Dashboard, overwrite bool) (bool, error) {
	return true, nil
}

func (d *DashboardStore) DeleteACLByUser(ctx context.Context, userID int64) error {
	return nil
}

func (d *DashboardStore) GetFolderByTitle(ctx context.Context, orgID int64, title string) (*folder.Folder, error) {
	ctx = middleware.NewQuerierContext(ctx, "GetFolderByTitle")

	if title == "" {
		return nil, dashboards.ErrFolderTitleEmpty
	}

	dashboard, err := d.opStorage.Dashboard.GetDashboard(ctx, &opstorage.GetDashboardQuery{
		Title: title,
		Type:  DashboardTypeFolder,
		OrgID: orgID,
	})
	switch {
	case errors.Is(err, opstorage.ErrNotFound):
		return nil, dashboards.ErrFolderNotFound
	case err != nil:
		return nil, err
	default:
		return dashboards.FromDashboard(dashboard.ToModel()), nil
	}
}

func (d *DashboardStore) GetFolderByID(ctx context.Context, orgID int64, id int64) (*folder.Folder, error) {
	ctx = middleware.NewQuerierContext(ctx, "GetFolderByID")

	if id == 0 {
		return nil, dashboards.ErrDashboardIdentifierNotSet
	}

	dashboard, err := d.opStorage.Dashboard.GetDashboard(ctx, &opstorage.GetDashboardQuery{
		ID:    id,
		Type:  DashboardTypeFolder,
		OrgID: orgID,
	})
	switch {
	case errors.Is(err, opstorage.ErrNotFound):
		return nil, dashboards.ErrFolderNotFound
	case err != nil:
		return nil, err
	default:
		return dashboards.FromDashboard(dashboard.ToModel()), nil
	}
}

func (d *DashboardStore) GetFolderByUID(ctx context.Context, orgID int64, uid string) (*folder.Folder, error) {
	ctx = middleware.NewQuerierContext(ctx, "GetFolderByUID")

	if uid == "" {
		return nil, dashboards.ErrDashboardIdentifierNotSet
	}

	dashboard, err := d.opStorage.Dashboard.GetDashboard(ctx, &opstorage.GetDashboardQuery{
		UID:   uid,
		Type:  DashboardTypeFolder,
		OrgID: orgID,
	})
	switch {
	case errors.Is(err, opstorage.ErrNotFound):
		return nil, dashboards.ErrFolderNotFound
	case err != nil:
		return nil, err
	default:
		return dashboards.FromDashboard(dashboard.ToModel()), nil
	}
}

func (d *DashboardStore) GetProvisionedDataByDashboardID(ctx context.Context, dashboardID int64) (*dashboards.DashboardProvisioning, error) {
	return nil, nil
}

func (d *DashboardStore) GetProvisionedDataByDashboardUID(ctx context.Context, orgID int64, dashboardUID string) (*dashboards.DashboardProvisioning, error) {
	return nil, nil
}

func (d *DashboardStore) GetProvisionedDashboardData(ctx context.Context, name string) ([]*dashboards.DashboardProvisioning, error) {
	return nil, nil
}

func (d *DashboardStore) SaveProvisionedDashboard(ctx context.Context, cmd dashboards.SaveDashboardCommand, provisioning *dashboards.DashboardProvisioning) (*dashboards.Dashboard, error) {
	return nil, nil
}

func (d *DashboardStore) SaveDashboard(ctx context.Context, cmd dashboards.SaveDashboardCommand) (*dashboards.Dashboard, error) {
	ctx = middleware.NewQuerierContext(ctx, "SaveDashboard")

	dashboard, err := d.opStorage.Dashboard.SaveDashboard(ctx, &opstorage.SaveDashboardQuery{
		Dashboard:    cmd.Dashboard,
		UserID:       cmd.UserID,
		Overwrite:    cmd.Overwrite,
		Message:      cmd.Message,
		OrgID:        cmd.OrgID,
		RestoredFrom: cmd.RestoredFrom,
		PluginID:     cmd.PluginID,
		FolderID:     cmd.FolderID,
		FolderUID:    cmd.FolderUID,
		IsFolder:     cmd.IsFolder,
		UpdatedAt:    cmd.UpdatedAt,
	})
	if err != nil {
		return nil, err
	}
	return dashboard.ToModel(), nil
}

func (d *DashboardStore) UpdateDashboardACL(ctx context.Context, dashboardID int64, items []*dashboards.DashboardACL) error {
	return nil
}

func (d *DashboardStore) SaveAlerts(ctx context.Context, dashID int64, alerts []*alertmodels.Alert) error {
	return nil
}

// UnprovisionDashboard removes row in dashboard_provisioning for the dashboard making it seem as if manually created.
// The dashboard will still have `created_by = -1` to see it was not created by any particular user.
func (d *DashboardStore) UnprovisionDashboard(ctx context.Context, id int64) error {
	return nil
}

func (d *DashboardStore) DeleteOrphanedProvisionedDashboards(ctx context.Context, cmd *dashboards.DeleteOrphanedProvisionedDashboardsCommand) error {
	return nil
}

func (d *DashboardStore) Count(ctx context.Context, scopeParams *quota.ScopeParameters) (*quota.Map, error) {
	ctx = middleware.NewQuerierContext(ctx, "Count")

	u := &quota.Map{}

	var (
		count int64
		tag   quota.Tag
		err   error
	)

	count, err = d.opStorage.Dashboard.Count(ctx, &opstorage.CountDashboardsQuery{
		OrgID:  scopeParams.OrgID,
		UserID: scopeParams.UserID,
	})
	if err != nil {
		return nil, err
	}

	if scopeParams != nil && scopeParams.OrgID != 0 {
		tag, err = quota.NewTag(dashboards.QuotaTargetSrv, dashboards.QuotaTarget, quota.OrgScope)
		if err != nil {
			return u, err
		}
		u.Set(tag, count)
	} else {
		tag, err = quota.NewTag(dashboards.QuotaTargetSrv, dashboards.QuotaTarget, quota.GlobalScope)
		if err != nil {
			return u, err
		}
		u.Set(tag, count)
	}
	return u, nil
}

func (d *DashboardStore) GetDashboardsByPluginID(ctx context.Context, query *dashboards.GetDashboardsByPluginIDQuery) ([]*dashboards.Dashboard, error) {
	ctx = middleware.NewQuerierContext(ctx, "GetDashboardsByPluginID")

	list, err := d.opStorage.Dashboard.GetDashboardsByPluginID(ctx, &opstorage.GetDashboardsByPluginIDQuery{
		PluginID: query.PluginID,
		OrgID:    query.OrgID,
	})
	if err != nil {
		return nil, err
	}

	dbds := make([]*dashboards.Dashboard, 0, len(list))
	for _, item := range list {
		dbds = append(dbds, item.ToModel())
	}

	return dbds, nil
}

func (d *DashboardStore) DeleteDashboard(ctx context.Context, cmd *dashboards.DeleteDashboardCommand) error {
	ctx = middleware.NewQuerierContext(ctx, "DeleteDashboard")
	err := d.opStorage.Dashboard.DeleteDashboard(ctx, &opstorage.DeleteDashboardQuery{
		ID:    cmd.ID,
		OrgID: cmd.OrgID,
	})
	if errors.Is(err, opstorage.ErrNotFound) {
		return dashboards.ErrDashboardNotFound
	}
	return err
}

func (d *DashboardStore) GetDashboard(ctx context.Context, query *dashboards.GetDashboardQuery) (*dashboards.Dashboard, error) {
	ctx = middleware.NewQuerierContext(ctx, "GetDashboard")

	q := &opstorage.GetDashboardQuery{
		ID:    query.ID,
		UID:   query.UID,
		OrgID: query.OrgID,
	}
	if query.Title != nil {
		q.Title = *query.Title
	}
	if query.FolderID != nil {
		q.FolderID = query.FolderID
	}

	dashboard, err := d.opStorage.Dashboard.GetDashboard(ctx, q)
	switch {
	case errors.Is(err, opstorage.ErrNotFound):
		return nil, dashboards.ErrDashboardNotFound
	case err != nil:
		return nil, err
	default:
		return dashboard.ToModel(), nil
	}
}

func (d *DashboardStore) GetDashboardUIDByID(ctx context.Context, query *dashboards.GetDashboardRefByIDQuery) (*dashboards.DashboardRef, error) {
	ctx = middleware.NewQuerierContext(ctx, "GetDashboardUIDByID")

	dashboardRef, err := d.opStorage.Dashboard.GetDashboardRef(ctx, &opstorage.GetDashboardRefQuery{ID: query.ID})
	switch {
	case errors.Is(err, opstorage.ErrNotFound):
		return nil, dashboards.ErrDashboardNotFound
	case err != nil:
		return nil, err
	default:
		return &dashboards.DashboardRef{
			UID:  dashboardRef.UID,
			Slug: dashboardRef.Slug,
		}, nil
	}
}

func (d *DashboardStore) GetDashboards(ctx context.Context, query *dashboards.GetDashboardsQuery) ([]*dashboards.Dashboard, error) {
	ctx = middleware.NewQuerierContext(ctx, "GetDashboards")

	list, err := d.opStorage.Dashboard.GetDashboards(ctx, &opstorage.GetDashboardsQuery{
		DashboardIDs:  query.DashboardIDs,
		DashboardUIDs: query.DashboardUIDs,
		OrgID:         query.OrgID,
	})
	if err != nil {
		return nil, err
	}

	dbds := make([]*dashboards.Dashboard, 0, len(list))
	for _, item := range list {
		dbds = append(dbds, item.ToModel())
	}

	return dbds, nil
}

func (d *DashboardStore) FindDashboards(ctx context.Context, query *dashboards.FindPersistedDashboardsQuery) ([]dashboards.DashboardSearchProjection, error) {
	ctx = middleware.NewQuerierContext(ctx, "FindDashboards")

	// queried by GetUserVisibleNamespaces on original dashboard load in loop for alerts
	if query.Type == DashboardTypeAlertFolder {
		return []dashboards.DashboardSearchProjection{}, nil
	}

	// fetch folder uids that required to show folder names in web
	folderUIDs := d.fetchFolderUIDs(ctx)

	list, err := d.opStorage.Dashboard.FindDashboards(ctx, &opstorage.FindDashboardsQuery{
		Title:         query.Title,
		OrgID:         query.OrgId,
		DashboardIDs:  query.DashboardIds,
		DashboardUIDs: query.DashboardUIDs,
		Type:          query.Type,
		FolderIDs:     query.FolderIds,
		Tags:          query.Tags,
		Limit:         query.Limit,
		Page:          query.Page,
	})
	if err != nil {
		return nil, err
	}

	dbds := make([]dashboards.DashboardSearchProjection, 0, len(list))
	for _, item := range list {
		dbds = append(dbds, dashboards.DashboardSearchProjection{
			ID:        item.ID,
			UID:       item.UID,
			Title:     item.Title,
			Slug:      item.Slug,
			IsFolder:  item.IsFolder,
			FolderID:  item.FolderID,
			FolderUID: folderUIDs[item.FolderID],
		})
	}

	return dbds, nil
}

func (d *DashboardStore) fetchFolderUIDs(ctx context.Context) map[int64]string {
	uids := make(map[int64]string)

	folders, err := d.opStorage.Dashboard.FindDashboards(ctx,
		&opstorage.FindDashboardsQuery{Type: DashboardTypeFolder},
	)
	if err != nil {
		d.logger.Error("failed to fetch folder uids", "error", err)
		return uids
	}

	for _, folderItem := range folders {
		uids[folderItem.ID] = folderItem.UID
	}
	return uids
}

func (d *DashboardStore) GetDashboardTags(ctx context.Context, query *dashboards.GetDashboardTagsQuery) ([]*dashboards.DashboardTagCloudItem, error) {
	ctx = middleware.NewQuerierContext(ctx, "GetDashboardTags")

	uniqueTags, err := d.opStorage.Dashboard.GetDashboardTags(ctx, &opstorage.GetDashboardTagsQuery{OrgID: query.OrgID})
	if err != nil {
		return nil, err
	}
	queryResult := make([]*dashboards.DashboardTagCloudItem, 0)
	for _, uniqueTag := range uniqueTags {
		queryResult = append(queryResult, &dashboards.DashboardTagCloudItem{
			Term:  uniqueTag.Name,
			Count: uniqueTag.Count,
		})
	}
	return queryResult, nil
}

func (d *DashboardStore) CountDashboardsInFolder(
	ctx context.Context, req *dashboards.CountDashboardsInFolderRequest) (int64, error) {
	ctx = middleware.NewQuerierContext(ctx, "CountDashboardsInFolder")
	return d.opStorage.Dashboard.CountDashboardsInFolder(ctx, &opstorage.CountDashboardsInFolderQuery{
		FolderID: req.FolderID,
		OrgID:    req.OrgID,
	})
}
