//go:build amd64
// +build amd64

// Package hash provides SIMD-optimized batch hash processing for AMD64/x86-64.
// This file contains the AMD64-specific implementation using AVX2 and SSE2.
package xxhash

import (
	"github.com/shaia/cuckoofilter/internal/hash/types"
	"github.com/shaia/cuckoofilter/internal/simd/cpu"
)

// BatchHashProcessor handles SIMD-optimized batch hashing for AMD64.
// It automatically selects the best SIMD instruction set (AVX2 or SSE2)
// based on CPU capabilities.
type BatchHashProcessor struct {
	simdType cpu.SIMDType
}

// NewBatchHashProcessor creates a new batch hash processor
func NewBatchHashProcessor(simdType cpu.SIMDType) *BatchHashProcessor {
	return &BatchHashProcessor{
		simdType: simdType,
	}
}

// ProcessBatchXXHash processes multiple items using SIMD-optimized XXHash.
//
// Performance characteristics:
//   - AVX2: Processes 4 items in parallel using 256-bit registers
//   - SSE2: Processes 2 items in parallel using 128-bit registers
//   - Fallback: Uses optimized single-item assembly hash
//
// The function automatically dispatches to the appropriate SIMD implementation
// based on the simdType specified during construction.
func (p *BatchHashProcessor) ProcessBatchXXHash(items [][]byte, fingerprintBits, numBuckets uint) []types.HashResult {
	results := make([]types.HashResult, len(items))

	switch p.simdType {
	case cpu.SIMDAVX2:
		// 4-way parallel processing with AVX2
		processBatchXXHashAVX2(items, results, fingerprintBits, numBuckets)
		return results
	case cpu.SIMDSSE2:
		// 2-way parallel processing with SSE2
		processBatchXXHashSSE2(items, results, fingerprintBits, numBuckets)
		return results
	default:
		// Scalar fallback using optimized single-item hash
		xxh := &XXHash{fingerprintBits: fingerprintBits}
		for i, item := range items {
			i1, i2, fp := xxh.GetIndices(item, numBuckets)
			results[i] = types.HashResult{I1: i1, I2: i2, Fp: fp}
		}
		return results
	}
}

// processBatchXXHashAVX2 is implemented in batch_avx2_amd64.s
// Processes 4 items in parallel using AVX2 256-bit SIMD instructions.
//
//go:noescape
func processBatchXXHashAVX2(items [][]byte, results []types.HashResult, fingerprintBits, numBuckets uint)

// processBatchXXHashSSE2 is implemented in batch_sse2_amd64.s
// Processes 2 items in parallel using SSE2 128-bit SIMD instructions.
//
//go:noescape
func processBatchXXHashSSE2(items [][]byte, results []types.HashResult, fingerprintBits, numBuckets uint)
