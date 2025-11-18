# SIMD-Optimized Bucket Operations

This directory contains SIMD-optimized implementations of bucket operations for the Cuckoo Filter.

## Overview

The bucket is the core data structure that holds fingerprints. Optimizing bucket operations can significantly improve overall filter performance, especially for larger bucket sizes.

## Supported Operations

### SIMD-Optimized Methods
- **ContainsSIMD(fp byte)** - Check if fingerprint exists
- **IsFullSIMD()** - Check if bucket has no empty slots
- **CountSIMD()** - Count non-zero fingerprints
- **FindFirstZeroSIMD()** - Find first empty slot
- **InsertSIMD(fp byte)** - Insert using SIMD-accelerated search

## Architecture Support

### AMD64 (x86-64)
- **Implementation**: Inline scalar loops (Go compiler optimized)
- **File**: `bucket_simd_amd64.go`
- **Rationale**: Buckets are typically small (4-64 bytes), where SIMD overhead exceeds benefits
- **Benefit**: Go compiler auto-vectorization, excellent branch prediction

### ARM64 (Apple Silicon, etc.)
- **Implementation**: Unrolled scalar loops optimized for ARM64
- **File**: `bucket_simd_arm64.s`
- **Benefit**: Eliminates loop overhead, better branch prediction for larger buckets (16-64 bytes)

### Other Platforms
- **Implementation**: Pure Go inline loops (same as AMD64)
- **File**: `bucket_simd_inline.go`
- **Benefit**: Portable, works everywhere, compiler-optimized

## Performance Characteristics

### ARM64 Benchmarks (Apple Silicon M-series)

| Operation | Size | Scalar (ns) | SIMD (ns) | Speedup |
|-----------|------|-------------|-----------|---------|
| Contains  | 2    | 0.79        | 2.08      | 0.38x ❌ |
| Contains  | 4    | 1.31        | 3.82      | 0.34x ❌ |
| Contains  | 8    | 2.34        | 3.82      | 0.61x ❌ |
| Contains  | 16   | 4.40        | 4.16      | **1.06x ✓** |
| Contains  | 32   | 8.56        | 5.99      | **1.43x ✓** |
| Contains  | 64   | 17.28       | 8.20      | **2.11x ✓** |
| IsFull    | 2    | 1.90        | 3.18      | 0.60x ❌ |
| IsFull    | 4    | 2.61        | 3.97      | 0.66x ❌ |
| IsFull    | 8    | 4.26        | 4.46      | 0.96x ≈ |
| IsFull    | 16   | 7.99        | 5.41      | **1.48x ✓** |
| IsFull    | 32   | 15.70       | 8.35      | **1.88x ✓** |
| IsFull    | 64   | 31.05       | 12.72     | **2.44x ✓** |
| Count     | 2    | 1.98        | 3.18      | 0.62x ❌ |
| Count     | 4    | 2.54        | 3.83      | 0.66x ❌ |
| Count     | 8    | 3.81        | 4.45      | 0.86x ❌ |
| Count     | 16   | 6.38        | 5.42      | **1.18x ✓** |
| Count     | 32   | 11.43       | 8.06      | **1.42x ✓** |
| Count     | 64   | 28.35       | 12.76     | **2.22x ✓** |

**Key Findings (ARM64 only):**
- **Size 64**: ARM64 assembly shows **111% improvement** for Contains, **144% improvement** for IsFull, **122% improvement** for Count - **BEST PERFORMANCE!**
- **Size 32**: ARM64 assembly shows **43% improvement** for Contains, **88% improvement** for IsFull, **42% improvement** for Count
- **Size 16**: ARM64 assembly shows **6% improvement** for Contains, **48% improvement** for IsFull, **18% improvement** for Count
- **Size 8**: ARM64 assembly shows marginal performance (no significant gains)
- **Sizes 2 & 4**: Go compiler optimizations are superior to assembly
- **Recommendation**: Use SIMD methods for bucket sizes **64** (optimal), **32**, and **16** on ARM64

### AMD64 Performance
- **Implementation**: Uses inline scalar loops (same as generic Go code)
- **Performance**: Relies on Go compiler auto-vectorization and branch prediction
- **Benefit**: Simplicity and portability without assembly maintenance burden
- **Note**: For bucket operations specifically, SIMD assembly overhead exceeds benefits for small data (4-64 bytes)

## Usage

```go
// Create SIMD-optimized bucket
bucket := NewSIMDBucket(64) // Size 64 recommended for optimal cache line fit and maximum SIMD benefits

// SIMD operations automatically use best implementation
if bucket.ContainsSIMD(fingerprint) {
    // Found
}

// Check if bucket is full
if bucket.IsFullSIMD() {
    // Need to kick out an item
}

// Fast insert with SIMD search
bucket.InsertSIMD(fingerprint)
```

## When to Use SIMD Buckets

### ✅ Use SIMD methods when:
- Bucket size is **64** (ARM64: 111-144% faster! Perfect cache line alignment)
- Bucket size is **32** (ARM64: 42-88% faster)
- Bucket size is **16** (ARM64: 6-48% faster)
- Running on **ARM64** (Apple Silicon, AWS Graviton) where assembly provides real benefits
- Doing many bucket lookups (hot path)
- Optimizing for throughput with larger buckets

### ❌ Don't use SIMD methods when:
- Bucket size is **2 or 4** (Go scalar code is faster even on ARM64)
- Bucket size is **8** (marginal gains, not worth complexity)
- Running on **AMD64** (uses same scalar code as regular methods)
- Optimizing for code simplicity over performance

## Implementation Details

### AMD64 Inline Scalar
```go
// Uses standard Go loops - compiler optimizes with auto-vectorization
func inlineContains(data []byte, fp byte) bool {
	for _, b := range data {
		if b == fp {
			return true
		}
	}
	return false
}
```

### ARM64 Unrolled Assembly
```asm
; Manually unrolled for specific sizes (e.g., 64 bytes)
MOVBU 0(R0), R3
CMP R2, R3
BEQ found
MOVBU 1(R0), R3
CMP R2, R3
BEQ found
; ... continues for all 64 bytes (unrolled)
```

## Future Improvements

1. **AMD64 SSE2/AVX2**: True SIMD for bucket operations if profiling shows benefits
2. **ARM NEON SIMD**: Use vector instructions once Go assembler supports them fully
3. **Batch Operations**: Process multiple buckets in parallel
4. **Auto-tuning**: Dynamically select scalar vs assembly based on runtime benchmarks

## Testing

Run tests:
```bash
go test ./internal/bucket/... -v
```

Run benchmarks:
```bash
go test ./internal/bucket/... -bench=. -benchmem
```

Compare scalar vs SIMD:
```bash
go test ./internal/bucket/... -bench=BucketContains -benchmem
```

## References

- [Go ARM64 Assembly](https://golang.org/doc/asm)
- [Intel SSE2 Intrinsics Guide](https://www.intel.com/content/www/us/en/docs/intrinsics-guide/)
- [ARM Architecture Reference Manual](https://developer.arm.com/documentation/)
