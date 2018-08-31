TEST?=$$(go list ./... |grep -v 'vendor')

test:
	go test $(TEST) -timeout 30s

integration:
	go get github.com/cenkalti/backoff
	go get golang.org/x/oauth2
	go test -tags=integration -v $(TEST)/test/integration -timeout 120m

.PHONY: test integration
