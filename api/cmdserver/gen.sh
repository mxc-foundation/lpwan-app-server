#!/usr/bin/env bash

PROTOBUF_PATH=`go list -f '{{ .Dir }}' github.com/golang/protobuf/ptypes`

# generate the gRPC code
protoc -I. -I${PROTOBUF_PATH} --go_out=paths=source_relative,plugins=grpc:. \
  ensure_default.proto
