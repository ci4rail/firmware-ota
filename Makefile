BIN_DIR ?= ./bin

all: build

build: io4edge-devsim io4edge-cli

test:
	go test ./...

clean:
	${MAKE} -C io4edge-devsim clean
	rm -rf bin/example

io4edge-devsim:
	${MAKE} -C cmd/io4edge-devsim build

io4edge-cli:
	${MAKE} -C cmd/io4edge-cli build


.PHONY: all build clean test io4edge-devsim io4edge-cli
