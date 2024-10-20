REPORT_COVERAGE=".reports/coverage.out"
REPORT_TESTS=".reports/tests.json"

.PHONY: help

help: ## Show this help
	@echo "Usage: make [target]"
	@echo ""
	@echo "Targets:"
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z0-9_-]+:.*?## / {printf "  %-15s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

download: ## Download dependencies
	@echo "Downloading dependencies..."
	@go mod download

generate: download ## Generate code
	@echo "Generating code..."
	@go generate ./...

test: generate ## Run tests
	@echo "Running tests..."
	@mkdir -p .reports
	@go test -v ./... -coverpkg=./...

test-coverage: generate ## Run tests and generate coverage report
	@echo "Running tests with coverage..."
	@mkdir -p .reports
	@go test -v ./... -coverpkg=./... -coverprofile=$(REPORT_COVERAGE) -json ./... 2>&1 | tee $(REPORT_TESTS)
