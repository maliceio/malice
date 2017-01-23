NAME=malice
ARCH=$(shell uname -m)
# VERSION=$(shell gorram github.com/maliceio/malice/version GetHumanVersion)
VERSION=0.2.0-alpha

GIT_COMMIT=$(git rev-parse HEAD)
GIT_DIRTY=$(test -n "`git status --porcelain`" && echo "+CHANGES" || true)
GIT_DESCRIBE=$(git describe --tags)

all: deps test validate

bindata:
	rm -f config/bindata.go plugins/bindata.go
	go-bindata -pkg config -ignore=load.go config/...
	mv bindata.go config/bindata.go
	go-bindata -pkg plugins -ignore="^.*.go" plugins/...
	mv bindata.go plugins/bindata.go

docker:
	docker build -t malice/build-linux-binaries -f .docker/Dockerfile.binaries .

build: bindata docker
	@echo "[Building Binaries]"
	docker run --rm -v `pwd`:/go/src/github.com/maliceio/malice:rw -e NAME=$(NAME) -e VERSION=$(VERSION) malice/build-linux-binaries

deps:
	go get -u github.com/progrium/gh-release/...
	go get -u -f github.com/tools/godep
	go get github.com/golang/lint/golint
	go get -u github.com/jteeuwen/go-bindata/...
	go get -u npf.io/gorram
	go get -t ./... || true

test:
	go test -race -cover ./...

validate: lint
	go vet ./...
	test -z "$(gofmt -s -l . | tee /dev/stderr)"

lint:
	out="$$(golint ./...)"; \
	if [ -n "$$(golint ./...)" ]; then \
		echo "$$out"; \
		exit 1; \
	fi

release: build
	rm -rf release && mkdir release
	go get github.com/progrium/gh-release/...
	cp build/*.zip release
	gh-release create maliceio/$(NAME) $(VERSION) $(shell git rev-parse --abbrev-ref HEAD) v$(VERSION)

destroy:
	rm -rf release
	rm -rf build
	gh-release destroy maliceio/$(NAME) $(VERSION) $(shell git rev-parse --abbrev-ref HEAD) v$(VERSION)

.PHONY: all release build destroy
