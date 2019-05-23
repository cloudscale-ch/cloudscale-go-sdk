TEST?=$$(go list ./... |grep -v 'vendor')

test:
	go test $(TEST) $(TESTARGS) -timeout 30s

integration:
	go test -tags=integration -v $(TEST)/test/integration/... $(TESTARGS) -timeout 120m

fmt:
	go fmt
	gofmt -l -w test/integration

.PHONY: test integration
