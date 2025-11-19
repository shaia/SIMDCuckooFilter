# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added
- **ARM64 NEON SIMD assembly** for bucket lookup operations (`internal/lookup/bucket_lookup_neon_arm64.s`)
  - 16-byte parallel processing using ARM64 NEON instructions
  - ~2-3x performance improvement over scalar implementation for buckets ≥16 bytes
  - Complete parity with AMD64 AVX2 SIMD optimizations
- Batch processing optimizations for CRC32 and FNV hash functions
- Comprehensive test suite organization (SIMD, validation, benchmark tests)
- Makefile for build automation and cross-platform compilation
- TESTING.md documentation for test guidance
- Comprehensive README.md with feature highlights, API reference, and examples
- Generic fallback for SIMD filter on non-amd64/arm64 platforms

### Changed
- **Refactored filter implementation** to share code between AMD64 and ARM64
  - Consolidated duplicate code in `filter.go` (shared implementation)
  - Reduced from 986 to 539 lines in test files (45% reduction)
  - Reduced from 425 to 259 lines in implementation files (39% reduction)
- **Renamed internal directories** for clarity:
  - `internal/simd/` → `internal/lookup/` (better describes purpose)
- Reorganized xxhash package structure for better maintainability
- Refactored hash package with improved organization and documentation
- Refactored filter files with consistent naming convention:
  - `cuckoofilter.go` → `filter.go`
  - `scalar_filter.go` → `filter_scalar.go`
  - `simd_filter.go` → `filter_amd64.go` / `filter_arm64.go`
- Updated examples to reflect NEON implementation
- Improved package documentation across hash implementations

### Fixed
- **Critical: ARM64 assembly calling convention** - Fixed return value offset (32 not 25) due to 8-byte alignment
- **Relocation algorithm bug** - Now uses `bucketSize` instead of `count` for standard cuckoo hashing behavior
- **Race condition in TestFilterConcurrentLookup** - Fixed using buffered channel for thread-safe error collection
- **Overflow in nextPowerOf2** - Added check for values ≥ 2^63 to prevent overflow on 64-bit platforms
- Race condition in random number generator usage
- Potential division by zero panic in `relocate()` function
- Assembly instruction typo: PCMPREQB → PCMPEQB in SSE2 implementation
- Misleading comment in SSE2 bucket contains implementation
- Uninitialized memory access in AVX2 batch hash processing

### Removed
- TODO comments for NEON implementation (now complete)

## [0.1.0] - Previous Release

### Added
- Initial Cuckoo Filter implementation
- SIMD optimizations for AMD64 (SSE2, AVX2)
- Multiple hash strategies (XXHash64, CRC32C, FNV-1a)
- Batch operations support
- Configurable bucket sizes (4, 8, 16, 32, 64)
- Configurable fingerprint sizes (8, 16, 32 bits)
- Thread-safe operations with RWMutex
- Comprehensive test coverage
- Benchmark suite

### Features
- High-performance probabilistic set membership testing
- Low false positive rates
- Support for deletions
- Efficient memory usage
- Hardware-accelerated hashing (CRC32C with SSE4.2, SIMD batch processing)

---

## Release Notes

### Version Numbering
This project follows [Semantic Versioning](https://semver.org/):
- MAJOR version for incompatible API changes
- MINOR version for backwards-compatible functionality additions
- PATCH version for backwards-compatible bug fixes

### Migration Guides

#### Upgrading to Unreleased
If you're upgrading from a previous version, note the following changes:

**File Renames** (internal only, no API changes):
- Filter implementation files have been renamed for consistency
- No changes to public API

**New Features**:
- ARM64 SIMD support is now automatically enabled on ARM64 platforms
- Batch processing for CRC32 and FNV hash strategies
- New test organization for easier test execution

**Bug Fixes**:
- Random number generation is now thread-safe
- AVX2 batch processing no longer accesses uninitialized memory
