.PHONY: build clean test lint sec package package-deb api statics requirements ui-requirements serve update-vendor internal/statics internal/migrations static/swagger/api.swagger.json
PKGS := $(shell go list ./... | grep -v /vendor |grep -v lora-app-server/api | grep -v /migrations | grep -v /static | grep -v /ui)
VERSION := $(shell git describe --tags --always --long |sed -e "s/^v//")

build: internal/statics internal/migrations
	mkdir -p build cache
	go build $(GO_EXTRA_BUILD_ARGS) -ldflags "-s -w -X main.version=$(VERSION)" -o build/lora-app-server cmd/lora-app-server/main.go

clean:
	@echo "Cleaning up workspace"
	@rm -rf build dist internal/migrations/migrations_gen.go internal/static/static_gen.go ui/build static/static
	@rm -f static/index.html static/icon.png static/mxc_icon.png static/manifest.json static/asset-manifest.json static/service-worker.js static/precache-manifest.*.js
	@rm -rf static/logo
	@rm -rf static/img
	@rm -f static/swagger/*.json
	@rm -rf docs/public
	@rm -rf dist

test: internal/statics internal/migrations
	# we only have non-generated code in ./internal, so we only count coverage for it
	go test -cover -coverprofile coverage.out -coverpkg ./internal/... ./...
	# IMPORTANT: required coverage can only be increased
	go tool cover -func coverage.out | \
		awk 'END { print "Coverage: " $$3; if ($$3+0 < 12.9) { print "Insufficient coverage"; exit 1; } }'

lint:
	@echo "Running code syntax check"
	@go get -u golang.org/x/lint/golint
	@golint -set_exit_status $(PKGS)

golangci-lint-new:
	docker pull golangci/golangci-lint:v1.36.0
	docker run --rm -v $$(pwd):/app -w /app golangci/golangci-lint:v1.36.0 golangci-lint run --new-from-rev master ./...

golangci-lint:
	docker pull golangci/golangci-lint:v1.36.0
	docker run --rm -v $$(pwd):/app -w /app golangci/golangci-lint:v1.36.0 golangci-lint run ./...

sec:
	@echo "Running code security check"
	@go get github.com/securego/gosec/cmd/gosec
	@gosec ./...

dist: ui/build internal/statics internal/migrations
	@goreleaser
	mkdir -p dist/upload/tar
	mkdir -p dist/upload/deb
	mkdir -p dist/upload/rpm
	mv dist/*.tar.gz dist/upload/tar
	mv dist/*.deb dist/upload/deb
	mv dist/*.rpm dist/upload/rpm

snapshot: ui/build internal/statics internal/migrations
	@goreleaser --snapshot

api:
	@echo "Generating API code from .proto files"
	@rm -rf /tmp/chirpstack-api
	@git clone https://github.com/brocaar/chirpstack-api.git /tmp/chirpstack-api
	@cp -rf /tmp/chirpstack-api/protobuf/* api/extapi/
	@go generate api/extapi/api.go
	@go generate api/appserver-serves-gateway/api.go
	@go generate api/appserver-serves-m2m/api.go
	@go generate api/networkserver/api.go

internal/statics internal/migrations: static/swagger/api.swagger.json
	@echo "Generating static files"
	@go generate internal/migrations/migrations.go
	@go generate internal/static/static.go


static/swagger/api.swagger.json:
	@echo "Generating combined Swagger JSON"
	@GOOS="" GOARCH="" go run api/extapi/swagger/main.go api/extapi/swagger > static/swagger/api.swagger.json
	@cp api/extapi/swagger/*.json static/swagger


# shortcuts for development
debug:
	go get github.com/go-delve/delve/cmd/dlv

dev-requirements:
	go get -u github.com/golang/protobuf/protoc-gen-go
	go get -u github.com/grpc-ecosystem/grpc-gateway/protoc-gen-grpc-gateway
	go get -u github.com/grpc-ecosystem/grpc-gateway/protoc-gen-swagger
	go get -u github.com/go-bindata/go-bindata/...

serve: build
	@echo "Starting LPWAN App Server"
	./build/lora-app-server

dep-graph:
	goda graph -short 'github.com/mxc-foundation/lpwan-app-server/...:root' | dot -Tpdf -o dep-graph.pdf
