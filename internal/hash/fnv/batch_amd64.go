//go:build amd64
// +build amd64

package fnv

import (
	"github.com/shaia/cuckoofilter/internal/hash/types"
)

// BatchProcessor handles optimized batch FNV hashing for AMD64.
// Uses parallel processing to maximize throughput.
type BatchProcessor struct{}

// NewBatchProcessor creates a new FNV batch processor
func NewBatchProcessor() *BatchProcessor {
	return &BatchProcessor{}
}

// ProcessBatch processes multiple items using optimized FNV-1a.
//
// Performance characteristics:
//   - Uses parallel goroutines for better throughput on large batches
//   - FNV-1a is a simple XOR+multiply operation, making it fast
//   - Sequential processing for small batches to avoid goroutine overhead
//
// Future optimization: Could implement SIMD vectorization of FNV-1a
// using AVX2 to hash 4 items simultaneously with vector operations.
func (p *BatchProcessor) ProcessBatch(items [][]byte, fingerprintBits, numBuckets uint) []types.HashResult {
	// For small-to-medium batches, sequential processing is faster
	// Goroutine overhead is ~1-2µs per goroutine, which exceeds the benefit for small batches
	if len(items) < 32 {
		return processSequential(items, fingerprintBits, numBuckets)
	}

	// For larger batches, process in parallel
	return processParallel(items, fingerprintBits, numBuckets)
}
