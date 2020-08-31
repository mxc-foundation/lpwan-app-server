import sys

localDevelopTemplate = \
    '''
    version: "2"
    services:  
      network-server:
        build:
          context: ../lpwan-server
          dockerfile: Dockerfile-devel
        volumes:
          - ../lpwan-server/configuration:/etc/network-server
          - ../lpwan-server:/network-server
        tty: true
          
      mxprotocol-server:
        build:
          context: ../mxprotocol-server
          dockerfile: Dockerfile-devel
        volumes:
          - ../mxprotocol-server/configuration:/etc/mxprotocol-server
          - ../mxprotocol-server:/mxprotocol-server
        links:
          - network-server
          - postgresql
          - redis
          - mosquitto
        environment:
          - APPSERVER=http://localhost:8080
          - MXPROTOCOL_SERVER=http://localhost:4000
        tty: true
        ports:
          - 4000:4000
        security_opt:
          - seccomp:unconfined
        cap_add:
          - SYS_PTRACE
      
      gatewaybridge:
        image: mxcdocker/chirpstack-gateway-bridge
        ports:
          - 1700:1700/udp
        volumes:
          - ./configuration/chirpstack-gateway-bridge:/etc/chirpstack-gateway-bridge
        restart: always
    
      geoserver:
        image: chirpstack/chirpstack-geolocation-server:3
        volumes:
          - ./configuration/chirpstack-geolocation-server:/etc/chirpstack-geolocation-server
        restart: always 
          
      appserver:
        build:
          context: .
          dockerfile: Dockerfile-devel
        volumes:
          - ./configuration:/etc/lora-app-server
          - ./:/lora-app-server
        links:
          - network-server
          - mxprotocol-server
          - postgres
          - redis
          - mosquitto
          - rabbitmq
          - zookeeper
          - kafka
        ports:
          - 8080:8080
          - 8004:8004
          - 8005:8005
        environment:
          - SUPERNODE_DATA_SERVICE=local
          - TEST_POSTGRES_DSN=postgres://chirpstack_as:chirpstack_as@postgres/chirpstack_as?sslmode=disable
          - TEST_REDIS_URL=redis://redis:6379
          - TEST_MQTT_SERVER=tcp://mosquitto:1883
          - TEST_RABBITMQ_URL=amqp://guest:guest@rabbitmq:5672/
          - TEST_KAFKA_BROKER=kafka:9092
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
        environment:
          - POSTGRES_PASSWORD=local_superuser_pass
    
      redis:
        image: redis:5-alpine
    
      mosquitto:
        image: eclipse-mosquitto
    
      postgres:
        image: postgres:9.6-alpine
        environment:
          - POSTGRES_HOST_AUTH_METHOD=trust
        volumes:
          - ./.docker-compose/postgresql/initdb:/docker-entrypoint-initdb.d
    
      rabbitmq:
        image: rabbitmq:3-alpine
    
      zookeeper:
        image: 'bitnami/zookeeper:3'
        environment:
          - ALLOW_ANONYMOUS_LOGIN=yes
          
      kafka:
        image: 'bitnami/kafka:2'
        environment:
          - KAFKA_CFG_ZOOKEEPER_CONNECT=zookeeper:2181
          - ALLOW_PLAINTEXT_LISTENER=yes
        depends_on:
          - zookeeper
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
          - 8004:8004
          - 8005:8005
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
        image: mxcdocker/supernode:network-server.2.0.0-6-g956f51f
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
          - 8004:8004
          - 8005:8005
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
        environment:
          - POSTGRES_PASSWORD=local_superuser_pass
    
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
          - 8004:8004
          - 8005:8005
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
