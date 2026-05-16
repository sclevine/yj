#!/bin/bash

set -e

version=${1:-0.0.0}
out=build-v${version}
cd "$(dirname "${BASH_SOURCE[0]}")"
mkdir -p "$out"

create_deb () {
  local IMAGE=$1
  local PLATFORM=$2

  # We need to build the DEB inside a docker container for each OS. That
  # docker container runs as root user yet the volume mounted source directory
  # is (in all likelihood) not owned by root. This presents a minor problem for
  # the Makefile which is used to create the DEB because that relies on Git to
  # determine the version. Git is quite picky about directory ownership
  # mismatches. So, we need to copy the source code over to another directory
  # inside the docker container before building. This happens to be a good
  # solution for another problem, which is identifying the resulting DEB, which
  # is always placed in the .. directory (for details, see
  # https://groups.google.com/g/linux.debian.bugs.dist/c/1KiGKfuFH3Y), and needs
  # to be renamed with an OS specific name anyways.
  docker run --platform=${PLATFORM} --rm -v $(pwd):/src ${IMAGE} /bin/bash -c "
    mkdir /build && \
    cp -R /src /build && \
    cd /build/src && \
    export DEBIAN_FRONTEND=noninteractive && \
    apt-get update && \
    apt-get install -y --no-install-recommends \
      build-essential ca-certificates debhelper dpkg-dev git golang-go && \
    make deb && \
    cd .. && \
    source /etc/os-release && \
    mv yj*.deb /src/\$(ls yj*.deb | sed 's/yj_/yj_'"\${ID}_\${VERSION_ID}_\${VERSION_CODENAME}_"'/')
  "
  mv yj*.deb "${out}"
}

for PLAT in linux/arm64 linux/amd64; do
  create_deb debian:10-slim ${PLAT}
  create_deb debian:11-slim ${PLAT}
  create_deb debian:12-slim ${PLAT}
  create_deb ubuntu:20.04 ${PLAT}
  create_deb ubuntu:22.04 ${PLAT}
done

GOOS=darwin GOARCH=amd64 go build -ldflags "-X main.Version=$version" -o "$out/yj-macos-amd64" .
GOOS=darwin GOARCH=arm64 go build -ldflags "-X main.Version=$version" -o "$out/yj-macos-arm64" .
GOOS=linux GOARCH=amd64 go build -ldflags "-X main.Version=$version" -o "$out/yj-linux-amd64" .
GOOS=linux GOARCH=arm64 go build -ldflags "-X main.Version=$version" -o "$out/yj-linux-arm64" .
GOOS=linux GOARCH=arm GOARM=5 go build -ldflags "-X main.Version=$version" -o "$out/yj-linux-arm-v5" .
GOOS=linux GOARCH=arm GOARM=7 go build -ldflags "-X main.Version=$version" -o "$out/yj-linux-arm-v7" .
GOOS=windows GOARCH=amd64 go build -ldflags "-X main.Version=$version" -o "$out/yj.exe" .

docker build . --build-arg "version=$version" -t "sclevine/yj:$version"
docker tag "sclevine/yj:$version" "sclevine/yj:latest"
