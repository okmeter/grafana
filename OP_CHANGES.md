# OP Changes

## Architecture

User being authenticated through custom auth proxy, which then adds the following headers:

- `X-WEBAUTH-USER` (https://grafana.com/docs/grafana/latest/setup-grafana/configure-security/configure-authentication/auth-proxy/)
- `X-WEBAUTH-ROLE` (https://grafana.com/docs/grafana/latest/setup-grafana/configure-security/configure-authentication/auth-proxy/)
- `X-REQUEST-CONTEXT`
- `X-USER-SESSION`

next, Grafana backend checks `X-WEBAUTH-USER` and `X-WEBAUTH-ROLE` headers (creates user with required permissions if no user/permissions found)  
then (by applied modifications):

1. adds `X-REQUEST-CONTEXT` and `X-USER-SESSION` headers to every upcoming request to OPStorage (get datasources, get dashboards..)  
   (so OPStorage can authenticate user and respond with user-related data)

2. while serving `Index` and `FrontendSettings` modifies corresponding DTOs by rewriting `AppURL` and `AppSubURL` to `AppURL/{X-REQUEST-CONTEXT}` and `/{X-REQUEST-CONTEXT}`  
   (so Grafana frontend pages can be accessed via `appURL/{X-REQUEST-CONTEXT}/` for example `http://grafana.uri/tenantUID/tenantSubUID`)

```text
location ~ ^/{X-REQUEST-CONTEXT}/(.*) {
    proxy_pass http://{GRAFANA-ENDPOINT}/$1;
}
               +---------+    X-WEBAUTH-USER         +---------+                          +---------+
               | Reverse |    X-WEBAUTH-ROLE         | Grafana |                          | OP      |
               | Proxy   +-------------------------->| Backend +------------------------->| Storage |
               |         |    X-REQUEST-CONTEXT      |         |   X-REQUEST-CONTEXT      |         |
               +---------+    X-USER-SESSION         +----+----+   X-USER-SESSION         +---------+
                    ^                                     |
                    |                                     |
                    |                    X-REQUEST-CONTEXT|
                    |                                     v
                    |                                +----------+
                    |                                | Grafana  |
                    |                                | Frontend |
                    |                                +----------+
                    |
                    +------------------- http://grafana.uri/{X-REQUEST-CONTEXT}/
```

### Dev Env

Quick start:

1. Start or port-forward OPStorage API
2. Execute:

```shell
OPSTORAGE_BASEURL="opStorageBaseURL" \
REQUEST_CONTEXT="{X-REQUEST-CONTEXT}" \
USER_ROLE="Editor" \
USER_SESSION="{X-USER-SESSION}" \
make op-develop -f op.mk
```

3. Open http://localhost:8080/{X-REQUEST-CONTEXT}/ page in your browser

Used environmental variables:

##### OPStorage (Grafana)

| Parameter         | Source                                              | Description               | Example                             |
|-------------------|-----------------------------------------------------|---------------------------|-------------------------------------|
| OPSTORAGE_BASEURL | Environment variable                                | OPStorage Base URL        | `http://host.docker.internal:10000` |
| OPSTORAGE_APIKEY  | Environment variable                                | OPStorage API Key         | apiKeyValue                         |
| USER_SESSION      | OP Middleware (Request Header: `X-USER-SESSION`)    | User Session cookie value | cookieValue                         |
| REQUEST_CONTEXT   | OP Middleware (Request Header: `X-REQUEST-CONTEXT`) | User Request context*     | `tenantRootUID/tenantSubRootUID`    |

*required for both authentication (OPStorage) and frontend routing (Grafana)

##### Grafana

| Parameter | Source                                                                            | Description                         | Example                       |
|-----------|-----------------------------------------------------------------------------------|-------------------------------------|-------------------------------|
| USER_ROLE | Grafana original Middleware (Request Headers: `X-WEBAUTH-ROLE`, `X-WEBAUTH-USER`) | Header required for authentication* | `Viewer`, `Editor` or `Admin` |

*`X-WEBAUTH-USER` can be different, but we use the role name for simplicity

##### Docker Compose

| Parameter     | Source                                                                     | Description            | Default                        | Example              |
|---------------|----------------------------------------------------------------------------|------------------------|--------------------------------|----------------------|
| GRAFANA_IMAGE | Service config (File `docker-compose.yml`)                                 | Local Docker image tag | op-grafana:develop             | `op-grafana:develop` |
| GRAFANA_IP    | Reverse proxy config (File `nginx.conf`)                                   | Grafana local IP*      | EVALUATED IN ENTRYPOINT SCRIPT | `172.17.0.1`         |
| GRAFANA_PORT  | Service and Reverse proxy config (Files `docker-compose.yml`, nginx.conf`) | Grafana local port     | 3030                           | `3030`               |

*required to support `host.docker.internal` for reverse proxy on both linux and macOS

## Updating Grafana from upstream

1. Lookup required grafana release at: https://github.com/grafana/grafana/releases
2. Clone required grafana release into custom directory `grafana-upstream`
```shell
git clone --depth=1 --branch $releaseTag https://github.com/grafana/grafana.git grafana-upstream
```
3. Remove .git related files from `grafana-upstream`
```shell
cd grafana-upstream
rm -rf .git
rm -f .gitignore
```
4. Checkout <this> grafana branch into two custom directories `grafana-update-branch-for-push` and `grafana-update-branch-for-diff`
5. Remove all codebase-related files and directories from directory `grafana-update-branch-for-push` (everything except `.helm`, ci, and `op`-prefixed files and directories from root)
6. Copy all `grafana-upstream` codebase-related files to `grafana-update-branch-for-push`
```shell
cp -a ~/Documents/grafana-upstream/. ~/Documents/grafana-update-branch-for-push/
```
7. List files that were previously affected inside `grafana-update-branch-for-diff` and apply corresponding changes inside: `grafana-update-branch-for-push`
```shell
make op-list-changes -f op.mk
```
```
./packages/grafana-data/src/themes/palette.ts:20:       // OP_CHANGES.md: added new color lightGray
./packages/grafana-data/src/themes/createColors.ts:217: // OP_CHANGES.md: added new color lightGray
./pkg/api/index.go:73:                                  // OP_CHANGES.md: override appURL and appSubURL for frontend
./pkg/api/frontendsettings.go:97:                       // OP_CHANGES.md: override appURL and appSubURL for frontend
...
```
8. Run `go mod vendor` inside `grafana-update-branch-for-push`
9. Launch Dev Env to check nothing is broken at the moment
10. Push applied changes to `grafana-update-branch-for-push` and create the merge request:
```shell
git add -f .
git commit -m "update version"
git push
```

## July 2023

### <a name="july2023Builds"></a> Builds

1. initial repository was cloned with currently running version commit:

```shell
git clone --depth=1 --branch v9.5.2 https://github.com/grafana/grafana.git
```

2. switched to `go mod vendor` to avoid long-running builds

3. following files and directories were added:

- `op.Dockerfile` (modified version of original Dockerfile to build both custom grafana backend/frontend)
- `op.mk` and `/op-develop` (files required to launch developer environment)

4. files that were modified:

- removed `/vendor` from `.gitignore` (to use it in Dockerfile)
- removed `/.github` to disable github actions and dependencies check
- added `/vendor` to `.gitattributes` (hide diffs in gitlab mrs)

### <a name="july2023Code"></a> Code

Added `op-pkg` package to both codebase and `op.Dockerfile` (require to be copied) with:

- service and store (to mimic internal logic with custom implementations)
- sdk (http sdk with client libraries and middlewares)
- opstorage (opstorage client library made with sdk)

#### Grafana internal codebase changes

API:

- `/pkg/api/http_server.go` (added authentication middlewares from `op-pkg/sdk`)
- `/pkg/api/accesscontrol.go` (added required rights for all dashboards and folders for `Viewer` and `Editor` by default)
- `/pkg/server/wire.go` (replaced original services requirements and stores with modified ones from `op-pkg`)

Services and service stores:

- `/pkg/services/datasources/service/datasource.go` (initial datasource Store implementation replacement)
- `/pkg/services/dashboards/database/database.go` (initial dashboard Store implementation replacement)
- `/pkg/services/folder/folderImpl/dashboard_folder_store.go` (initial dashboard Store implementation replacement)
- `/pkg/services/secrets/manager.go` (changes to use modified version of encryption service from `op-pkg` only)
- `/pkg/services/ngalert/api/util.go` (use op middlewares in alerting service)

Frontend:

- `/pkg/api/frontendsettings.go` (override `appURL` and `appSubURL` to use dynamic url sub-paths like `localhost:3000/sub1/sub2.../dashboards`)
- `/pkg/api/index.go` (override `appURL` and `appSubURL` to use dynamic url sub-paths like `localhost:3000/sub1/sub2.../dashboards`
- `/pkg/services/navtree/navtreeimpl/navtree.go` (disable navigation on alerting page)
- `/packages/grafana-data/src/themes/palette.ts` (add new color `lightGray`)
- `/packages/grafana-data/src/themes/createColors.ts` (override `background.canvas` color with `lightGray`)
- `/public/app/core/components/AppChrome/AppChrome.tsx` (override `searchBarHidden` with `true`, replace `NavToolbar` with actions, add css for actions, disable `NavToolbar` and `MegaMenu`)
- `/public/app/core/components/PageNew/Page.tsx` (disable `padding` in css)
- `/public/app/features/dashboard/components/DashNav/DashNav.tsx` (override `canStar`, `canShare` and `isStarred` with `false`)
- `/public/app/features/dashboard/components/SaveDashboard/forms/SaveDashboardAsForm.tsx` (disable `showRoot` in `FolderPicker`)
- `/public/app/features/alerting/routes.tsx` (removed unnecessary routes: `/alerting/ng/list`, `/alerting/notifications`, `/alerting/notifications/*`, `/alerting/notification/*`, `/alerting`, `/alerting/home`, `/alerting/routes*`, `/alerting/groups/`)
- `/public/app/features/alerting/unified/RuleList.tsx` (removed button export `url:/api/v1/provisioning/alert-rules/export`)
- `/public/app/features/alerting/unified/components/rule-editor/AlertRuleForm.tsx` (set type: `RuleFormType.cloudAlerting`)
- `/public/app/features/alerting/unified/components/rule-editor/query-and-alert-condition/AlertType.tsx` (set `RuleFormType.cloudAlerting` as defaultRuleType)
- `/public/app/features/alerting/unified/components/rule-editor/rule-types/RuleTypePicker.tsx` (disable selected `GrafanaManagedRuleType`)
- `/public/app/features/alerting/unified/components/rule-viewer/RuleViewerVisualization.tsx` (removed `view in Explore` button)
- `/public/app/features/alerting/unified/components/rules/RuleActionsButtons.tsx` (removed `copy to clipboard` button)
- `/public/app/features/alerting/unified/components/rules/RuleDetailsActionButtons.tsx` (removed `copy to clipboard` button)
- `/public/app/features/alerting/unified/hooks/useAlertManagerSourceName.ts` (use first available alert manager name)
- `/public/app/features/alerting/unified/utils/datasource.ts` (remove `grafanaAlertManagerDataSource` from available datasources)
- `/public/app/features/alerting/unified/utils/rule-form.ts` (removed logic related to `RuleFormType.grafana`)
- `/public/app/features/alerting/unified/Home.tsx` (removed `WelcomeCTABox`es)
- `/public/app/features/explore/ExploreToolbar.tsx` (hide `<AppChromeUpdate actions={[shareButton}}`)
- `/public/app/features/explore/RichHistory/RichHistoryCard.tsx` (`hide onCreateShortLink`)
