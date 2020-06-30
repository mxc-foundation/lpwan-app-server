#!/usr/bin/env bash

# generate the gRPC code
protoc -I. --go_out=plugins=grpc:. \
  heartbeat.proto
