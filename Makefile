.PHONY: default build test

default: test build

build:
	go build -o ./dist/lambda-build .

test:
	go test -v ./...
