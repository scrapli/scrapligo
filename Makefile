.DEFAULT_GOAL := help

NET_EXAMPLES?=examples/network_driver

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
	-@go run ${NET_EXAMPLES}/simple/main.go -argument=${NET_EXAMPLES}/simple/commandsfile

net-log-example: ## Run Network logging example
	-@go run ${NET_EXAMPLES}/logging/main.go

net-interactive-example: ## Run Network Interactive example
	-@go run ${NET_EXAMPLES}/interactive/main.go

net-factory-example: ## Run Network Factory example
	-@go run ${NET_EXAMPLES}/factory/main.go

net-onopen-example: ## Run Network On Open example
	-@go run ${NET_EXAMPLES}/custom_onopen/main.go

net-channellog-example: ## Run Network Channel Log example
	-@go run ${NET_EXAMPLES}/channellog/main.go

help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'