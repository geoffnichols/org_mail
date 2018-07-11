

.DEFAULT_GOAL := help

help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(firstword $(MAKEFILE_LIST)) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

fmt: ## Format all code
	go fmt *.go

build: ## Build the go program
	go build *.go

clean: ## Remove artiacts
	rm -f org_mail

run: ## Run program
	go run *.go

.PHONY: help fmt clean
