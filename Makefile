package = github.com/cultureamp/cfparams

.PHONY: release
release: cfparams-darwin-amd64.gz cfparams-linux-amd64.gz

%.gz: %
	gzip $<

%-darwin-amd64:
	GOOS=darwin GOARCH=amd64 go build -o $@ $(package)

%-linux-amd64:
	GOOS=linux GOARCH=amd64 go build -o $@ $(package)
