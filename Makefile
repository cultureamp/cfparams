PACKAGE = github.com/cultureamp/cfparams
VERSION = $(shell git describe --tags --candidates=1 --dirty)
FLAGS=-X main.Version=$(VERSION) -s -w

cfparams: main.go parameters.go tags.go template.go parameterstore/store.go
	go build -ldflags="$(FLAGS)"

.PHONY: install
install:
	go install -ldflags="$(FLAGS)" $(PACKAGE)

.PHONY: release
release: \
	build/cfparams-$(VERSION)-darwin-amd64.tar.gz \
	build/cfparams-$(VERSION)-linux-amd64.tar.gz

%.tar.gz: %
	cp $< build/cfparams
	chmod 0755 build/cfparams
	tar czf $<.tar.gz -C build cfparams
	rm build/cfparams

%-darwin-amd64:
	GOOS=darwin GOARCH=amd64 go build -ldflags="$(FLAGS)" -o $@ $(PACKAGE)

%-linux-amd64:
	GOOS=linux GOARCH=amd64 go build -ldflags="$(FLAGS)" -o $@ $(PACKAGE)
