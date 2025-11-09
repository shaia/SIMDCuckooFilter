//go:build amd64
// +build amd64

package fnv

import (
	"hash/fnv"

	"github.com/shaia/cuckoofilter/internal/hash/types"
	"github.com/shaia/cuckoofilter/internal/simd/cpu"
)

// BatchProcessor handles optimized batch FNV hashing for AMD64.
// Uses parallel processing to maximize throughput.
type BatchProcessor struct {
	simdType cpu.SIMDType
}

// NewBatchProcessor creates a new FNV batch processor
func NewBatchProcessor(simdType cpu.SIMDType) *BatchProcessor {
	return &BatchProcessor{
		simdType: simdType,
	}
}

// ProcessBatch processes multiple items using optimized FNV-1a.
//
// Performance characteristics:
//   - Uses parallel goroutines for better throughput
//   - FNV-1a is a simple XOR+multiply operation, making it fast
//   - Processes items in parallel chunks to maximize CPU utilization
//
// Future optimization: Could implement SIMD vectorization of FNV-1a
// using SSE2/AVX2 to hash 2-4 items simultaneously with vector operations.
func (p *BatchProcessor) ProcessBatch(items [][]byte, fingerprintBits, numBuckets uint) []types.HashResult {
	results := make([]types.HashResult, len(items))

	// For small-to-medium batches, sequential processing is faster
	// Goroutine overhead is ~1-2µs per goroutine, which exceeds the benefit for small batches
	if len(items) < 32 {
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

	// For larger batches, process in parallel
	// Use larger chunk sizes to amortize goroutine overhead
	chunkSize := (len(items) + 3) / 4 // Process in 4 chunks
	if chunkSize < 8 {
		chunkSize = 8 // Minimum chunk size to justify goroutine overhead
	}

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
				item := items[i]

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
			done <- struct{}{}
		}(c.start, c.end)
	}

	// Wait for all chunks to complete
	for range chunks {
		<-done
	}

	return results
}
