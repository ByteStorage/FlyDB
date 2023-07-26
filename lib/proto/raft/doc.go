package raftPB

// note: run ```go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28```
// note: run ```go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2```
//go:generate protoc --go_out=. --go-grpc_out=. ./raft.proto
