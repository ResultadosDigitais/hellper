.PHONY: all build lint test vet
CHECK_FILES?=$$(go list ./... | grep -v /vendor/)
APP_NAME=notify-api

help: ## Show this help.
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {sub("\\\\n",sprintf("\n%22c"," "), $$2);printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)

all: lint vet test build ## Run the tests and build the binary.

build: ## Build the binary.
	go build -o bin/$(APP_NAME) *.go

lint: ## Lint the code.
	golint $(CHECK_FILES)

test: ## Run tests.
	go test -race -cover -p 1 -v $(CHECK_FILES)

vet: ## Vet the code
	go vet $(CHECK_FILES)

run: ## Run application
	docker-compose up

install: ## Install application on local machine or container
	go install *.go