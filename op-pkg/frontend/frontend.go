package frontend

import (
	"context"
	"path"

	op_middleware "github.com/grafana/grafana/op-pkg/sdk/middleware"
)

// BuildAppURLOverrides creates new values for appURL and appSubURL only for fronted
// This function is required to replace existing approach to run behind reverse proxy https://grafana.com/tutorials/run-grafana-behind-a-proxy/
// But instead of using constant endpoint, this allows grafana to serve under dynamic sub-paths like localhost:3000/sub1/sub2/...subN
// The context contains the required sub-path to be resolved for every request, so we take it and rewrite values of appURL and appSubURL
// for both FrontendSettings and IndexViewData to pass this to frontend application
// It also takes to have corresponding location configured on reverse-proxy side (eq. ~ ^/[0-9a-z-]+/[0-9a-z-]+/(.*) for /sub1/sub2/ replacement)
// Note: this function completely ignores configured appSubURL value, so it's recommended not to use [server] section from config
// Note: appSubURL MUST start with "/", otherwise it will affect the frontend, forcing it to prefix url of every request with url of page that was open in browser
func BuildAppURLOverrides(ctx context.Context, appURL string) (string, string) {
	requestContextData := op_middleware.GetRequestContextData(ctx)
	return path.Join(appURL, requestContextData), path.Join("/", requestContextData)
}
