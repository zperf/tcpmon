ROOT_DIR := $(dir $(realpath $(lastword $(MAKEFILE_LIST))))
BUILD_FLAGS := -trimpath -gcflags "all=-N -l"

.PHONY: all
all: build

.PHONY: build
build:
	go build $(BUILD_FLAGS) -o bin/tcpmon main.go

.PHONY: build-linux
build-linux:
	GOOS=linux go build $(BUILD_FLAGS) -o bin/tcpmon main.go

.PHONY: release
release: proto
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build $(BUILD_FLAGS) -o bin/tcpmon-x86_64 github.com/zperf/tcpmon
	GOOS=linux GOARCH=arm64 CGO_ENABLED=0 go build $(BUILD_FLAGS) -o bin/tcpmon-aarch64 github.com/zperf/tcpmon

.PHONY: proto
proto:
	protoc -Iproto --gogofaster_out=tcpmon/ tcpmon.proto

.PHONY: check
check: build
ifeq (, $(shell which golangci-lint))
	$(MAKE) tools
endif
	golangci-lint run
	go test -race -v ./...

.PHONY: tools
tools:
	$(MAKE) -C tools

.PHONY: package
package: release rpm

.PHONY: rpm
rpm:
	$(MAKE) -C rpm

.PHONY: clean
clean:
	find rpm -name "tcpmon*rpm" -type f -exec rm -f {} \;

.PHONY: docker
docker:
	$(MAKE) -C docker builder

# list all tests
.PHONY: tests
tests:
	go test -v ./... -list Test | grep -v "^ok\|no test files"
