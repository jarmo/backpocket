BINARY = backpocket
GOARCH = amd64
GO_BUILD = GOARCH=${GOARCH} go build -mod=vendor

all: clean linux darwin windows

clean:
	rm -rf bin/

vendor:
	go mod vendor
	go mod tidy

examples:
	mv -f `go run . http://www.cnn.com/2017/06/12/politics/donald-trump-cabinet-meeting/index.html` examples/successful.html
	mv -f `go run . http://failed-article-url.com` examples/failed.html

linux: vendor
	GOOS=linux ${GO_BUILD} -o bin/linux_${GOARCH}/${BINARY}

darwin: vendor
	GOOS=darwin ${GO_BUILD} -o bin/darwin_${GOARCH}/${BINARY}

windows: vendor
	GOOS=windows ${GO_BUILD} -o bin/windows_${GOARCH}/${BINARY}.exe

release: examples all
	script/release.sh

.PHONY: all examples clean release vendor linux darwin windows install
