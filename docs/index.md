# LPWAN App Server

# Content
[Setup](#Start LPWAN App Server)

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

## Stop any existing processes in other terminal tabs

Check that you haven't left any existing instances of lpwan-app-server running (i.e. if you're already running the UI with `npm start`).

## Change to compatible version of Node.js

Since the Docker container will use Node.js v10.16.3, switch to that version on your local machine.
[Node Version Manager](https://github.com/nvm-sh/nvm#install--update-script) makes this convenient to do:

```bash
nvm install v10.16.3 &&
nvm use v10.16.3
```

## Install dependencies

We'll install dependencies as follows **OUTSIDE** the Docker container:

```bash
cd lpwan-app-server/ &&
cd ui/ &&
rm package-lock.json &&
rm -rf ./node_modules/node-sass/ &&
rm -rf ~/.npm/node-sass &&
npm install &&
cd ./node_modules/node-sass &&
npm install &&
cd ../../ &&
npm rebuild node-sass --force &&
rm package-lock.json &&
cd ../
```

> Note: The above strange sequence of steps should avoid encountering error `npm ERR! Cannot read property 'match' of undefined` or `Failed to compile. ... Node Sass could not find a binding for your current environment ...` by removing the node-sass from cache and re-installing it after installing other dependencies, and then removing the package-lock.json again to avoid encountering the same error in the Docker container. See https://stackoverflow.com/questions/37986800/node-sass-couldnt-find-a-binding-for-your-current-environment. If you accidently ran `yarn`, then remove `yarn.lock` too (i.e. `rm yarn.lock`), otherwise you'll encounter errors due to conflicts between package-lock.json and yarn.lock.

## Build Docker container and start container shell session:

If you want to use __local__ postgresql and mqtt service, do following command in directory where Makefile is:
```bash
$ make server_local
```

If you want to use __remote__ postgresql and mqtt service, do following command in directory where Makefile is, and insert remote server domain name after the prompt:
```bash
$ make server_remote
Start docker container with remote database and mqtt service
Insert remote server domain name (excluding the port): 

```

## Start LPWAN App Server:

Run the following **INSIDE** the Docker container that was started from the previous step:

```bash
make ui-requirements &&
make dev-requirements &&
make clean &&
make build &&
./build/lora-app-server
```

Then keep the Docker container running,
and **OUTSIDE** the Docker container, in terminal run:

```bash
cd lpwan-app-server/ &&
cd ui/ &&
cd node_modules/node-sass &&
npm install &&
cd ../../ &&
npm start
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

> Note: An alternative to running `git rebase -i feature/MCL-117` is to merge and resolve conflicts instead with `git merge feature/MCL-117` 

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

## Development

### Library Requirements

All libraries used in the UI should provide React Native support

### Database Access

Try using a PostgreSQL GUI to easily resolve issues in the test database

Example:
* Download http://www.psequel.com/
* Enter connection details that are either in your .env file, or in your custom configuration file: configuration/lora-app-server.toml

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
