#!/bin/bash
# shellcheck disable=SC2164
cd cmd/cli
go build flydb.go
./flydb

