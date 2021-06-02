.PHONY: default clean build test deps docker docker-build docker-push fmt fmt-ci

default: clean test build

clean:
	@rm ./dist/wavelength

build: ./dist/wavelength

./dist/wavelength:
	go build -o ./dist/wavelength .

test:
	go test -v ./...

deps:
	go mod download

image-name := ghcr.io/lewis-od/wavelength
image-version := $(shell git describe --tags)
docker: docker-build docker-push

docker-build:
	docker build -t $(image-name):$(image-version) .

docker-push:
	docker push $(image-name):$(image-version)

fmt:
	go fmt ./...

fmt-ci:
	if [ -z $$( go fmt ./... ) ]; then exit 0; else exit 1; fi
