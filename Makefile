.PHONY: build run test fmt

build:
	go build -o bin/cln ./cmd/cln

run:
	go run ./cmd/cln $(ARGS)

test:
	go test ./...

fmt:
	gofmt -s -w .
