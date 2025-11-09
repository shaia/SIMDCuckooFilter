# ARM64 CRC32 Hardware Implementation

## Overview

This directory contains an optimized ARM64 implementation of CRC32C (Castagnoli) hashing using ARMv8 hardware CRC32 instructions.

## Files

- **[batch_arm64.go](batch_arm64.go)** - ARM64-specific batch processor
- **[batch_arm64_asm.go](batch_arm64_asm.go)** - Go wrapper for assembly function
- **[batch_arm64.s](batch_arm64.s)** - ARMv8 assembly implementation using hardware CRC32C instructions
- **[crc32_arm64_test.go](crc32_arm64_test.go)** - Comprehensive tests for ARM64 implementation

## Hardware Instructions Used

The implementation uses ARMv8 CRC32 hardware instructions:

| Instruction | Description | Data Size |
|-------------|-------------|-----------|
| `CRC32CB` | CRC32C byte | 8-bit |
| `CRC32CH` | CRC32C halfword | 16-bit |
| `CRC32CW` | CRC32C word | 32-bit |
| `CRC32CX` | CRC32C doubleword | 64-bit |

These instructions compute the Castagnoli polynomial (CRC32C), which matches `hash/crc32.Castagnoli` exactly.

## Implementation Details

### Algorithm

The assembly implementation processes data in chunks for optimal performance:

1. **Initialize**: CRC32 starts at `0xFFFFFFFF`
2. **Process 8-byte chunks**: Use `CRC32CX` for 64-bit data
3. **Process 4-byte chunks**: Use `CRC32CW` for remaining 32-bit data
4. **Process 2-byte chunks**: Use `CRC32CH` for remaining 16-bit data
5. **Process remaining bytes**: Use `CRC32CB` for final bytes
6. **Finalize**: Invert all bits with `MVN` (bitwise NOT)

### Performance Characteristics

- **Hardware acceleration**: 3-5x faster than pure software CRC32 (when not using stdlib)
- **Note**: Go's standard library already uses hardware CRC32 on ARM64, so the performance benefit is marginal for this use case
- **Optimized for**: Small to medium batch sizes (4-64 items)

### Build Tags

```go
//go:build arm64
// +build arm64
```

Only builds on ARM64 platforms (Apple Silicon, AWS Graviton, etc.)

## Testing

The implementation includes comprehensive tests:

- **Correctness**: Validates against Go's stdlib `hash/crc32`
- **Edge cases**: Empty data, large data, unaligned data, mixed sizes
- **Consistency**: Hardware vs software implementation comparison
- **Batch sizes**: Tests with 1-64 items

Run tests:
```bash
go test -v ./internal/hash/crc32 -run TestARM64
```

Run benchmarks:
```bash
go test -bench=ARM64 -benchmem ./internal/hash/crc32
```

## Platform Support

✅ **Supported Platforms:**
- Apple Silicon (M1, M2, M3)
- AWS Graviton (ARM64)
- Any ARMv8+ processor with CRC32 extension

❌ **Not supported:**
- ARMv7 and earlier (no CRC32 instructions)
- Non-ARM platforms (use AMD64 or generic implementation)

## Comparison with Other Platforms

| Platform | Implementation | Instructions |
|----------|----------------|--------------|
| **ARM64** | Hardware CRC32C | CRC32CB/CH/CW/CX |
| **AMD64** | Hardware CRC32C | SSE4.2 CRC32 |
| **Generic** | Software | Pure Go fallback |

## Future Improvements

Potential optimizations for future versions:

1. **NEON parallel processing**: Process multiple items simultaneously using NEON SIMD
2. **Prefetching**: Use ARM64 prefetch instructions for large data
3. **Batch optimization**: Optimize for specific batch sizes
4. **Memory alignment**: Ensure optimal memory access patterns

## References

- [ARMv8 Architecture Reference Manual](https://developer.arm.com/documentation/ddi0487/latest)
- [ARM CRC32 Instructions](https://developer.arm.com/architectures/instruction-sets/intrinsics/#f:@navigationhierarchiessimdisa=[Neon]&f:@navigationhierarchiesinstructionsetsarchitectures=[A64]&q=crc32)
- [CRC32C (Castagnoli) Polynomial](https://en.wikipedia.org/wiki/Cyclic_redundancy_check#CRC-32C_(Castagnoli))
