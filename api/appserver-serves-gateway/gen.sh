#!/usr/bin/env bash

GRPC_PATH=`go list -f '{{ .Dir }}' github.com/grpc-ecosystem/grpc-gateway/protoc-gen-grpc-gateway`
GRPC_PATH="${GRPC_PATH}/../third_party/googleapis"

# generate the gRPC code
protoc -I. -I${GRPC_PATH} --go_out=paths=source_relative,plugins=grpc:. \
    heartbeat.proto
