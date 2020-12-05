.PHONY: all build lint test vet
CHECK_FILES?=$$(go list ./... | grep -v /vendor/)
APP_NAME=hellper

GO ?= go
GORUN ?= $(GO) run
GOIMPORTS ?= $(GORUN) golang.org/x/tools/cmd/goimports
GIT ?= git
GITDIFF ?= $(GIT) diff

help: ## Show this help.
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {sub("\\\\n",sprintf("\n%22c"," "), $$2);printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST) | sort

all: lint vet test build ## Run the tests and build the binary.

build: ## Build the binary.
	[ -d bin ] || mkdir bin
	go build -o ./bin/$(APP_NAME) ./cmd/http

lint: ## Lint the code.
	golint $(CHECK_FILES)

test: ## Run tests.
	go test -race -cover -p 1 -v $(CHECK_FILES)

vet: ## Vet the code
	go vet $(CHECK_FILES)

install: ## Install application on local machine or container
	go install *.go

run: ## Run application
	docker-compose up

run-local: ## Run application locally (without docker)
	go run ./cmd/http -v

run-db: ## Start the database on docker
	docker-compose up -d hellper_db

migrate: ## Migrate the database
	docker-compose exec hellper sh -c "go run ./cmd/migrations -v"

migrate-local: ## Migrate the database using local binaries
	go run ./cmd/migrations -v

clean: ## Remove all resources
	docker-compose rm -sf

goimports:
	@$(GOIMPORTS) -w $(SOURCES)

git/diff:
	@if ! $(GITDIFF) --quiet; then \
		printf 'Found changes on local workspace. Please run this target and commit the changes\n' ; \
		exit 1; \
	fi