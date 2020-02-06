import sys

localDevelopTemplate = \
'''
version: "2"
services:  
  network-server:
    image: mxcdocker/supernode:network-server.v0.0.1-19-g030527a
    ports:
      - 8000:8000
    volumes:
      - ./configuration/loraserver:/etc/loraserver
  
  appserver:
    build:
      context: .
      dockerfile: Dockerfile-devel
    volumes:
      - ./configuration:/etc/lora-app-server
      - ./:/lora-app-server
    links:
      - postgresql
      - redis
      - mosquitto
      - network-server
    ports:
      - 8080:8080
    environment:
      - SUPERNODE_DATA_SERVICE=local
    env_file:
      - {}
    security_opt:
      - seccomp:unconfined
    cap_add:
      - SYS_PTRACE
    tty: true

  postgresql:
    image: postgres:9.6-alpine
    volumes:
      - ./.docker-compose/postgresql/initdb:/docker-entrypoint-initdb.d

  redis:
    image: redis:5-alpine

  mosquitto:
    image: eclipse-mosquitto
'''

remoteDevelopTemplate = \
'''
version: "2"
services:
  appserver:
    build:
      context: .
      dockerfile: Dockerfile-devel
    volumes:
      - ./configuration:/etc/lora-app-server
      - ./:/lora-app-server  
    links:
      - redis
    ports:
      - 8080:8080
    environment:
      - SUPERNODE_DATA_SERVICE=remote
    env_file:
      - {}
    security_opt:
      - seccomp:unconfined
    cap_add:
      - SYS_PTRACE
    tty: true
    
  redis:
    image: redis:5-alpine
'''

localTestingTemplate = \
'''
version: "3"

services:
  network-server:
    image: mxcdocker/supernode:network-server.v0.0.1-19-g030527a
    ports:
      - 8000:8000
    volumes:
      - ./configuration/loraserver:/etc/loraserver

  appserver:
    image: registry.gitlab.com/mxcfoundation/cloud/lpwan-app-server:latest
    volumes:
      - ./configuration/lora-app-server:/etc/lora-app-server
    environment:
      - SUPERNODE_DATA_SERVICE=local
    env_file:
      - {}
    ports:
      - 8080:8080
    depends_on:
      - postgresql
      - redis
      - mosquitto

  gatewaybridge:
    image: loraserver/lora-gateway-bridge:3
    ports:
      - 1700:1700/udp
    volumes:
      - ./configuration/lora-gateway-bridge:/etc/lora-gateway-bridge

  geoserver:
    image: loraserver/lora-geo-server:3
    volumes:
      - ./configuration/lora-geo-server:/etc/lora-geo-server

  postgresql:
    image: postgres:9.6-alpine
    env_file:
      - {}
    volumes:
      - ./configuration/postgresql/initdb:/docker-entrypoint-initdb.d
      - postgresqldata:/var/lib/postgresql/data

  redis:
    image: redis:5-alpine
    volumes:
      - redisdata:/data

  mosquitto:
    image: eclipse-mosquitto
    ports:
      - 1883:1883

  mxprotocol-server:
    image: mxcdocker/mxprotocol-server:0.3.0-14-ga8be3482
    volumes:
      - ./configuration/mxprotocol-server:/etc/mxprotocol-server
    environment:
      - SUPERNODE_DATA_SERVICE=local
    env_file:
      - {}
    depends_on:
      - postgresql
      - redis
      - mosquitto
    ports:
      - 4000:4000

volumes:
  postgresqldata:
  redisdata:
'''

imageTemplate = \
'''
version: "2"
services:
  appserver:
    build:
      context: .
      dockerfile: Dockerfile
    volumes:
      - ./configuration:/etc/lora-app-server
    ports:
      - 8080:8080
    environment:
      - APPSERVER=http://localhost:8080
      - MXPROTOCOL_SERVER=http://localhost:4000
'''

inputList = {
    "env_file": "default",
    "supernode_data_service": "default",
    "supernode_mode": "default",
}

if __name__ == "__main__":
    for item in sys.argv[1:]:
        if len(item.split('=')) != 2:
            print("invalid argument: ", item)
            exit()

        key, value = item.split('=')
        inputList[key] = value

    if (inputList["supernode_mode"] == "development") and (inputList["supernode_data_service"] == "local"):
        print(localDevelopTemplate.format(inputList["env_file"]))
        exit()

    elif (inputList["supernode_mode"] == "development") and (inputList["supernode_data_service"] == "remote"):
        print(remoteDevelopTemplate.format(inputList["env_file"]))
        exit()

    elif (inputList["supernode_mode"] == "testing") and (inputList["supernode_data_service"] == "local"):
        print(localTestingTemplate.format(inputList["env_file"], inputList["env_file"], inputList["env_file"]))
        exit()

    elif (inputList["supernode_mode"] == "image") and (inputList["supernode_data_service"] == "image"):
        print(imageTemplate)
        exit()

    else:
        print(inputList["supernode_mode"], inputList["supernode_data_service"])
        print("invalid argument")

