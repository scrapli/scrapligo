.DEFAULT_GOAL := help

lint: ## Run linters
	gofmt -w -s .
	goimports -w .
	golines -w .
	golangci-lint run

test: ## Test execution
	go test -cover -v ./channel/... ./driver/... ./netconf/... ./logging/... ./transport/...

test_race: ## Test execution with race flag
	go test -cover -race -v ./channel/... ./driver/... ./netconf/... ./logging/... ./transport/...

help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'