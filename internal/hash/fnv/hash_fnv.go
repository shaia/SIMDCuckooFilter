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

// GetIndices computes the two bucket indices and fingerprint for an item in a cuckoo filter.
//
// This method hashes the input item using FNV-1a (Fowler-Noll-Vo) and derives:
//   - i1: The primary bucket index, computed as fnv1a(item) % numBuckets
//   - i2: The alternative bucket index, computed as (i1 ^ fnv1a(fp)) % numBuckets
//   - fp: A non-zero fingerprint (1-255) extracted from the hash, used to identify the item
//
// FNV-1a is a simple, fast non-cryptographic hash function with good distribution properties.
// It's implemented in pure Go without assembly, making it portable across all architectures.
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
// Each call creates a new hasher instance, avoiding any shared state.
//
// Example:
//
//	fnv := NewFNVHash(8, nil)
//	i1, i2, fp := fnv.GetIndices([]byte("example"), 1024)
//	// i1 and i2 are candidate buckets where the item could be stored
//	// fp identifies the item within those buckets
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

// GetAltIndex computes the alternative bucket index given a current index and fingerprint.
//
// This method implements the core cuckoo hashing property where each item has exactly two
// possible bucket locations. The calculation uses XOR with the FNV-1a hash of the fingerprint:
//
//	altIndex = (index ^ fnv1a(fp)) % numBuckets
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
// It uses stack-allocated buffers and creates a new hasher instance per call,
// avoiding any shared state.
//
// Example:
//
//	fnv := NewFNVHash(8, nil)
//	i1, _, fp := fnv.GetIndices([]byte("example"), 1024)
//	i2 := fnv.GetAltIndex(i1, fp, 1024)  // Get alternative location
//	i1Back := fnv.GetAltIndex(i2, fp, 1024)  // Returns to i1 (symmetry property)
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

// GetIndicesBatch computes indices and fingerprints for multiple items efficiently.
//
// This method processes multiple items in a single call. When a batch processor is configured,
// it can leverage optimizations like:
//   - Amortized function call overhead
//   - Better cache utilization through sequential processing
//   - Potential for parallel processing across multiple cores
//
// FNV-1a is a pure Go implementation, so batch processing primarily benefits from reduced
// overhead rather than SIMD parallelization. It serves as a reliable fallback when
// architecture-specific optimizations aren't available.
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
// Different goroutines can call this method simultaneously on the same FNVHash instance.
//
// Performance considerations:
//   - FNV-1a is simple and fast, but not as optimized as XXHash or CRC32C
//   - Batch processing provides moderate benefits through reduced overhead
//   - Best used as a portable fallback when other hash functions aren't suitable
//
// Example:
//
//	fnv := NewFNVHash(8, batchProcessor)
//	items := [][]byte{[]byte("key1"), []byte("key2"), []byte("key3")}
//	results := fnv.GetIndicesBatch(items, 1024)
//	for i, result := range results {
//	    fmt.Printf("Item %d: i1=%d, i2=%d, fp=%d\n", i, result.I1, result.I2, result.Fp)
//	}
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
