#
# chart-releaser
#

BIN_NAME    := chart-releaser
BIN_VERSION := v0.1.4
IMG_NAME    := chartreleaser/chart-releaser

GIT_COMMIT  ?= $(shell git rev-parse --short HEAD 2> /dev/null || true)
GIT_TAG     ?= $(shell git describe --tags 2> /dev/null || true)
BUILD_DATE  := $(shell date -u +%Y-%m-%dT%T 2> /dev/null)
GO_VERSION  := $(shell go version | awk '{ print $$3 }')

PKG_CTX := github.com/edaniszewski/chart-releaser/pkg
LDFLAGS := -w \
	-X ${PKG_CTX}.BuildDate=${BUILD_DATE} \
	-X ${PKG_CTX}.Commit=${GIT_COMMIT} \
	-X ${PKG_CTX}.Tag=${GIT_TAG} \
	-X ${PKG_CTX}.GoVersion=${GO_VERSION} \
	-X ${PKG_CTX}.Version=${BIN_VERSION}

.PHONY: build build-linux clean cover docker fmt github-tag lint test version help
.DEFAULT_GOAL := help


build:  ## Build the binary
	CGO_ENABLED=0 go build -ldflags "${LDFLAGS}" -o ${BIN_NAME} cmd/chart_releaser.go

build-linux:  ## Build the binary for linux amd64
	GO_OS=linux GO_ARCH=amd64 CGO_ENABLED=0 go build -ldflags "${LDFLAGS}" -o ${BIN_NAME} cmd/chart_releaser.go

clean:  ## Clean build and test artifacts
	@rm ${BIN_NAME}
	@rm coverage.out

cover:  ## Open a coverage report
	go tool cover -html=coverage.out

docker:  ## Build the docker image
	docker build -t ${IMG_NAME} .

fmt:  ## Run goimports formatting on all go files
	@find . -name '*.go' -not -wholename './vendor/*' | while read -r file; do goimports -w "$$file"; done

github-tag:  ## Create and push a tag with the current version
	git tag -a ${BIN_VERSION} -m "${BIN_NAME} version ${BIN_VERSION}"
	git push -u origin ${BIN_VERSION}

lint:  ## Lint project source files
	golangci-lint run ./...

test:  ## Run project unit tests
	go test --race -coverprofile=coverage.out -covermode=atomic ./...

version:  ## Print the version of the project
	@echo "${BIN_VERSION}"

help:  ## Print usage information
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "\033[36m%-15s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST) | sort
