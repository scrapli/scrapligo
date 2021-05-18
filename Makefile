lint:
	gofmt -w -s .
	goimports -w .
	golines -w .
	golangci-lint run

test:
	go test -cover -v ./channel/... ./driver/... ./netconf/... ./logging/... ./transport/...
