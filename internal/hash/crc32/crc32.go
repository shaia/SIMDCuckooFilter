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

func (h *CRC32Hash) GetAltIndex(index uint, fp byte, numBuckets uint) uint {
	// Hash the fingerprint to get alternative index
	// Use stack-allocated buffer for thread safety
	fpBuf := [1]byte{fp}
	fpHash := crc32.Checksum(fpBuf[:], h.Table)
	altIndex := (uint64(index) ^ uint64(fpHash)) % uint64(numBuckets)
	return uint(altIndex)
}

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
