# LPWAN App Server

[![CircleCI](https://circleci.com/gh/brocaar/lora-app-server.svg?style=svg)](https://circleci.com/gh/brocaar/lora-app-server)

LPWAN App Server is an open-source LoRaWAN application-server, part of the
[LPWAN Server](https://www.loraserver.io/) project. It is responsible
for the node "inventory" part of a LoRaWAN infrastructure, handling of received
application payloads and the downlink application payload queue. It comes
with a web-interface and API (RESTful JSON and gRPC) and supports authorization
by using JWT tokens (optional). Received payloads are published over MQTT
and payloads can be enqueued by using MQTT or the API.

## Architecture

![architecture](https://www.loraserver.io/img/architecture.png)

### Component links

* [LPWAN Gateway Bridge](https://www.loraserver.io/lora-gateway-bridge)
* [LPWAN Gateway Config](https://www.loraserver/lora-gateway-config)
* [LPWAN Server](https://www.loraserver.io/loraserver/)
* [LPWAN App Server](https://www.loraserver.io/lora-app-server/)

## Links

* [Downloads](https://www.loraserver.io/lora-app-server/overview/downloads/)
* [Docker image](https://hub.docker.com/r/loraserver/lora-app-server/)
* [Documentation & screenshots](https://www.loraserver.io/lora-app-server/) and [Getting started](https://www.loraserver.io/lora-app-server/getting-started/)
* [Building from source](https://www.loraserver.io/lora-app-server/community/source/)
* [Contributing](https://www.loraserver.io/lora-app-server/community/contribute/)
* Support
  * [Support forum](https://forum.loraserver.io)
  * [Bug or feature requests](https://github.com/mxc-foundation/lpwan-app-server/issues)

## License

LPWAN App Server is distributed under the MIT license. See also
[LICENSE](https://github.com/mxc-foundation/lpwan-app-server/blob/master/LICENSE).
