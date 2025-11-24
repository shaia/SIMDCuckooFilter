// Package hash provides hash function implementations for the Cuckoo Filter.
//
// The package supports multiple hash algorithms optimized for different use cases:
//   - XXHash: Fastest, best overall performance with SIMD optimizations
//   - CRC32C: Hardware-accelerated on modern CPUs (SSE4.2)
//   - FNV-1a: Simple, good distribution, pure Go fallback
//
// All hash implementations support batch processing for improved throughput.
//
// # Architecture
//
// The package is organized into subpackages for each hash algorithm:
//   - xxhash: XXHash64 with SIMD optimization (AVX2)
//   - crc32: CRC32C with hardware acceleration
//   - fnv: FNV-1a hash
//
// Each subpackage follows a consistent structure:
//   - Main implementation file (xxhash.go, hash_crc32.go, hash_fnv.go)
//   - Platform-specific batch processors (batch_amd64.go, batch_generic.go)
//   - Assembly optimizations where applicable (*.s files)
package hash

import "github.com/shaia/simdcuckoofilter/internal/hash/types"

// HashResult is an alias to types.HashResult for convenience
type HashResult = types.HashResult

// HashInterface defines the interface for hash functions used in the cuckoo filter.
type HashInterface interface {
	// GetIndices returns the two bucket indices and fingerprint for an item
	GetIndices(item []byte, numBuckets uint) (i1, i2 uint, fp uint16)

	// GetAltIndex calculates the alternative index given one index and fingerprint
	GetAltIndex(index uint, fp uint16, numBuckets uint) uint

	// GetIndicesBatch processes multiple items in batch (SIMD-optimized when available)
	// Returns results in the same order as input items
	GetIndicesBatch(items [][]byte, numBuckets uint) []HashResult
}

// HashStrategy represents different hash function options
type HashStrategy int

const (
	// HashStrategyFNV uses FNV-1a hash (default, good distribution)
	HashStrategyFNV HashStrategy = iota
	// HashStrategyCRC32 uses CRC32C (hardware accelerated on modern CPUs)
	HashStrategyCRC32
	// HashStrategyXXHash uses XXHash (fastest, best performance)
	HashStrategyXXHash
)

// String returns the name of the hash strategy
func (s HashStrategy) String() string {
	switch s {
	case HashStrategyFNV:
		return "FNV-1a"
	case HashStrategyCRC32:
		return "CRC32C"
	case HashStrategyXXHash:
		return "XXHash64"
	default:
		return "Unknown"
	}
}
