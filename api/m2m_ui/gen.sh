#!/usr/bin/env bash
GRPC_GW_PATH=`go list -f '{{ .Dir }}' github.com/grpc-ecosystem/grpc-gateway/protoc-gen-grpc-gateway`
GRPC_GW_PATH="${GRPC_GW_PATH}/../third_party/googleapis"

LS_PATH_APP=`go list -f '{{ .Dir }}' github.com/mxc-foundation/lpwan-app-server/api`
LS_PATH_APP="${LS_PATH_APP}"

# generate the gRPC code
protoc -I. -I${GRPC_GW_PATH} -I${LS_PATH_APP} --go_out=plugins=grpc:. \
    profile.proto \
    ext_account.proto \
    super_node.proto \
    topup.proto \
    wallet.proto \
    withdraw.proto \
    m2mserver_device.proto \
    m2mserver_gateway.proto \
    settings.proto \
    server.proto \
    staking.proto

# generate the JSON interface code
protoc -I. -I${GRPC_GW_PATH} -I${LS_PATH_APP} --grpc-gateway_out=logtostderr=true:. \
    profile.proto \
    ext_account.proto \
    super_node.proto \
    topup.proto \
    wallet.proto \
    withdraw.proto \
    m2mserver_device.proto \
    m2mserver_gateway.proto \
    server.proto \
    staking.proto \
    settings.proto

# generate the swagger definitions
protoc -I. -I${GRPC_GW_PATH} -I${LS_PATH_APP} --swagger_out=json_names_for_fields=true:./swagger \
    profile.proto \
    ext_account.proto \
    super_node.proto \
    topup.proto \
    wallet.proto \
    withdraw.proto \
    m2mserver_device.proto \
    m2mserver_gateway.proto \
    settings.proto \
    server.proto \
    staking.proto