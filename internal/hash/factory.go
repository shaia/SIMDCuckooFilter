package hash

import (
	stdcrc32 "hash/crc32"

	crc32hash "github.com/shaia/cuckoofilter/internal/hash/crc32"
	fnvhash "github.com/shaia/cuckoofilter/internal/hash/fnv"
	"github.com/shaia/cuckoofilter/internal/hash/xxhash"
)

// NewHashFunction creates a hash function based on the strategy and fingerprint bits.
// Automatically uses the best SIMD implementation available for the platform at compile time.
func NewHashFunction(strategy HashStrategy, fingerprintBits uint) HashInterface {
	switch strategy {
	case HashStrategyCRC32:
		crcTable := stdcrc32.MakeTable(stdcrc32.Castagnoli)
		crcBatchProcessor := crc32hash.NewBatchProcessor(crcTable)
		return crc32hash.NewCRC32Hash(crcTable, fingerprintBits, crcBatchProcessor)
	case HashStrategyXXHash:
		xxhashBatchProcessor := xxhash.NewBatchHashProcessor()
		return xxhash.NewXXHash(fingerprintBits, xxhashBatchProcessor)
	default: // HashStrategyFNV
		fnvBatchProcessor := fnvhash.NewBatchProcessor()
		return fnvhash.NewFNVHash(fingerprintBits, fnvBatchProcessor)
	}
}
