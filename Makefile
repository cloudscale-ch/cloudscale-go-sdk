TEST?=$$(go list ./... |grep -v 'vendor')
VERSION ?= $(shell cat VERSION)

test:
	go test -v $(TEST) $(TESTARGS) -timeout 30s

clean-testcache:
	@go clean -testcache  # Force retesting of code

integration: clean-testcache
	go test -tags=integration -v $(TEST)/test/integration/... $(TESTARGS) -parallel 4 -timeout 120m

integration-short: clean-testcache
	go test -tags=integration -short -v $(TEST)/test/integration/... $(TESTARGS) -parallel 4 -timeout 120m

vet:
	go vet ./...
	go vet -tags=integration ./test/integration/

fmt:
	go fmt
	gofmt -l -w test/integration

bump-version:
	@[ "${NEW_VERSION}" ] || ( echo "NEW_VERSION must be set (ex. make NEW_VERSION=v1.x.x bump-version)"; exit 1 )
	@(echo ${NEW_VERSION} | grep -E "^v") || ( echo "NEW_VERSION must be a semver ('v' prefix is required)"; exit 1 )
	@echo "Bumping VERSION from $(VERSION) to $(NEW_VERSION)"
	@echo $(NEW_VERSION) > VERSION
	@sed -i.bak -e 's/${VERSION}/${NEW_VERSION}/g' cloudscale.go
	@rm cloudscale.go.bak

.PHONY: test vet integration integration-short clean-testcache
