#!/usr/bin/env bash

VERSION=v1.1.0

docker run \
  --mount "type=bind,source=$HOME/.aws,target=/root/.aws" \
  --mount "type=bind,source=$PWD,target=/project" \
  ghcr.io/lewis-od/wavelength:$VERSION "$@"
