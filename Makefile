# CuckooFilter Makefile
# Build automation for the CuckooFilter project

.PHONY: help all build build-all build-linux-amd64 build-linux-arm64 build-darwin-amd64 build-darwin-arm64
.PHONY: test test-verbose test-short test-race test-simd test-amd64 test-arm64
.PHONY: bench bench-all bench-insert bench-lookup bench-batch bench-hash
.PHONY: coverage coverage-html
.PHONY: vet fmt check-fmt lint
.PHONY: clean clean-all ci pre-commit

# Default target - run standard checks
all: fmt vet test
	@echo "All checks passed"

# Help target
help:
	@echo "CuckooFilter Build Automation"
	@echo ""
	@echo "Available targets:"
	@echo "  all            - Run fmt, vet, and test (default)"
	@echo "  build          - Build the project"
	@echo "  test           - Run all tests"
	@echo "  test-verbose   - Run tests with verbose output"
	@echo "  test-short     - Run short tests only"
	@echo "  test-race      - Run tests with race detector"
	@echo "  test-simd      - Run SIMD-specific tests"
	@echo "  bench          - Run all benchmarks"
	@echo "  bench-all      - Run all benchmarks with multiple runs"
	@echo "  bench-insert   - Run insert-related benchmarks"
	@echo "  bench-lookup   - Run lookup-related benchmarks"
	@echo "  bench-batch    - Run batch operation benchmarks"
	@echo "  bench-hash     - Run hash strategy benchmarks"
	@echo "  coverage       - Generate test coverage report"
	@echo "  coverage-html  - Generate HTML coverage report"
	@echo "  lint           - Run golangci-lint (requires golangci-lint)"
	@echo "  vet            - Run go vet"
	@echo "  fmt            - Format code with go fmt"
	@echo "  check-fmt      - Check if code is formatted"
	@echo "  ci             - Run all CI checks (check-fmt, vet, test, coverage)"
	@echo "  pre-commit     - Quick pre-commit checks (fmt, vet, test-short)"
	@echo "  clean          - Clean build artifacts (coverage files)"
	@echo "  clean-all      - Deep clean including Go caches (affects all projects)"
	@echo ""
	@echo "Cross-compilation targets:"
	@echo "  build-linux-amd64   - Build for Linux AMD64"
	@echo "  build-linux-arm64   - Build for Linux ARM64"
	@echo "  build-darwin-amd64  - Build for macOS AMD64"
	@echo "  build-darwin-arm64  - Build for macOS ARM64 (Apple Silicon)"
	@echo "  build-all           - Build for all platforms"
	@echo ""
	@echo "Platform-specific tests:"
	@echo "  test-amd64     - Run tests with GOARCH=amd64"
	@echo "  test-arm64     - Run tests with GOARCH=arm64"

# Build targets
build:
	@echo "Building CuckooFilter..."
	go build ./...

build-linux-amd64:
	@echo "Building for Linux AMD64..."
	GOOS=linux GOARCH=amd64 go build ./...

build-linux-arm64:
	@echo "Building for Linux ARM64..."
	GOOS=linux GOARCH=arm64 go build ./...

build-darwin-amd64:
	@echo "Building for macOS AMD64..."
	GOOS=darwin GOARCH=amd64 go build ./...

build-darwin-arm64:
	@echo "Building for macOS ARM64 (Apple Silicon)..."
	GOOS=darwin GOARCH=arm64 go build ./...

build-all: build-linux-amd64 build-linux-arm64 build-darwin-amd64 build-darwin-arm64
	@echo "Built for all platforms"

# Test targets
test:
	@echo "Running tests..."
	go test ./... -count=1

test-verbose:
	@echo "Running tests with verbose output..."
	go test ./... -v -count=1

test-short:
	@echo "Running short tests..."
	go test ./... -short -count=1

test-race:
	@echo "Running tests with race detector..."
	go test -race ./... -count=1

test-simd:
	@echo "Running SIMD-specific tests..."
	go test -v -run=.*SIMD.* -count=1 ./...

test-amd64:
	@echo "Running tests for AMD64..."
	GOARCH=amd64 go test ./... -count=1

test-arm64:
	@echo "Running tests for ARM64..."
	GOARCH=arm64 go test ./... -count=1

# Benchmark targets
bench:
	@echo "Running benchmarks..."
	go test -bench=. -benchmem -count=1 ./...

bench-all:
	@echo "Running benchmarks (5 runs)..."
	go test -bench=. -benchmem -count=5 -benchtime=2s ./...

bench-insert:
	@echo "Running insert benchmarks..."
	go test -bench=BenchmarkInsert.* -benchmem -count=1 ./...

bench-lookup:
	@echo "Running lookup benchmarks..."
	go test -bench=BenchmarkLookup.* -benchmem -count=1 ./...

bench-batch:
	@echo "Running batch benchmarks..."
	go test -bench=Benchmark.*Batch.* -benchmem -count=1 ./...

bench-hash:
	@echo "Running hash strategy benchmarks..."
	go test -bench=BenchmarkHashStrategies.* -benchmem -count=1 ./...

# Coverage targets
coverage:
	@echo "Generating coverage report..."
	go test ./... -coverprofile=coverage.out -count=1
	go tool cover -func=coverage.out

coverage-html:
	@echo "Generating HTML coverage report..."
	go test ./... -coverprofile=coverage.out -count=1
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

# Code quality targets
vet:
	@echo "Running go vet..."
	go vet ./...

fmt:
	@echo "Formatting code..."
	go fmt ./...

check-fmt:
	@echo "Checking code formatting..."
	@test -z "$$(gofmt -l .)" || (echo "Code is not formatted. Run 'make fmt'" && exit 1)

lint:
	@echo "Running golangci-lint..."
	@which golangci-lint > /dev/null || (echo "golangci-lint not installed. Install from https://golangci-lint.run/usage/install/" && exit 1)
	golangci-lint run

# Clean target
clean:
	@echo "Cleaning build artifacts..."
	rm -f coverage.out coverage.html
	go clean

# Deep clean - clears global Go caches (affects all projects)
clean-all:
	@echo "Deep cleaning (clears global Go caches)..."
	rm -f coverage.out coverage.html
	go clean -cache -testcache
	@echo "Warning: Global Go build and test caches have been cleared"

# CI target - runs all checks
ci: check-fmt vet test coverage
	@echo "All CI checks passed"

# Quick check before commit
pre-commit: fmt vet test-short
	@echo "Pre-commit checks passed"
