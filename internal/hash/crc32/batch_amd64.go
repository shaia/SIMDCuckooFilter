//go:build amd64
// +build amd64

package crc32hash

import (
	"hash/crc32"

	"github.com/shaia/simdcuckoofilter/internal/hash/types"
)

// BatchProcessor handles optimized batch CRC32 hashing for AMD64.
// Uses hardware-accelerated CRC32C and parallel processing.
type BatchProcessor struct {
	table *crc32.Table
}

// NewBatchProcessor creates a new CRC32 batch processor.
// Uses hardware-accelerated CRC32C (SSE4.2) with parallel processing.
func NewBatchProcessor(table *crc32.Table) *BatchProcessor {
	return &BatchProcessor{
		table: table,
	}
}

// ProcessBatch processes multiple items using hardware-accelerated CRC32C.
//
// Performance characteristics:
//   - Uses Go's hardware-accelerated CRC32C (SSE4.2 instruction)
//   - Processes items in parallel using separate goroutines for better throughput
//   - Automatically balances work across available CPU cores
//
// The Go standard library's crc32 already uses SSE4.2 hardware acceleration,
// so we focus on parallel processing rather than custom SIMD assembly.
func (p *BatchProcessor) ProcessBatch(items [][]byte, fingerprintBits, numBuckets uint) []types.HashResult {
	results := make([]types.HashResult, len(items))

	// For small-to-medium batches, sequential processing is faster
	// Goroutine overhead is ~1-2µs per goroutine, which exceeds the benefit for small batches
	if len(items) < 32 {
		for i, item := range items {
			hashVal := crc32.Checksum(item, p.table)
			fp := fingerprint(uint64(hashVal), fingerprintBits)
			i1 := uint(hashVal % uint32(numBuckets))
			fpHash := crc32.Checksum([]byte{fp}, p.table)
			i2 := (uint64(i1) ^ uint64(fpHash)) % uint64(numBuckets)
			results[i] = types.HashResult{I1: i1, I2: uint(i2), Fp: fp}
		}
		return results
	}

	// For larger batches, process in parallel
	// CRC32 computation is fast, so we use larger chunk sizes to amortize goroutine overhead
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
				hashVal := crc32.Checksum(item, p.table)
				fp := fingerprint(uint64(hashVal), fingerprintBits)
				i1 := uint(hashVal % uint32(numBuckets))
				fpHash := crc32.Checksum([]byte{fp}, p.table)
				i2 := (uint64(i1) ^ uint64(fpHash)) % uint64(numBuckets)
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
