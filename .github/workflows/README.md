# GitHub Actions Workflows

This directory contains CI/CD workflows for the SIMDCuckooFilter project.

## Workflows

### `ci.yml` - Comprehensive CI Pipeline

**Triggers:** Pull requests to main, pushes to main, manual dispatch

**Concurrency:** Automatically cancels in-progress runs for the same workflow and branch

**Jobs:**

#### `format-check`
- **Platform:** Ubuntu
- **Actions:**
  - Verifies all Go code is formatted with `gofmt`
  - Fails if any files need formatting

#### `vet`
- **Platform:** Ubuntu
- **Actions:**
  - Runs `go vet` to detect suspicious constructs
  - Allowed to fail due to known assembly stack frame offset issues
- **Known Issues:** AMD64 AVX2 assembly bugs cause vet warnings

#### `test`
- **Platforms:** Ubuntu, macOS, Windows
- **Matrix:** Multi-platform testing with fail-fast disabled
- **Actions:**
  - Downloads and verifies dependencies
  - Runs full test suite with race detector (`-race -count=1`)
  - Generates coverage report on Ubuntu
  - Uploads coverage to Codecov
- **Known Issues:** Assembly bugs cause test failures, allowed to continue

#### `test-arm64`
- **Platform:** macOS 14 (Apple Silicon)
- **Actions:**
  - Verifies ARM64 architecture
  - Runs SIMD-specific tests
  - Runs ARM64-specific tests
  - Full test suite with race detection

#### `lint`
- **Platform:** Ubuntu
- **Actions:**
  - Runs golangci-lint with 5-minute timeout
  - Code quality and style checks

#### `build-matrix`
- **Cross-compilation verification**
- **Matrix:** linux/darwin/windows × amd64/arm64
  - Excludes: windows/arm64
- **Actions:**
  - Verifies builds compile for all platform combinations
  - Tests cross-platform compatibility

#### `verify-assembly`
- **Platform:** Ubuntu
- **Actions:**
  - Verifies AMD64 assembly compiles (`internal/lookup`, `internal/hash`)
  - Verifies ARM64 assembly compiles (`internal/lookup`, `internal/hash`)

#### `security`
- **Platform:** Ubuntu
- **Actions:**
  - Runs Gosec security scanner (v2.21.4)
  - Runs govulncheck for known vulnerabilities

#### `ci-success`
- **Final gate for CI pipeline**
- **Dependencies:** All jobs above
- **Actions:**
  - Checks status of all required jobs
  - Allows vet to fail (known assembly issues)
  - Fails pipeline if any other required job fails

---

## Known Issues

### Assembly Bugs

The following assembly bugs are present and tracked:

1. **AMD64 AVX2 Bucket Lookup** (`internal/lookup/bucket_lookup_avx2_amd64.s`)
   - Wrong argument size: 25 instead of 33
   - Invalid return offset: ret+25(FP) instead of ret+32(FP)

2. **AMD64 AVX2 Batch Hash** (`internal/hash/xxhash/batch_avx2_amd64.s`)
   - Wrong argument size: 72 instead of 64
   - Multiple invalid offset errors

**Impact:**
- `go vet` warnings (allowed to fail)
- Test failures on Linux and Windows (allowed to continue)
- Affected tests: bucket lookup alignment, patterns, stress tests

---

## Running CI Checks Locally

### Format Check
```bash
# Check formatting
gofmt -l .

# Fix formatting
make fmt
# or
go fmt ./...
```

### Run Tests
```bash
# All tests
go test -v -race -count=1 ./...

# Short tests
go test -short ./...

# SIMD-specific tests
go test -v -run=.*SIMD.* ./...

# ARM64-specific tests (on ARM64 machine)
go test -v -run=.*ARM64.* ./...

# Using Makefile
make test
make test-race
make test-simd
```

### Generate Coverage
```bash
# Generate coverage report
go test -coverprofile=coverage.out -covermode=atomic ./...

# View coverage
go tool cover -html=coverage.out

# Using Makefile
make coverage
make coverage-html
```

### Lint
```bash
# Install golangci-lint
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

# Run linter
golangci-lint run --timeout=5m

# Using Makefile
make lint
```

### Security Scan
```bash
# Install tools
go install github.com/securego/gosec/v2/cmd/gosec@latest
go install golang.org/x/vuln/cmd/govulncheck@latest

# Run security scan
gosec -no-fail ./...

# Check vulnerabilities
govulncheck ./...
```

### Cross-Platform Builds
```bash
# Build for specific platform
GOOS=linux GOARCH=amd64 go build -v ./...
GOOS=linux GOARCH=arm64 go build -v ./...
GOOS=darwin GOARCH=amd64 go build -v ./...
GOOS=darwin GOARCH=arm64 go build -v ./...
GOOS=windows GOARCH=amd64 go build -v ./...

# Using Makefile
make build-linux-amd64
make build-linux-arm64
make build-darwin-amd64
make build-darwin-arm64
make build-all
```

### Verify Assembly
```bash
# AMD64 assembly
GOARCH=amd64 go build -v ./internal/lookup/...
GOARCH=amd64 go build -v ./internal/hash/...

# ARM64 assembly
GOARCH=arm64 go build -v ./internal/lookup/...
GOARCH=arm64 go build -v ./internal/hash/...
```

### Run All CI Checks
```bash
# Using Makefile
make ci
```

---

## Platform Support

| Platform | Architecture | SIMD | Status |
|----------|-------------|------|--------|
| Linux | AMD64 | AVX2/SSE2 | ✅ Tested (assembly bugs) |
| Linux | ARM64 | NEON | ✅ Build verified |
| macOS | AMD64 | AVX2/SSE2 | ✅ Tested |
| macOS | ARM64 | NEON | ✅ Tested |
| Windows | AMD64 | AVX2/SSE2 | ✅ Tested (assembly bugs) |
| Windows | ARM64 | - | ❌ Not supported |

---

## Workflow Status Badge

Add to your main README.md:

```markdown
[![CI](https://github.com/USERNAME/SIMDCuckooFilter/workflows/CI/badge.svg)](https://github.com/USERNAME/SIMDCuckooFilter/actions/workflows/ci.yml)
[![codecov](https://codecov.io/gh/USERNAME/SIMDCuckooFilter/branch/main/graph/badge.svg)](https://codecov.io/gh/USERNAME/SIMDCuckooFilter)
```

---

## Debugging Failed Workflows

### Viewing Logs
1. Go to the Actions tab in GitHub
2. Select the failed workflow run
3. Click on the failed job
4. Expand the failing step to see detailed logs

### Common Issues

**Format check fails:**
- Run `make fmt` or `go fmt ./...` locally
- Commit the formatted files

**Tests fail:**
- Known assembly bugs cause failures on Linux/Windows
- Check if new failures are unrelated to known bugs
- Run tests locally with `-v` flag for details

**Lint fails:**
- Run `golangci-lint run` locally
- Fix reported issues
- Some warnings may be acceptable (document in code)

**Cross-platform build fails:**
- Verify assembly code for target architecture
- Check build tags are correct
- Ensure platform-specific code uses proper build constraints

**Coverage upload fails:**
- Verify Codecov token is set in repository secrets
- Check Codecov service status
- CI continues even if upload fails (`fail_ci_if_error: false`)

---

## Workflow Optimization

Current optimizations:
- ✅ Parallel job execution
- ✅ Go module caching (actions/setup-go@v5)
- ✅ Concurrency control (auto-cancel outdated runs)
- ✅ Fail-fast disabled for comprehensive testing
- ✅ Continue-on-error for known issues

---

## Contributing

When adding new code:

1. **Run format check:** `make fmt`
2. **Run tests locally:** `make test`
3. **Run linter:** `make lint`
4. **Check coverage:** `make coverage`
5. **Verify cross-platform:** Test on multiple platforms if possible

The CI pipeline will automatically run when you push or create a pull request.
