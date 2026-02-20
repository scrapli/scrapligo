.DEFAULT_GOAL := help

## Show this help
help:
	@awk -f build/makefile-doc.awk $(MAKEFILE_LIST)

##@ Development
## Format all go files
fmt:
	go run mvdan.cc/gofumpt -w .
	go run github.com/daixiang0/gci write .
	go run github.com/golangci/golines --base-formatter="gofmt" -w .

## Lint all go files
lint: fmt
	golangci-lint run

##@ Testing
## Run unit tests
test:
	go test -v -coverprofile=cover.out `go list ./... | grep -v e2e`

## Run unit tests with race flag
test-race:
	go test -v -coverprofile=cover.out -race `go list ./... | grep -v e2e`

## Run e2e tests against "full" test topology (count to never cache e2e tests)
test-e2e:
	go test -v ./e2e/... -count=1

## Run e2e tests against "ci" test topology with race flag (count to never cache e2e tests)
test-e2e-ci:
	go test -v ./e2e/... -platforms nokia_srl -race -count=1 -skip-slow

##@ Testing Coverage
## Produce html coverage report
cov:
	go tool cover -html=cover.out

##@ Test Environment
## Runs the clab functional testing topo; uses the clab launcher to run nicely on darwin
run-clab:
	rm -r .clab/* || true
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
		--platform=linux/arm64 \
		--privileged \
		--pid=host \
		--stop-signal=SIGINT \
		-v /var/run/docker.sock:/var/run/docker.sock \
		-v /run/netns:/run/netns \
		-v "$$(pwd):$$(pwd)" \
		-e "WORKDIR=$$(pwd)/.clab" \
		-e "HOST_ARCH=$$(uname -m)" \
		ghcr.io/scrapli/scrapli_clab/launcher:0.0.7

## Runs the clab functional testing topo with the ci specific topology - omits ceos
run-clab-ci:
	mkdir .clab || true
	rm -r .clab/* || true
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
        -v "$$(pwd):$$(pwd)" \
        -e "WORKDIR=$$(pwd)/.clab" \
        -e "HOST_ARCH=$$(uname -m)" \
        -e "CLAB_TOPO=topo.ci.$$(uname -m).yaml" \
        ghcr.io/scrapli/scrapli_clab/launcher:0.0.7
