#!/bin/bash

set -e

version=${1:-0.0.0}
out=build-v${version}
cd "$(dirname "${BASH_SOURCE[0]}")"
mkdir -p "$out"

GOOS=darwin GOARCH=amd64 go build -ldflags "-X main.Version=$version" -o "$out/yj-macos-amd64" .
GOOS=darwin GOARCH=arm64 go build -ldflags "-X main.Version=$version" -o "$out/yj-macos-arm64" .
GOOS=linux GOARCH=amd64 go build -ldflags "-X main.Version=$version" -o "$out/yj-linux-amd64" .
GOOS=linux GOARCH=arm64 go build -ldflags "-X main.Version=$version" -o "$out/yj-linux-arm64" .
GOOS=linux GOARCH=arm GOARM=5 go build -ldflags "-X main.Version=$version" -o "$out/yj-linux-arm-v5" .
GOOS=linux GOARCH=arm GOARM=7 go build -ldflags "-X main.Version=$version" -o "$out/yj-linux-arm-v7" .
GOOS=windows GOARCH=amd64 go build -ldflags "-X main.Version=$version" -o "$out/yj.exe" .

docker build . --build-arg "version=$version" -t "sclevine/yj:$version"
docker tag "sclevine/yj:$version" "sclevine/yj:latest"