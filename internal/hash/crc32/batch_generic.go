//go:build !amd64 && !arm64
// +build !amd64,!arm64

package crc32hash

import (
	"hash/crc32"

	"github.com/shaia/cuckoofilter/internal/hash/types"
	"github.com/shaia/cuckoofilter/internal/simd/cpu"
)

// BatchProcessor handles batch CRC32 hashing for generic platforms.
type BatchProcessor struct {
	table    *crc32.Table
	simdType cpu.SIMDType
}

// NewBatchProcessor creates a new CRC32 batch processor
func NewBatchProcessor(table *crc32.Table, simdType cpu.SIMDType) *BatchProcessor {
	return &BatchProcessor{
		table:    table,
		simdType: simdType,
	}
}

// ProcessBatch processes multiple items using CRC32.
// Generic implementation uses sequential processing.
func (p *BatchProcessor) ProcessBatch(items [][]byte, fingerprintBits, numBuckets uint) []types.HashResult {
	results := make([]types.HashResult, len(items))

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
