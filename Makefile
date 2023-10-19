ROOT_DIR := $(dir $(realpath $(lastword $(MAKEFILE_LIST))))
BUILD_FLAGS := -trimpath -gcflags "all=-N -l"
DEV_BUILD_FLAGS := -race

VERSION := $(shell git describe --always --tags | sed 's/^v//g' | awk -F- '{print $$1}')
RELEASE := $(shell git describe --always --tags | awk -F- '{if(!$$2){$$2=0};print $$2".el7"}')

.PHONY: all
all: build

.PHONY: build
build:
	go build $(DEV_BUILD_FLAGS) $(BUILD_FLAGS) -o bin/tcpmon main.go

.PHONY: build-linux
build-linux:
	GOOS=linux go build $(DEV_BUILD_FLAGS) $(BUILD_FLAGS) -o bin/tcpmon main.go

.PHONY: release
release: proto
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build $(BUILD_FLAGS) -o bin/x86_64/tcpmon github.com/zperf/tcpmon
	GOOS=linux GOARCH=arm64 CGO_ENABLED=0 go build $(BUILD_FLAGS) -o bin/aarch64/tcpmon github.com/zperf/tcpmon

.PHONY: proto
proto:
	mkdir -p tcpmon/tproto
	protoc -Iproto --gogofaster_out=tcpmon/tproto tcpmon.proto
	sed -i 's/package tcpmon/package tproto/g' tcpmon/tproto/tcpmon.pb.go

.PHONY: gproto
gproto:
	mkdir -p tcpmon/gproto
	protoc --go_out=. --go_opt=Mproto/tcpmon.proto=./tcpmon/gproto proto/tcpmon.proto

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
ifeq (, $(shell which nfpm))
	$(MAKE) tools
endif
	VERSION=$(VERSION) RELEASE=$(RELEASE) nfpm package -p rpm -t rpm/

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
