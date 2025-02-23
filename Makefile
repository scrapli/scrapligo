.DEFAULT_GOAL := help

help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

fmt: ## Run formatters
	gofumpt -w .
	gci write .
	golines -w .

lint: ## Run linters
	golangci-lint run

test: ## Run unit tests
	gotestsum --format testname --hide-summary=skipped -- -coverprofile=cover.out ./...

test-race: ## Run unit tests with race flag
	gotestsum --format testname --hide-summary=skipped -- -coverprofile=cover.out ./... -race

test-e2e: ## Run e2e tests against "full" test topology
	gotestsum --format testname --hide-summary=skipped -- ./e2e/...

test-e2e-ci: ## Run e2e tests against "ci" test topology with race flag
	gotestsum --format testname --hide-summary=skipped -- ./e2e/... -platforms nokia_srl -race

cov:  ## Produce html coverage report
	go tool cover -html=cover.out

deploy-clab-full: ## Deploy "full" test topology
	cd .clab && sudo clab deploy -t topo-full.yaml

destroy-clab-full: ## Destroy "full" test topology
	cd .clab && sudo clab destroy -t topo-full.yaml

deploy-clab-ci: ## Deploy "ci" test topology
	cd .clab && sudo clab deploy -t topo-ci.yaml

destroy-clab-ci: ## Destroy "ci" test topology
	cd .clab && sudo clab destroy -t topo-ci.yaml