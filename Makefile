NAME=malice
ARCH=$(shell uname -m)
VERSION=0.1.0-alpha

build:
	# mkdir -p build/Linux  && GOOS=linux  go build -ldflags "-X main.Version=$(VERSION)" -o build/Linux/$(NAME)
	mkdir -p build/Darwin && GOOS=darwin go build -ldflags "-X main.Version=$(VERSION)" -o build/Darwin/$(NAME)

deps:
	go get -u github.com/progrium/gh-release/...
	go get -u -f github.com/tools/godep
	go get || true

test: deps build
	echo "Built" || true
# 	tests/shunit2 tests/*.sh

release: build
	rm -rf release && mkdir release
	# tar -zcf release/$(NAME)_$(VERSION)_linux_$(ARCH).tgz -C build/Linux $(NAME)
	tar -zcf release/$(NAME)_$(VERSION)_darwin_$(ARCH).tgz -C build/Darwin $(NAME)
	gh-release create maliceio/$(NAME) $(VERSION) $(shell git rev-parse --abbrev-ref HEAD) v$(VERSION)

destroy:
	rm -rf release
	rm -rf build
	gh-release destroy maliceio/$(NAME) $(VERSION) $(shell git rev-parse --abbrev-ref HEAD) v$(VERSION)

.PHONY: release build destroy
