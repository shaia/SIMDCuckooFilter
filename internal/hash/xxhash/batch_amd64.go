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

// ProcessBatchXXHash processes multiple items using AVX2-optimized XXHash.
// Processes 4 items in parallel using 256-bit AVX2 registers for maximum performance.
func (p *BatchHashProcessor) ProcessBatchXXHash(items [][]byte, fingerprintBits, numBuckets uint) []types.HashResult {
	results := make([]types.HashResult, len(items))
	processBatchXXHashAVX2(items, results, fingerprintBits, numBuckets)
	return results
}

// processBatchXXHashAVX2 is implemented in batch_avx2_amd64.s
// Processes 4 items in parallel using AVX2 256-bit SIMD instructions.
//
//go:noescape
func processBatchXXHashAVX2(items [][]byte, results []types.HashResult, fingerprintBits, numBuckets uint)
