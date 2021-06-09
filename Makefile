.DEFAULT_GOAL := help

lint: ## Run linters
	gofmt -w -s .
	goimports -w .
	golines -w .
	golangci-lint run

test: ## Test execution
	go test -cover -v ./channel/... ./driver/... ./netconf/... ./logging/... ./transport/...

base-example: ## Run Base example
	-@go run examples/base_driver/main.go

net-simple-example: ## Run Network simple example
	-@go run examples/network_driver/simple/main.go -argument=examples/network_driver/simple/commandsfile

net-log-example: ## Run Network logging example
	-@go run examples/network_driver/logging/main.go

net-interactive-example: ## Run Network Interactive example
	-@go run examples/network_driver/interactive/main.go

net-factory-example: ## Run Network Factory example
	-@go run examples/network_driver/factory/main.go

net-onopen-example: ## Run Network On Open example
	-@go run examples/network_driver/custom_onopen/main.go

net-channellog-example: ## Run Network Channel Log example
	-@go run examples/network_driver/channellog/main.go

help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'