#!/usr/bin/env bash

GRPC_GW_PATH=`go list -f '{{ .Dir }}' github.com/grpc-ecosystem/grpc-gateway/protoc-gen-grpc-gateway`
GRPC_GW_PATH="${GRPC_GW_PATH}/../third_party/googleapis"

LS_PATH=`go list -f '{{ .Dir }}' github.com/mxc-foundation/lpwan-app-server/api`
LS_PATH="${LS_PATH}/m2m_ui/"

# generate the gRPC code
protoc -I=. -I${GRPC_GW_PATH} -I${LS_PATH} --go_out=paths=source_relative,plugins=grpc:. \
    appserver.proto