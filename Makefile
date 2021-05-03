.PHONY: default build test

default: test build

build:
	go build -o ./dist/wavelength .

test:
	go test -v ./...
