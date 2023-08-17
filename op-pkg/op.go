package op_pkg

import (
	"os"
	"sync"

	"github.com/grafana/grafana/op-pkg/opstorage"
	"github.com/grafana/grafana/op-pkg/store"

	"github.com/grafana/grafana/pkg/infra/log"
)

var (
	opStorageOnce sync.Once
	opStorage     *opstorage.Storage
)

func getOPStorage() *opstorage.Storage {
	opStorageOnce.Do(func() {
		var (
			baseURL = os.Getenv("OPSTORAGE_BASEURL")
			apiKey  = os.Getenv("OPSTORAGE_APIKEY")
		)
		opStorage = opstorage.New(baseURL, apiKey)
	})
	return opStorage
}

var (
	datasourceOnce  sync.Once
	datasourceStore *store.DatasourceStore
)

func GetDatasourceStore(logger log.Logger) *store.DatasourceStore {
	datasourceOnce.Do(func() {
		datasourceStore = store.NewDatasourceStore(logger, getOPStorage())
	})
	return datasourceStore
}

var (
	dashboardOnce  sync.Once
	dashboardStore *store.DashboardStore
)

func GetDashboardStore(logger log.Logger) *store.DashboardStore {
	dashboardOnce.Do(func() {
		dashboardStore = store.NewDashboardStore(logger, getOPStorage())
	})
	return dashboardStore
}
