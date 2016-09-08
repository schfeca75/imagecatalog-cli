BINARY=imagecatalog

VERSION=1.0
BUILD_TIME=$(shell date +%FT%T)
LDFLAGS=-ldflags "-X github.com/hortonworks/imagecatalog-cli/cli.Version=${VERSION} -X github.com/hortonworks/imagecatalog-cli/cli.BuildTime=${BUILD_TIME}"

format:
	gofmt -w .

build: format build-darwin build-linux

build-darwin:
	GOOS=darwin go build -a -installsuffix cgo ${LDFLAGS} -o build/Darwin/${BINARY} main.go

build-linux:
	GOOS=linux go build -a -installsuffix cgo ${LDFLAGS} -o build/Linux/${BINARY} main.go

.PHONY: build
