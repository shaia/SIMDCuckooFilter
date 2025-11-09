// Package xxhash provides SIMD-optimized XXHash64 implementation.
// This is the primary hash function used by the Cuckoo Filter with
// multi-architecture SIMD support (AMD64 AVX2/SSE2, ARM64 NEON).
package xxhash

import "github.com/shaia/cuckoofilter/internal/hash/types"

// XXHash is a SIMD-optimized XXHash64 implementation.
// For production use, this provides excellent performance across architectures.
// Alternatively, consider github.com/cespare/xxhash/v2 for a pure Go version.
type XXHash struct {
	fingerprintBits uint
	batchProcessor  *BatchHashProcessor
	fpBuf           [1]byte // Reusable buffer for GetAltIndex to avoid allocations
}

// NewXXHash creates a new XXHash instance
func NewXXHash(fingerprintBits uint, batchProcessor *BatchHashProcessor) *XXHash {
	return &XXHash{
		fingerprintBits: fingerprintBits,
		batchProcessor:  batchProcessor,
	}
}

const (
	prime64_1 = 11400714785074694791
	prime64_2 = 14029467366897019727
	prime64_3 = 1609587929392839161
	prime64_4 = 9650029242287828579
	prime64_5 = 2870177450012600261
)

func (h *XXHash) GetIndices(item []byte, numBuckets uint) (i1, i2 uint, fp byte) {
	hashVal := h.hash64(item)

	// Extract fingerprint
	fp = fingerprint(hashVal, h.fingerprintBits)

	// Calculate first index
	i1 = uint(hashVal % uint64(numBuckets))

	// Calculate second index
	i2 = h.GetAltIndex(i1, fp, numBuckets)

	return i1, i2, fp
}

func (h *XXHash) GetAltIndex(index uint, fp byte, numBuckets uint) uint {
	h.fpBuf[0] = fp
	fpHash := h.hash64(h.fpBuf[:])
	altIndex := (uint64(index) ^ fpHash) % uint64(numBuckets)
	return uint(altIndex)
}

func (h *XXHash) GetIndicesBatch(items [][]byte, numBuckets uint) []types.HashResult {
	// Use SIMD batch processor if available
	if h.batchProcessor != nil {
		return h.batchProcessor.ProcessBatchXXHash(items, h.fingerprintBits, numBuckets)
	}

	// Scalar fallback
	results := make([]types.HashResult, len(items))
	for i, item := range items {
		i1, i2, fp := h.GetIndices(item, numBuckets)
		results[i] = types.HashResult{I1: i1, I2: i2, Fp: fp}
	}
	return results
}

func (h *XXHash) hash64(data []byte) uint64 {
	return hash64XXHashInternal(data)
}

// hash64XXHashGo is the Go fallback implementation
func hash64XXHashGo(data []byte) uint64 {
	var hash uint64

	if len(data) >= 8 {
		hash = prime64_5 + uint64(len(data))
		for len(data) >= 8 {
			k := uint64(data[0]) | uint64(data[1])<<8 | uint64(data[2])<<16 | uint64(data[3])<<24 |
				uint64(data[4])<<32 | uint64(data[5])<<40 | uint64(data[6])<<48 | uint64(data[7])<<56
			k *= prime64_2
			k = (k << 31) | (k >> 33)
			k *= prime64_1
			hash ^= k
			hash = ((hash << 27) | (hash >> 37)) * prime64_1
			hash += prime64_4
			data = data[8:]
		}
	} else {
		hash = prime64_5 + uint64(len(data))
	}

	for len(data) > 0 {
		hash ^= uint64(data[0]) * prime64_5
		hash = ((hash << 11) | (hash >> 53)) * prime64_1
		data = data[1:]
	}

	hash ^= hash >> 33
	hash *= prime64_2
	hash ^= hash >> 29
	hash *= prime64_3
	hash ^= hash >> 32

	return hash
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
