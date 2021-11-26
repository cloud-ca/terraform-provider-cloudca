#!/usr/bin/env bash

set -o errexit
set -o nounset
set -o pipefail

BUILD_DIR="${1:-bin}"
VERSION=$2
GOOS="${3:-"linux darwin windows freebsd openbsd solaris"}"
GOARCH="${4:-"386 amd64 arm"}"
GOLDFLAGS="$5"

if [ -z "${VERSION}" ]; then
    echo "Error: VERSION is missing. e.g. ./build-all.sh <build_dir> <version> <build_os_list> <build_arch_list> <build_ldflag>"
    exit 1
fi

if [ -z "${GOLDFLAGS}" ]; then
    echo "Error: GOLDFLAGS is missing. e.g. ./build-all.sh <build_dir> <version> <build_os_list> <build_arch_list> <build_ldflag>"
    exit 1
fi

PWD=$(cd $(dirname "$0") && pwd -P)
BUILD_DIR="${PWD}/../../${BUILD_DIR}"

CGO_ENABLED=0 gox \
    -verbose \
    -ldflags "${GOLDFLAGS}" \
    -gcflags=-trimpath=$(go env GOPATH) \
    -os="${GOOS}" \
    -arch="${GOARCH}" \
    -osarch="!darwin/arm !darwin/386" \
    -output="${BUILD_DIR}/{{.OS}}-{{.Arch}}/{{.Dir}}_${VERSION}" ${PWD}/../../
