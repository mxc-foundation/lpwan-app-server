FROM registry.gitlab.com/mxcfoundation/cloud/lpwan-app-server:baseimage AS development

ENV PROJECT_PATH=/lora-app-server
ENV PATH=$PATH:$PROJECT_PATH/build
ENV CGO_ENABLED=0
ENV GO_EXTRA_BUILD_ARGS="-a -installsuffix cgo"

COPY . $PROJECT_PATH
WORKDIR $PROJECT_PATH

RUN make clean ui/build build

FROM alpine:latest AS production

WORKDIR /root/
RUN apk --no-cache add ca-certificates
RUN mkdir /etc/lora-app-server
COPY --from=development /lora-app-server/build/ .
COPY --from=development /lora-app-server/configuration/ .
COPY --from=development /lora-app-server/scripts/init .
RUN ["chmod", "+x", "./start"]
ENTRYPOINT ["./start"]
