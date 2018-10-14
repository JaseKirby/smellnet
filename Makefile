# run from git bash only on windows
APPNAME=smellnet
VERSION ?= vlatest
GOBIN=$(shell pwd)/bin
GOFILES=$(wildcard *.go)

build:
	GOOS=windows GOARCH=amd64 GOBIN=$(GOBIN) go build -o bin/$(APPNAME).exe

clean:
	echo "Cleaning"
	GOBIN=$(GOBIN) go clean
	rm -rf $(GOBIN)/*

windows:
	GOOS=windows GOARCH=amd64 GOBIN=$(GOBIN) go build -o bin/release/$(APPNAME)-$(VERSION)-windows-amd64/$(APPNAME).exe
	zip -r bin/release/$(APPNAME)-$(VERSION)-windows-amd64.zip bin/release/$(APPNAME)-$(VERSION)-windows-amd64
linux:
	GOOS=windows GOARCH=amd64 GOBIN=$(GOBIN) go build -o bin/release/$(APPNAME)-$(VERSION)-linux-amd64/$(APPNAME)
	tar -zcvf bin/release/$(APPNAME)-$(VERSION)-linux-amd64.tar.gz bin/release/$(APPNAME)-$(VERSION)-linux-amd64

release: windows linux

.PHONY: build clean windows linux release
