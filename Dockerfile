FROM golang:1.12-alpine AS development

ENV PROJECT_PATH=/lora-app-server
ENV PATH=$PATH:$PROJECT_PATH/build
ENV CGO_ENABLED=0
ENV GO_EXTRA_BUILD_ARGS="-a -installsuffix cgo"

RUN apk add --no-cache ca-certificates make git bash protobuf alpine-sdk nodejs nodejs-npm python

RUN mkdir -p $PROJECT_PATH
COPY . $PROJECT_PATH
WORKDIR $PROJECT_PATH

RUN make dev-requirements ui-requirements clean ui/build_dep ui/build build

FROM alpine:latest AS production

WORKDIR /root/
RUN apk --no-cache add ca-certificates
RUN mkdir /etc/lora-app-server
COPY --from=development /lora-app-server/build/ .
COPY --from=development /lora-app-server/configuration/ .
COPY --from=development /lora-app-server/scripts/init .
RUN ["chmod", "+x", "./start"]
ENTRYPOINT ["./start"]
