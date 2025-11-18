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
// TODO: Implement NEON-parallel batch processing (2-4 items in parallel)
type BatchHashProcessor struct{}

// NewBatchHashProcessor creates a new batch hash processor for ARM64.
func NewBatchHashProcessor() *BatchHashProcessor {
	return &BatchHashProcessor{}
}

// ProcessBatchXXHash processes multiple items using optimized XXHash.
//
// Performance characteristics:
//   - Uses optimized ARM64 assembly for each item (~32% faster than Go)
//   - Future: NEON parallel processing for 2-4 items simultaneously
//
// The current implementation processes items sequentially but uses
// the optimized ARM64 assembly hash function for each item.
func (p *BatchHashProcessor) ProcessBatchXXHash(items [][]byte, fingerprintBits, numBuckets uint) []types.HashResult {
	results := make([]types.HashResult, len(items))

	// TODO: Implement NEON-optimized parallel batch processing
	// The challenge is ARM64 assembly syntax differs from x86,
	// requiring careful handling of NEON instructions in Go assembly.
	//
	// For now, use scalar loop with optimized single-item hash
	// which still provides ~32% speedup over pure Go.
	xxh := &XXHash{fingerprintBits: fingerprintBits}
	for i, item := range items {
		i1, i2, fp := xxh.GetIndices(item, numBuckets)
		results[i] = types.HashResult{I1: i1, I2: i2, Fp: fp}
	}
	return results
}
