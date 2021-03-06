import sys

localDevelopTemplate = \
    '''
    version: "3"
    services:  
      web-ui:
        image: mxcdocker/webui:latest  
        container_name: web-ui
        ports:
          - 3001:3001
        restart: always
      
      appserver:
        build:
          context: .
          dockerfile: Dockerfile-devel
        volumes:
          - ./configuration:/etc/lora-app-server
          - ./:/lora-app-server
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
        
      network-server:
        build:
          context: ../../lpwan-server
          dockerfile: Dockerfile
        volumes:
          - ../../lpwan-server/configuration:/etc/loraserver
          - ../../lpwan-server:/network-server
        tty: true
          
      mxprotocol-server:
        build:
          context: ../../mxprotocol-server
          dockerfile: Dockerfile-devel
        volumes:
          - ../../mxprotocol-server/configuration:/etc/mxprotocol-server
          - ../../mxprotocol-server/configuration/ecc:/etc/ecc
          - ../../mxprotocol-server:/mxprotocol-server
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

      dhx-center:
        build:
          context: ../../dhx-center
          dockerfile: Dockerfile-devel
        volumes:
          - ../../dhx-center/configuration:/etc/dhx-center
        tty: true
      
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

      postgresql:
        image: postgres:9.6-alpine
        volumes:
          - ./.docker-compose/postgresql/initdb:/docker-entrypoint-initdb.d
        environment:
          - POSTGRES_PASSWORD=local_superuser_pass
    
      redis:
        image: redis:5-alpine
    
      mosquitto:
        build:
          context: ../../mosquitto-with-auth-plug
          dockerfile: Dockerfile
        container_name: mosquitto
        volumes:
          - ../../mosquitto-with-auth-plug/configuration/:/etc/mosquitto
        ports:
          - 1883:1883
          - 8883:8883
        restart: always
    
      rabbitmq:
        image: rabbitmq:3-alpine
      
      testsmtpd:
        image: registry.gitlab.com/shaydo/testsmtpd
        ports:
          - 8100:80

      influxdb:
        image: influxdb:1.8
        environment:
          - INFLUXDB_DB=mxprotocol
          - INFLUXDB_ADMIN_USER=admin
          - INFLUXDB_ADMIN_PASSWORD=admin
          - INFLUXDB_USER=influx
          - INFLUXDB_USER_PASSWORD=influx
    '''

inputList = {
    "env_file": "default",
    "supernode_mode": "default",
}

if __name__ == "__main__":
    for item in sys.argv[1:]:
        if len(item.split('=')) != 2:
            print("invalid argument: ", item)
            exit()

        key, value = item.split('=')
        inputList[key] = value

    if (inputList["supernode_mode"] == "development"):
        print(localDevelopTemplate.format(inputList["env_file"]))
        exit()

    else:
        print("invalid argument")
