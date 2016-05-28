OBSERVER_GO_EXECUTABLE ?= go
VERSION := $(shell git describe --tags)
DIST_DIRS := find * -type d -exec

build:
	CGO_ENABLED=0 go install -ldflags "-X main.version=${VERSION}" -v ./...

install:
	go get ./...

clean:
	go clean -i ./...
	rm -rf pkg bin

style:
	go vet ./...

test:
	go get github.com/axw/gocov/gocov
	go get github.com/AlekSi/gocov-xml
	go get gopkg.in/matm/v1/gocov-html
	rm -rf coverage.html coverage.xml
	sh coverage.sh --html

