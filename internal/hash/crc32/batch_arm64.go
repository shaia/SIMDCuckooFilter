//go:build arm64
// +build arm64

package crc32hash

import (
	"hash/crc32"

	"github.com/shaia/cuckoofilter/internal/hash/types"
	"github.com/shaia/cuckoofilter/internal/simd/cpu"
)

// BatchProcessor handles optimized batch CRC32 hashing for ARM64.
// Uses hardware CRC32 instructions available on ARMv8 and later.
type BatchProcessor struct {
	table    *crc32.Table
	simdType cpu.SIMDType
	useSIMD  bool // Whether to use hardware CRC32 instructions
}

// NewBatchProcessor creates a new CRC32 batch processor for ARM64
func NewBatchProcessor(table *crc32.Table, simdType cpu.SIMDType) *BatchProcessor {
	return &BatchProcessor{
		table:    table,
		simdType: simdType,
		useSIMD:  true, // Enable hardware CRC32 by default
	}
}

// NewBatchProcessorNoSIMD creates a batch processor without hardware CRC32 (for comparison/testing)
func NewBatchProcessorNoSIMD(table *crc32.Table, simdType cpu.SIMDType) *BatchProcessor {
	return &BatchProcessor{
		table:    table,
		simdType: simdType,
		useSIMD:  false, // Disable hardware CRC32
	}
}

// ProcessBatch processes multiple items using optimized CRC32.
//
// Performance characteristics:
//   - Uses ARMv8 hardware CRC32 instructions (CRC32C variant)
//   - Processes items sequentially but with optimized assembly
//   - Hardware acceleration provides ~3-5x speedup over software implementation
//
// ARM64 has dedicated CRC32C instructions that match the Castagnoli polynomial,
// which is exactly what we need for this hash function.
func (p *BatchProcessor) ProcessBatch(items [][]byte, fingerprintBits, numBuckets uint) []types.HashResult {
	results := make([]types.HashResult, len(items))

	// Use hardware CRC32 assembly if enabled and batch is large enough
	if p.useSIMD && len(items) >= 4 {
		return p.processBatchHardware(items, fingerprintBits, numBuckets)
	}

	// For small batches or when SIMD is disabled, use stdlib
	for i, item := range items {
		hashVal := crc32.Checksum(item, p.table)
		fp := fingerprint(uint64(hashVal), fingerprintBits)
		i1 := uint(hashVal % uint32(numBuckets))
		fpHash := crc32.Checksum([]byte{fp}, p.table)
		i2 := (uint64(i1) ^ uint64(fpHash)) % uint64(numBuckets)
		results[i] = types.HashResult{I1: i1, I2: uint(i2), Fp: fp}
	}

	return results
}

// processBatchHardware uses ARM64 CRC32C hardware instructions
func (p *BatchProcessor) processBatchHardware(items [][]byte, fingerprintBits, numBuckets uint) []types.HashResult {
	results := make([]types.HashResult, len(items))
	crc32Results := make([]uint32, len(items))

	// Call assembly function to compute CRC32C values using hardware instructions
	batchCRC32Hardware(items, crc32Results)

	// Convert CRC32 values to hash results
	for i, hashVal := range crc32Results {
		fp := fingerprint(uint64(hashVal), fingerprintBits)
		i1 := uint(hashVal % uint32(numBuckets))
		fpHash := crc32.Checksum([]byte{fp}, p.table)
		i2 := (uint64(i1) ^ uint64(fpHash)) % uint64(numBuckets)
		results[i] = types.HashResult{I1: i1, I2: uint(i2), Fp: fp}
	}

	return results
}
