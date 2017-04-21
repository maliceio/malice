REPO=malice
NAME=malice
VERSION=$(shell cat VERSION)

SOURCE_FILES?=$$(go list ./... | grep -v '/vendor/\|/templates/\|/api')
TEST_PATTERN?=.
TEST_OPTIONS?=

GIT_COMMIT=$(git rev-parse HEAD)
GIT_DIRTY=$(test -n "`git status --porcelain`" && echo "+CHANGES" || true)
GIT_DESCRIBE=$(git describe --tags)


bindata: ## Embed binary data in malice program
	@echo "===> Embedding Binary Data"
	rm -f config/bindata.go plugins/bindata.go
	go-bindata -pkg config -ignore=load.go config/...
	mv bindata.go config/bindata.go
	go-bindata -pkg plugins -ignore="^.*.go|\\.DS_Store" plugins/...
	mv bindata.go plugins/bindata.go

docker:
	docker build -t malice/build-linux-binaries -f .docker/Dockerfile.binaries .

osx: ## Install OSX dev dependencies
	brew tap homebrew/bundle
	brew bundle
	gem install --no-ri --no-rdoc fpm

setup: ## Install all the build and lint dependencies
	@echo "===> Installing deps"
	go get -u github.com/alecthomas/gometalinter
	go get -u github.com/jteeuwen/go-bindata/...
	go get -u github.com/kardianos/govendor
	go get -u github.com/pierrre/gotestcover
	go get -u golang.org/x/tools/cmd/cover
	govendor sync
	gometalinter --install

test: ## Run all the tests
	@echo "===> Running Tests"
	gotestcover $(TEST_OPTIONS) -covermode=count -coverprofile=coverage.out $(SOURCE_FILES) -run $(TEST_PATTERN) -timeout=30s

cover: test ## Run all the tests and opens the coverage report
	@echo "===> Running Cover"
	go tool cover -html=coverage.out

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

release: bindata ## Create a new release from the VERSION
	@echo "===> Creating Release"
	git tag ${VERSION}
	git push origin ${VERSION}
	goreleaser

destroy: ## Remove release from the VERSION
	@echo "===> Deleting Release"
	rm -rf dist
	gh-release destroy maliceio/$(NAME) $(VERSION) $(shell git rev-parse --abbrev-ref HEAD) v$(VERSION)

ci: lint test ## Run all the tests and code checks

build: bindata ## Build a beta version of malice
	@echo "===> Building Binaries"
	go build -o malice-beta

# Absolutely awesome: http://marmelab.com/blog/2016/02/29/auto-documented-makefile.html
help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

.DEFAULT_GOAL := help
