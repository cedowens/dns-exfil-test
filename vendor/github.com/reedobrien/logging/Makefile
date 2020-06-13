.SILENT: ; # no need for @

PROJECT			=logging
PROJECT_DIR		=$(shell pwd)
GOFILES         :=$(shell find . -name '*.go' -not -path './vendor/*')
GOPACKAGES      :=$(shell go list ./... | grep -v /vendor/| grep -v /checkers)
OS              := $(shell go env GOOS)
ARCH            := $(shell go env GOARCH)

GITHASH         :=$(shell git rev-parse --short HEAD)
GITBRANCH       :=$(shell git rev-parse --abbrev-ref HEAD)
GITTAGORBRANCH 	:=$(shell sh -c 'git describe --always --dirty 2>/dev/null')
BUILDDATE      	:=$(shell date -u +%Y%m%d%H%M)
GO_LDFLAGS		?= -s -w
GO_BUILD_FLAGS  :=-ldflags "${GOLDFLAGS} -X main.BuildVersion=${GITTAGORBRANCH} -X main.GitHash=${GITHASH} -X main.GitBranch=${GITBRANCH} -X main.BuildDate=${BUILDDATE}"

GOFLAGS		:="-mod=vendor"
TOOLSDIR 	:=$(PROJECT_DIR)/_tools

CI_LINT_VERSION := 1.22.2

## What if there's no CIRCLE_BUILD_NUM
ifeq ($$CIRCLE_BUILD_NUM, "")
		BUILD_NUM:=""
else
		CB:=$$CIRCLE_BUILD_NUM
		BUILD_NUM:=$(CB)/
endif

WORKDIR         :=$(PROJECT_DIR)/_workdir

default: build-linux

build-linux:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o $(WORKDIR)/$(PROJECT)_linux_amd64 $(GO_BUILD_FLAGS)

build:
	CGO_ENABLED=0 go build -o $(WORKDIR)/$(PROJECT)_$(OS)_$(ARCH) $(GO_BUILD_FLAGS)

clean:
	rm -f $(WORKDIR)/*
	rm -rf .cover
	go clean -r

circle-ready: dependencies
	GOFLAGS=$(GOFLAGS) GOPRIVATE=$(GOPRIVATE) go mod tidy
	GOFLAGS=$(GOFLAGS) GOPRIVATE=$(GOPRIVATE) go mod vendor

coverage:
	./_misc/coverage.sh

coverage-html:
	./_misc/coverage.sh --html

dependencies:
	echo "Installing/Upgrading golangci-lint..."
	CGO_ENABLED=1 GOFLAGS=$(GOFLAGS) $(TOOLSDIR)/golangci-lint --version | grep -q " $(CI_LINT_VERSION) " || curl -sfL https://install.goreleaser.com/github.com/golangci/golangci-lint.sh | sh -s -- -b $(TOOLSDIR) v$(CI_LINT_VERSION)

develop: dependencies
	(cd .git/hooks && ln -sf ../../_misc/pre-push.bash pre-push )

lint:
	echo "golangci-lint..."
	CGO_ENABLED=1 GOFLAGS=$(GOFLAGS) $(TOOLSDIR)/golangci-lint run --deadline=5m --enable-all --exclude-use-default=false ./...

mod-update:
	GOPRIVATE=$(GOPRIVATE) go get -u ./...
	GOFLAGS=$(GOFLAGS) GOPRIVATE=$(GOPRIVATE) go mod tidy
	GOFLAGS=$(GOFLAGS) GOPRIVATE=$(GOPRIVATE) go mod vendor


test:
	CGO_ENABLED=0 go test $(GOPACKAGES)

test-race:
	CGO_ENABLED=1 go test -race $(GOPACKAGES)
