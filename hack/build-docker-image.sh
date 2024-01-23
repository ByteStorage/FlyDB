#!/bin/bash

DOCKER_USERNAME="${DOCKER_USERNAME:-bytestorage}"
IMAGE_NAME="${IMAGE_NAME:-flydb}"
TAG="${TAG:-latest}"

# check docker command
if ! [ -x "$(command -v docker)" ]; then
    echo "Docker is not installed"
    exit 1
fi

# check go command
if ! [ -x "$(command -v go)" ]; then
    echo "Go is not installed"
    exit 1
fi

if [ -z "$IMAGE_NAME" ] || [ -z "$TAG" ]; then
    echo "Please set DOCKER_USERNAME, IMAGE_NAME and TAG environment variables"
    exit 1
fi

go build -o bin/flydb-server cmd/server/cli/flydb-server.go

echo "IMAGE_NAME: $IMAGE_NAME"
echo "TAG: $TAG"

# build docker image
docker build -t "$DOCKER_USERNAME/$IMAGE_NAME:$TAG" -f docker/Dockerfile .

