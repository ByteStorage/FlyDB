#!/bin/bash

# 提示用户输入 Docker Hub 相关信息
# shellcheck disable=SC2162
read -p "Please enter your Docker Hub username: " DOCKER_USERNAME
# shellcheck disable=SC2162
read -p "Please enter image name: " IMAGE_NAME
# shellcheck disable=SC2162
read -p "Please enter image tag: " TAG

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

# check environment variables
if [ -z "$DOCKER_USERNAME" ] || [ -z "$IMAGE_NAME" ] || [ -z "$TAG" ]; then
    echo "Please set DOCKER_USERNAME, IMAGE_NAME and TAG environment variables"
    exit 1
fi

go build -o bin/flydb-server cmd/server/cli/flydb-server.go

echo "DOCKER_USERNAME: $DOCKER_USERNAME"
echo "IMAGE_NAME: $IMAGE_NAME"
echo "TAG: $TAG"

# build docker image
docker build -t "$DOCKER_USERNAME/$IMAGE_NAME:$TAG" -f docker/Dockerfile .

# docker login

# push docker image
docker push "$DOCKER_USERNAME/$IMAGE_NAME:$TAG"

# docker rmi $DOCKER_USERNAME/$IMAGE_NAME:$TAG
