.PHONY: build clean test package package-deb ui api statics requirements ui-requirements serve update-vendor internal/statics internal/migrations static/swagger/api.swagger.json
PKGS := $(shell go list ./... | grep -v /vendor |grep -v lora-app-server/api | grep -v /migrations | grep -v /static | grep -v /ui)
VERSION := $(shell git describe --always |sed -e "s/^v//")
M2M_SERVER=$(shell cat lora-app-server.toml | grep mxp_server=| sed 's/^mxp_server=//g')
M2M_SERVER_DEV=$(shell cat lora-app-server.toml | grep mxp_server_development=| sed 's/^mxp_server_development=//g')
DEMO_USER=$(shell cat lora-app-server.toml | grep demo_user=| sed 's/^demo_user=//g')
DEMO_USER_PASSWORD=$(shell cat lora-app-server.toml | grep demo_user_password=| sed 's/^demo_user_password=//g')

build: ui/build internal/statics internal/migrations
	mkdir -p build
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
	@rm -f ui/.env.*

test: internal/statics internal/migrations
	@echo "Running tests"
	@rm -f coverage.out
	@for pkg in $(PKGS) ; do \
		golint $$pkg ; \
	done
	@go vet $(PKGS)
	@go test -p 1 -v $(PKGS) -cover -coverprofile coverage.out

dist: ui/build internal/statics internal/migrations
	@goreleaser
	mkdir -p dist/upload/tar
	mkdir -p dist/upload/deb
	mv dist/*.tar.gz dist/upload/tar
	mv dist/*.deb dist/upload/deb

snapshot: ui/build internal/statics internal/migrations
	@goreleaser --snapshot

ui/build:
	@echo "Building ui"
	@cd ui && printf 'REACT_APP_M2M_SERVER=$(M2M_SERVER)\nREACT_APP_DEMO_USER=$(DEMO_USER)\nREACT_APP_DEMO_USER_PASSWORD=$(DEMO_USER_PASSWORD)'  >> .env.production 
	@cd ui && printf 'REACT_APP_M2M_SERVER=$(M2M_SERVER_DEV)\nREACT_APP_DEMO_USER=$(DEMO_USER)\nREACT_APP_DEMO_USER_PASSWORD=$(DEMO_USER_PASSWORD)' >> .env.development 
	@cd ui && npm run build
	@mv ui/build/* static

api:
	@echo "Generating API code from .proto files"
	@go mod vendor
	@go generate api/api.go
	@rm -rf vendor/

internal/statics internal/migrations: static/swagger/api.swagger.json
	@echo "Generating static files"
	@go generate internal/migrations/migrations.go
	@go generate internal/static/static.go


static/swagger/api.swagger.json:
	@echo "Generating combined Swagger JSON"
	@GOOS="" GOARCH="" go run api/swagger/main.go api/swagger > static/swagger/api.swagger.json
	@cp api/swagger/*.json static/swagger


# shortcuts for development

dev-requirements:
	go mod download
	go install golang.org/x/lint/golint
	go install github.com/grpc-ecosystem/grpc-gateway/protoc-gen-grpc-gateway
	go install github.com/grpc-ecosystem/grpc-gateway/protoc-gen-swagger
	go install github.com/golang/protobuf/protoc-gen-go
	go install github.com/elazarl/go-bindata-assetfs/go-bindata-assetfs
	go install github.com/jteeuwen/go-bindata/go-bindata
	go install golang.org/x/tools/cmd/stringer
	go install github.com/goreleaser/goreleaser
	go install github.com/goreleaser/nfpm

ui-requirements:
	@echo "Installing UI requirements"
	@cd ui && npm install

serve: build
	@echo "Starting Lora App Server"
	./build/lora-app-server

update-vendor:
	@echo "Updating vendored packages"
	@govendor update +external

run-compose-test:
	docker-compose run --rm appserver make test
