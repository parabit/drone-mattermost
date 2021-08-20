#!/bin/bash

SRC=$(realpath $(cd -P "$(dirname "${BASH_SOURCE[0]}")" && pwd))

set -e

pushd $SRC &> /dev/null

# disable cgo
(set -x;
  GO111MODULE=on \
  CGO_ENABLED=0 \
    go build .
  docker build -t kenshaw/drone-mattermost:latest .
  docker push kenshaw/drone-mattermost:latest
)

popd &> /dev/null
