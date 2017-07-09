REPO=malice
NAME=engine
VERSION=$(shell cat .release/VERSION)
MESSAGE?="New release"

# TODO remove \|/templates/\|/api
SOURCE_FILES?=$$(go list ./... | grep -v '/vendor/\|/templates/\|/api')
TEST_PATTERN?=.
TEST_OPTIONS?=

GIT_COMMIT=$(git rev-parse HEAD)
GIT_DIRTY=$(test -n "`git status --porcelain`" && echo "+CHANGES" || true)
GIT_DESCRIBE=$(git describe --tags)


bindata: ## Embed binary data in malice program
	@echo "===> Embedding Binary Data"
	# tomlupdate --path config/config.toml $(VERSION)
	rm -f config/bindata.go plugins/bindata.go
	go-bindata -pkg config -ignore=load.go config/...
	mv bindata.go config/bindata.go
	go-bindata -pkg plugins -ignore="^.*.go|\\.DS_Store" plugins/...
	mv bindata.go plugins/bindata.go

docker: ## Build docker image
	cd .docker; docker build -t $(REPO)/$(NAME):$(VERSION) .

size: docker ## Add docker image size to READMEs
	sed -i.bu 's/docker%20image-.*-blue/docker%20image-$(shell docker images --format "{{.Size}}" $(REPO)/$(NAME):$(VERSION)| cut -d' ' -f1)%20MB-blue/' README.md
	sed -i.bu 's/docker%20image-.*-blue/docker%20image-$(shell docker images --format "{{.Size}}" $(REPO)/$(NAME):$(VERSION)| cut -d' ' -f1)%20MB-blue/' .docker/README.md

osx: ## Install OSX dev dependencies
	brew tap homebrew/bundle
	brew bundle
	gem install --no-ri --no-rdoc fpm

# TODO switch to golang/dep
setup: ## Install all the build and lint dependencies
	@echo "===> Installing deps"
	go get -u github.com/alecthomas/gometalinter
	# go get -u github.com/shurcooL/markdownfmt
	go get -u github.com/jteeuwen/go-bindata/...
	go get -u github.com/kardianos/govendor
	go get -u golang.org/x/tools/cmd/cover
	go get -u github.com/fatih/gomodifytags
	go get -u github.com/maliceio/malice/utils/tomlupdate
	govendor sync
	gometalinter --install

test: ## Run all the tests
	@echo "===> Running Tests"
	echo 'mode: atomic' > coverage.tmp
	$(SOURCE_FILES) | xargs -n1 -I{} sh -c 'go test -covermode=atomic -coverprofile=coverage.tmp {} && tail -n +2 coverage.tmp >> coverage.txt' && rm coverage.tmp

cover: test ## Run all the tests and opens the coverage report
	@echo "===> Running Cover"
	go tool cover -html=coverage.txt

fmt: ## gofmt and goimports all go files
	@echo "===> Formatting Go Files"
	find . -name '*.go' -not -wholename './vendor/*' | while read -r file; do gofmt -w -s "$$file"; goimports -w "$$file"; done

lint: ## Run all the linters
	@echo "===> Lintting"
	gometalinter --vendor --disable-all \
		--enable=deadcode \
		--enable=ineffassign \
		--enable=gosimple \
		--enable=staticcheck \
		--enable=gofmt \
		--enable=goimports \
		--enable=dupl \
		--enable=misspell \
		--enable=errcheck \
		--enable=vet \
		--enable=vetshadow \
		--deadline=10m \
		./...
	markdownfmt -w README.md
	markdownfmt -w CHANGELOG.md
	markdownfmt -w .release/RELEASE.md

release: ## Create a new release from the VERSION
	@echo "===> Creating Release"
	git tag -a $(VERSION) -m ${MESSAGE}
	git push origin $(VERSION)
	# goreleaser --release-notes .release/RELEASE.md
	goreleaser --rm-dist

destroy: ## Remove release from the VERSION
	@echo "===> Deleting Release"
	rm -rf dist
	git tag -d $(VERSION)
	git push origin :refs/tags/$(VERSION)

ci: lint test ## Run all the tests and code checks

build: bindata size ## Build a beta version of malice
	@echo "===> Building Binaries"
	go build -ldflags "-X main.version=${VERSION}-beta" -o malice-beta

# Absolutely awesome: http://marmelab.com/blog/2016/02/29/auto-documented-makefile.html
help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

.DEFAULT_GOAL := help
