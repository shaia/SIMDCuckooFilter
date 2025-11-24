# Versioning Guide

This document describes the versioning strategy and release process for the SIMDCuckooFilter project.

## Overview

SIMDCuckooFilter follows [Semantic Versioning 2.0.0](https://semver.org/). Version numbers use the format:

```
vMAJOR.MINOR.PATCH
```

Where:
- **MAJOR**: Incremented for incompatible API changes (breaking changes)
- **MINOR**: Incremented for new functionality in a backward-compatible manner
- **PATCH**: Incremented for backward-compatible bug fixes

## Version Information

The current version is defined in `version.go` with the following constants:

```go
const (
    Major = 0
    Minor = 1
    Patch = 0
    Version = "v0.1.0"
)
```

### Accessing Version Information

You can access version information programmatically:

```go
import "github.com/shaia/simdcuckoofilter"

// Get the full version string
version := simdcuckoofilter.FullVersion()

// Get structured version information
info := simdcuckoofilter.VersionInfo()
fmt.Printf("Version: %s\n", info["version"])
fmt.Printf("Major: %d, Minor: %d, Patch: %d\n",
    info["major"], info["minor"], info["patch"])
```

## Release Process

### Automated Release (Recommended)

The project uses GitHub Actions for automated releases:

1. **Create a Release PR**:
   ```bash
   # Update version in version.go
   # Update CHANGELOG.md
   git checkout -b release/v0.2.0
   git add version.go CHANGELOG.md
   git commit -m "chore: prepare release v0.2.0"
   git push origin release/v0.2.0
   ```

2. **Submit PR**: Create a pull request to merge into `main`

3. **Review and Approval**: Get PR approved and merged

4. **Tag the Release**:
   ```bash
   git checkout main
   git pull origin main
   git tag -a v0.2.0 -m "Release v0.2.0: Description of changes"
   git push origin v0.2.0
   ```

5. **Automated Actions**: GitHub Actions will:
   - Validate the tag format
   - Run all tests and checks
   - Create a GitHub release with auto-generated notes
   - Notify Go module proxy

### Manual Release Process

If you need to create a release manually:

1. **Prepare the Release**:
   ```bash
   # Create a feature branch
   git checkout -b release/vX.Y.Z

   # Update version.go
   # Update CHANGELOG.md
   # Update README.md if needed

   # Commit changes
   git add .
   git commit -m "chore: prepare release vX.Y.Z"
   ```

2. **Run Pre-Release Checks**:
   ```bash
   make ci              # Run all CI checks
   make test-race       # Test with race detector
   make coverage        # Check test coverage
   ```

3. **Create and Push PR**:
   ```bash
   git push origin release/vX.Y.Z
   # Create PR on GitHub
   ```

4. **After PR Approval and Merge**:
   ```bash
   git checkout main
   git pull origin main

   # Create annotated tag
   git tag -a vX.Y.Z -m "Release vX.Y.Z

   Changes:
   - Feature/fix 1
   - Feature/fix 2
   - Feature/fix 3
   "

   # Push tag
   git push origin vX.Y.Z
   ```

5. **Create GitHub Release**:
   - Go to GitHub Releases page
   - Click "Draft a new release"
   - Select the tag you just created
   - Add release notes (can use "Generate release notes")
   - Publish release

6. **Verify Publication**:
   ```bash
   # Wait a few minutes, then verify on Go proxy
   curl https://proxy.golang.org/github.com/shaia/simdcuckoofilter/@v/vX.Y.Z.info
   ```

## Pre-Release Versions

For alpha, beta, or release candidate versions, use the format:

```
vMAJOR.MINOR.PATCH-prerelease.N
```

Examples:
- `v0.2.0-alpha.1` - First alpha release
- `v0.2.0-beta.2` - Second beta release
- `v1.0.0-rc.1` - First release candidate

Update the `PreRelease` constant in `version.go`:

```go
const PreRelease = "alpha.1"
```

## Release Checklist

Before creating a release, ensure:

- [ ] All tests pass (`make test`)
- [ ] Race detector passes (`make test-race`)
- [ ] Code is formatted (`make fmt`)
- [ ] Linter passes (`make lint`)
- [ ] Coverage is acceptable (`make coverage`)
- [ ] `CHANGELOG.md` is updated with all changes
- [ ] `version.go` is updated with the new version
- [ ] Documentation is up to date
- [ ] Breaking changes are clearly documented
- [ ] Migration guide is provided (for major versions)

## Version Examples

### Patch Release (v0.1.0 → v0.1.1)

**When**: Bug fixes, performance improvements, documentation updates

```go
const (
    Major = 0
    Minor = 1
    Patch = 1
    Version = "v0.1.1"
)
```

**Example changes**:
- Fixed race condition in bucket allocation
- Improved hash distribution for small filters
- Updated documentation examples

### Minor Release (v0.1.1 → v0.2.0)

**When**: New features, backward-compatible API additions

```go
const (
    Major = 0
    Minor = 2
    Patch = 0
    Version = "v0.2.0"
)
```

**Example changes**:
- Added support for custom hash functions
- New `BatchContains()` method
- Additional configuration options

### Major Release (v0.2.0 → v1.0.0)

**When**: Breaking API changes, architectural changes

```go
const (
    Major = 1
    Minor = 0
    Patch = 0
    Version = "v1.0.0"
)
```

**Example changes**:
- Changed function signatures
- Removed deprecated APIs
- Restructured package layout

## Best Practices

1. **Always run tests**: Never tag a release with failing tests
2. **Update CHANGELOG**: Keep a detailed changelog for users
3. **Use annotated tags**: Include meaningful commit messages in tags
4. **Document breaking changes**: Clearly explain what changed and why
5. **Provide migration guides**: Help users upgrade smoothly
6. **Follow semantic versioning**: Be consistent and predictable
7. **Test pre-releases**: Use alpha/beta versions for major changes

## Troubleshooting

### Wrong Tag Version

If you created a tag with the wrong version:

```bash
# Delete local tag
git tag -d vX.Y.Z

# Delete remote tag
git push origin :refs/tags/vX.Y.Z

# Create correct tag
git tag -a vX.Y.Z -m "Release vX.Y.Z"
git push origin vX.Y.Z
```

### Update Existing Tag

**WARNING**: Only do this for unpublished releases!

```bash
# Delete and recreate tag
git tag -d vX.Y.Z
git push origin :refs/tags/vX.Y.Z
git tag -a vX.Y.Z -m "Release vX.Y.Z"
git push origin vX.Y.Z
```

### Go Module Not Updated

If the Go module proxy doesn't show your new version:

```bash
# Request explicit update
curl "https://proxy.golang.org/github.com/shaia/simdcuckoofilter/@v/vX.Y.Z.info"

# Check available versions
go list -m -versions github.com/shaia/simdcuckoofilter

# Clear local cache and refetch
go clean -modcache
go get github.com/shaia/simdcuckoofilter@vX.Y.Z
```

## For Users

### Installing a Specific Version

```bash
# Install latest version
go get github.com/shaia/simdcuckoofilter@latest

# Install specific version
go get github.com/shaia/simdcuckoofilter@v0.1.0

# Install from main branch
go get github.com/shaia/simdcuckoofilter@main
```

### Checking Version in Your Code

```go
import (
    "fmt"
    "github.com/shaia/simdcuckoofilter"
)

func main() {
    fmt.Println("Using SIMDCuckooFilter version:", simdcuckoofilter.Version)

    // Get full version info
    info := simdcuckoofilter.VersionInfo()
    fmt.Printf("Version details: %+v\n", info)
}
```

## Version History

See [CHANGELOG.md](CHANGELOG.md) for a complete version history and detailed change information.

## References

- [Semantic Versioning 2.0.0](https://semver.org/)
- [Go Modules Reference](https://go.dev/ref/mod)
- [GitHub Releases](https://docs.github.com/en/repositories/releasing-projects-on-github)
