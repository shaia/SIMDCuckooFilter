package fnv

import (
	"hash/fnv"

	"github.com/shaia/simdcuckoofilter/internal/hash/types"
)

// processItemFNV computes hash result for a single item using FNV-1a.
// This is the core hashing logic shared across all platforms.
func processItemFNV(item []byte, fingerprintBits, numBuckets uint) types.HashResult {
	// Hash the item
	hasher := fnv.New64a()
	hasher.Write(item)
	hashVal := hasher.Sum64()

	// Extract fingerprint
	fp := fingerprint(hashVal, fingerprintBits)
	i1 := uint(hashVal % uint64(numBuckets))

	// Calculate alternative index using fingerprint hash
	// Use stack-allocated buffer to avoid heap allocation
	fpHasher := fnv.New64a()
	fpBuf := [2]byte{byte(fp), byte(fp >> 8)}
	len := 1
	if fingerprintBits > 8 {
		len = 2
	}
	fpHasher.Write(fpBuf[:len])
	fpHash := fpHasher.Sum64()
	i2 := (uint64(i1) ^ fpHash) % uint64(numBuckets)

	return types.HashResult{I1: i1, I2: uint(i2), Fp: fp}
}

// processSequential processes items sequentially without goroutines.
// Used for small batches where goroutine overhead exceeds any benefit.
func processSequential(items [][]byte, fingerprintBits, numBuckets uint) []types.HashResult {
	results := make([]types.HashResult, len(items))
	for i, item := range items {
		results[i] = processItemFNV(item, fingerprintBits, numBuckets)
	}
	return results
}

// processParallel processes items in parallel using goroutines.
// Splits work into chunks to maximize CPU utilization for large batches.
func processParallel(items [][]byte, fingerprintBits, numBuckets uint) []types.HashResult {
	results := make([]types.HashResult, len(items))

	// Calculate optimal chunk size
	chunkSize := (len(items) + 3) / 4 // Process in 4 chunks
	if chunkSize < 8 {
		chunkSize = 8 // Minimum chunk size to justify goroutine overhead
	}

	// Build chunk list
	type chunk struct {
		start, end int
	}
	chunks := make([]chunk, 0, 4)
	for i := 0; i < len(items); i += chunkSize {
		end := i + chunkSize
		if end > len(items) {
			end = len(items)
		}
		chunks = append(chunks, chunk{start: i, end: end})
	}

	// Process chunks in parallel
	// Buffer the channel to avoid blocking goroutines waiting to send completion signal
	done := make(chan struct{}, len(chunks))
	for _, c := range chunks {
		go func(start, end int) {
			for i := start; i < end; i++ {
				results[i] = processItemFNV(items[i], fingerprintBits, numBuckets)
			}
			done <- struct{}{}
		}(c.start, c.end)
	}

	// Wait for all chunks to complete
	for range chunks {
		<-done
	}

	return results
}
