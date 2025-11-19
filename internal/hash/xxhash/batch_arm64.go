//go:build arm64
// +build arm64

// Package hash provides optimized batch hash processing for ARM64.
// This file contains the ARM64-specific implementation.
package xxhash

import (
	"github.com/shaia/simdcuckoofilter/internal/hash/types"
)

// BatchHashProcessor handles batch hashing on ARM64.
// Uses optimized ARM64 assembly for single-item hashing (~32% faster than pure Go).
//
// Note: NEON-parallel batch processing (2-4 items in parallel like AVX2 on AMD64)
// is a potential future optimization that would require complex SIMD implementation.
type BatchHashProcessor struct{}

// NewBatchHashProcessor creates a new batch hash processor for ARM64.
func NewBatchHashProcessor() *BatchHashProcessor {
	return &BatchHashProcessor{}
}

// ProcessBatchXXHash processes multiple items using optimized XXHash.
//
// Performance characteristics:
//   - Uses optimized ARM64 assembly for each item (~32% faster than Go)
//   - Processes items sequentially (unlike AVX2 which processes 4 in parallel)
//
// The current implementation processes items sequentially but uses
// the optimized ARM64 assembly hash function for each item, providing
// significant speedup over pure Go implementation.
func (p *BatchHashProcessor) ProcessBatchXXHash(items [][]byte, fingerprintBits, numBuckets uint) []types.HashResult {
	results := make([]types.HashResult, len(items))

	// Process each item using optimized ARM64 assembly hash
	xxh := &XXHash{fingerprintBits: fingerprintBits}
	for i, item := range items {
		i1, i2, fp := xxh.GetIndices(item, numBuckets)
		results[i] = types.HashResult{I1: i1, I2: i2, Fp: fp}
	}
	return results
}
