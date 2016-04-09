GOTOOLS=\
	github.com/tools/godep \
	github.com/mitchellh/gox \
	golang.org/x/tools/cmd/cover \
	golang.org/x/tools/cmd/vet \
	github.com/jteeuwen/go-bindata/... \
	github.com/elazarl/go-bindata-assetfs/...
# DEPS = $(shell go list -f '{{range .TestImports}}{{.}} {{end}}' ./...)
PACKAGES = $(shell go list ./...)
VETARGS?=-asmdecl -atomic -bool -buildtags -copylocks -methods \
         -nilfunc -printf -rangeloops -shift -structtags -unsafeptr
VERSION?=$(shell awk -F\" '/^const Version/ { print $$2; exit }' version.go)

all: deps static-assets

# bin generates the releaseable binaries for Vault
bin: generate
	@sh -c "'$(CURDIR)/scripts/build.sh'"

# dev creates binaries for testing Vault locally. These are put
# into ./bin/ as well as $GOPATH/bin
dev: generate
	@MALICE_DEV=1 sh -c "'$(CURDIR)/scripts/build.sh'"

start:
	echo Starting malice machine

deps:
	@echo "--> Installing build dependencies"
	@go get -v $(GOTOOLS)
	# @go get -d -v ./... $(DEPS)

vet:
	@go tool vet 2>/dev/null ; if [ $$? -eq 3 ]; then \
		go get golang.org/x/tools/cmd/vet; \
	fi
	@echo "--> Running go tool vet $(VETARGS) ."
	@go tool vet $(VETARGS) . ; if [ $$? -eq 1 ]; then \
		echo ""; \
		echo "Vet found suspicious constructs. Please check the reported constructs"; \
		echo "and fix them if necessary before submitting the code for reviewal."; \
	fi

# generate runs `go generate` to build the dynamically generated source files
generate: deps
	find . -type f -name '.DS_Store' -delete
	go generate ./...

# generates the static web ui
static-assets: deps
	@echo "--> Generating static assets"
	@go-bindata -pkg data ./data/
	@mv bindata.go data/

.PHONY: all bin dev dist cov deps test vet web web-push generate test-nodep static-assets
