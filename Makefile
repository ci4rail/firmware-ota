BIN_DIR ?= ./bin

all: build

build: io4edge-devsim io4edge-cli

test:
	go test ./...

clean:
	${MAKE} -C io4edge-devsim clean
	rm -rf bin/example

io4edge-devsim: proto
	${MAKE} -C cmd/io4edge-devsim build

io4edge-cli: proto
	${MAKE} -C cmd/io4edge-cli build


proto:
	go get google.golang.org/protobuf/cmd/protoc-gen-go
	go install google.golang.org/protobuf/cmd/protoc-gen-go
	protoc -I=./proto ./proto/io4edge_base_function.proto --go_out=.

.PHONY: all build clean test io4edge-devsim io4edge-cli proto
