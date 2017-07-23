.PHONY: build dev size tags tar test run ssh circle node push dockerfile

REPO=maliceio/engine
ORG=malice
NAME=engine
VERSION=$(shell cat VERSION)

all: gotest build size test

build: ## Build docker image
	docker build -t $(ORG)/$(NAME):$(VERSION) .

size: tags ## Update docker image size in README.md
	sed -i.bu 's/docker%20image-.*-blue/docker%20image-$(shell docker images --format "{{.Size}}" $(ORG)/$(NAME):$(VERSION)| cut -d' ' -f1)-blue/' README.md

tags: ## Show all docker image tags
	docker images --format "table {{.Repository}}\t{{.Tag}}\t{{.Size}}" $(ORG)/$(NAME)

daemon: stop ## Run malice engine daemon
	@echo "===> Starting kibana elasticsearch..."
	@docker run --init -d --name engine $(ORG)/$(NAME):$(VERSION)

ssh: ## SSH into docker image
	@docker run --init -it --rm --entrypoint=sh $(ORG)/$(NAME):$(VERSION)

tar: ## Export tar of docker image
	docker save $(ORG)/$(NAME):$(VERSION) -o $(NAME).tar

gotest: ## Run go tests
	go test -v

test: ## Test malice engine
	docker-compose -f ./docker-compose.ci.yml up -d
	docker-compose -f docker-compose.ci.yml run httpie http://engine:3333/login username=admin password=admin

push: build ## Push docker image to docker registry
	@echo "===> Pushing $(ORG)/$(NAME):$(VERSION) to docker hub..."
	@docker push $(ORG)/$(NAME):$(VERSION)

circle: ci-size ## Get docker image size from CircleCI
	@sed -i.bu 's/docker%20image-.*-blue/docker%20image-$(shell cat .circleci/SIZE)-blue/' README.md
	@echo "===> Image size is: $(shell cat .circleci/SIZE)"

ci-build:
	@echo "===> Getting CircleCI build number"
	@http https://circleci.com/api/v1.1/project/github/${REPO} | jq '.[0].build_num' > .circleci/build_num

ci-size: ci-build
	@echo "===> Getting image build size from CircleCI"
	@http "$(shell http https://circleci.com/api/v1.1/project/github/${REPO}/$(shell cat .circleci/build_num)/artifacts circle-token==${CIRCLE_TOKEN} | jq '.[].url')" > .circleci/SIZE

clean: ## Clean docker image and stop all running containers
	docker-clean stop
	docker rmi maliceengine_httpie || true
	docker rmi $(ORG)/$(NAME) || true
	docker rmi $(ORG)/$(NAME):$(VERSION) || true

stop: ## Kill running malice-engine docker containers
	@docker rm -f engine || true

# Absolutely awesome: http://marmelab.com/blog/2016/02/29/auto-documented-makefile.html
help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

.DEFAULT_GOAL := help
