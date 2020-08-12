#!/usr/bin/env bash

PROTOBUF_PATH=$(go list -f '{{ .Dir }}' github.com/golang/protobuf/ptypes)

# generate the gRPC code
protoc -I. -I.. -I${PROTOBUF_PATH} --go_out=plugins=grpc:. \
  m2m_device.proto \
  m2m_gateway.proto \
  notification.proto
