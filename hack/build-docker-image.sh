#!/bin/bash

DOCKER_USERNAME="$DOCKER_USERNAME"
IMAGE_NAME="$IMAGE_NAME"
TAG="$TAG"

if [ -z "$DOCKER_USERNAME" ] || [ -z "$IMAGE_NAME" ] || [ -z "$TAG" ]; then
    echo "Please set DOCKER_USERNAME, IMAGE_NAME and TAG environment variables"
    exit 1
fi

docker build -t $DOCKER_USERNAME/$IMAGE_NAME:$TAG .

# docker login

docker push $DOCKER_USERNAME/$IMAGE_NAME:$TAG

# docker rmi $DOCKER_USERNAME/$IMAGE_NAME:$TAG
