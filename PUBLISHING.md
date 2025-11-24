# Publishing Guide

This guide provides step-by-step instructions for publishing a new release of SIMDCuckooFilter.

## Pre-Publishing Checklist

Before you start the publishing process, ensure:

### 1. Code Quality

- [ ] All tests pass locally and on CI
  ```bash
  make test
  make test-race
  ```

- [ ] No test failures or flaky tests
  ```bash
  make test-verbose
  ```

- [ ] Code is properly formatted
  ```bash
  make fmt
  make check-fmt
  ```

- [ ] Linter passes (if available)
  ```bash
  make lint  # or golangci-lint run
  ```

- [ ] Code has been reviewed and approved
  ```bash
  make vet
  ```

### 2. Documentation

- [ ] `CHANGELOG.md` is up to date with all changes
- [ ] `README.md` reflects current features and usage
- [ ] `version.go` has the correct version numbers
- [ ] API documentation (godoc comments) is complete
- [ ] Examples are working and up to date
  ```bash
  cd examples/basic_usage && go run main.go
  cd examples/custom_config && go run main.go
  ```

### 3. Performance and Benchmarks

- [ ] Benchmarks have been run and results documented
  ```bash
  make bench
  ```

- [ ] No performance regressions identified
- [ ] SIMD optimizations are working correctly
  ```bash
  make test-simd
  ```

### 4. Breaking Changes

- [ ] Breaking changes (if any) are documented in `CHANGELOG.md`
- [ ] Migration guide is provided for major version changes
- [ ] Deprecation warnings added for APIs that will be removed

## Publishing Steps

### Step 1: Prepare the Release Branch

```bash
# Ensure you're on main and up to date
git checkout main
git pull origin main

# Create a release branch
git checkout -b release/v0.2.0
```

### Step 2: Update Version Information

Edit `version.go`:

```go
const (
    Major = 0
    Minor = 2
    Patch = 0
    Version = "v0.2.0"
    PreRelease = ""           // Empty for stable release
    BuildMetadata = ""        // Empty for official release
)
```

### Step 3: Update CHANGELOG.md

Add a new section at the top of `CHANGELOG.md`:

```markdown
## [0.2.0] - 2024-01-15

### Added
- New feature X that provides Y functionality
- Support for custom hash functions with Z interface
- Batch processing optimization for large datasets

### Changed
- Improved memory efficiency by 15%
- Updated hash distribution algorithm for better performance

### Fixed
- Race condition in concurrent bucket access
- Memory leak in filter reset operation

### Performance
- 20% faster insertion for filters with >100k elements
- Reduced memory footprint by 12% for small filters
```

### Step 4: Update README.md (if needed)

Update any sections that reference:
- Installation instructions
- Feature lists
- Performance benchmarks
- Usage examples
- Version compatibility

### Step 5: Run Full Test Suite

```bash
# Run all checks
make ci

# Run tests on different architectures (if applicable)
make test-amd64
make test-arm64

# Run race detector
make test-race

# Generate and review coverage
make coverage-html
# Open coverage.html in browser to review
```

**IMPORTANT**: Do NOT proceed if any tests fail!

### Step 6: Commit and Push Release Branch

```bash
# Commit all changes
git add version.go CHANGELOG.md README.md
git commit -m "chore: prepare release v0.2.0"

# Push to GitHub
git push origin release/v0.2.0
```

### Step 7: Create Pull Request

1. Go to GitHub repository
2. Create a new Pull Request from `release/v0.2.0` to `main`
3. Title: "Release v0.2.0"
4. Description should include:
   - Summary of changes
   - Link to relevant issues/PRs
   - Testing performed
   - Breaking changes (if any)

Example PR description:
```markdown
## Release v0.2.0

### Summary
This release includes significant performance improvements and new features
for batch operations.

### Changes
- Added batch insert/lookup operations (#15)
- Improved SIMD utilization for ARM64 (#18)
- Fixed race condition in concurrent access (#20)

### Performance Improvements
- 20% faster insertions for large filters
- 15% reduced memory usage

### Breaking Changes
None

### Testing
- [x] All tests pass
- [x] Race detector passes
- [x] Benchmarks show expected improvements
- [x] Tested on AMD64 and ARM64

Closes #15, #18, #20
```

### Step 8: Get PR Approved and Merge

1. Wait for CI checks to pass
2. Request review from maintainers
3. Address any feedback
4. Once approved, merge the PR (use "Squash and merge" or "Merge commit")

### Step 9: Create and Push Git Tag

```bash
# Switch to main and pull the merged changes
git checkout main
git pull origin main

# Create an annotated tag
git tag -a v0.2.0 -m "Release v0.2.0

## Highlights
- Performance improvements for batch operations
- Memory efficiency enhancements
- ARM64 SIMD optimizations

## Changes
- Added batch insert/lookup operations
- Improved SIMD utilization for ARM64
- Fixed race condition in concurrent access

## Performance
- 20% faster insertions for large filters
- 15% reduced memory usage

See CHANGELOG.md for complete details.
"

# Push the tag
git push origin v0.2.0
```

### Step 10: Create GitHub Release

1. Go to: `https://github.com/shaia/simdcuckoofilter/releases/new`

2. Select the tag: `v0.2.0`

3. Set release title: `v0.2.0 - Performance Improvements and Batch Operations`

4. Add release notes (can auto-generate and then edit):

```markdown
## What's Changed

### Features
* Batch insert and lookup operations for improved throughput by @username in #15
* Enhanced ARM64 SIMD support for 2x faster operations by @username in #18

### Bug Fixes
* Fixed race condition in concurrent bucket access by @username in #20
* Corrected fingerprint calculation for edge cases by @username in #21

### Performance Improvements
* 20% faster insertions for filters with >100k elements
* 15% reduced memory footprint for small filters
* Improved hash distribution reducing collision rate by 8%

### Documentation
* Updated README with batch operation examples
* Added performance comparison charts
* Improved API documentation

## Benchmarks

| Operation | v0.1.0 | v0.2.0 | Improvement |
|-----------|--------|--------|-------------|
| Insert    | 45 ns  | 36 ns  | 20% faster  |
| Lookup    | 38 ns  | 32 ns  | 16% faster  |
| Batch(16) | 520 ns | 380 ns | 27% faster  |

**Full Changelog**: https://github.com/shaia/simdcuckoofilter/compare/v0.1.0...v0.2.0
```

5. Check "Set as latest release" (or "Set as pre-release" if applicable)

6. Click "Publish release"

### Step 11: Verify Publication

Wait a few minutes for the Go module proxy to update, then verify:

```bash
# Check if version is available on Go proxy
curl https://proxy.golang.org/github.com/shaia/simdcuckoofilter/@v/v0.2.0.info

# Should return something like:
# {"Version":"v0.2.0","Time":"2024-01-15T10:30:00Z"}

# List all available versions
go list -m -versions github.com/shaia/simdcuckoofilter

# Test installation
cd /tmp
mkdir test-install && cd test-install
go mod init test
go get github.com/shaia/simdcuckoofilter@v0.2.0
```

### Step 12: Update Dependent Projects (if any)

If you maintain projects that depend on SIMDCuckooFilter:

```bash
# In dependent project
go get github.com/shaia/simdcuckoofilter@v0.2.0
go mod tidy
```

## Post-Publishing Tasks

### 1. Announce the Release

Consider announcing on:

- [ ] GitHub Discussions (if enabled)
- [ ] Project README or website
- [ ] Twitter/social media
- [ ] Reddit (r/golang)
- [ ] Hacker News (for major releases)
- [ ] Go Forum
- [ ] Company/team Slack or communication channel

Example announcement:

```markdown
🎉 SIMDCuckooFilter v0.2.0 is now available!

New in this release:
✨ Batch operations for 27% better throughput
🚀 20% faster insertions
💾 15% reduced memory usage
🔧 ARM64 SIMD optimizations

Get it: go get github.com/shaia/simdcuckoofilter@v0.2.0

Release notes: https://github.com/shaia/simdcuckoofilter/releases/tag/v0.2.0
```

### 2. Update Repository Metadata

On GitHub repository settings:

- [ ] Update repository description if features changed
- [ ] Update repository topics/tags (e.g., "simd", "cuckoo-filter", "go", "high-performance")
- [ ] Update website URL if applicable
- [ ] Ensure LICENSE is correct
- [ ] Check repository visibility settings

### 3. Monitor Issues and Feedback

After release:

- [ ] Monitor GitHub issues for bug reports
- [ ] Watch for build failures in projects using your library
- [ ] Respond to questions and feedback promptly
- [ ] Track performance reports from users

### 4. Prepare for Next Release

- [ ] Create a milestone for next version
- [ ] Label issues/PRs with target version
- [ ] Start planning next features
- [ ] Update project roadmap (if you have one)

## Version-Specific Guidelines

### Patch Release (0.1.0 → 0.1.1)

Focus on:
- Bug fixes only
- No new features
- No API changes
- Minimal risk

Quick checklist:
1. Fix the bug
2. Add test case
3. Update CHANGELOG
4. Bump patch version
5. Release

### Minor Release (0.1.0 → 0.2.0)

Can include:
- New features
- Backward-compatible API additions
- Performance improvements
- Deprecations (with warnings)

Additional steps:
- Document new features thoroughly
- Add examples for new functionality
- Update benchmarks
- Consider blog post for significant features

### Major Release (0.9.0 → 1.0.0)

Requires:
- Complete documentation review
- Migration guide
- Extensive testing
- Beta/RC releases
- Community feedback period

Additional considerations:
- Breaking changes clearly documented
- Deprecated APIs removed
- Stability guarantees established
- Long-term support plan

## Pre-Release Process

For alpha, beta, or RC releases:

### Alpha Release (v0.2.0-alpha.1)

```go
const PreRelease = "alpha.1"
```

Use when:
- Feature is implemented but needs testing
- API might still change
- Not production-ready

### Beta Release (v0.2.0-beta.1)

```go
const PreRelease = "beta.1"
```

Use when:
- Feature is complete
- API is stable
- Needs real-world testing
- Bug fixes only from here

### Release Candidate (v0.2.0-rc.1)

```go
const PreRelease = "rc.1"
```

Use when:
- All features complete
- All tests passing
- Documentation complete
- Final validation before stable release

## Troubleshooting

### Tests Failing on CI but Passing Locally

```bash
# Clean everything and retry
make clean-all
go clean -cache -testcache -modcache
make test

# Check for race conditions
make test-race

# Check for platform-specific issues
GOOS=linux GOARCH=amd64 go test ./...
GOOS=linux GOARCH=arm64 go test ./...
```

### Tag Already Exists

```bash
# Delete local tag
git tag -d v0.2.0

# Delete remote tag (⚠️  CAUTION: Only if not yet published!)
git push origin :refs/tags/v0.2.0

# Recreate and push
git tag -a v0.2.0 -m "Release v0.2.0"
git push origin v0.2.0
```

### Go Proxy Not Updating

```bash
# Request explicit update
curl https://proxy.golang.org/github.com/shaia/simdcuckoofilter/@v/v0.2.0.info

# If still not working, wait up to 30 minutes
# The proxy updates on first request, but has caching

# Check proxy directly
curl https://sum.golang.org/lookup/github.com/shaia/simdcuckoofilter@v0.2.0
```

### Version Not Appearing in `go list`

```bash
# Clear local cache
go clean -modcache

# Try direct request
go get github.com/shaia/simdcuckoofilter@v0.2.0

# Check available versions
GOPROXY=https://proxy.golang.org go list -m -versions github.com/shaia/simdcuckoofilter
```

## Emergency Rollback

If a critical bug is discovered immediately after release:

1. **Acknowledge the issue** publicly on GitHub
2. **Do NOT delete the tag** (breaks reproducibility)
3. **Create a patch release** (e.g., v0.2.1) with the fix
4. **Deprecate the broken version** in release notes
5. **Announce the fix** to all channels used for original announcement

## References

- [Semantic Versioning](https://semver.org/)
- [Go Modules Reference](https://go.dev/ref/mod)
- [Publishing Go Modules](https://go.dev/doc/modules/publishing)
- [GitHub Releases Guide](https://docs.github.com/en/repositories/releasing-projects-on-github)
- [Go Module Proxy Protocol](https://go.dev/ref/mod#goproxy-protocol)

## Contact

For questions about the release process:
- Open an issue on GitHub
- Check existing documentation in VERSIONING.md
- Review past releases for examples
