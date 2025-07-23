#!/bin/bash
# Export go binary path (put it first in PATH to override system version)
export PATH=$(go env GOPATH)/bin:$PATH

# Generate Go code from pb files
# First generate user.proto
protoc --go_out=../grpc --go-grpc_out=../grpc --proto_path=. user.proto

# Then generate todo.proto with gRPC service
protoc --go_out=../grpc --go-grpc_out=../grpc --proto_path=. todo.proto