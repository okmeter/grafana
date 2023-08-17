#!/bin/sh
GRAFANA_IP=$(cat /etc/hosts | grep "host.docker.internal" | tr -d "[:space:]" | sed -e "s|host.docker.internal||")
export GRAFANA_IP
envsubst '$$GRAFANA_IP $$GRAFANA_PORT $$USER_ROLE $$REQUEST_CONTEXT $$USER_SESSION' < /etc/nginx/conf.d/templates/default.conf > /etc/nginx/conf.d/default.conf
cat /etc/nginx/conf.d/default.conf
exec nginx -g "daemon off;"
