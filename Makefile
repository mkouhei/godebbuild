#!/usr/bin/make -f
# -*- makefile -*-

BIN := debbuild
SRC := *.go
GOPKG := github.com/mkouhei/godebbuild/
GOPATH := $(CURDIR)/_build
export GOPATH
PATH := $(CURDIR)/_build/bin:$(PATH)
export PATH


all: precheck clean test format build

precheck:
	@if [ -d .git ]; then \
	set -e; \
	diff -u .git/hooks/pre-commit utils/pre-commit.txt ;\
	[ -x .git/hooks/pre-commit ] ;\
	fi

prebuild:
	go get -d -v ./...
	install -d $(CURDIR)/_build/src/$(GOPKG)
	cp -a $(CURDIR)/*.go $(CURDIR)/examples $(CURDIR)/_build/src/$(GOPKG)


build: prebuild
	go build -ldflags "-X main.version $(shell git describe --always)" -o _build/$(BIN)

build-only:
	go build -ldflags "-X main.version $(shell git describe --always)" -o _build/$(BIN)

clean:
	@rm -f _build/$(BIN)

format:
	for src in $(SRC); do \
		gofmt -w $$src ;\
		goimports -w $$src ;\
	done


test: prebuild
	go get github.com/golang/lint/golint
	golint
	go vet
	go test -v -coverprofile=c.out $(GOPKG)
	go tool cover -func=c.out
	unlink c.out
