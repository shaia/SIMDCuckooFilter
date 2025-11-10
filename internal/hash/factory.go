package hash

import (
	stdcrc32 "hash/crc32"

	crc32hash "github.com/shaia/cuckoofilter/internal/hash/crc32"
	fnvhash "github.com/shaia/cuckoofilter/internal/hash/fnv"
	"github.com/shaia/cuckoofilter/internal/hash/xxhash"
	"github.com/shaia/cuckoofilter/internal/simd/cpu"
)

// NewHashFunction creates a hash function based on the strategy and fingerprint bits.
// Uses default configuration without SIMD optimizations.
func NewHashFunction(strategy HashStrategy, fingerprintBits uint) HashInterface {
	return NewHashFunctionWithSIMD(strategy, fingerprintBits, cpu.SIMDNone)
}

// NewHashFunctionWithSIMD creates a hash function with SIMD support.
// The simdType parameter enables SIMD-optimized batch processing when available.
func NewHashFunctionWithSIMD(strategy HashStrategy, fingerprintBits uint, simdType cpu.SIMDType) HashInterface {
	switch strategy {
	case HashStrategyCRC32:
		crcTable := stdcrc32.MakeTable(stdcrc32.Castagnoli)
		var crcBatchProcessor *crc32hash.BatchProcessor
		if simdType != cpu.SIMDNone {
			crcBatchProcessor = crc32hash.NewBatchProcessor(crcTable, simdType)
		}
		return crc32hash.NewCRC32Hash(crcTable, fingerprintBits, crcBatchProcessor)
	case HashStrategyXXHash:
		var xxhashBatchProcessor *xxhash.BatchHashProcessor
		if simdType != cpu.SIMDNone {
			xxhashBatchProcessor = xxhash.NewBatchHashProcessor(simdType)
		}
		return xxhash.NewXXHash(fingerprintBits, xxhashBatchProcessor)
	default: // HashStrategyFNV
		var fnvBatchProcessor *fnvhash.BatchProcessor
		if simdType != cpu.SIMDNone {
			fnvBatchProcessor = fnvhash.NewBatchProcessor(simdType)
		}
		return fnvhash.NewFNVHash(fingerprintBits, fnvBatchProcessor)
	}
}
