// Package crc32hash provides CRC32C (Castagnoli) hash implementation.
// This implementation leverages hardware acceleration when available:
// - AMD64: SSE4.2 CRC32 instructions
// - ARM64: ARMv8 CRC32C instructions
// - Other platforms: Optimized Go implementation
package crc32hash

import (
	"hash/crc32"

	"github.com/shaia/cuckoofilter/internal/hash/types"
)

// CRC32Hash implements the CRC32C (Castagnoli) hash function.
// Uses hardware-accelerated CRC32 instructions when available:
// SSE4.2 on AMD64, ARMv8 CRC32 on ARM64.
//
// CRC32Hash instances are safe for concurrent use by multiple goroutines.
type CRC32Hash struct {
	Table           *crc32.Table
	FingerprintBits uint
	batchProcessor  *BatchProcessor
}

// NewCRC32Hash creates a new CRC32Hash instance
func NewCRC32Hash(table *crc32.Table, fingerprintBits uint, batchProcessor *BatchProcessor) *CRC32Hash {
	return &CRC32Hash{
		Table:           table,
		FingerprintBits: fingerprintBits,
		batchProcessor:  batchProcessor,
	}
}

// GetIndices computes the two bucket indices and fingerprint for an item in a cuckoo filter.
//
// This method hashes the input item using CRC32C (Castagnoli) and derives:
//   - i1: The primary bucket index, computed as crc32c(item) % numBuckets
//   - i2: The alternative bucket index, computed as (i1 ^ crc32c(fp)) % numBuckets
//   - fp: A non-zero fingerprint (1-255) extracted from the hash, used to identify the item
//
// CRC32C is hardware-accelerated on modern CPUs (SSE4.2 on AMD64, ARMv8 CRC32 on ARM64),
// providing excellent performance with minimal CPU overhead.
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
//	crc := NewCRC32Hash(crc32.MakeTable(crc32.Castagnoli), 8, nil)
//	i1, i2, fp := crc.GetIndices([]byte("example"), 1024)
//	// i1 and i2 are candidate buckets where the item could be stored
//	// fp identifies the item within those buckets
func (h *CRC32Hash) GetIndices(item []byte, numBuckets uint) (i1, i2 uint, fp byte) {
	// CRC32C checksum (hardware accelerated on modern CPUs)
	hashVal := crc32.Checksum(item, h.Table)

	// Extract fingerprint
	fp = fingerprint(uint64(hashVal), h.FingerprintBits)

	// Calculate first index
	i1 = uint(hashVal % uint32(numBuckets))

	// Calculate second index
	i2 = h.GetAltIndex(i1, fp, numBuckets)

	return i1, i2, fp
}

// GetAltIndex computes the alternative bucket index given a current index and fingerprint.
//
// This method implements the core cuckoo hashing property where each item has exactly two
// possible bucket locations. The calculation uses XOR with the CRC32C hash of the fingerprint:
//
//	altIndex = (index ^ crc32c(fp)) % numBuckets
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
//	crc := NewCRC32Hash(crc32.MakeTable(crc32.Castagnoli), 8, nil)
//	i1, _, fp := crc.GetIndices([]byte("example"), 1024)
//	i2 := crc.GetAltIndex(i1, fp, 1024)  // Get alternative location
//	i1Back := crc.GetAltIndex(i2, fp, 1024)  // Returns to i1 (symmetry property)
func (h *CRC32Hash) GetAltIndex(index uint, fp byte, numBuckets uint) uint {
	// Hash the fingerprint to get alternative index
	// Use stack-allocated buffer for thread safety
	fpBuf := [1]byte{fp}
	fpHash := crc32.Checksum(fpBuf[:], h.Table)
	altIndex := (uint64(index) ^ uint64(fpHash)) % uint64(numBuckets)
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
// While CRC32C is inherently sequential (each byte depends on the previous state),
// batch processing still provides performance benefits through reduced overhead and
// better memory access patterns.
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
// Different goroutines can call this method simultaneously on the same CRC32Hash instance.
//
// Performance considerations:
//   - CRC32C is hardware-accelerated and already very fast for single items
//   - Batch processing provides diminishing returns compared to SIMD-capable hashes
//   - Main benefit is reduced function call overhead for many small items
//
// Example:
//
//	crc := NewCRC32Hash(crc32.MakeTable(crc32.Castagnoli), 8, batchProcessor)
//	items := [][]byte{[]byte("key1"), []byte("key2"), []byte("key3")}
//	results := crc.GetIndicesBatch(items, 1024)
//	for i, result := range results {
//	    fmt.Printf("Item %d: i1=%d, i2=%d, fp=%d\n", i, result.I1, result.I2, result.Fp)
//	}
func (h *CRC32Hash) GetIndicesBatch(items [][]byte, numBuckets uint) []types.HashResult {
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
