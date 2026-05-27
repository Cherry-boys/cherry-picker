.PHONY: build test lint install clean

build:
	go build -o bin/ares-pp-cli ./cmd/ares-pp-cli

test:
	go test ./...

lint:
	golangci-lint run

install:
	go install ./cmd/ares-pp-cli

clean:
	rm -rf bin/

build-mcp:
	go build -o bin/ares-pp-mcp ./cmd/ares-pp-mcp

install-mcp:
	go install ./cmd/ares-pp-mcp

build-all: build build-mcp
