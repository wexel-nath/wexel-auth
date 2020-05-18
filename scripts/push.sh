#!/bin/bash
set -e

PROJECT_NAME="${PROJECT_NAME:-auth}"

push() {
	local image="wexel/$PROJECT_NAME:$1"
	docker inspect "$image"
	docker push "$image"
}

push api
push db-init
push keygen
