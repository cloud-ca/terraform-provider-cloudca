
default: build

init:
	curl https://glide.sh/get | sh
	glide install

build:
	go build .

build-all:
		gox -verbose \
		-ldflags "-X main.version=${VERSION}" \
		-os="linux darwin windows freebsd openbsd" \
		-arch="amd64 386 armv5 armv6 armv7 arm64" \
		-osarch="!darwin/arm64" \
		-output="dist/{{.OS}}-{{.Arch}}/{{.Dir}}" .

clean:
	rm -rf dist terraform-provider-cloudca

.PHONY: init build vet build-all
