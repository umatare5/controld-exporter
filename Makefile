.PHONY: help build lint test-unit test-unit-coverage clean image pre-commit-install pre-commit-test pre-commit-uninstall

# Binary name and paths
BINARY_NAME := controld-exporter
BUILD_DIR := ./tmp
BINARY_PATH := $(BUILD_DIR)/$(BINARY_NAME)
COVERAGE_DIR := ./coverage
IMAGE_DIR := $(BUILD_DIR)/image
GOARCH := $(shell go env GOARCH)

# Go build flags
LDFLAGS := -X github.com/umatare5/controld-exporter/internal/cli.version=$(shell cat VERSION)
BUILD_FLAGS := -trimpath -ldflags "$(LDFLAGS)"

# Default target
.DEFAULT_GOAL := help

# Show available targets
help:
	@echo "Available targets:"
	@echo "  build                - Build the binary"
	@echo "  lint                 - Run linters (golangci-lint)"
	@echo "  test-unit            - Run unit tests with colored output"
	@echo "  test-unit-coverage   - Generate HTML coverage report"
	@echo "  clean                - Remove build artifacts and backup files"
	@echo "  image                - Build Docker image"
	@echo "  pre-commit-install   - Install the pre-commit hooks"
	@echo "  pre-commit-test      - Run every hook across the whole tree"
	@echo "  pre-commit-uninstall - Remove the pre-commit hooks"
	@echo ""
	@echo "Requirements:"
	@echo "  - gotestsum: go install gotest.tools/gotestsum@latest"
	@echo "  - golangci-lint: https://golangci-lint.run/docs/welcome/install/local/"
	@echo "  - pre-commit: https://pre-commit.com/#install"
	@echo "  - gitleaks: https://github.com/gitleaks/gitleaks#installing"

build: $(BINARY_PATH)

# Build the binary
$(BINARY_PATH):
	mkdir -p $(BUILD_DIR)
	go build $(BUILD_FLAGS) -o $(BINARY_PATH) ./cmd

# Lint the code
lint:
	golangci-lint run
	go mod tidy

# Run unit tests with gotestsum (shows individual test results with color)
test-unit:
	@command -v gotestsum >/dev/null 2>&1 || { echo "Error: gotestsum is not installed. Run: go install gotest.tools/gotestsum@latest"; exit 1; }
	mkdir -p $(COVERAGE_DIR)
	gotestsum --format testname -- -race -coverprofile=$(COVERAGE_DIR)/report.out ./...

# Generate coverage report (HTML)
test-unit-coverage: test-unit
	go tool cover -html=$(COVERAGE_DIR)/report.out -o $(COVERAGE_DIR)/report.html
	@echo "Coverage report generated: $(COVERAGE_DIR)/report.html"

# Clean build artifacts and backup files
clean:
	rm -rf $(BUILD_DIR) $(COVERAGE_DIR)
	find . -name "*.bak*" -type f -delete 2>/dev/null || true

# Docker targets
# The Dockerfile is written for GoReleaser, which hands docker a context holding
# the binary and the licence notices. Building from the repository root instead
# cannot work: the binary is not there, and .dockerignore excludes it by name if
# it is. This target assembles the same context.
image:
	mkdir -p $(IMAGE_DIR)/linux/$(GOARCH)
	CGO_ENABLED=0 GOOS=linux go build $(BUILD_FLAGS) -o $(IMAGE_DIR)/linux/$(GOARCH)/$(BINARY_NAME) ./cmd
	cp LICENSE NOTICE $(IMAGE_DIR)/
	docker build --platform linux/$(GOARCH) -f Dockerfile -t $(USER)/$(BINARY_NAME) $(IMAGE_DIR)

# Pre-commit targets
# Install the hooks declared in .pre-commit-config.yaml
pre-commit-install:
	@command -v pre-commit >/dev/null 2>&1 || { echo "Error: pre-commit is not installed. See: https://pre-commit.com/#install"; exit 1; }
	@pre-commit install --allow-missing-config

# Run every hook across the whole tree without committing
pre-commit-test:
	@pre-commit run --all-files

# Remove the hooks
pre-commit-uninstall:
	@pre-commit uninstall
