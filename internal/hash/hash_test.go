package hash

import (
	"hash/crc32"
	"testing"

	crc32hash "github.com/shaia/simdcuckoofilter/internal/hash/crc32"
	"github.com/shaia/simdcuckoofilter/internal/hash/fnv"
	"github.com/shaia/simdcuckoofilter/internal/hash/xxhash"
)

// TestHashStrategyString tests the String method for HashStrategy
func TestHashStrategyString(t *testing.T) {
	testCases := []struct {
		strategy HashStrategy
		expected string
	}{
		{HashStrategyFNV, "FNV-1a"},
		{HashStrategyCRC32, "CRC32C"},
		{HashStrategyXXHash, "XXHash64"},
		{HashStrategy(999), "Unknown"},
	}

	for _, tc := range testCases {
		result := tc.strategy.String()
		if result != tc.expected {
			t.Errorf("Strategy %d: got %q, want %q", tc.strategy, result, tc.expected)
		}
	}
}

// TestAllHashImplementations tests all hash implementations against the interface
func TestAllHashImplementations(t *testing.T) {
	testData := []byte("test data for all implementations")
	numBuckets := uint(1024)
	fingerprintBits := uint(8)

	implementations := []struct {
		name string
		hash HashInterface
	}{
		{"FNV", fnv.NewFNVHash(fingerprintBits, nil)},
		{"CRC32", crc32hash.NewCRC32Hash(crc32.MakeTable(crc32.Castagnoli), fingerprintBits, nil)},
		{"XXHash", xxhash.NewXXHash(fingerprintBits, nil)},
	}

	for _, impl := range implementations {
		t.Run(impl.name, func(t *testing.T) {
			// Test GetIndices
			i1, i2, fp := impl.hash.GetIndices(testData, numBuckets)

			// Verify indices are in range
			if i1 >= numBuckets {
				t.Errorf("%s: i1=%d exceeds numBuckets=%d", impl.name, i1, numBuckets)
			}
			if i2 >= numBuckets {
				t.Errorf("%s: i2=%d exceeds numBuckets=%d", impl.name, i2, numBuckets)
			}

			// Verify fingerprint is not zero
			if fp == 0 {
				t.Errorf("%s: fingerprint is zero", impl.name)
			}

			// Test GetAltIndex
			altIdx := impl.hash.GetAltIndex(i1, fp, numBuckets)
			if altIdx >= numBuckets {
				t.Errorf("%s: altIdx=%d exceeds numBuckets=%d", impl.name, altIdx, numBuckets)
			}

			// Verify reversibility
			originalIdx := impl.hash.GetAltIndex(altIdx, fp, numBuckets)
			if originalIdx != i1 {
				t.Errorf("%s: GetAltIndex not reversible: %d -> %d -> %d", impl.name, i1, altIdx, originalIdx)
			}

			// Test consistency
			i1_2, i2_2, fp2 := impl.hash.GetIndices(testData, numBuckets)
			if i1 != i1_2 || i2 != i2_2 || fp != fp2 {
				t.Errorf("%s: inconsistent results", impl.name)
			}
		})
	}
}

// TestHashImplementationsBatch tests batch operations for all implementations
func TestHashImplementationsBatch(t *testing.T) {
	items := [][]byte{
		[]byte("item1"),
		[]byte("item2"),
		[]byte("item3"),
		[]byte("item4"),
	}
	numBuckets := uint(1024)
	fingerprintBits := uint(8)

	implementations := []struct {
		name string
		hash HashInterface
	}{
		{"FNV", fnv.NewFNVHash(fingerprintBits, nil)},
		{"CRC32", crc32hash.NewCRC32Hash(crc32.MakeTable(crc32.Castagnoli), fingerprintBits, nil)},
		{"XXHash", xxhash.NewXXHash(fingerprintBits, nil)},
	}

	for _, impl := range implementations {
		t.Run(impl.name, func(t *testing.T) {
			// Get batch results
			batchResults := impl.hash.GetIndicesBatch(items, numBuckets)

			// Verify result count
			if len(batchResults) != len(items) {
				t.Errorf("%s: got %d results, want %d", impl.name, len(batchResults), len(items))
			}

			// Verify each result matches individual processing
			for i, item := range items {
				i1, i2, fp := impl.hash.GetIndices(item, numBuckets)

				if batchResults[i].I1 != i1 {
					t.Errorf("%s: item %d: batch I1=%d, individual I1=%d", impl.name, i, batchResults[i].I1, i1)
				}
				if batchResults[i].I2 != i2 {
					t.Errorf("%s: item %d: batch I2=%d, individual I2=%d", impl.name, i, batchResults[i].I2, i2)
				}
				if batchResults[i].Fp != fp {
					t.Errorf("%s: item %d: batch Fp=%d, individual Fp=%d", impl.name, i, batchResults[i].Fp, fp)
				}

				// Verify fingerprint is not zero
				if batchResults[i].Fp == 0 {
					t.Errorf("%s: item %d: fingerprint is zero", impl.name, i)
				}

				// Verify indices are in range
				if batchResults[i].I1 >= numBuckets || batchResults[i].I2 >= numBuckets {
					t.Errorf("%s: item %d: indices out of range", impl.name, i)
				}
			}
		})
	}
}

// TestHashDistribution tests that hash implementations distribute values reasonably
func TestHashDistribution(t *testing.T) {
	numBuckets := uint(1024)
	fingerprintBits := uint(8)
	numSamples := 10000

	implementations := []struct {
		name string
		hash HashInterface
	}{
		{"FNV", fnv.NewFNVHash(fingerprintBits, nil)},
		{"CRC32", crc32hash.NewCRC32Hash(crc32.MakeTable(crc32.Castagnoli), fingerprintBits, nil)},
		{"XXHash", xxhash.NewXXHash(fingerprintBits, nil)},
	}

	for _, impl := range implementations {
		t.Run(impl.name, func(t *testing.T) {
			bucketCounts := make([]int, numBuckets)
			fpCounts := make([]int, 256) // 8-bit fingerprints

			// Generate samples
			for i := 0; i < numSamples; i++ {
				// Create unique input
				data := make([]byte, 8)
				for j := 0; j < 8; j++ {
					data[j] = byte((i >> (j * 8)) & 0xFF)
				}

				i1, _, fp := impl.hash.GetIndices(data, numBuckets)
				bucketCounts[i1]++
				fpCounts[fp]++
			}

			// Check that distribution is reasonable
			// Expected count per bucket
			expectedPerBucket := float64(numSamples) / float64(numBuckets)

			// Count buckets that are significantly over/under utilized
			emptyBuckets := 0
			overloadedBuckets := 0

			for _, count := range bucketCounts {
				if count == 0 {
					emptyBuckets++
				}
				if float64(count) > expectedPerBucket*3 {
					overloadedBuckets++
				}
			}

			// With a good hash, very few buckets should be empty or heavily overloaded
			if emptyBuckets > int(numBuckets)/20 { // Allow up to 5% empty
				t.Errorf("%s: too many empty buckets: %d/%d", impl.name, emptyBuckets, numBuckets)
			}
			if overloadedBuckets > int(numBuckets)/20 { // Allow up to 5% overloaded
				t.Errorf("%s: too many overloaded buckets: %d/%d", impl.name, overloadedBuckets, numBuckets)
			}

			// Check fingerprint distribution
			emptyFingerprints := 0
			for i, count := range fpCounts {
				if i == 0 {
					// Fingerprint 0 should never occur
					if count > 0 {
						t.Errorf("%s: fingerprint 0 occurred %d times (should be 0)", impl.name, count)
					}
				} else if count == 0 {
					emptyFingerprints++
				}
			}

			t.Logf("%s distribution: %d empty buckets, %d overloaded buckets, %d empty fingerprints",
				impl.name, emptyBuckets, overloadedBuckets, emptyFingerprints)
		})
	}
}

// TestHashImplementationsConsistency tests that each implementation is internally consistent
func TestHashImplementationsConsistency(t *testing.T) {
	numBuckets := uint(1024)
	fingerprintBits := uint(8)

	testCases := [][]byte{
		[]byte(""),
		[]byte("a"),
		[]byte("hello world"),
		[]byte("The quick brown fox jumps over the lazy dog"),
		make([]byte, 1000),
	}

	implementations := []struct {
		name string
		hash HashInterface
	}{
		{"FNV", fnv.NewFNVHash(fingerprintBits, nil)},
		{"CRC32", crc32hash.NewCRC32Hash(crc32.MakeTable(crc32.Castagnoli), fingerprintBits, nil)},
		{"XXHash", xxhash.NewXXHash(fingerprintBits, nil)},
	}

	for _, impl := range implementations {
		t.Run(impl.name, func(t *testing.T) {
			for i, tc := range testCases {
				// Hash multiple times
				i1_1, i2_1, fp1 := impl.hash.GetIndices(tc, numBuckets)
				i1_2, i2_2, fp2 := impl.hash.GetIndices(tc, numBuckets)
				i1_3, i2_3, fp3 := impl.hash.GetIndices(tc, numBuckets)

				// Verify all results are identical
				if i1_1 != i1_2 || i1_1 != i1_3 {
					t.Errorf("%s: test %d: i1 inconsistent: %d, %d, %d", impl.name, i, i1_1, i1_2, i1_3)
				}
				if i2_1 != i2_2 || i2_1 != i2_3 {
					t.Errorf("%s: test %d: i2 inconsistent: %d, %d, %d", impl.name, i, i2_1, i2_2, i2_3)
				}
				if fp1 != fp2 || fp1 != fp3 {
					t.Errorf("%s: test %d: fp inconsistent: %d, %d, %d", impl.name, i, fp1, fp2, fp3)
				}
			}
		})
	}
}

// TestHashImplementationsDifferentInputs verifies different inputs produce different hashes
func TestHashImplementationsDifferentInputs(t *testing.T) {
	numBuckets := uint(1024)
	fingerprintBits := uint(8)

	// Create clearly different inputs
	inputs := [][]byte{
		[]byte("input1"),
		[]byte("input2"),
		[]byte("input3"),
		[]byte("completely different"),
		[]byte("another unique value"),
	}

	implementations := []struct {
		name string
		hash HashInterface
	}{
		{"FNV", fnv.NewFNVHash(fingerprintBits, nil)},
		{"CRC32", crc32hash.NewCRC32Hash(crc32.MakeTable(crc32.Castagnoli), fingerprintBits, nil)},
		{"XXHash", xxhash.NewXXHash(fingerprintBits, nil)},
	}

	for _, impl := range implementations {
		t.Run(impl.name, func(t *testing.T) {
			seen := make(map[uint]int)
			collisions := 0

			for idx, input := range inputs {
				i1, _, _ := impl.hash.GetIndices(input, numBuckets)

				if prevIdx, exists := seen[i1]; exists {
					collisions++
					t.Logf("%s: collision: inputs %d and %d both hash to %d", impl.name, prevIdx, idx, i1)
				} else {
					seen[i1] = idx
				}
			}

			// Some collisions are acceptable, but all inputs shouldn't collide
			if collisions >= len(inputs)-1 {
				t.Errorf("%s: too many collisions: %d out of %d inputs", impl.name, collisions, len(inputs))
			}
		})
	}
}

// TestHashImplementationsEdgeCases tests edge cases for all implementations
func TestHashImplementationsEdgeCases(t *testing.T) {
	fingerprintBits := uint(8)

	implementations := []struct {
		name string
		hash HashInterface
	}{
		{"FNV", fnv.NewFNVHash(fingerprintBits, nil)},
		{"CRC32", crc32hash.NewCRC32Hash(crc32.MakeTable(crc32.Castagnoli), fingerprintBits, nil)},
		{"XXHash", xxhash.NewXXHash(fingerprintBits, nil)},
	}

	edgeCases := []struct {
		name       string
		data       []byte
		numBuckets uint
	}{
		{"empty data, 1 bucket", []byte(""), 1},
		{"empty data, many buckets", []byte(""), 1024},
		{"1 byte, 1 bucket", []byte("x"), 1},
		{"1 bucket", []byte("test"), 1},
		{"2 buckets", []byte("test"), 2},
		{"large power of 2 buckets", []byte("test"), 65536},
		{"non-power of 2 buckets", []byte("test"), 1000},
		{"large data", make([]byte, 10000), 1024},
	}

	for _, impl := range implementations {
		for _, tc := range edgeCases {
			t.Run(impl.name+"/"+tc.name, func(t *testing.T) {
				i1, i2, fp := impl.hash.GetIndices(tc.data, tc.numBuckets)

				if i1 >= tc.numBuckets {
					t.Errorf("i1=%d exceeds numBuckets=%d", i1, tc.numBuckets)
				}
				if i2 >= tc.numBuckets {
					t.Errorf("i2=%d exceeds numBuckets=%d", i2, tc.numBuckets)
				}
				if fp == 0 {
					t.Error("fingerprint is zero")
				}
			})
		}
	}
}

// BenchmarkAllHashImplementations benchmarks all hash implementations
func BenchmarkAllHashImplementations(b *testing.B) {
	data := []byte("benchmark test data with reasonable length for comparison")
	numBuckets := uint(1024)
	fingerprintBits := uint(8)

	implementations := []struct {
		name string
		hash HashInterface
	}{
		{"FNV", fnv.NewFNVHash(fingerprintBits, nil)},
		{"CRC32", crc32hash.NewCRC32Hash(crc32.MakeTable(crc32.Castagnoli), fingerprintBits, nil)},
		{"XXHash", xxhash.NewXXHash(fingerprintBits, nil)},
	}

	for _, impl := range implementations {
		b.Run(impl.name, func(b *testing.B) {
			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				_, _, _ = impl.hash.GetIndices(data, numBuckets)
			}
		})
	}
}

// BenchmarkAllHashImplementationsBatch benchmarks batch operations
func BenchmarkAllHashImplementationsBatch(b *testing.B) {
	items := make([][]byte, 32)
	for i := range items {
		items[i] = []byte("benchmark test data item")
	}
	numBuckets := uint(1024)
	fingerprintBits := uint(8)

	implementations := []struct {
		name string
		hash HashInterface
	}{
		{"FNV", fnv.NewFNVHash(fingerprintBits, nil)},
		{"CRC32", crc32hash.NewCRC32Hash(crc32.MakeTable(crc32.Castagnoli), fingerprintBits, nil)},
		{"XXHash", xxhash.NewXXHash(fingerprintBits, nil)},
	}

	for _, impl := range implementations {
		b.Run(impl.name, func(b *testing.B) {
			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				_ = impl.hash.GetIndicesBatch(items, numBuckets)
			}
		})
	}
}
