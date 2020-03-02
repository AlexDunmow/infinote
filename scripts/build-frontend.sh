#!/bin/bash

set -e -u -x
VERSION=$(cat version/version)

echo
echo "Installing system dependencies..."

echo
echo "Building..."

cd repo/web
npm install
npm run build

cd ../..
tar cvf "./web/release-web-$(echo $VERSION).tar.gz" ./repo/web/dist