REPO=malice
NAME=engine
BUILD ?= 0.1
LATEST ?= 0.1
NEWSIZE = $(shell docker images --format "{{.Size}}" $(REPO)/$(NAME):$(BUILD))

build:
	cd $(BUILD); docker build -t $(REPO)/$(NAME):$(BUILD) .

size: build
ifeq "$(BUILD)" "$(LATEST)"
	sed -i.bu 's/docker image-.*-blue/docker image-$(NEWSIZE)-blue/' README.md
	sed -i.bu '/latest/ s/[0-9.]\{3,5\} MB/$(NEWSIZE)/' README.md
endif
	sed -i.bu '/$(BUILD)/ s/[0-9.]\{3,5\} MB/$(NEWSIZE)/' README.md

tags:
	docker images --format "table {{.Repository}}\t{{.Tag}}\t{{.Size}}" $(REPO)/$(NAME)

.PHONY: build size tags
