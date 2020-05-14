#!/bin/sh
set -e

PROJECT_NAME="${PROJECT_NAME:-auth}"

run() {
	COMPOSE_FILE='docker/docker-compose.yml:docker/docker-compose.dev.yml' \
	COMPOSE_PROJECT_NAME="$PROJECT_NAME" \
		docker-compose up --force-recreate -d "$@"
}

run api db-init
