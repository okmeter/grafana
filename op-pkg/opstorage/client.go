package opstorage

import (
	"errors"
	"net/http"

	"github.com/grafana/grafana/op-pkg/sdk/client"
	"github.com/grafana/grafana/op-pkg/sdk/roundtripper"
)

var (
	// ErrNotFound is a custom error to make it easier to differ proxy misconfiguration (default 404 response)
	// from missing an actual item (query by id failed)
	ErrNotFound = errors.New("not found")
	// ErrEmptyUserSession is a custom error to identify empty user session
	ErrEmptyUserSession = errors.New("empty user session")
)

type Storage struct {
	Datasource *datasourceStorage
	Dashboard  *dashboardStorage
}

func New(baseURL, apiKey string) *Storage {
	c := client.New(
		baseURL,
		client.OptionHeader("X-API-Key", apiKey),
		client.OptionTransport(
			roundtripper.NewLoggingRoundTripper(
				http.DefaultTransport, "op-storage"),
		),
	)
	s := &Storage{}
	s.Datasource = &datasourceStorage{client: c}
	s.Dashboard = &dashboardStorage{client: c}
	return s
}
