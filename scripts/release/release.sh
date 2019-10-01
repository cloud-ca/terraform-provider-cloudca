#!/usr/bin/env bash

set -o errexit
set -o pipefail

CURRENT_BRANCH=$(git rev-parse --abbrev-ref HEAD)
if [ -z "${CURRENT_BRANCH}" -o "${CURRENT_BRANCH}" != "master" ]; then
    echo "Error: The current branch is '${CURRENT_BRANCH}', switch to 'master' to do the release."
    exit 1
fi

if [ -n "$(git status --short)" ]; then
    echo "Error: There are untracked/modified changes, commit or discard them before the release."
    exit 1
fi

RELEASE_VERSION=$1
PUSH=$2
CURRENT_VERSION=$3
FROM_MAKEFILE=$4

if [ -z "${RELEASE_VERSION}" ]; then
    if [ -z "${FROM_MAKEFILE}" ]; then
        echo "Error: VERSION is missing. e.g. ./release.sh <version> <push>"
    else
        echo "Error: missing value for 'version'. e.g. 'make release version=x.y.z'"
    fi
    exit 1
fi

if [ -z "${PUSH}" ]; then
    echo "Error: PUSH is missing. e.g. ./release.sh <version> <push>"
    exit 1
fi

if [ -z "${CURRENT_VERSION}" ]; then
    CURRENT_VERSION=$(git describe --tags --exact-match 2>/dev/null || git describe --tags 2>/dev/null || echo "v0.0.1-$(COMMIT_HASH)")
fi

if [ "v${RELEASE_VERSION}" = "${CURRENT_VERSION}" ]; then
    echo "Error: provided version (v${RELEASE_VERSION}) already exists."
    exit 1
fi

if [ $(git describe --tags "v${RELEASE_VERSION}" 2>/dev/null) ]; then
    echo "Error: provided version (v${RELEASE_VERSION}) already exists."
    exit 1
fi

PWD=$(cd $(dirname "$0") && pwd -P)
CLOSEST_VERSION=$(git describe --tags --abbrev=0)

# Bump the released version in README and version.go
sed -i -E 's|'${CLOSEST_VERSION}'|v'${RELEASE_VERSION}'|g' README.md
sed -i -E 's|v'${RELEASE_VERSION}'-alpha|v'${RELEASE_VERSION}'|g' cmd/rostamctl/version/version.go

# Commit changes
printf "\033[36m==> %s\033[0m\n" "Commit changes for release version v${RELEASE_VERSION}"
git add README.md cmd/cca/version/version.go TODO
git commit -m "Release version v${RELEASE_VERSION}"

if [ "${PUSH}" == "true" ]; then
    printf "\033[36m==> %s\033[0m\n" "Push commits for v${RELEASE_VERSION}"
    git push origin master
fi

# Generate Changelog
make --no-print-directory -f ${PWD}/../../Makefile changelog push="${PUSH}" next="--next-tag v${RELEASE_VERSION}"

# Tag the release
printf "\033[36m==> %s\033[0m\n" "Tag release v${RELEASE_VERSION}"
git tag --annotate --message "v${RELEASE_VERSION} Release" "v${RELEASE_VERSION}"

if [ "${PUSH}" == "true" ]; then
    printf "\033[36m==> %s\033[0m\n" "Push tag release v${RELEASE_VERSION}"
    git push origin v${RELEASE_VERSION}
fi

# Bump the next version in version.go
NEXT_VERSION=$(echo "${RELEASE_VERSION}" | sed 's/^v//' | awk -F'[ .]' '{print $1"."$2+1".0"}')
sed -i -E 's|v'${RELEASE_VERSION}'|v'${NEXT_VERSION}'-alpha|g' cmd/cca/version/version.go TODO

# Commit changes
printf "\033[36m==> %s\033[0m\n" "Bump version to v${NEXT_VERSION}-alpha"
git add cmd/rostamctl/version/version.go TODO
git commit -m "Bump version to v${NEXT_VERSION}-alpha" TODO

if [ "${PUSH}" == "true" ]; then
    printf "\033[36m==> %s\033[0m\n" "Push commits for v${NEXT_VERSION}-alpha"
    git push origin master
fi
