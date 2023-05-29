#!/bin/sh
echo "build start"
VERSION_INFO="-X 'main.version=$(git describe --tags --abbrev=0)' -X 'main.date=$(date +%F)' -X 'main.commit=$(git rev-parse --short HEAD)'"
go build -mod vendor --ldflags "${VERSION_INFO}" main.go
echo "complete"