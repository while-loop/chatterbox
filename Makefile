# service specific vars
TARGET	 := chatterbox
VERSION	 := 0.0.1

.PHONY: deps test build cont all deploy help clean lint
.DEFAULT_GOAL := help

help: ## help
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "\033[36m%-15s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)

all: deps gen lint test build ## get && test && build

build: clean gen lint ## build service binary file
	@echo "[build] building go binary"
	@go build

clean: ## remove service bin from $GOPATH/bin
	@echo "[clean] removing service files"
	rm -f ${TARGET}

cont: ## build a non-cached service container
	docker build -t ${TARGET} -t ${TARGET}:${VERSION} . --no-cache

deploy: ## deploy lastest built container to docker hub
	docker push ${TARGET}

deps: ## get service pkg + test deps
	@echo "[deps] getting go deps"
	go get -v -t ./...

gen: ## generate static www files
	@echo "[gen] generating binddata"
	go generate ./...

lint: ## apply golint
	@echo "[lint] applying go fmt & vet"
	go fmt ./...
	go vet ./...

release: test cont deploy ## build and deploy a docker container

test: gen lint ## test service code
	@echo "[test] running tests w/ cover"
	go test ./... -cover