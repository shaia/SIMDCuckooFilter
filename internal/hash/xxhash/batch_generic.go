//go:build !amd64 && !arm64
// +build !amd64,!arm64

// Package hash provides batch hash processing for generic platforms.
// This file contains the fallback implementation for platforms without
// architecture-specific optimizations (i.e., not AMD64 or ARM64).
package xxhash

import (
	"github.com/shaia/cuckoofilter/internal/hash/types"
	"github.com/shaia/cuckoofilter/internal/simd/cpu"
)

// BatchHashProcessor handles batch hashing using pure Go implementation.
// This is the fallback for platforms without SIMD optimizations.
type BatchHashProcessor struct {
	simdType cpu.SIMDType
}

// NewBatchHashProcessor creates a new batch hash processor.
// On generic platforms, this uses pure Go implementation regardless of simdType.
func NewBatchHashProcessor(simdType cpu.SIMDType) *BatchHashProcessor {
	return &BatchHashProcessor{
		simdType: simdType,
	}
}

// ProcessBatchXXHash processes multiple items using pure Go XXHash.
//
// This is a scalar fallback implementation for platforms without
// architecture-specific optimizations. It uses the pure Go implementation
// of XXHash64 for maximum portability.
func (p *BatchHashProcessor) ProcessBatchXXHash(items [][]byte, fingerprintBits, numBuckets uint) []types.HashResult {
	results := make([]types.HashResult, len(items))
	xxh := &XXHash{fingerprintBits: fingerprintBits}
	for i, item := range items {
		i1, i2, fp := xxh.GetIndices(item, numBuckets)
		results[i] = types.HashResult{I1: i1, I2: i2, Fp: fp}
	}
	return results
}
