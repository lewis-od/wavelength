FROM golang:1.15.12 as build

WORKDIR /wavelength

COPY Makefile .
COPY go.mod .
COPY go.sum .
RUN make deps

COPY . .
RUN make build

FROM debian:buster-slim

COPY --from=build /wavelength/dist/wavelength /usr/bin/wavelength
RUN chmod +x /usr/bin/wavelength

RUN mkdir -p /root/.aws && chown -R root /root/.aws
VOLUME /root/.aws

RUN mkdir -p /project && chown -R root /project
WORKDIR /project
VOLUME /project

ENTRYPOINT ["/usr/bin/wavelength"]

