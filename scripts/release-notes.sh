#!/bin/bash

VERSION="$(git describe --tags)"
function getUrl() {
    echo $(cat ./dist/terraform-provider-cloudca_${VERSION}_SWIFTURLS | grep $1)
}

echo "[DESCRIPTION HERE]

# Issues fixed
[ISSUES HERE]

# Downloads

**macOS**
- 64-bit: $(getUrl darwin)

**Linux**
- 64-bit:  $(getUrl linux-amd64)
- 32-bit: $(getUrl linux-386)
- Arm: $(getUrl linux-arm)

**Windows**
- 64-bit: $(getUrl windows-amd64)
- 32-bit: $(getUrl windows-386)

**FreeBSD**
- 64-bit: $(getUrl freebsd-amd64)
- 32-bit: $(getUrl freebsd-386)
- Arm: $(getUrl freebsd-arm)

**OpenBSD**
- 64-bit: $(getUrl openbsd-amd64)
- 32-bit: $(getUrl openbsd-386)

**Solaris**
- 64-bit: $(getUrl solaris-amd64)
"