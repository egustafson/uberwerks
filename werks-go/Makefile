# I am a -*-Makefile-*-
#
# ##########################################################
#

GO_FLAGS =

.PHONY: all
all: build test

.PHONY: build
build:
	go build ${GO_FLAGS} ./...

.PHONY: test
test:
	go test ./...

.PHONY: clean
clean:
	go clean ./...
