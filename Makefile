# Project variables
ORG         := cloud-ca
NAME        := terraform-provider-cloudca
DESCRIPTION := Terraform provider to interact with cloud.ca infrastructure
AUTHOR      := cloud.ca
URL         := https://github.com/cloud-ca/terraform-provider-cloudca
LICENSE     := MIT

# Repository variables
PACKAGE     := github.com/$(ORG)/$(NAME)

# Build variables
BUILD_DIR   := bin
COMMIT_HASH ?= $(shell git rev-parse --short HEAD 2>/dev/null)
VERSION     ?= $(shell git describe --tags --exact-match 2>/dev/null || git describe --tags 2>/dev/null || echo "v0.0.1-$(COMMIT_HASH)")
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
GITCHGLOG_VERSION := 0.8.0
GOLANGCI_VERSION  := v1.18.0

.PHONY: default
default: help

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

.PHONY: tidy
tidy: ## Tidy up 'vendor' dependencies
	@ $(MAKE) --no-print-directory log-$@
	$(GOCMD) mod tidy

.PHONY: lint
lint: ## Run linter
	@ $(MAKE) --no-print-directory log-$@
	GO111MODULE=on golangci-lint run ./...

.PHONY: fmt
fmt: ## Format go files
	@ $(MAKE) --no-print-directory log-$@
	goimports -w $(GOFILES)

.PHONY: checkfmt
checkfmt: RESULT ?= $(shell goimports -l $(GOFILES) | tee >(if [ "$$(wc -l)" = 0 ]; then echo "OK"; fi))
checkfmt: SHELL  := /usr/bin/env bash
checkfmt: ## Check formatting of go files
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
	GOOS=$(GOOS) GOARCH=$(GOARCH) CGO_ENABLED=0 $(GOBUILD) -o ./$(BUILD_DIR)/$(GOOS)-$(GOARCH)/$(NAME)_$(VERSION)

.PHONY: build-all
build-all: clean ## Build binaries for all OS/ARCH
	@ $(MAKE) --no-print-directory log-$@
	@ ./scripts/build/build-all.sh "$(BUILD_DIR)" "$(VERSION)" "$(GOOS)" "$(GOARCH)" $(GOLDFLAGS)
	@ ./scripts/build/compress.sh "$(BUILD_DIR)" "$(NAME)" "$(VERSION)"

#####################
## Release targets ##
#####################
.PHONY: release patch minor major
PATTERN =

release: version ?= $(shell echo $(VERSION) | sed 's/^v//' | awk -F'[ .]' '{print $(PATTERN)}')
release: push    ?= false
release: ## Prepare release
	@ $(MAKE) --no-print-directory log-$@
	@ ./scripts/release/release.sh "$(version)" "$(push)" "$(VERSION)" "1"

patch: PATTERN = '\$$1\".\"\$$2\".\"\$$3+1'
patch: release ## Prepare Patch release

minor: PATTERN = '\$$1\".\"\$$2+1\".0\"'
minor: release ## Prepare Minor release

major: PATTERN = '\$$1+1\".0.0\"'
major: release ## Prepare Major release

####################
## Helper targets ##
####################
.PHONY: changelog
changelog: push ?= false
changelog: next ?=
changelog: ## Generate Changelog
	@ $(MAKE) --no-print-directory log-$@
	git-chglog --config ./scripts/chglog/config-full-history.yml --tag-filter-pattern v[0-9]+.[0-9]+.[0-9]+$$ --output CHANGELOG.md $(next)
	@ git add CHANGELOG.md
	@ git commit -m "Update Changelog"
	@ if $(push) = "true"; then git push origin master; fi

.PHONY: tools git-chglog goimports golangci gox

git-chglog:
	curl -sfL https://github.com/git-chglog/git-chglog/releases/download/$(GITCHGLOG_VERSION)/git-chglog_$(shell go env GOOS)_$(shell go env GOARCH) -o $(shell go env GOPATH)/bin/git-chglog && chmod +x $(shell go env GOPATH)/bin/git-chglog

goimports:
	GO111MODULE=off go get -u golang.org/x/tools/cmd/goimports

golangci:
	curl -sfL https://install.goreleaser.com/github.com/golangci/golangci-lint.sh | sh -s  -- -b $(shell go env GOPATH)/bin $(GOLANGCI_VERSION)

gox:
	GO111MODULE=off go get -u github.com/mitchellh/gox

tools: ## Install required tools
	@ $(MAKE) --no-print-directory log-$@
	@ $(MAKE) --no-print-directory git-chglog goimports golangci gox

####################################
## Self-Documenting Makefile Help ##
####################################
.PHONY: help
help:
	@ grep -h -E '^[a-zA-Z0-9_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

log-%:
	@ grep -h -E '^$*:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m==> %s\033[0m\n", $$2}'
