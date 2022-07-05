SHELL:=/usr/bin/env bash
.EXPORT_ALL_VARIABLES:

all: test

build: clean copy_configs ./build/server ## Build the project
	@echo 'Done'

./build/.running:
	touch ./build/.running
	make server

run: ./build/.running ## Run the server and scanner

./build/server: cmd/server/main.go cmd/server/appConfig.go
	@go build -o build/server cmd/server/main.go cmd/server/appConfig.go

server: ./build/server ## Run just the webserver
	rm -f /tmp/kb.txt
	@cd ./build && \
	./server  --keyboard-file /tmp/kb.txt --config-path=$(shell pwd)/build/configs/config.json

copy_configs: configs
	@mkdir -p ./build/configs
	@cp -r ./configs/*.json ./build/configs/ || :

.PHONY: test
test: ## Run all tests
	@go test ./...

.PHONY: vtest
vtest: ## Run all tests with verbose flag set
	@go test -v -count=1 ./...

.PHONY: clean
clean: ## Remove generated build files
	@rm -rf ./build
	@go clean -testcache

.PHONY: postSomeText
postSomeText: ## Post some text for testing
	curl --request POST \
		--include \
		--location \
		--header "Content-Type: text/plain" \
		--data 'this is a String' \
		localhost:8282/write/string
	@echo ''

.PHONY: help
help: ## Show help message
	@grep -E '^[[:alnum:]_-]+[[:blank:]]?:.*##' $(MAKEFILE_LIST) \
		| sort \
		| awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
