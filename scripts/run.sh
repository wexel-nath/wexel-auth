#!/bin/sh
set -e

PROJECT_NAME="${PROJECT_NAME:-auth}"
NETWORK_NAME="${NETWORK_NAME:-wexel}"
[ -z $(docker network ls --filter name=^${NETWORK_NAME}$ --format="{{ .Name }}") ] \
	&& docker network create "$NETWORK_NAME"

run() {
	COMPOSE_FILE='docker/docker-compose.yml:docker/docker-compose.dev.yml' \
	COMPOSE_PROJECT_NAME="$PROJECT_NAME" \
		docker-compose up --force-recreate -d "$@"
}

run api db-init keygen
