## Makefile commands base
# https://marmelab.com/blog/2016/02/29/auto-documented-makefile.html
# https://www.gnu.org/software/make/manual/html_node/Phony-Targets.html

COLOR_BLUE=\033[0;34m
COLOR_GREEN=\033[0;32m
COLOR_RESET=\033[0m
TARGET_REGEX="^[a-zA-Z0-9_.-]+:.*?\#\#"

.DEFAULT_GOAL := help

.PHONY: help
help: ## Show this help
	@echo "Usage: make ${COLOR_BLUE}[target]${COLOR_RESET}"
	@echo ""
	@echo "Targets:"
	@grep -E ${TARGET_REGEX} $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "  ${COLOR_BLUE}%-36s${COLOR_RESET} %s\n", $$1, $$2}'

# ----------------------------------------------
# Go commands

GO_REPORT_DIR ?= .reports
GO_REPORT_COVERAGE_FILE ?= "$(GO_REPORT_DIR)"/coverage.out
GO_REPORT_TEST_FILE ?= "$(GO_REPORT_DIR)"/tests.json

.PHONY: download
download: ## Download dependencies
	@echo "Downloading dependencies..."
	@go mod download

.PHONY: generate
generate: download ## Generate code
	@echo "Generating code..."
	@go generate ./...

# ---------------------- Tests ------------------------

.PHONY: tests
tests: generate ## Run tests
	@echo "Running tests..."
	@go test ./... -coverpkg=./...

.PHONY: tests.coverage
tests.coverage: generate ## Run tests and generate coverage report
	@echo "Running tests with coverage..."
	@mkdir -p "$$(dirname $(GO_REPORT_COVERAGE_FILE))"
	@mkdir -p "$$(dirname $(GO_REPORT_TEST_FILE))"
	@go test -v ./... -coverpkg=./... -coverprofile=$(GO_REPORT_COVERAGE_FILE) -json ./... 2>&1 | tee $(GO_REPORT_TEST_FILE)
	@echo "Coverage report generated in ${COLOR_START}$(GO_REPORT_COVERAGE_FILE)${COLOR_RESET}"
	@echo "Test report generated in ${COLOR_START}$(GO_REPORT_TEST_FILE)${COLOR_RESET}"

# ---------------------- Dependencies ------------------------

.PHONY: deps.update
deps.update: ## Update package (go.mod and go.sum) and launch tests
	@echo "Tidying up go modules..."
	@go get -u ./...
	@go mod tidy
	@$(MAKE) generate
	@$(MAKE) tests

.PHONY: deps.revert
deps.revert: ## Revert changes in go.mod and go.sum
	@echo "Reverting changes in go modules..."
	@git checkout go.mod go.sum
	@go mod tidy
	@$(MAKE) generate
	@$(MAKE) tests

# ---------------------- Documentation ------------------------

.PHONY: godoc
godoc: ## Run a local godoc server with your package documentation
	@echo "Running godoc server..."
	@go install golang.org/x/pkgsite/cmd/pkgsite@676c19eae995f25cccb1e097a1308caecf93d08a
	@echo "${COLOR_GREEN}Press CTRL+C to stop the server${COLOR_RESET}"
	@pkgsite -http "localhost:8880" -open .
