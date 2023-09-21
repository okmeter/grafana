op-develop: op-build op-compose

## build grafana image

GRAFANA_IMAGE := "op-grafana:develop"
GRAFANA_PORT  := "3030"

op-build:
	docker build \
	-f op.Dockerfile \
	--build-arg COMMIT_SHA=$$(git rev-parse --short HEAD) \
	--build-arg BUILD_BRANCH=$$(git rev-parse --abbrev-ref HEAD) \
	--tag ${GRAFANA_IMAGE} .

## start grafana

op-compose:
	GRAFANA_IMAGE=${GRAFANA_IMAGE} \
	GRAFANA_PORT=${GRAFANA_PORT} \
	OPSTORAGE_BASEURL=${OPSTORAGE_BASEURL} \
	OPSTORAGE_APIKEY=${OPSTORAGE_APIKEY} \
	docker compose -f op-develop/docker-compose.yml up

## update grafana version

op-list-changes:
	grep -rnw '.' -e 'OP_CHANGES.md'
