# vim: ft=make ts=4

SRC=$(shell find . -name *.go)
BINARY=perfApp
GO_BUILD_RECIPE:=CGO_ENABLED=0 go build
ARCHS=amd64 ppc64le arm64 s390x

build: clean
	mkdir -p build
	for arch in $(ARCHS); do \
	  CGO_ENABLED=0 GOOS=linux GOARCH=$$arch go build -o build/$(BINARY)-$$arch cmd/perfApp/perfApp.go; \
	done

clean:
	rm -rf build
