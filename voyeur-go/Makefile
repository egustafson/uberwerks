# I am a -*-Makefile-*-
#
# --------------------------------------------------
#

GIT_SUMMARY := $(shell git describe --tags --dirty --always)
BUILD_DATE  := $(shell date -u "+%Y-%m-%dT%H:%M:%SZ")

# --------------------------------------------------
GO_FLAGS = 

.PHONY: all
all: build test

.PHONY: build
build: voyeur-go

.PHONY: voyeur-go  # force a rebuild always
voyeur-go:
	go build -ldflags "-X main.GitSummary=$(GIT_SUMMARY) -X main.BuildDate=${BUILD_DATE}" ${GO_FLAGS}

.PHONY: test
test:
	go test ./...

.PHONY: lint
lint:

.PHONY: package
package:
	echo "package as (possibly) multile distros including deb"

.PHONY: clean
clean:
	go clean -cache
	go clean ./...
