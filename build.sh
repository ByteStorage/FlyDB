#!/bin/bash

export CGO_ENABLED=0
export GOPATH=$(go env GOPATH)

# shellcheck disable=SC2164
go build -gcflags=-trimpath=$GOPATH -asmflags=-trimpath=$GOPATH -ldflags "-s -w -extldflags '-static'" -o ./bin/flydb-server cmd/server/cli/flydb-server.go
go build  -o ./bin/flydb-client cmd/client/cli/flydb-client.go
echo "build success"

echo "Now you can run the follow command:
      start server: ./bin/flydb-server
      start client: ./bin/flydb-client 127.0.0.1:8999"