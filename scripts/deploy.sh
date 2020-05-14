#!/bin/bash
set -e

PROJECT_NAME="${PROJECT_NAME:-auth}"

compose() {
	COMPOSE_FILE='docker/docker-compose.yml:docker/docker-compose.prod.yml' \
	COMPOSE_PROJECT_NAME="$PROJECT_NAME" \
	DOCKER_HOST='ssh://ec2-user@ec2-13-210-184-139.ap-southeast-2.compute.amazonaws.com' \
		docker-compose "$@"
}

compose pull api db-init
compose up -d api db-init
