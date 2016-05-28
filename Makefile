OBSERVER_GO_EXECUTABLE ?= go
VERSION := $(shell git describe --tags)
DIST_DIRS := find * -type d -exec

build:
	${OBSERVER_GO_EXECUTABLE} get ./...
	${OBSERVER_GO_EXECUTABLE} build -o observer -ldflags "-X main.version=${VERSION}" observer.go

install: build
	install -d ${DESTDIR}/usr/local/bin/
	install -m 755 ./observer ${DESTDIR}/usr/local/bin/observer

test:
	${OBSERVER_GO_EXECUTABLE} get github.com/axw/gocov/gocov
	${OBSERVER_GO_EXECUTABLE} get github.com/AlekSi/gocov-xml
	${OBSERVER_GO_EXECUTABLE} get gopkg.in/matm/v1/gocov-html
	rm -rf coverage.html coverage.xml
	sh coverage.sh --html

style:
	${OBSERVER_GO_EXECUTABLE} vet ./...

clean:
	${OBSERVER_GO_EXECUTABLE} clean -i ./...
	rm -f ./observer.test
	rm -f ./observer
	rm -rf ./dist

bootstrap-dist:
	${OBSERVER_GO_EXECUTABLE} get -u github.com/mitchellh/gox

build-all:
	gox -verbose \
	-ldflags "-X main.version=${VERSION}" \
	-os="linux darwin windows " \
	-arch="amd64 386" \
	-output="dist/{{.OS}}-{{.Arch}}/{{.Dir}}" .

dist: build-all
	cd dist && \
	$(DIST_DIRS) cp ../LICENSE {} \; && \
	$(DIST_DIRS) cp ../README.md {} \; && \
	$(DIST_DIRS) tar -zcf observer-${VERSION}-{}.tar.gz {} \; && \
	$(DIST_DIRS) zip -r observer-${VERSION}-{}.zip {} \; && \
	cd ..


.PHONY: build test install clean bootstrap-dist build-all dist