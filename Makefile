all: build

.PHONY: build
build: proto
	go build -o bin/tcpmon main.go

.PHONY: release
release: proto
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -trimpath -gcflags "all=-N -l" -o bin/tcpmon github.com/zperf/tcpmon

.PHONY: proto
proto:
	protoc --go_out=. --go_opt=Mproto/tcpmon.proto=./tcpmon proto/tcpmon.proto

.PHONY: check
check: build
	go test -race -v ./...
