CWD=$(shell pwd)
GOPATH := $(CWD)

prep:
	if test -d pkg; then rm -rf pkg; fi

self:   prep
	if test -s src; then rm -rf src; fi
	mkdir -p src/github.com/sfomuseum/go-html-offline
	cp *.go src/github.com/sfomuseum/go-html-offline/
	cp -r http src/github.com/sfomuseum/go-html-offline/
	cp -r server src/github.com/sfomuseum/go-html-offline/
	cp -r vendor/* src/

rmdeps:
	if test -d src; then rm -rf src; fi 

build:	fmt bin

deps:   
	@GOPATH=$(GOPATH) go get -u "github.com/facebookgo/atomicfile"
	@GOPATH=$(GOPATH) go get -u "github.com/whosonfirst/walk"
	@GOPATH=$(GOPATH) go get -u "github.com/whosonfirst/algnhsa"
	@GOPATH=$(GOPATH) go get -u "github.com/whosonfirst/go-whosonfirst-cli"
	@GOPATH=$(GOPATH) go get -u "golang.org/x/net/html"

vendor-deps: rmdeps deps
	if test ! -d vendor; then mkdir vendor; fi
	if test -d vendor; then rm -rf vendor; fi
	cp -r src vendor
	find vendor -name '.git' -print -type d -exec rm -rf {} +
	rm -rf src

fmt:
	go fmt cmd/*.go
	go fmt http/*.go
	go fmt server/*.go
	go fmt *.go

bin: 	rmdeps self
	rm -rf bin/*
	@GOPATH=$(shell pwd) go build -o bin/add-service-worker cmd/add-service-worker.go
	@GOPATH=$(shell pwd) go build -o bin/service-worker-inventoryd cmd/service-worker-inventoryd.go

lambda:
	@make self
	if test -f main; then rm -f main; fi
	if test -f deployment.zip; then rm -f deployment.zip; fi
	@GOPATH=$(GOPATH) GOOS=linux go build -o main cmd/service-worker-inventoryd.go
	zip deployment.zip main
	rm -f main

dist-build:
	OS=darwin make dist-os
	OS=windows make dist-os
	OS=linux make dist-os

dist-os:
	mkdir -p dist/$(OS)
	GOOS=$(OS) GOPATH=$(GOPATH) GOARCH=386 go build -o dist/$(OS)/add-service-worker cmd/add-service-worker.go
