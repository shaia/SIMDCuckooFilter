//go:build amd64
// +build amd64

// Package hash provides SIMD-optimized batch hash processing for AMD64/x86-64.
// This file contains the AMD64-specific implementation using AVX2.
package xxhash

import (
	"github.com/shaia/simdcuckoofilter/internal/hash/types"
)

// BatchHashProcessor handles AVX2-optimized batch hashing for AMD64.
// Processes 4 items in parallel using 256-bit AVX2 registers.
type BatchHashProcessor struct{}

// NewBatchHashProcessor creates a new batch hash processor
func NewBatchHashProcessor() *BatchHashProcessor {
	return &BatchHashProcessor{}
}

// ProcessBatchXXHash processes multiple items using XXHash.
// Currently falls back to scalar implementation to support 16-bit fingerprints.
// TODO: Update AVX2 assembly to support 16-bit fingerprints and re-enable.
func (p *BatchHashProcessor) ProcessBatchXXHash(items [][]byte, fingerprintBits, numBuckets uint) []types.HashResult {
	results := make([]types.HashResult, len(items))

	// Use scalar implementation for now
	xxh := &XXHash{fingerprintBits: fingerprintBits}
	for i, item := range items {
		i1, i2, fp := xxh.GetIndices(item, numBuckets)
		results[i] = types.HashResult{I1: i1, I2: i2, Fp: fp}
	}

	return results
}

// processBatchXXHashAVX2 is implemented in batch_avx2_amd64.s
// Processes 4 items in parallel using AVX2 256-bit SIMD instructions.
// Note: Currently unused as we fallback to scalar for 16-bit support
/*
//go:noescape
func processBatchXXHashAVX2(items [][]byte, results []types.HashResult, fingerprintBits, numBuckets uint)
*/
