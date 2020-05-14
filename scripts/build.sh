#!/bin/bash
set -e

PROJECT_NAME="${PROJECT_NAME:-auth}"

build() {
	dockerfile="./docker/Dockerfile.$1"

	docker build \
		-t "wexel/$PROJECT_NAME:$1" \
		-f "$dockerfile" \
		.
}

build api
build db-init
