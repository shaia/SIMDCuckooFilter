// Package fnv provides FNV-1a (Fowler-Noll-Vo) hash implementation.
// FNV-1a is a simple, fast hash function with good distribution properties.
// It's implemented in pure Go and serves as a reliable fallback option.
package fnv

import (
	"hash/fnv"

	"github.com/shaia/cuckoofilter/internal/hash/types"
)

// FNVHash implements the FNV-1a (Fowler-Noll-Vo) hash function.
// FNV-1a provides good distribution and is implemented in pure Go.
//
// FNVHash instances are safe for concurrent use by multiple goroutines.
type FNVHash struct {
	FingerprintBits uint
	batchProcessor  *BatchProcessor
}

// NewFNVHash creates a new FNVHash instance
func NewFNVHash(fingerprintBits uint, batchProcessor *BatchProcessor) *FNVHash {
	return &FNVHash{
		FingerprintBits: fingerprintBits,
		batchProcessor:  batchProcessor,
	}
}

func (h *FNVHash) GetIndices(item []byte, numBuckets uint) (i1, i2 uint, fp byte) {
	// Hash the item
	hasher := fnv.New64a()
	hasher.Write(item)
	hashVal := hasher.Sum64()

	// Extract fingerprint from hash
	fp = fingerprint(hashVal, h.FingerprintBits)

	// Calculate first index
	i1 = uint(hashVal % uint64(numBuckets))

	// Calculate second index using fingerprint
	i2 = h.GetAltIndex(i1, fp, numBuckets)

	return i1, i2, fp
}

func (h *FNVHash) GetAltIndex(index uint, fp byte, numBuckets uint) uint {
	// Use fingerprint to compute alternative index
	// This ensures i2 != i1 for the same fingerprint
	// Use stack-allocated buffer for thread safety
	hasher := fnv.New64a()
	fpBuf := [1]byte{fp}
	hasher.Write(fpBuf[:])
	fpHash := hasher.Sum64()

	altIndex := (uint64(index) ^ fpHash) % uint64(numBuckets)
	return uint(altIndex)
}

func (h *FNVHash) GetIndicesBatch(items [][]byte, numBuckets uint) []types.HashResult {
	// Use batch processor if available
	if h.batchProcessor != nil {
		return h.batchProcessor.ProcessBatch(items, h.FingerprintBits, numBuckets)
	}

	// Fallback to sequential processing
	results := make([]types.HashResult, len(items))
	for i, item := range items {
		i1, i2, fp := h.GetIndices(item, numBuckets)
		results[i] = types.HashResult{I1: i1, I2: i2, Fp: fp}
	}
	return results
}

// fingerprint extracts a fingerprint from a hash value
func fingerprint(hashVal uint64, bits uint) byte {
	fp := byte(hashVal & ((1 << bits) - 1))
	// Ensure fingerprint is never zero (0 means empty slot)
	if fp == 0 {
		fp = 1
	}
	return fp
}
