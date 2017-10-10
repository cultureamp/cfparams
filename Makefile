PACKAGE = github.com/cultureamp/cfparams
VERSION = $(shell git describe --tags --candidates=1 --dirty)
FLAGS=-X main.Version=$(VERSION) -s -w

cfparams: main.go
	go build -ldflags="$(FLAGS)"

.PHONY: install
install:
	go install -ldflags="$(FLAGS)" $(PACKAGE)

.PHONY: release
release: cfparams-$(VERSION)-darwin-amd64.gz cfparams-$(VERSION)-linux-amd64.gz

%.gz: %
	gzip $<

%-darwin-amd64:
	GOOS=darwin GOARCH=amd64 go build -ldflags="$(FLAGS)" -o $@ $(PACKAGE)

%-linux-amd64:
	GOOS=linux GOARCH=amd64 go build -ldflags="$(FLAGS)" -o $@ $(PACKAGE)
