package hash

import (
	"testing"

	crc32hash "github.com/shaia/simdcuckoofilter/internal/hash/crc32"
	fnvhash "github.com/shaia/simdcuckoofilter/internal/hash/fnv"
	"github.com/shaia/simdcuckoofilter/internal/hash/xxhash"
)

// TestNewHashFunction verifies that NewHashFunction creates the correct hash type
// for each strategy without SIMD support.
func TestNewHashFunction(t *testing.T) {
	testCases := []struct {
		name     string
		strategy HashStrategy
		wantType string
	}{
		{"CRC32", HashStrategyCRC32, "*crc32.CRC32Hash"},
		{"XXHash", HashStrategyXXHash, "*xxhash.XXHash"},
		{"FNV", HashStrategyFNV, "*fnvhash.FNVHash"},
		{"Default", HashStrategy(99), "*fnvhash.FNVHash"}, // Unknown defaults to FNV
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			h := NewHashFunction(tc.strategy, 8)
			if h == nil {
				t.Fatal("NewHashFunction returned nil")
			}

			// Verify the correct type was created
			switch tc.strategy {
			case HashStrategyCRC32:
				if _, ok := h.(*crc32hash.CRC32Hash); !ok {
					t.Errorf("Expected *crc32.CRC32Hash, got %T", h)
				}
			case HashStrategyXXHash:
				if _, ok := h.(*xxhash.XXHash); !ok {
					t.Errorf("Expected *xxhash.XXHash, got %T", h)
				}
			default: // FNV or unknown
				if _, ok := h.(*fnvhash.FNVHash); !ok {
					t.Errorf("Expected *fnvhash.FNVHash, got %T", h)
				}
			}
		})
	}
}

// TestNewHashFunctionWithSIMD verifies that NewHashFunction creates
// the correct hash type with automatic SIMD support.
func TestNewHashFunctionWithSIMD(t *testing.T) {
	testCases := []struct {
		name     string
		strategy HashStrategy
	}{
		{"CRC32", HashStrategyCRC32},
		{"XXHash", HashStrategyXXHash},
		{"FNV", HashStrategyFNV},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			h := NewHashFunction(tc.strategy, 8)
			if h == nil {
				t.Fatal("NewHashFunction returned nil")
			}

			// Verify the hash can perform basic operations
			item := []byte("test")
			numBuckets := uint(1024)
			i1, i2, fp := h.GetIndices(item, numBuckets)

			// Verify indices are within bounds
			if i1 >= numBuckets {
				t.Errorf("i1 (%d) >= numBuckets (%d)", i1, numBuckets)
			}
			if i2 >= numBuckets {
				t.Errorf("i2 (%d) >= numBuckets (%d)", i2, numBuckets)
			}

			// Verify fingerprint is non-zero
			if fp == 0 {
				t.Error("fingerprint is zero (should be non-zero)")
			}

			// Verify alternative index calculation is symmetric
			i2Calculated := h.GetAltIndex(i1, fp, numBuckets)
			if i2Calculated != i2 {
				t.Errorf("GetAltIndex(%d, %d, %d) = %d, want %d", i1, fp, numBuckets, i2Calculated, i2)
			}

			// Verify XOR symmetry: GetAltIndex(GetAltIndex(i1, fp, n), fp, n) == i1
			i1Back := h.GetAltIndex(i2, fp, numBuckets)
			if i1Back != i1 {
				t.Errorf("XOR symmetry broken: GetAltIndex(GetAltIndex(%d, %d, %d), %d, %d) = %d, want %d",
					i1, fp, numBuckets, fp, numBuckets, i1Back, i1)
			}
		})
	}
}

// TestFactoryBatchProcessing verifies that batch processing works correctly
// for hashes created by the factory.
func TestFactoryBatchProcessing(t *testing.T) {
	testCases := []struct {
		name     string
		strategy HashStrategy
	}{
		{"CRC32", HashStrategyCRC32},
		{"XXHash", HashStrategyXXHash},
		{"FNV", HashStrategyFNV},
	}

	items := [][]byte{
		[]byte("test1"),
		[]byte("test2"),
		[]byte("test3"),
		[]byte("test4"),
	}
	numBuckets := uint(1024)

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Create hash with automatic SIMD support
			h := NewHashFunction(tc.strategy, 8)

			// Get individual results
			individualResults := make([]struct{ i1, i2 uint; fp byte }, len(items))
			for i, item := range items {
				i1, i2, fp := h.GetIndices(item, numBuckets)
				individualResults[i] = struct{ i1, i2 uint; fp byte }{i1, i2, fp}
			}

			// Get batch results
			batchResults := h.GetIndicesBatch(items, numBuckets)

			// Verify batch results match individual results
			if len(batchResults) != len(items) {
				t.Fatalf("Expected %d batch results, got %d", len(items), len(batchResults))
			}

			for i := range items {
				if batchResults[i].I1 != individualResults[i].i1 {
					t.Errorf("Item %d: batch i1=%d, individual i1=%d", i, batchResults[i].I1, individualResults[i].i1)
				}
				if batchResults[i].I2 != individualResults[i].i2 {
					t.Errorf("Item %d: batch i2=%d, individual i2=%d", i, batchResults[i].I2, individualResults[i].i2)
				}
				if batchResults[i].Fp != individualResults[i].fp {
					t.Errorf("Item %d: batch fp=%d, individual fp=%d", i, batchResults[i].Fp, individualResults[i].fp)
				}
			}
		})
	}
}

// TestFactoryFingerprintBits verifies that the factory respects the
// fingerprintBits parameter for all hash strategies.
func TestFactoryFingerprintBits(t *testing.T) {
	testCases := []struct {
		name     string
		strategy HashStrategy
		bits     uint
	}{
		{"CRC32_8bits", HashStrategyCRC32, 8},
		{"CRC32_4bits", HashStrategyCRC32, 4},
		{"XXHash_8bits", HashStrategyXXHash, 8},
		{"XXHash_4bits", HashStrategyXXHash, 4},
		{"FNV_8bits", HashStrategyFNV, 8},
		{"FNV_4bits", HashStrategyFNV, 4},
	}

	item := []byte("test")
	numBuckets := uint(1024)

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			h := NewHashFunction(tc.strategy, tc.bits)
			_, _, fp := h.GetIndices(item, numBuckets)

			// Verify fingerprint fits within specified bits
			maxValue := byte((1 << tc.bits) - 1)
			if fp > maxValue {
				t.Errorf("Fingerprint %d exceeds max value %d for %d bits", fp, maxValue, tc.bits)
			}

			// Verify fingerprint is non-zero
			if fp == 0 {
				t.Error("Fingerprint is zero (should be non-zero)")
			}
		})
	}
}

// TestFactoryConsistency verifies that hashes created by the factory
// produce consistent results across multiple calls.
func TestFactoryConsistency(t *testing.T) {
	strategies := []HashStrategy{
		HashStrategyCRC32,
		HashStrategyXXHash,
		HashStrategyFNV,
	}

	item := []byte("consistency_test")
	numBuckets := uint(1024)

	for _, strategy := range strategies {
		t.Run(strategy.String(), func(t *testing.T) {
			h := NewHashFunction(strategy, 8)

			// Hash the same item multiple times
			results := make([]struct{ i1, i2 uint; fp byte }, 100)
			for i := range results {
				i1, i2, fp := h.GetIndices(item, numBuckets)
				results[i] = struct{ i1, i2 uint; fp byte }{i1, i2, fp}
			}

			// Verify all results are identical
			for i := 1; i < len(results); i++ {
				if results[i] != results[0] {
					t.Errorf("Inconsistent result at iteration %d: got %+v, want %+v",
						i, results[i], results[0])
				}
			}
		})
	}
}
