BINARY = backpocket
GOARCH = amd64
GO_BUILD = GOARCH=${GOARCH} go build -mod=vendor
PREFIX ?= ${GOPATH}

all: clean linux darwin windows

clean:
	rm -rf bin/

vendor:
	go mod vendor
	go mod tidy

linux: vendor
	GOOS=linux ${GO_BUILD} -o bin/linux_${GOARCH}/${BINARY}

darwin: vendor
	GOOS=darwin ${GO_BUILD} -o bin/darwin_${GOARCH}/${BINARY}

windows: vendor
	GOOS=windows ${GO_BUILD} -o bin/windows_${GOARCH}/${BINARY}.exe

install:
	cp -Rf bin/ "${PREFIX}/bin"

release: all
	script/release.sh

.PHONY: all clean release vendor linux darwin windows install
