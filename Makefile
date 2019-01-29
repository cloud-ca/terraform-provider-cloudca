VERSION := $(shell git describe --tags)
VERSION_COMMIT := $(shell git describe --always --long)
ifeq ($(VERSION),)
VERSION:=$(VERSION_COMMIT)
endif

default: build

vendor:
	GO111MODULE=on go mod vendor

build:
	go build .

build-all: clean
	@gox -verbose \
		-ldflags "-X main.version=${VERSION}" \
		-gcflags=-trimpath=${GOPATH} \
		-os="linux darwin windows freebsd openbsd solaris" \
		-arch="386 amd64 arm" \
		-osarch="!darwin/arm !darwin/386" \
		-output="dist/{{.OS}}-{{.Arch}}/{{.Dir}}_${VERSION}" .

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

upload:
	rm -f ./dist/terraform-provider-cloudca_${VERSION}_SWIFTURLS ;
	SWIFT_ACCOUNT=`swift stat | grep Account: | sed s/Account:// | tr -d '[:space:]'` ; \
	SWIFT_URL=https://objects-qc.cloud.ca/v1 ; \
	SWIFT_CONTAINER=terraform-provider-cloudca ; \
	for FILE in `ls ./dist | grep -i terraform.*\.zip` ; do \
		echo "Uploading $$FILE to swift" ; \
		swift upload $${SWIFT_CONTAINER} ./dist/$$FILE --object-name ${VERSION}/$$FILE ; \
		echo "$${SWIFT_URL}/$${SWIFT_ACCOUNT}/$${SWIFT_CONTAINER}/${VERSION}/$$FILE" >> ./dist/terraform-provider-cloudca_${VERSION}_SWIFTURLS ; \
	done

release-notes: 
	./release-notes.sh > ./dist/release.md ;

release: build-all upload release-notes

clean:
	rm -rf dist terraform-provider-cloudca

.PHONY: default vendor build build-all upload release-notes release clean
