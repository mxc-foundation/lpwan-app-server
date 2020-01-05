# LPWAN App Server

# Setup

See MXC Developer Handbook for further information.

Note: UI part from m2m has been merged into lpwan-app-server, m2m no longer contains UI part.  
However part of the APIs get data from m2m service, you need to start m2m service for accessing all features correctly.

## Environment

#### Set up docker
- [Install Docker](https://docs.docker.com/install/linux/docker-ce/ubuntu/)  
Just follow __Install using the repository / SET UP THE REPOSITORY__, no need to install docker engine community

- [Install docker-compose](https://docs.docker.com/compose/install/)
Just follow __Install Compose on Linux systems__

- Add user to docker group
```bash
$ sudo usermod -aG docker $USER
```

## Clone the repo:

```bash
git clone git@gitlab.com:MXCFoundation/cloud/lpwan-app-server.git &&
cd lpwan-app-server
```

## Fetch latest develop branch:

```bash
git fetch origin develop:develop &&
git checkout develop &&
git pull --rebase origin develop
```

## Existing or new feature branch

* New feature branch required?

```
git checkout -b feature/MCL-XXX
```

* Existing feature branch?

> Example: If there is a "feature" branch that you are working on in Jira
(i.e. feature/MCL-117) and you are working on a task of that feature,
then create a branch from that feature that is prefixed with your name
(i.e. luke/MCL-118-page-network-servers)

```bash
git fetch origin feature/MCL-117:feature/MCL-117 &&
git checkout feature/MCL-117 &&
git pull --rebase origin feature/MCL-117
```

## Create task branch from feature branch:

```bash
git checkout -b luke/MCL-118-page-network-servers
```

## Install dependencies:

```bash
cd ui/ &&
npm install
```

## Build Docker container and start container shell session:

If you want to use __local__ postgresql and mqtt service, do following command in directory where Makefile is:
```bash
$ make server_local
```

If you want to use __remote__ postgresql and mqtt service, do following command in directory where Makefile is, and insert remote server domain name after the prompt:
```bash
$ make server_remote
Start docker container with remote database and mqtt service
Insert remote server domain name: 

```

## Start LPWAN App Server:

```bash
make ui-requirements &&
make dev-requirements &&
make clean &&
make build &&
./build/lora-app-server
```

**HACK**
If it then gives a `Failed to compile` error as shown below:
```
Failed to compile.

./src/assets/scss/DefaultTheme.scss
Error: Missing binding /lora-app-server/ui/node_modules/node-sass/vendor/linux_musl-x64-64/binding.node
Node Sass could not find a binding for your current environment: Linux/musl 64-bit with Node.js 10.x

Found bindings for the following environments:
  - OS X 64-bit with Node.js 12.x

This usually happens because your environment has changed since running `npm install`.
Run `npm rebuild node-sass` to download the binding for your current environment.
```

Then keep the Docker container running,
and outside the Docker container, in terminal run:

```bash
cd ui/ &&
cd node_modules/node-sass &&
sudo npm install &&
cd ../../../
```

Then back in the Docker container run the following commands again, and it should compile successfully and run:

```
make build &&
./build/lora-app-server
```

Open web browser at: http://localhost:8080

Enter credentials to login: admin, password: admin

See below how to enable debugging and live reload.

## Frequently apply latest from feature branch into your task branch:

```
git checkout feature/MCL-117 &&
git pull --rebase origin feature/MCL-117
git checkout luke/MCL-118-page-network-servers
git rebase -i feature/MCL-117
```

## Debugging with live reload:

After the LPWAN App Server is built and running from the Docker container,
if you just go to http://localhost:8080, then you won't get live reload support.
So to enable debugging and live reload, additionally run the following (outside the
Docker container):

```bash
cd lpwan-app-server &&
cd ui/ &&
npm start
```

Then open in your web browser: http://localhost:3000

Now when you make changes it will automatically refresh

## Configuration

##### - redirect database
For sharing testing data during development, set postgresql service server wherever it is needed.
Change in configuration/lora-app-server.toml
```toml
[postgresql]
dsn="postgres://USERNAME:PASSWORD@SERVICE_SERVER_DOMAIN_NAME:5432/DATABASE_NAME?sslmode=disable"
```

After changing config file, simply restart the service in docker container again

```bash
$ ./build/lora-app-server -c configuration/lora-app-server.toml
```

# Intro

LPWAN App Server is an open-source LoRaWAN application-server, part of the
[LPWAN Server](https://www.loraserver.io/) project. It is responsible
for the node "inventory" part of a LoRaWAN infrastructure, handling of received
application payloads and the downlink application payload queue. It comes
with a web-interface and API (RESTful JSON and gRPC) and supports authorization
by using JWT tokens (optional). Received payloads are published over MQTT
and payloads can be enqueued by using MQTT or the API.

## Architecture


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
