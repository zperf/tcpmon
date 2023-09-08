.PHONY: all
all:
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.53.3
	go install github.com/bep/punused@latest
	go install golang.org/x/tools/gopls@latest
	go install honnef.co/go/tools/cmd/staticcheck@latest
	go install github.com/go-delve/delve/cmd/dlv@latest
	go install github.com/chzchzchz/goword@latest
	go run codegen.go
