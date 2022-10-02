.PHONY: build
build:
	go build -v ./

test:
	go test -v -race -timeout 30s ./...

.DEFAULT_GOAL := build