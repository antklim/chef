.PHONY: lint
lint: go-lint ## Run linter

.PHONY: test
test: go-test ## Run tests

.PHONY: cov-report
cov-report: go-cov-report ## View coverage report (HTML version)

.PHONY: cov-report-ci
cov-report-ci: go-cov-report-ci ## View coverage report (text version for CI)

.PHONY: help
help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

go-lint:
	golangci-lint run -v

go-test:
	go test -race -cover -coverprofile=coverage.out -count=1 ./...

go-cov-report:
	go tool cover -html=coverage.out

go-cov-report-ci:
	go tool cover -func=coverage.out

.DEFAULT_GOAL := help