# Project variables
NAME        := terraform-provider-cloudca
DESCRIPTION := Terraform provider to interact with cloud.ca infrastructure
VENDOR      := cloud.ca
URL         := https://github.com/cloud-ca/terraform-provider-cloudca
LICENSE     := MIT

# Build variables
BUILD_DIR   := bin
VERSION     := $(shell git describe --always --tags --dirty 2>/dev/null || echo "0.0.0")
COMMIT_HASH := $(shell git rev-parse --short HEAD 2>/dev/null)
BUILD_DATE  := $(shell date +%FT%T%z)

# Go variables
GOOS        ?= linux darwin windows freebsd openbsd solaris
GOARCH      ?= 386 amd64 arm
GOPKGS      := $(shell go list ./... | grep -v /vendor)
GOFILES     := $(shell find . -type f -name '*.go' -not -path "./vendor/*")
LDFLAGS     := -X main.version=$(VERSION) -X main.commitHash=${COMMIT_HASH} -X main.buildDate=${BUILD_DATE}
GOBUILD     := go build -ldflags "$(LDFLAGS)"

# General variables
BLUE_COLOR  := \033[36m
NO_COLOR    := \033[0m

.PHONY: version
version:
	@echo $(VERSION)

#########################
## Development targets ##
#########################
.PHONY: clean
clean: log-clean ## Clean builds
	rm -rf ./$(BUILD_DIR) $(NAME)

.PHONY: vendor
vendor: log-vendor ## Install 'vendor' dependencies
	GO111MODULE=on go mod vendor

.PHONY: verify
verify: log-verify ## Verify 'vendor' dependencies
	GO111MODULE=on go mod verify

.PHONY: lint
lint: log-lint ## Run linter
	gometalinter ./...

.PHONY: format
format: log-format ## Format all go files
	go fmt $(GOPKGS)

.PHONY: checkfmt
checkfmt: SHELL :=/bin/bash
checkfmt: RESULT = $(shell gofmt -l $(GOFILES) | tee >(if [ "$$(wc -l)" = 0 ]; then echo "OK"; fi))
checkfmt: log-checkfmt ## Check formatting of all go files
	@echo $(RESULT)
	@if [ "$(RESULT)" != "OK" ]; then exit 1; fi

.PHONY: test
test: log-test ## Run tests
	go test -v $(GOPKGS)

.PHONY: tools
tools: log-tools ## Install required tools
	@curl -L https://git.io/vp6lP | sh # gometalinter
	@cd /tmp && go get -v -u github.com/mitchellh/gox # gox
	@cd /tmp && go get -v -u github.com/git-chglog/git-chglog/cmd/git-chglog # git-chglog

###################
## Build targets ##
###################
.PHONY: build
build: GOOS   = $(shell go env GOOS)
build: GOARCH = $(shell go env GOARCH)
build: clean log-build ## Build binary for current OS/ARCH
	@ GOOS=$(GOOS) GOARCH=$(GOARCH) $(GOBUILD) -o ./$(BUILD_DIR)/$(GOOS)-$(GOARCH)/$(NAME)_$(VERSION) && echo "./$(BUILD_DIR)/$(GOOS)-$(GOARCH)/$(NAME)_$(VERSION)"

.PHONY: build-all
build-all: SHELL :=/bin/bash
build-all: clean log-build-all ## Build binary for all OS/ARCH
	@gox -verbose \
		-ldflags "${LDFLAGS}" \
		-gcflags=-trimpath=${GOPATH} \
		-os="${GOOS}" \
		-arch="${GOARCH}" \
		-osarch="!darwin/arm !darwin/386" \
		-output="$(BUILD_DIR)/{{.OS}}-{{.Arch}}/{{.Dir}}_${VERSION}" .

	@printf "\n"
	@for platform in `find ./$(BUILD_DIR) -mindepth 1 -maxdepth 1 -type d` ; do \
		OSARCH=`basename $${platform}` ; \
		printf -- "--> %15s: Done\n" "$${OSARCH}" ; \
		pushd $${platform} >/dev/null 2>&1 ; \
		zip -q ../$(NAME)_$(VERSION)_$${OSARCH}.zip ./* ; \
		popd >/dev/null 2>&1 ; \
	done

	@ pushd ./$(BUILD_DIR) >/dev/null 2>&1 ; \
	shasum -a256 *.zip > ./$(NAME)_${VERSION}_SHA256SUMS ; \
	popd >/dev/null 2>&1 ; \
	printf -- "\n--> %15s: Done\n" "sha256sum"

#####################
## Release targets ##
#####################
.PHONY: authors
authors: log-authors ## Generate list of Authors
	@echo "# Authors\n" > AUTHORS.md
	git log --all --format='- %aN \<%aE\>' | sort -u | egrep -v noreply >> AUTHORS.md

.PHONY: changelog
changelog: log-changelog ## Generate content of Changelog
	git-chglog --output CHANGELOG.md

.PHONY: release patch minor major
PATTERN =

release: version ?= $(shell echo $(VERSION) | sed 's/^v//' | awk -F'[ .]' '{print $(PATTERN)}')
release: push ?= false
release: SHELL :=/bin/bash
release: log-release ## Prepare Module release
	@ if [ -z "$(version)" ]; then \
		echo "Error: missing value for 'version'. e.g. 'make release version=x.y.z'"; \
	elif [ "v$(version)" = "$(VERSION)" ] ; then \
		echo "Error: provided version (v$(version)) exists."; \
	else \
		git tag --annotate --message "v$(version) Release" v$(version); \
		echo "Tag v$(version) Release"; \
		if [ $(push) = "true" ]; then \
			git push origin v$(version); \
			echo "Push v$(version) Release"; \
		fi \
	fi

patch: PATTERN = '\$$1\".\"\$$2\".\"\$$3+1'
patch: release ## Prepare Module Patch release

minor: PATTERN = '\$$1\".\"\$$2+1\".0\"'
minor: release ## Prepare Module Minor release

major: PATTERN = '\$$1+1\".0.0\"'
major: release ## Prepare Module Major release

####################################
## Self-Documenting Makefile Help ##
####################################
.PHONY: help
.DEFAULT_GOAL := help
help:
	@grep -h -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "$(BLUE_COLOR)%-20s$(NO_COLOR) %s\n", $$1, $$2}'

log-%:
	@grep -h -E '^$*:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "$(BLUE_COLOR)==> %s$(NO_COLOR)\n", $$2}'
