# SIMD Cuckoo Filter

A high-performance, SIMD-optimized Cuckoo filter implementation in Go with native assembly for AMD64 (AVX2) and ARM64 (NEON).

## Features

- **SIMD Acceleration**: Native assembly implementations for maximum performance
  - AMD64: AVX2 instructions (32-byte parallel processing)
  - ARM64: NEON instructions (16-byte parallel processing)
- **Concurrent-Safe**: Thread-safe operations with mutex protection
- **Batch Operations**: Optimized batch insert, lookup, and delete
- **Configurable**: Customizable fingerprint size, bucket size, and relocation parameters
- **Zero Dependencies**: Pure Go with assembly optimizations

## Performance

SIMD optimizations provide significant speedups for bucket operations:

- **ARM64 NEON**: ~2-3x faster than scalar for buckets ≥16 bytes
- **AMD64 AVX2**: ~3-4x faster than scalar for buckets ≥32 bytes
- **Batch Operations**: Additional speedup through parallel hash computation

## Installation

```bash
go get github.com/shaia/simdcuckoofilter
```

## Quick Start

```go
package main

import (
    "fmt"
    "github.com/shaia/simdcuckoofilter"
)

func main() {
    // Create a filter with capacity for 10,000 items
    filter, err := cuckoofilter.New(10000)
    if err != nil {
        panic(err)
    }

    // Insert items
    filter.Insert([]byte("apple"))
    filter.Insert([]byte("banana"))

    // Lookup items
    found := filter.Lookup([]byte("apple")) // true
    notFound := filter.Lookup([]byte("grape")) // false (probably)

    // Delete items
    filter.Delete([]byte("banana"))

    // Check statistics
    fmt.Printf("Count: %d\n", filter.Count())
    fmt.Printf("Load Factor: %.2f%%\n", filter.LoadFactor()*100)
}
```

## Advanced Configuration

```go
filter, err := cuckoofilter.New(10000,
    cuckoofilter.WithFingerprintSize(8),  // 8-bit fingerprints (max supported)
    cuckoofilter.WithBucketSize(32),      // Larger buckets for better SIMD utilization
    cuckoofilter.WithMaxKicks(500),       // Relocation attempts before failure
    cuckoofilter.WithXXHash(),            // Use XXHash64 for better performance
)
```

### Configuration Options

| Option | Description | Default | Notes |
|--------|-------------|---------|-------|
| `WithFingerprintSize(bits)` | Bits per fingerprint | 8 | Valid: 1, 2, 4, 8 |
| `WithBucketSize(size)` | Fingerprints per bucket | 4 | Valid: 2, 4, 8, 16, 32, 64 |
| `WithMaxKicks(kicks)` | Relocation attempts | 500 | Range: 1-1000 |
| `WithFNVHash()` | Use FNV-1a hash (default) | ✓ | Moderate speed, good distribution |
| `WithXXHash()` | Use XXHash64 | | Fast, excellent distribution |
| `WithCRC32Hash()` | Use CRC32C | | Fastest, hardware-accelerated |
| `WithBatchSize(size)` | Batch processing size | 32 | Range: 1-256 |

## Batch Operations

Batch operations provide better performance through parallel hash computation:

```go
items := [][]byte{
    []byte("item1"),
    []byte("item2"),
    []byte("item3"),
}

// Batch insert
results := filter.InsertBatch(items)

// Batch lookup
found := filter.LookupBatch(items)

// Batch delete
deleted := filter.DeleteBatch(items)
```

## API Reference

### Creation

- `New(capacity uint, opts ...Option) (*CuckooFilter, error)` - Create a new filter

### Operations

- `Insert(item []byte) bool` - Insert an item
- `Lookup(item []byte) bool` - Check if item exists
- `Delete(item []byte) bool` - Remove an item
- `InsertBatch(items [][]byte) []bool` - Batch insert
- `LookupBatch(items [][]byte) []bool` - Batch lookup
- `DeleteBatch(items [][]byte) []bool` - Batch delete

### Statistics

- `Count() uint` - Number of items in filter
- `Capacity() uint` - Maximum capacity
- `LoadFactor() float64` - Current load (0.0 to 1.0)
- `OptimalBatchSize() int` - Recommended batch size
- `Reset()` - Clear all items

## Architecture

### SIMD Implementations

**AMD64 (AVX2)**
- Bucket lookup: 32 bytes processed in parallel
- Batch hashing: 4 items processed simultaneously
- File: `internal/lookup/bucket_lookup_avx2_amd64.s`

**ARM64 (NEON)**
- Bucket lookup: 16 bytes processed in parallel
- Optimized assembly hashing: ~32% faster than Go
- File: `internal/lookup/bucket_lookup_neon_arm64.s`

### Hash Strategies

| Strategy | Speed | Quality | Use Case |
|----------|-------|---------|----------|
| FNV-1a | Moderate | Good | Default, compatibility |
| XXHash64 | Fast | Excellent | General purpose, better distribution |
| CRC32C | Fastest | Good | High-throughput scenarios |

## Memory Usage

Memory usage depends on configuration:

```
Memory = numBuckets × bucketSize × (fingerprintBits / 8)
```

Example for 10,000 item capacity:
- 8-bit fingerprints, 4-entry buckets: ~10 KB
- 8-bit fingerprints, 32-entry buckets: ~80 KB

## False Positive Rate

False positive probability depends on fingerprint size and load factor:

| Fingerprint Bits | Load 50% | Load 75% | Load 90% |
|------------------|----------|----------|----------|
| 1 bit | ~25% | ~37% | ~45% |
| 2 bits | ~12% | ~19% | ~25% |
| 4 bits | ~3% | ~6% | ~9% |
| 8 bits | ~0.4% | ~0.8% | ~1.2% |

## Examples

See the `examples/` directory for complete examples:

- `examples/basic_usage/` - Simple insert/lookup/delete operations
- `examples/custom_config/` - Advanced configuration options

## Testing

Run the comprehensive test suite:

```bash
# All tests
go test ./...

# With race detection
go test -race ./...

# Benchmarks
go test -bench=. ./...
```

See `TESTING.md` for detailed testing documentation.

## Platform Support

- **AMD64**: Requires AVX2 support (Intel Haswell 2013+, AMD Excavator 2015+)
- **ARM64**: Requires NEON support (all ARM64 processors)
- **Other platforms**: Falls back to Go implementation (not recommended)

## Limitations

- **Fingerprint size**: Maximum 8 bits (1 byte)
- **No resizing**: Filter capacity is fixed at creation
- **False positives**: Small probability of false positives (no false negatives)
- **Delete caveat**: Deleting non-existent items may cause false negatives

## Contributing

Contributions are welcome! Please ensure:

1. All tests pass: `go test ./...`
2. Code is formatted: `go fmt ./...`
3. Assembly changes include tests and benchmarks

## License

[License information to be added]

## References

- [Cuckoo Filter: Practically Better Than Bloom](https://www.cs.cmu.edu/~dga/papers/cuckoo-conext2014.pdf) (Fan et al., 2014)
- [Intel AVX2 Intrinsics Guide](https://www.intel.com/content/www/us/en/docs/intrinsics-guide/index.html)
- [ARM NEON Intrinsics Reference](https://developer.arm.com/architectures/instruction-sets/intrinsics/)
