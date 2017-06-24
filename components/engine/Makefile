REPO=malice
NAME=engine
VERSION=$(shell cat VERSION)

all: gotest build size test

build:
	docker build -t $(REPO)/$(NAME):$(VERSION) .

size:
	sed -i.bu 's/docker%20image-.*-blue/docker%20image-$(shell docker images --format "{{.Size}}" $(REPO)/$(NAME):$(VERSION)| cut -d' ' -f1)-blue/' README.md

tags:
	docker images --format "table {{.Repository}}\t{{.Tag}}\t{{.Size}}" $(REPO)/$(NAME)

tar:
	docker save $(REPO)/$(NAME):$(VERSION) -o $(NAME).tar

gotest:
	go test -v

test:
	docker-compose -f ./docker-compose.ci.yml up -d
	docker-compose -f docker-compose.ci.yml run httpie http://engine:3333/login username=admin password=admin

clean:
	docker-clean stop
	docker rmi maliceengine_httpie
	docker rmi $(REPO)/$(NAME)
	docker rmi $(REPO)/$(NAME):$(VERSION)

.PHONY: build dev size tags test gotest clean
