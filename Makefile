# setup variables
GOBIN ?= $(shell go env GOPATH)/bin

.PHONY: lint
lint: go-lint ## Run linter

.PHONY: test
test: go-test ## Run tests

.PHONY: build
build: go-build ## Build Chef

.PHONY: cov-report
cov-report: go-cov-report ## View coverage report (HTML version)

.PHONY: cov-report-ci
cov-report-ci: go-cov-report-ci ## View coverage report (text version for CI)

.PHONY: integration
integration: cmd-integration ## Chef command integration tests

.PHONY: help
help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

go-lint:
	golangci-lint run -v

go-test:
	go test -race -cover -coverprofile=coverage.out -count=1 ./...

go-build:
	go build -o $(GOBIN)/chef cmd/chef/main.go

go-cov-report:
	go tool cover -html=coverage.out

go-cov-report-ci:
	go tool cover -func=coverage.out

$(GOBIN)/commander:
	cd && GO111MODULE=auto go get github.com/commander-cli/commander/cmd/commander

run-integration:
	rm -rf tmp && mkdir -p tmp/subdir
	commander test commander.yml
.PHONY: run-integration

integration-cleanup:
	./scripts/integration/test-cleanup.sh
.PHONY: integration-cleanup

cmd-integration: build $(GOBIN)/commander
	@$(MAKE) run-integration; \
	ret=$$?; \
	$(MAKE) integration-cleanup; \
	exit $$ret

.DEFAULT_GOAL := help
