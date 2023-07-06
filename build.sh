#!/bin/bash
# shellcheck disable=SC2164
go build -o ./bin/flydb-server cmd/server/cli/flydb-server.go
go build -o ./bin/flydb-client cmd/client/cli/flydb-client.go
echo "build success"

echo "Now you can run the follow command:
      start server: ./bin/flydb-server
      start client: ./bin/flydb-client 127.0.0.1:8999"

