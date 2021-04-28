#!/usr/bin/env bash

GRPC_GW_PATH=$(go list -f '{{ .Dir }}' github.com/grpc-ecosystem/grpc-gateway/protoc-gen-grpc-gateway)
GRPC_GW_PATH="${GRPC_GW_PATH}/../third_party/googleapis"

# generate the gRPC code
protoc -I. -I${GRPC_GW_PATH} --go_out=paths=source_relative,plugins=grpc:. \
  device.proto \
  application.proto \
  deviceQueue.proto \
  as_common.proto \
  user.proto \
  gateway.proto \
  organization.proto \
  as_profiles.proto \
  networkServer.proto \
  serviceProfile.proto \
  deviceProfile.proto \
  gatewayProfile.proto \
  multicastGroup.proto \
  fuotaDeployment.proto \
  internal.proto \
  serverInfo.proto \
  topup.proto \
  wallet.proto \
  withdraw.proto \
  settings.proto \
  staking.proto \
  dhx.proto \
  external_user.proto \
  shopify_integration.proto \
  dfi_service.proto \
  download.proto

# generate the JSON interface code
protoc -I. -I${GRPC_GW_PATH} --grpc-gateway_out=paths=source_relative,logtostderr=true:. \
  device.proto \
  application.proto \
  deviceQueue.proto \
  as_common.proto \
  user.proto \
  gateway.proto \
  organization.proto \
  as_profiles.proto \
  networkServer.proto \
  serviceProfile.proto \
  deviceProfile.proto \
  gatewayProfile.proto \
  multicastGroup.proto \
  fuotaDeployment.proto \
  internal.proto \
  serverInfo.proto \
  topup.proto \
  wallet.proto \
  withdraw.proto \
  settings.proto \
  staking.proto \
  dhx.proto \
  external_user.proto \
  shopify_integration.proto \
  dfi_service.proto \
  download.proto

# generate the swagger definitions
protoc -I. -I${GRPC_GW_PATH} --swagger_out=json_names_for_fields=true,simple_operation_ids=true:./swagger \
  device.proto \
  application.proto \
  deviceQueue.proto \
  as_common.proto \
  user.proto \
  gateway.proto \
  organization.proto \
  as_profiles.proto \
  networkServer.proto \
  serviceProfile.proto \
  deviceProfile.proto \
  gatewayProfile.proto \
  multicastGroup.proto \
  fuotaDeployment.proto \
  internal.proto \
  serverInfo.proto \
  topup.proto \
  wallet.proto \
  withdraw.proto \
  settings.proto \
  staking.proto \
  dhx.proto \
  external_user.proto \
  shopify_integration.proto \
  dfi_service.proto \
  download.proto

# merge the swagger code into one file
#go run swagger/main.go swagger > ../static/swagger/api.swagger.json
