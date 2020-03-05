#!/usr/bin/env bash

GRPC_GW_PATH=`go list -f '{{ .Dir }}' github.com/grpc-ecosystem/grpc-gateway/protoc-gen-grpc-gateway`
GRPC_GW_PATH="${GRPC_GW_PATH}/../third_party/googleapis"

LS_PATH_NS=`go list -f '{{ .Dir }}' github.com/mxc-foundation/lpwan-server/api/ns`
LS_PATH_NS="${LS_PATH_NS}/../.."

PROTOBUF_PATH=`go list -f '{{ .Dir }}' github.com/golang/protobuf/ptypes`

# generate the gRPC code
protoc -I. -I${LS_PATH_NS} -I${GRPC_GW_PATH} -I${PROTOBUF_PATH} --go_out=plugins=grpc:. \
    device.proto \
    application.proto \
    deviceQueue.proto \
    common.proto \
    user.proto \
    gateway.proto \
    organization.proto \
    profiles.proto \
    networkServer.proto \
    serviceProfile.proto \
    deviceProfile.proto \
    gatewayProfile.proto \
    multicastGroup.proto \
	  fuotaDeployment.proto \
    internal.proto \
    serverInfo.proto \
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
protoc -I. -I${LS_PATH_NS} -I${GRPC_GW_PATH} -I${PROTOBUF_PATH} --grpc-gateway_out=logtostderr=true:. \
    device.proto \
    application.proto \
    deviceQueue.proto \
    common.proto \
    user.proto \
    gateway.proto \
    organization.proto \
    profiles.proto \
    networkServer.proto \
    serviceProfile.proto \
    deviceProfile.proto \
    gatewayProfile.proto \
    multicastGroup.proto \
	  fuotaDeployment.proto \
    internal.proto \
    serverInfo.proto \
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

# generate the swagger definitions
protoc -I. -I${LS_PATH_NS} -I${GRPC_GW_PATH} -I${PROTOBUF_PATH} --swagger_out=json_names_for_fields=true:./swagger \
    device.proto \
    application.proto \
    deviceQueue.proto \
    common.proto \
    user.proto \
    gateway.proto \
    organization.proto \
    profiles.proto \
    networkServer.proto \
    serviceProfile.proto \
    deviceProfile.proto \
    gatewayProfile.proto \
    multicastGroup.proto \
	  fuotaDeployment.proto \
    internal.proto \
    serverInfo.proto \
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

# merge the swagger code into one file
#go run swagger/main.go swagger > ../static/swagger/api.swagger.json
