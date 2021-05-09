.PHONY: default build test deps fmt fmt-ci

default: test build

build:
	go build -o ./dist/wavelength .

test:
	go test -v ./...

deps:
	go get ./...

fmt:
	go fmt ./...

fmt-ci:
	if [ -z $$( go fmt ./... ) ]; then exit 0; else exit 1; fi
