ROOT_DIR := $(dir $(realpath $(lastword $(MAKEFILE_LIST))))

all: build

.PHONY: build
build:
	go build -o bin/tcpmon main.go

.PHONY: build-linux
build-linux:
	GOOS=linux GOARCH=arm64 go build -o bin/tcpmon-linux main.go

.PHONY: release
release: proto
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -trimpath -gcflags "all=-N -l" -o bin/tcpmon-x86_64 github.com/zperf/tcpmon
	GOOS=linux GOARCH=arm64 CGO_ENABLED=0 go build -trimpath -gcflags "all=-N -l" -o bin/tcpmon-aarch64 github.com/zperf/tcpmon

.PHONY: proto
proto: docker
	docker run -t --rm --user ${shell id -u}:${shell id -g} -v ${ROOT_DIR}:${ROOT_DIR} -w ${ROOT_DIR} tcpmon:builder protoc --go_out=. --go_opt=Mproto/tcpmon.proto=./tcpmon proto/tcpmon.proto

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
	$(MAKE) -C docker

# list all tests
.PHONY: tests
tests:
	go test -v ./... -list Test | grep -v "^ok\|no test files"
