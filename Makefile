.PHONY: default clean build test deps fmt fmt-ci

default: clean test build

clean:
	@rm ./dist/wavelength

build: ./dist/wavelength

./dist/wavelength:
	go build -o ./dist/wavelength .

test:
	go test -v ./...

deps:
	go get ./...

fmt:
	go fmt ./...

fmt-ci:
	if [ -z $$( go fmt ./... ) ]; then exit 0; else exit 1; fi
