HOST?=0.0.0.0
PORT?=4000
FILEPATH?=./auth.json
REDIS_ADDR?=127.0.0.1:6379

build: ## Runs gb build on the project
	gb build

dev: ## Builds and runs the service with local environment
	gb build && HOST="${HOST}" PORT="${PORT}" FILEPATH="${FILEPATH}" REDIS_ADDR="${REDIS_ADDR}" ./bin/reception

lint: ## Runs all packages in the service through golint
	GOPATH=$(PWD):$(PWD)/vendor golint reception/...

run: ## Runs the service with local environment unless overridden
	HOST="${HOST}" PORT="${PORT}" REDIS_ADDR="${REDIS_ADDR}" FILEPATH="${FILEPATH}" ./bin/reception

test: ## Runs gb test with the -v verbose flag
	gb test -v

vet: ## Runs all packages in the service through go-vet
	GOPATH=$(PWD):$(PWD)/vendor go vet reception/...

.PHONY: help

help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

.DEFAULT_GOAL := help
