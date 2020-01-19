ROJECT_NAME := "lpwan-app-server"
PKG := "."
PKG_LIST := $(shell go list ${PKG}/... | grep -v /vendor/)
GO_FILES := $(shell find . -name '*.go' | grep -v /vendor/ | grep -v _test.go)
VERSION := $(shell git describe --tags |sed -e "s/^v//")

.PHONY: all dep build ui/test ui/build_dep ui/built clean test coverage coverhtml lint

all: ui/build build

lint: ## Lint the files
	@go get -u golang.org/x/lint/golint
	@golint -set_exit_status ${PKG_LIST}

test: ## Run unittests
	@go test -timeout 60s -short ${PKG_LIST}

race: dep ## Run data race detector
	@go test -timeout 60s -race -short ${PKG_LIST}

msan: dep ## Run memory sanitizer
	@apk add --no-cache clang
	@export CC=clang
	@go test -timeout 60s -msan -short ${PKG_LIST}

coverage: ## Generate global code coverage report
	@go test -timeout 60s -covermode=count -coverprofile /tmp/coverage.out ${PKG_LIST}

coverhtml: ## Generate global code coverage report in HTML
	@go tool cover -html=/tmp/coverage.out -o /tmp/coverage.html

dep: ## Get the backend dependencies
	@go get -u github.com/grpc-ecosystem/grpc-gateway/protoc-gen-grpc-gateway
	@go get -u github.com/grpc-ecosystem/grpc-gateway/protoc-gen-swagger
	@go get -u github.com/golang/protobuf/protoc-gen-go
	@go get -u github.com/elazarl/go-bindata-assetfs/go-bindata-assetfs
	@go get -u github.com/jteeuwen/go-bindata/go-bindata

ui/test: ## Run frontend syntax check
	@cd ui && npm install && npm run test

ui/build_dep: ## Get the frontend dependencies
	@echo "Building node-sass"
	@cd ui/node_modules/node-sass/ && npm install && npm run build

ui/build: ui/build_dep ## Build the frontend
	@echo "Building ui"
	@cd ui && npm run build
	@mv ui/build/* static

build: dep ## Build the backend binary file
	@export CGO_ENABLED=0
	@export GO_EXTRA_BUILD_ARGS="-a -installsuffix cgo"
	#@go build -i -v $(PKG)
	@go build $(GO_EXTRA_BUILD_ARGS) -ldflags "-s -w -X main.version=$(VERSION)" -o build/lora-app-server cmd/lora-app-server/main.go

clean: ## Remove previous build
	@rm -f $(PROJECT_NAME)

help: ## Display this help screen
	@grep -h -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
