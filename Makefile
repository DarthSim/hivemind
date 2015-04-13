NAME = hivemind
VENDOR = ${PWD}/_vendor
BIN = ${PWD}/bin/${NAME}

GOBIN ?= ${GOPATH}/bin
GOPATH := ${VENDOR}:${GOPATH}

export GOPATH

.PHONY: all clean build prepare

all: build

clean:
	rm -rf bin/

build: prepare
	go build -v -o ${BIN}

prepare:
	mkdir -p ${VENDOR}
	go get -d -v

install: build
	cp ${BIN} ${GOBIN}
