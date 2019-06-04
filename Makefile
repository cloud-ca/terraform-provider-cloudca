# Project variables
NAME        := terraform-provider-cloudca
DESCRIPTION := Terraform provider to interact with cloud.ca infrastructure
VENDOR      := cloud.ca
URL         := https://github.com/cloud-ca/terraform-provider-cloudca
LICENSE     := MIT

# Build variables
BUILD_DIR   := bin
COMMIT_HASH ?= $(shell git rev-parse --short HEAD 2>/dev/null)
VERSION     ?= $(shell git describe --tags --exact-match 2>/dev/null || git describe --tags 2>/dev/null || echo "v0.0.0-$(COMMIT_HASH)")
BUILD_DATE  ?= $(shell date +%FT%T%z)

# Go variables
GOOS        := linux darwin windows freebsd openbsd solaris
GOARCH      := 386 amd64 arm
GOCMD       := GO111MODULE=on go
MODVENDOR   := -mod=vendor
GOPKGS      ?= $(shell $(GOCMD) list $(MODVENDOR) ./... | grep -v /vendor)
GOFILES     ?= $(shell find . -type f -name '*.go' -not -path "./vendor/*")

GOLDFLAGS   :="
GOLDFLAGS   += -X $(PACKAGE)/main.version=$(VERSION)
GOLDFLAGS   += -X $(PACKAGE)/main.commitHash=$(COMMIT_HASH)
GOLDFLAGS   += -X $(PACKAGE)/main.buildDate=$(BUILD_DATE)
GOLDFLAGS   +="

GOBUILD     := $(GOCMD) build $(MODVENDOR) -ldflags $(GOLDFLAGS)

# Binary versions
GOLANGCI_VERSION  := v1.16.0
GITCHGLOG_VERSION := 0.8.0
GOX_VERSION       := v1.0.1

.DEFAULT_GOAL := help

.PHONY: version
version: ## Show version of provider
	@ echo "$(NAME) - $(VERSION) - $(BUILD_DATE)"

#########################
## Development targets ##
#########################
.PHONY: clean
clean: ## Clean builds
	@ $(MAKE) --no-print-directory log-$@
	rm -rf ./$(BUILD_DIR) $(NAME)

.PHONY: vendor
vendor: ## Install 'vendor' dependencies
	@ $(MAKE) --no-print-directory log-$@
	$(GOCMD) mod vendor

.PHONY: verify
verify: ## Verify 'vendor' dependencies
	@ $(MAKE) --no-print-directory log-$@
	$(GOCMD) mod verify

.PHONY: lint
lint: ## Run linter
	@ $(MAKE) --no-print-directory log-$@
	GO111MODULE=on golangci-lint run ./...

.PHONY: format
format: ## Format all go files
	@ $(MAKE) --no-print-directory log-$@
	$(GOCMD) fmt $(GOPKGS)

.PHONY: checkfmt
checkfmt: RESULT ?= $(shell gofmt -l $(GOFILES) | tee >(if [ "$$(wc -l)" = 0 ]; then echo "OK"; fi))
checkfmt: SHELL  := /bin/bash
checkfmt: ## Check formatting of all go files
	@ $(MAKE) --no-print-directory log-$@
	@ echo "$(RESULT)"
	@ if [ "$(RESULT)" != "OK" ]; then exit 1; fi

.PHONY: test
test: ## Run tests
	@ $(MAKE) --no-print-directory log-$@
	$(GOCMD) test $(MODVENDOR) -v $(GOPKGS)

###################
## Build targets ##
###################
.PHONY: build
build: GOOS   := $(shell go env GOOS)
build: GOARCH := $(shell go env GOARCH)
build: clean ## Build binary for current OS/ARCH
	@ $(MAKE) --no-print-directory log-$@
	@ GOOS=$(GOOS) GOARCH=$(GOARCH) CGO_ENABLED=0 $(GOBUILD) -o ./$(BUILD_DIR)/$(GOOS)-$(GOARCH)/$(NAME)_$(VERSION) && echo "./$(BUILD_DIR)/$(GOOS)-$(GOARCH)/$(NAME)_$(VERSION)"

.PHONY: build-all
build-all: SHELL := /bin/bash
build-all: clean ## Build binary for all OS/ARCH
	@ $(MAKE) --no-print-directory log-$@
	@ CGO_ENABLED=0 gox -verbose \
		-ldflags "$(LDFLAGS)" \
		-gcflags=-trimpath=$(GOPATH) \
		-os="$(GOOS)" \
		-arch="$(GOARCH)" \
		-osarch="!darwin/arm !darwin/386" \
		-output="$(BUILD_DIR)/{{.OS}}-{{.Arch}}/{{.Dir}}_$(VERSION)" .

	@ printf "\n"
	@ for platform in `find ./$(BUILD_DIR) -mindepth 1 -maxdepth 1 -type d` ; do \
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

####################
## Helper targets ##
####################
.PHONY: tools
tools: ## Install required tools
	@ $(MAKE) --no-print-directory log-$@
	@ curl -sfL https://install.goreleaser.com/github.com/golangci/golangci-lint.sh | sh -s  -- -b $(shell go env GOPATH)/bin $(GOLANGCI_VERSION)
	@ curl -sfL https://github.com/git-chglog/git-chglog/releases/download/$(GITCHGLOG_VERSION)/git-chglog_$(shell go env GOOS)_$(shell go env GOARCH) -o $(shell go env GOPATH)/bin/git-chglog && chmod +x $(shell go env GOPATH)/bin/git-chglog
	@ cd /tmp && go get -v -u github.com/mitchellh/gox

.PHONY: authors
authors: ## Generate list of Authors
	@ $(MAKE) --no-print-directory log-$@
	@ echo "# Authors\n" > AUTHORS.md
	git log --all --format='- %aN \<%aE\>' | sort -u | egrep -v noreply >> AUTHORS.md

.PHONY: changelog
changelog: ## Generate content of Changelog
	@ $(MAKE) --no-print-directory log-$@
	git-chglog --output CHANGELOG.md

#####################
## Release targets ##
#####################
.PHONY: release patch minor major
PATTERN =

release: version ?= $(shell echo $(VERSION) | sed 's/^v//' | awk -F'[ .]' '{print $(PATTERN)}')
release: push    := false
release: SHELL   := /bin/bash
release: ## Prepare release
	@ $(MAKE) --no-print-directory log-$@
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
help:
	@ grep -h -E '^[a-zA-Z0-9_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

log-%:
	@ grep -h -E '^$*:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m==> %s\033[0m\n", $$2}'
