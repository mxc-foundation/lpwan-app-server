.PHONY: format build clean test lint sec package package-deb ui/build ui/build_dep api statics requirements ui-requirements serve update-vendor internal/statics internal/migrations static/swagger/api.swagger.json
PKGS := $(shell go list ./... | grep -v /vendor |grep -v lora-app-server/api | grep -v /migrations | grep -v /static | grep -v /ui)
VERSION := $(shell git describe --tags --always --long |sed -e "s/^v//")

build: internal/statics internal/migrations
	mkdir -p build cache
	go build $(GO_EXTRA_BUILD_ARGS) -ldflags "-s -w -X main.version=$(VERSION)" -o build/lora-app-server cmd/lora-app-server/main.go

clean:
	@echo "Cleaning up workspace"
	@rm -rf build dist internal/migrations/migrations_gen.go internal/static/static_gen.go ui/build static/static
	@rm -f static/index.html static/icon.png static/manifest.json static/asset-manifest.json static/service-worker.js
	@rm -rf static/logo
	@rm -rf static/img
	@rm -f static/swagger/*.json
	@rm -rf docs/public
	@rm -rf dist

test: internal/statics internal/migrations
	@echo "Running tests"
	@rm -f coverage.out
	@for pkg in $(PKGS) ; do \
		golint $$pkg ; \
	done
	@go vet $(PKGS)
	@go test -p 1 -v $(PKGS) -cover -coverprofile coverage.out

lint:
	@echo "Running code syntax check"
	@go get -u golang.org/x/lint/golint
	@golint -set_exit_status $(PKGS)

golangci-lint-new:
	docker pull golangci/golangci-lint
	docker run --rm -v $$(pwd):/app -w /app golangci/golangci-lint:v1.26.0 golangci-lint run --new-from-rev master ./...

golangci-lint:
	docker pull golangci/golangci-lint
	docker run --rm -v $$(pwd):/app -w /app golangci/golangci-lint:v1.26.0 golangci-lint run ./...

sec:
	@echo "Running code security check"
	@go get github.com/securego/gosec/cmd/gosec
	@gosec ./...

dist: ui/build internal/statics internal/migrations
	@goreleaser
	mkdir -p dist/upload/tar
	mkdir -p dist/upload/deb
	mv dist/*.tar.gz dist/upload/tar
	mv dist/*.deb dist/upload/deb

snapshot: ui/build internal/statics internal/migrations
	@goreleaser --snapshot

ui/test:
	@echo "Running react tests"
	@cd ui && npm test

ui/build_dep:
	@echo "Building node-sass"
	@cd ui/node_modules/node-sass/ && npm install
	@echo "Running npm audit fix"
	@cd ui && npm audit fix

ui/build:
	@echo "BUilding ui"
	@cd ui && npm run build
	@mv ui/build/* static

api:
	@echo "Generating API code from .proto files"
	@go generate api/appserver-serves-ui/api.go
	@go generate api/appserver-serves-gateway/api.go

internal/statics internal/migrations: static/swagger/api.swagger.json
	@echo "Generating static files"
	@go generate internal/migrations/migrations.go
	@go generate internal/static/static.go


static/swagger/api.swagger.json:
	@echo "Generating combined Swagger JSON"
	@GOOS="" GOARCH="" go run api/appserver-serves-ui/swagger/main.go api/appserver-serves-ui/swagger > static/swagger/api.swagger.json
	@cp api/appserver-serves-ui/swagger/*.json static/swagger


# shortcuts for development

dev-requirements:
	go install github.com/grpc-ecosystem/grpc-gateway/protoc-gen-grpc-gateway
	go install github.com/grpc-ecosystem/grpc-gateway/protoc-gen-swagger
	go install github.com/golang/protobuf/protoc-gen-go
	go install github.com/elazarl/go-bindata-assetfs/go-bindata-assetfs
	go install github.com/jteeuwen/go-bindata/go-bindata

ui-requirements:
	@echo "Installing UI requirements"
	@cd ui && npm install

serve: build
	@echo "Starting LPWAN App Server"
	./build/lora-app-server

update-vendor:
	@echo "Updating vendored packages"
	@govendor update +external

run-compose-test:
	docker-compose run --rm appserver make test
