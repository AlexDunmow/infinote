#!/bin/bash

set -e -u -x
VERSION=$(cat version/version)

echo "Building version $(echo $VERSION)"

echo
echo "Installing system dependencies..."

apk add git -q

echo
echo "Building..."

cd repo/server
go build -o platform cmd/platform/main.go

cd ../..


tar cvf "./server/release-server-$(echo $VERSION).tar.gz" ./repo/server/platform