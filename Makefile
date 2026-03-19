.PHONY: run-debug 

SHELL = /bin/zsh

build-cli:
	go build -o bin/paytunnel ./cmd/paytunnel  

run-debug:
	go run examples/btcpay-basics/main.go

lint:
	golangci-lint version && golangci-lint run --verbose  -E  misspell   