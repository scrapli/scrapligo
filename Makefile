.DEFAULT_GOAL := help

help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

lint: ## Run linters
	gofumpt -w .
	goimports -w .
	golines -w .
	golangci-lint run

test: ## Run unit tests
	gotestsum --format testname -- -coverprofile=cover.out ./...

test-race: ## Run unit tests with race flag
	gotestsum --format testname -- -coverprofile=cover.out ./... -race

test-functional: ## Run functional tests against "full" test topology
	gotestsum --format testname -- ./... -functional

test-ci: ## Run functional tests against "ci" test topology with race flag
	gotestsum --format testname -- ./... -functional -platforms nokia_srlinux -race

cov:  ## Produce html coverage report
	go tool cover -html=cover.out

deploy-clab-full: ## Deploy "full" test topology
	cd clab && sudo clab deploy -t topo-full.yaml

destroy-clab-full: ## Destroy "full" test topology
	cd clab && sudo clab destroy -t topo-full.yaml

deploy-clab-ci: ## Deploy "ci" test topology
	cd clab && sudo clab deploy -t topo-ci.yaml

destroy-clab-ci: ## Destroy "ci" test topology
	cd clab && sudo clab destroy -t topo-ci.yaml