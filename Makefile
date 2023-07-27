all: build

.PHONY: build
build:
	go build -o bin/tcpmon main.go

.PHONY: release
release:
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -trimpath -gcflags "all=-N -l" -o bin/tcpmon github.com/zperf/tcpmon
