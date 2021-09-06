NAME = io4edge-devsim
BIN_DIR ?= ../../bin
VERSION ?= $(shell git describe --match=NeVeRmAtCh --always --abbrev=40 --dirty)
GO_LDFLAGS = -ldflags "-X github.com/ci4rail/firmware-ota/cmd/io4edge-devsim/pkg/version.Version=$(VERSION)"

all: test build-io4edge-devsim

build: build-io4edge-devsim

build-io4edge-devsim:
	go mod tidy
	GOOS=linux go build $(GO_LDFLAGS) -o ${BIN_DIR}/${NAME} main.go

test:
	go test ./...

clean:
	rm -f ${BIN_DIR}/${NAME}

.PHONY: all proto build test clean clean build-io4edge-devsim
