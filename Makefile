VERSION := $(shell git describe --tags)
VERSION_COMMIT := $(shell git describe --always --long)
ifeq ($(VERSION),)
VERSION:=$(VERSION_COMMIT)
endif

default: build

init:
	curl https://glide.sh/get | sh
	glide install

build:
	go build .

build-all: clean
	@gox -verbose \
		-ldflags "-X main.version=${VERSION}" \
		-gcflags=-trimpath=${GOPATH} \
		-os="linux darwin windows freebsd openbsd solaris" \
		-arch="386 amd64 arm" \
		-osarch="!darwin/arm !darwin/386" \
		-output="dist/{{.OS}}-{{.Arch}}/{{.Dir}}" .

	@for PLATFORM in `find ./dist -mindepth 1 -maxdepth 1 -type d` ; do \
		OSARCH=`basename $$PLATFORM` ; \
		echo "--> $$OSARCH" ; \
		pushd $$PLATFORM >/dev/null 2>&1 ; \
		zip ../terraform-provider-cloudca_$(VERSION)_$$OSARCH.zip ./* ; \
		popd >/dev/null 2>&1 ; \
	done

	pushd ./dist ; \
	shasum -a256 *.zip > ./terraform-provider-cloudca_${VERSION}_SHA256SUMS ; \
	popd >/dev/null 2>&1 ;
clean:
	rm -rf dist terraform-provider-cloudca

.PHONY: init build build-all clean
