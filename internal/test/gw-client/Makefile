.PHONY: build clean
HOSTNAME := $(shell hostname -f)

build:
	mkdir -p build
	go build $(GO_EXTRA_BUILD_ARGS) -ldflags "-s -w -X main.host=$(HOSTNAME)" -o build/gw-client cmd/main.go

clean:
	@echo "Cleaning up workspace"
	@rm -rf build
