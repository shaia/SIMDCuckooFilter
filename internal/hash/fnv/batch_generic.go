//go:build !amd64
// +build !amd64

package fnv

import (
	"hash/fnv"

	"github.com/shaia/cuckoofilter/internal/hash/types"
	"github.com/shaia/cuckoofilter/internal/simd/cpu"
)

// BatchProcessor handles batch FNV hashing for generic platforms.
type BatchProcessor struct {
	simdType cpu.SIMDType
}

// NewBatchProcessor creates a new FNV batch processor
func NewBatchProcessor(simdType cpu.SIMDType) *BatchProcessor {
	return &BatchProcessor{
		simdType: simdType,
	}
}

// ProcessBatch processes multiple items using FNV-1a.
// Generic implementation uses sequential processing.
func (p *BatchProcessor) ProcessBatch(items [][]byte, fingerprintBits, numBuckets uint) []types.HashResult {
	results := make([]types.HashResult, len(items))

	for i, item := range items {
		hasher := fnv.New64a()
		hasher.Write(item)
		hashVal := hasher.Sum64()

		fp := fingerprint(hashVal, fingerprintBits)
		i1 := uint(hashVal % uint64(numBuckets))

		fpHasher := fnv.New64a()
		fpHasher.Write([]byte{fp})
		fpHash := fpHasher.Sum64()
		i2 := (uint64(i1) ^ fpHash) % uint64(numBuckets)

		results[i] = types.HashResult{I1: i1, I2: uint(i2), Fp: fp}
	}

	return results
}
