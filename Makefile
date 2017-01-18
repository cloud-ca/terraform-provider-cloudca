
default: build

init:
	curl https://glide.sh/get | sh
	glide install

build:
	go build .

build-all:
	# compile for all OS/Arch using Gox
	gox -verbose \
		-ldflags "-X main.version=${VERSION}" \
		-os="linux darwin windows freebsd openbsd solaris" \
		-arch="386 amd64 arm" \
		-osarch="!darwin/arm !darwin/386" \
		-output="dist/{{.OS}}-{{.Arch}}/{{.Dir}}" .

	# zip the executables
	for PLATFORM in `find ./dist -mindepth 1 -maxdepth 1 -type d` ; do \
		OSARCH=`basename $$PLATFORM` ; \
		echo "--> $$OSARCH" ; \
		pushd $$PLATFORM >/dev/null 2>&1 ; \
		zip ../$$OSARCH.zip ./* ; \
		popd >/dev/null 2>&1 ; \
	done

clean:
	rm -rf dist terraform-provider-cloudca

.PHONY: init build vet build-all
