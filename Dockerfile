FROM golang:1.15.12 as build

WORKDIR /wavelength

COPY Makefile .
COPY go.mod .
COPY go.sum .
RUN make deps

COPY . .
RUN make build
RUN chmod +x dist/wavelength

