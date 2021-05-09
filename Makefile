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
	! go fmt ./... 2>&1 | read 2>/dev/null
