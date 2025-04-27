.DEFAULT_GOAL := help

help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

fmt: ## Run formatters
	gofumpt -w .
	gci write .
	golines -w .

lint: fmt ## Run linters
	golangci-lint run

test: ## Run unit tests
	gotestsum --format testname --hide-summary=skipped -- -coverprofile=cover.out `go list ./... | grep -v e2e`

test-race: ## Run unit tests with race flag
	gotestsum --format testname --hide-summary=skipped -- -coverprofile=cover.out -race `go list ./... | grep -v e2e`

test-e2e: ## Run e2e tests against "full" test topology (count to never cache e2e tests)
	gotestsum --format testname --hide-summary=skipped -- ./e2e/... -count=1

test-e2e-ci: ## Run e2e tests against "ci" test topology with race flag (count to never cache e2e tests)
	gotestsum --format testname --hide-summary=skipped -- ./e2e/... -platforms nokia_srl -race -count=1

cov:  ## Produce html coverage report
	go tool cover -html=cover.out

build-clab-launcher: ## Builds the netopeer and clab launcher images
	docker build \
		-f e2e/clab/netopeer/Dockerfile \
		-t libscrapli-netopeer2:latest \
		e2e/clab/netopeer
	docker build \
		-f e2e/clab/launcher/Dockerfile \
		-t clab-launcher:latest \
		e2e/clab/launcher

run-clab: ## Runs the clab functional testing topo; uses the clab launcher to run nicely on darwin
	docker network rm clab || true
	docker network create \
		--driver bridge \
		--subnet=172.20.20.0/24 \
		--gateway=172.20.20.1 \
		--ipv6 \
		--subnet=2001:172:20:20::/64 \
		--gateway=2001:172:20:20::1 \
		--opt com.docker.network.driver.mtu=65535 \
		--label containerlab \
		clab
	docker run \
		-d \
		--rm \
		--name clab-launcher \
		--privileged \
		--pid=host \
		--stop-signal=SIGINT \
		-v /var/run/docker.sock:/var/run/docker.sock \
		-v /run/netns:/run/netns \
		-v "$$(pwd)/e2e/clab:$$(pwd)/e2e/clab" \
		-e "LAUNCHER_WORKDIR=$$(pwd)/e2e/clab" \
		-e "HOST_ARCH=$$(uname -m)" \
		clab-launcher:latest

