#!/bin/bash

set -ex;

# 移至脚本目录
cd /build

echo `pwd`

echo "set go env GOFLAGS"
/usr/local/go/bin/go env -w GOFLAGS="-buildvcs=false"

/usr/local/go/bin/go mod tidy

echo "build start"

CGO_ENABLED=1 GOOS=linux GOARCH=amd64 /usr/local/go/bin/go \
build -ldflags="-s -X teamide/pkg/base.version="$1"" \
-o linux-amd64-electron-server .

CGO_ENABLED=1 GOOS=linux GOARCH=amd64 /usr/local/go/bin/go \
build -ldflags="-s -X teamide/pkg/base.version="$1" -X main.buildFlags=--isServer" \
-o linux-amd64-server .

CGO_ENABLED=1 GOOS=linux GOARCH=amd64 /usr/local/go/bin/go \
build -ldflags="-s -X teamide/pkg/base.version="$1"" \
-o linux-amd64-node teamide/pkg/node/main

echo "build success"


