all: build

.PHONY: build
build: proto
	go build -o bin/tcpmon main.go

.PHONY: release
release: proto
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -trimpath -gcflags "all=-N -l" -o bin/tcpmon github.com/zperf/tcpmon

.PHONY: proto
proto:
ifeq (, $(shell which protoc-gen-go))
	$(MAKE) install-deps
endif
	protoc --go_out=. --go_opt=Mproto/tcpmon.proto=./tcpmon proto/tcpmon.proto

.PHONY: check
check: build
ifeq (, $(shell which golangci-lint))
	$(MAKE) install-deps
endif
	golangci-lint run
	go test -race -v ./...

.PHONY: install-deps
install-deps:
	go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.53.3
