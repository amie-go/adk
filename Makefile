REPORT_COVERAGE=".reports/coverage.out"
REPORT_TESTS=".reports/tests.json"

# https://www.gnu.org/software/make/manual/html_node/Phony-Targets.html
.PHONY: help download generate test test-coverage deps-update deps-revert godoc

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
	@go test ./... -coverpkg=./...

test-coverage: generate ## Run tests and generate coverage report
	@echo "Running tests with coverage..."
	@mkdir -p .reports
	@go test -v ./... -coverpkg=./... -coverprofile=$(REPORT_COVERAGE) -json ./... 2>&1 | tee $(REPORT_TESTS)

deps-update: ## Update package (go.mod and go.sum) and launch tests
	@echo "Tidying up go modules..."
	@go get -u ./...
	@go mod tidy
	@go generate ./...
	@go test ./... -coverpkg=./...

deps-revert: ## Revert changes in go.mod and go.sum
	@echo "Reverting changes in go modules..."
	@git checkout go.mod go.sum
	@go mod tidy
	@go generate ./...
	@go test ./... -coverpkg=./...

godoc: ## Run a local godoc server with your package documentation
	@echo "Running godoc server..."
	@go install golang.org/x/pkgsite/cmd/pkgsite@676c19eae995f25cccb1e097a1308caecf93d08a
	@echo "\033[32mPress CTRL+C to stop the server\033[m"
	@pkgsite -http "localhost:8880" -open .
