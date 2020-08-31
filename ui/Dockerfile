FROM golang:1.15-alpine AS development

ENV PROJECT_PATH=/web-ui
ENV PATH=$PATH:$PROJECT_PATH/build
RUN apk add nodejs nodejs-npm python3

RUN mkdir -p $PROJECT_PATH
COPY . $PROJECT_PATH
WORKDIR $PROJECT_PATH

RUN npm install
RUN npm run build

FROM alpine:latest AS production

WORKDIR /root/
RUN apk add nodejs
RUN mkdir ./build ./node_modules
COPY --from=development /web-ui/build/ ./build
COPY --from=development /web-ui/node_modules ./node_modules
COPY --from=development /web-ui/server.js .
COPY --from=development /web-ui/start .
RUN ["chmod", "+x", "./start"]
ENTRYPOINT ["./start"]