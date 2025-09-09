#!/bin/bash

SRC=$(realpath $(cd -P "$(dirname "${BASH_SOURCE[0]}")" && pwd))

set -e

pushd $SRC &> /dev/null

# disable cgo
(set -x;
  GOOS=linux GOARCH=amd64 \
  GO111MODULE=on \
  CGO_ENABLED=0 \
    go build .
  docker build -t parabit/drone-mattermost:latest --platform linux/amd64 .
  docker push parabit/drone-mattermost:latest
)

popd &> /dev/null
