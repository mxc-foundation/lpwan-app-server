# Gw-client
This is a provisioning-server client written in Go. It will only call APIs defined in api_gateway.proto :
```proto
syntax = "proto3";

package appserver_serves_gateway;

service HeartbeatService {
    rpc Heartbeat (HeartbeatRequest) returns (HeartbeatResponse);
}

// MiningRequest sends gateway list to m2m
message HeartbeatRequest {
    string gateway_mac = 1;
    string model = 2;
    string config_hash = 3;
    string os_version = 4;
    string statistics = 5;
}

message HeartbeatResponse {
    string new_firmware_link = 1;
    string config = 2;
}
``` 

## How to build the client
```bash
$ git clone git@gitlab.com:MXCFoundation/cloud/lpwan-app-server.git
$ cd internal/test/gw-client
$ make clean; make
$ ls build           
gw-client
$ ./build/gw-client help
  Usage:
     [flags]
     [command]
  
  Available Commands:
    heartbeat   Send heartbeat to supernode
    help        Help about any command
  
  Flags:
        --client-tls-certificate-path string   client TLS certificate
        --client-tls-key-path string           client TLS key
    -h, --help                                 help for this command
        --mac string                           mac address of gateway
        --model string                         model of gateway
        --os-version string                    os version of gateway
        --root-ca-path string                  rootCA
        --server string                        address of provisioning server
        --sn string                            serial number of gateway
  
  Use " [command] --help" for more information about a command.

```

## How to call API: Heartbeat
Simulate new gateway:
```bash
$ ./build/gw-client --server SERVER_ADDR:8004 \
  --root-ca-path PATH_TO_ECC_ROOT_CA_PEM \
  --client-tls-certificate-path PATH_TO_CLIENT_ECC_TLS_CRT \
  --client-tls-key-path PATH_TO_CLIENT_ECC_TLS_KEY \
  --mac MAC_ADDRESS --model MODEL --os-version OS_VERSION \
  heartbeat
```
Simulate old gateway:
```bash
$ ./build/gw-client --server SERVER_ADDR:8005 \
  --root-ca-path PATH_TO_ECC_ROOT_CA_PEM \
  --client-tls-certificate-path PATH_TO_CLIENT_ECC_TLS_CRT \
  --client-tls-key-path PATH_TO_CLIENT_ECC_TLS_KEY \
  --mac MAC_ADDRESS --model MODEL --os-version OS_VERSION \
  heartbeat
```