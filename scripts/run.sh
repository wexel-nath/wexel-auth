#!/bin/bash
set -e

PROJECT_NAME="${PROJECT_NAME:-auth}"

run() {
	DATABASE_URL='postgresql://nathanw:bonnie@db:5432/auth?sslmode=disable' \
	DB_HOST='db' \
	DB_PASS='bonnie' \
	DB_USER='nathanw' \
	docker stack deploy \
		--compose-file 'docker/docker-stack.yml' \
		--compose-file 'docker/docker-stack.dev.yml' \
		"$PROJECT_NAME"
}

run
