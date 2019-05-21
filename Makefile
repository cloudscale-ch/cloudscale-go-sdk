TEST?=$$(go list ./... |grep -v 'vendor')

test:
	go test $(TEST) $(TESTARGS) -timeout 30s

integration:
	go get github.com/cenkalti/backoff
	go get golang.org/x/oauth2
	GOPATH=$(GOPATH):$(PWD) go test -tags=integration -v integration $(TESTARGS) -timeout 120m

fmt:
	go fmt
	gofmt -l -w test/integration

.PHONY: test integration

