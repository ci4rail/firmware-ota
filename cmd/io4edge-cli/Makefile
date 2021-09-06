NAME = io4edge-cli
BIN_DIR ?= ../../bin
VERSION ?= $(shell git describe --match=NeVeRmAtCh --always --abbrev=40 --dirty)
GO_LDFLAGS = -ldflags "-X github.com/ci4rail/firmware-ota/cmd/io4edge-cli/pkg/version.Version=$(VERSION)"

all: test build-io4edge-cli

build: build-io4edge-cli

build-io4edge-cli:
	go mod tidy
	GOOS=linux go build $(GO_LDFLAGS) -o ${BIN_DIR}/${NAME} main.go

test:
	go test ./...

clean:
	rm -f ${BIN_DIR}/${NAME}

.PHONY: all proto build build-io4edge-cli
