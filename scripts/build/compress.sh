#!/usr/bin/env bash

set -o errexit
set -o nounset
set -o pipefail

BUILD_DIR="${1:-bin}"
NAME="${2:-terraform-provider-cloudca}"
VERSION=$3

if [ -z "${NAME}" ]; then
    echo "Error: NAME is missing. e.g. ./compress.sh <name> <version>"
    exit 1
fi

if [ -z "${VERSION}" ]; then
    echo "Error: VERSION is missing. e.g. ./compress.sh <name> <version>"
    exit 1
fi

PWD=$(cd $(dirname "$0") && pwd -P)
BUILD_DIR="${PWD}/../../${BUILD_DIR}"

printf "\033[36m==> Compress binaries\033[0m\n"

for platform in $(find ${BUILD_DIR} -mindepth 1 -maxdepth 1 -type d); do
    OSARCH=$(basename ${platform})
    FULLNAME="${NAME}_${VERSION}_${OSARCH}"

    if ! command -v zip >/dev/null; then
        echo "Error: cannot compress, 'zip' not found"
        exit 1
    fi

    zip -q -j ${BUILD_DIR}/${FULLNAME}.zip ${platform}/*
    printf -- "--> %15s: bin/%s\n" "${OSARCH}" "${FULLNAME}.zip"
done

printf "\033[36m==> Generate checksum\033[0m\n"

cd ${BUILD_DIR}
touch ${NAME}_${VERSION}_SHA256SUMS

for binary in $(find . -mindepth 1 -maxdepth 1 -type f | grep -v "${NAME}_${VERSION}_SHA256SUMS" | sort); do
    binary=$(basename ${binary})

    if command -v sha256sum >/dev/null; then
        sha256sum ${binary} >>${NAME}_${VERSION}_SHA256SUMS
    elif command -v shasum >/dev/null; then
        shasum -a256 ${binary} >>${NAME}_${VERSION}_SHA256SUMS
    fi
done

cd - >/dev/null 2>&1
printf -- "--> %15s: bin/%s\n" "sha256sum" "${NAME}_${VERSION}_SHA256SUMS"
