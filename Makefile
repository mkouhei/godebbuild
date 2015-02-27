#!/usr/bin/make -f
# -*- makefile -*-

BIN := debbuild
SRC := *.go
GOPKG := github.com/mkouhei/godebbuild/
GOPATH := $(CURDIR)/_build
export GOPATH
PATH := $(CURDIR)/_build/bin:$(PATH)
export PATH
# "FLAGS=" when no update package
FLAGS := -u
# "FUNC=-html" when generate HTML coverage report
FUNC := -func


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
	@rm -rf _build/$(BIN) $(GOPATH)/src/$(GOPKG)

format:
	for src in $(SRC); do \
		gofmt -w $$src ;\
		goimports -w $$src ;\
	done


test: prebuild
	go get $(FLAGS) golang.org/x/tools/cmd/goimports
	go get $(FLAGS) github.com/golang/lint/golint
	go get $(FLAGS) golang.org/x/tools/cmd/vet
	go get $(FLAGS) golang.org/x/tools/cmd/cover
	_build/bin/golint
	go vet
	go test -v -covermode=count -coverprofile=c.out $(GOPKG)
	go tool cover $(FUNC)=c.out
	unlink c.out
	rm -f main.test
