# CRC32 Hash Implementation

Hardware-accelerated CRC32C (Castagnoli) hash implementation for the Cuckoo Filter, with platform-specific optimizations.

## Files Organization

### Core Implementation
- **`crc32.go`** - Main CRC32Hash implementation and interface
- **`crc32_test.go`** - Core functionality tests

### AMD64 (x86-64) Platform
- **`batch_amd64.go`** - AMD64-specific batch processor
- **`batch_amd64_asm.go`** - Assembly function declarations
- **`batch_amd64.s`** - SSE4.2 hardware CRC32 assembly implementation
- **`batch_amd64_test.go`** - AMD64-specific SIMD tests

### ARM64 Platform
- **`batch_arm64.go`** - ARM64-specific batch processor
- **`batch_arm64_asm.go`** - Assembly function declarations
- **`batch_arm64.s`** - ARMv8 hardware CRC32C assembly implementation
- **`batch_arm64_test.go`** - ARM64-specific tests
- **`ARM64.md`** - Detailed ARM64 implementation documentation

### Generic Fallback
- **`batch_generic.go`** - Pure Go implementation for other platforms

### Benchmarks & Tests
- **`benchmark_test.go`** - Performance benchmarks
- **`crc32_arm64_test.go`** - Additional ARM64 test cases

## Platform Support

| Platform | Implementation | Instructions | Status |
|----------|----------------|--------------|--------|
| **AMD64** | Hardware CRC32 | SSE4.2 `CRC32` | ✅ Production |
| **ARM64** | Hardware CRC32C | ARMv8 `CRC32C*` | ✅ Production |
| **Other** | Software | Pure Go | ✅ Fallback |

## Usage

The package automatically selects the best implementation based on the build platform:

```go
import crc32hash "github.com/shaia/cuckoofilter/internal/hash/crc32"

table := crc32.MakeTable(crc32.Castagnoli)
h := crc32hash.NewCRC32Hash(table, 8, processor)

// Single item
i1, i2, fp := h.GetIndices(item, numBuckets)

// Batch processing (SIMD-optimized)
results := h.GetIndicesBatch(items, numBuckets)
```

## Performance

- **AMD64**: 3-5x faster with SSE4.2 hardware acceleration
- **ARM64**: 3-5x faster with ARMv8 CRC32C instructions
- **Generic**: Uses Go stdlib's optimized CRC32

## Build Tags

Platform-specific files use build tags:
- `//go:build amd64` - AMD64 only
- `//go:build arm64` - ARM64 only
- `//go:build !amd64 && !arm64` - Generic fallback

## Testing

```bash
# Run all tests
go test ./internal/hash/crc32

# Run platform-specific tests
go test ./internal/hash/crc32 -run TestARM64  # ARM64 only
go test ./internal/hash/crc32 -run TestSIMD   # AMD64 only

# Run benchmarks
go test ./internal/hash/crc32 -bench=.
```

## References

- [CRC32C (Castagnoli) Polynomial](https://en.wikipedia.org/wiki/Cyclic_redundancy_check#CRC-32C_(Castagnoli))
- [Intel SSE4.2 CRC32 Instructions](https://www.intel.com/content/www/us/en/docs/intrinsics-guide/index.html#text=crc32)
- [ARM CRC32 Instructions](https://developer.arm.com/architectures/instruction-sets/intrinsics/#f:@navigationhierarchiessimdisa=[Neon]&q=crc32)
