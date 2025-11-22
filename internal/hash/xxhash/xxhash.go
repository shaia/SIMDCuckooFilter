// Package xxhash provides SIMD-optimized XXHash64 implementation.
// This is the primary hash function used by the Cuckoo Filter with
// multi-architecture SIMD support:
//   - AMD64: AVX2 SIMD batch processing
//   - ARM64: Optimized assembly (NEON batch processing planned for future)
package xxhash

import "github.com/shaia/simdcuckoofilter/internal/hash/types"

// XXHash is a SIMD-optimized XXHash64 implementation.
// For production use, this provides excellent performance across architectures.
// Alternatively, consider github.com/cespare/xxhash/v2 for a pure Go version.
//
// XXHash instances are safe for concurrent use by multiple goroutines.
type XXHash struct {
	fingerprintBits uint
	batchProcessor  *BatchHashProcessor
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

// GetIndices computes the two bucket indices and fingerprint for an item in a cuckoo filter.
//
// This method hashes the input item using XXHash64 and derives:
//   - i1: The primary bucket index, computed as hash(item) % numBuckets
//   - i2: The alternative bucket index, computed as (i1 ^ hash(fp)) % numBuckets
//   - fp: A non-zero fingerprint (1-255) extracted from the hash, used to identify the item
//
// Parameters:
//   - item: The data to hash (typically a key or value being inserted into the filter)
//   - numBuckets: The total number of buckets in the cuckoo filter
//
// Returns:
//   - i1: Primary bucket index (0 <= i1 < numBuckets)
//   - i2: Alternative bucket index (0 <= i2 < numBuckets), where i2 = GetAltIndex(i1, fp, numBuckets)
//   - fp: Fingerprint byte (1 <= fp <= 255, never 0 as that indicates an empty slot)
//
// Thread-safety: This method is safe for concurrent use by multiple goroutines.
//
// Example:
//
//	xxh := NewXXHash(8, nil)
//	i1, i2, fp := xxh.GetIndices([]byte("example"), 1024)
//	// i1 and i2 are candidate buckets where the item could be stored
//	// fp identifies the item within those buckets
func (h *XXHash) GetIndices(item []byte, numBuckets uint) (uint, uint, uint16) {
	hashVal := h.hash64(item)

	// Extract fingerprint
	fp := fingerprint(hashVal, h.fingerprintBits)

	// Calculate first index
	i1 := uint(hashVal % uint64(numBuckets))

	// Calculate second index
	i2 := h.GetAltIndex(i1, fp, numBuckets)

	return i1, i2, fp
}

// GetAltIndex computes the alternative bucket index given a current index and fingerprint.
//
// This method implements the core cuckoo hashing property where each item has exactly two
// possible bucket locations. The calculation uses XOR with the hash of the fingerprint:
//
//	altIndex = (index ^ hash(fp)) % numBuckets
//
// This formula has the important mathematical property that applying it twice returns the
// original index (since XOR is self-inverse):
//
//	GetAltIndex(GetAltIndex(i1, fp, n), fp, n) == i1
//
// This symmetry allows the cuckoo filter to efficiently swap between the two candidate
// bucket locations during insertion and lookup operations.
//
// Parameters:
//   - index: The current bucket index (typically i1 or i2)
//   - fp: The fingerprint byte associated with the item
//   - numBuckets: The total number of buckets in the cuckoo filter
//
// Returns:
//   - The alternative bucket index (0 <= altIndex < numBuckets)
//
// Thread-safety: This method is safe for concurrent use by multiple goroutines.
// It uses stack-allocated buffers to avoid shared state.
//
// Example:
//
//	xxh := NewXXHash(8, nil)
//	i1, _, fp := xxh.GetIndices([]byte("example"), 1024)
//	i2 := xxh.GetAltIndex(i1, fp, 1024)  // Get alternative location
//	i1Back := xxh.GetAltIndex(i2, fp, 1024)  // Returns to i1 (symmetry property)
func (h *XXHash) GetAltIndex(index uint, fp uint16, numBuckets uint) uint {
	// 0x5bd1e995 is the MurmurHash2 constant, used here for mixing
	// This is faster than a full hash and sufficient for alternative index
	hash := uint64(fp) * 0x5bd1e995

	// Ensure hash is non-zero modulo numBuckets (which is power of 2)
	// by making sure it is odd. This guarantees i2 != i1 when numBuckets > 1.
	if numBuckets > 1 {
		hash |= 1
	}

	return uint((uint64(index) ^ hash) % uint64(numBuckets))
}

// GetIndicesBatch computes indices and fingerprints for multiple items efficiently.
//
// This method processes multiple items in a single call, leveraging SIMD optimizations
// when available to achieve significant performance gains:
//   - AMD64: AVX2 (4-way parallel) batch processing
//   - ARM64: Optimized assembly (NEON batch processing planned for future)
//   - Fallback: Sequential scalar processing if no SIMD support
//
// The batch processing can be 2-4x faster than calling GetIndices repeatedly,
// especially for small to medium-sized items. SIMD instructions process multiple
// hash computations simultaneously, amortizing memory access costs.
//
// Parameters:
//   - items: Slice of byte slices to hash (can be variable length)
//   - numBuckets: The total number of buckets in the cuckoo filter
//
// Returns:
//   - Slice of HashResult structs, one per input item, in the same order.
//     Each result contains: I1 (primary index), I2 (alternative index), Fp (fingerprint)
//
// Thread-safety: This method is safe for concurrent use by multiple goroutines.
// Different goroutines can call this method simultaneously on the same XXHash instance.
//
// Performance considerations:
//   - Batch sizes of 4-16 items typically show the best SIMD speedup
//   - Very small batches (1-2 items) may not benefit from SIMD overhead
//   - Large batches (100+ items) amortize any setup costs effectively
//
// Example:
//
//	xxh := NewXXHash(8, batchProcessor)
//	items := [][]byte{[]byte("key1"), []byte("key2"), []byte("key3")}
//	results := xxh.GetIndicesBatch(items, 1024)
//	for i, result := range results {
//	    fmt.Printf("Item %d: i1=%d, i2=%d, fp=%d\n", i, result.I1, result.I2, result.Fp)
//	}
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
func fingerprint(hashVal uint64, bits uint) uint16 {
	fp := uint16(hashVal & ((1 << bits) - 1))
	// Ensure fingerprint is never zero (0 means empty slot)
	if fp == 0 {
		fp = 1
	}
	return fp
}
