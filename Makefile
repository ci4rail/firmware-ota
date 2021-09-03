BIN_DIR ?= ./bin

all: build 

build: netio-devsim netio-cli

test:
	go test ./...

clean:
	${MAKE} -C netio-devsim clean
	rm -rf bin/example

netio-devsim: proto
	${MAKE} -C cmd/netio-devsim build

netio-cli: proto
	${MAKE} -C cmd/netio-cli build


proto:
	go get google.golang.org/protobuf/cmd/protoc-gen-go
	go install google.golang.org/protobuf/cmd/protoc-gen-go
	protoc -I=./proto ./proto/netio_base_function.proto --go_out=.

.PHONY: all build clean test netio-devsim netio-cli proto