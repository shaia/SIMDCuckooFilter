//go:build amd64
// +build amd64

package crc32hash

import (
	"hash/crc32"

	"github.com/shaia/cuckoofilter/internal/hash/types"
	"github.com/shaia/cuckoofilter/internal/simd/cpu"
)

// BatchProcessor handles optimized batch CRC32 hashing for AMD64.
// Uses SIMD assembly for true parallel processing when available,
// falls back to goroutine-based parallelism otherwise.
type BatchProcessor struct {
	table    *crc32.Table
	simdType cpu.SIMDType
	useSIMD  bool // Whether to use SIMD assembly implementation
}

// NewBatchProcessor creates a new CRC32 batch processor
func NewBatchProcessor(table *crc32.Table, simdType cpu.SIMDType) *BatchProcessor {
	return &BatchProcessor{
		table:    table,
		simdType: simdType,
		useSIMD:  true, // Enable SIMD by default
	}
}

// NewBatchProcessorNoSIMD creates a batch processor without SIMD (for comparison/testing)
func NewBatchProcessorNoSIMD(table *crc32.Table, simdType cpu.SIMDType) *BatchProcessor {
	return &BatchProcessor{
		table:    table,
		simdType: simdType,
		useSIMD:  false, // Disable SIMD
	}
}

// ProcessBatch processes multiple items using optimized CRC32.
//
// Performance characteristics:
//   - Uses Go's hardware-accelerated CRC32 (SSE4.2 instruction)
//   - Processes items in parallel using separate goroutines for better throughput
//   - Automatically balances work across available CPU cores
//
// The Go standard library's crc32 already uses SSE4.2 hardware acceleration,
// so we focus on parallel processing rather than custom SIMD assembly.
func (p *BatchProcessor) ProcessBatch(items [][]byte, fingerprintBits, numBuckets uint) []types.HashResult {
	results := make([]types.HashResult, len(items))

	// Use SIMD assembly if enabled and batch is large enough
	if p.useSIMD && len(items) >= 4 {
		return p.processBatchSIMD(items, fingerprintBits, numBuckets)
	}

	// For small batches, sequential processing is faster (avoids goroutine overhead)
	if len(items) < 8 {
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
	// CRC32 computation is fast, so we use larger chunk sizes to reduce overhead
	chunkSize := (len(items) + 3) / 4 // Process in 4 chunks
	if chunkSize < 4 {
		chunkSize = 4
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
	done := make(chan struct{})
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

// processBatchSIMD uses custom assembly to process batches with SIMD CRC32C instructions
func (p *BatchProcessor) processBatchSIMD(items [][]byte, fingerprintBits, numBuckets uint) []types.HashResult {
	results := make([]types.HashResult, len(items))
	crc32Results := make([]uint32, len(items))

	// Call assembly function to compute CRC32 values
	batchCRC32SIMD(items, crc32Results)

	// Convert CRC32 values to hash results
	for i, hashVal := range crc32Results {
		fp := fingerprint(uint64(hashVal), fingerprintBits)
		i1 := uint(hashVal % uint32(numBuckets))
		fpHash := crc32.Checksum([]byte{fp}, p.table)
		i2 := (uint64(i1) ^ uint64(fpHash)) % uint64(numBuckets)
		results[i] = types.HashResult{I1: i1, I2: uint(i2), Fp: fp}
	}

	return results
}
