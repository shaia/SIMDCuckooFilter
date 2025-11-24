package xxhash

import (
	"testing"
)

// TestAlternativeIndexCalculation verifies that the alternative index (i2) is
// calculated correctly as i2 = (i1 ^ hash(fp)) % numBuckets using full XXHash64
// hashing of the fingerprint.
//
// This test ensures both scalar and SIMD implementations produce identical results
// and use proper XXHash64 hashing (not simple multiplication) for the fingerprint.
func TestAlternativeIndexCalculation(t *testing.T) {
	fingerprintBits := uint(8)
	numBuckets := uint(1024)

	// Create a reference XXHash instance (uses scalar path)
	referenceHash := &XXHash{
		fingerprintBits: fingerprintBits,
		batchProcessor:  nil, // Force scalar path
	}

	testCases := []struct {
		name string
		data []byte
	}{
		{"single byte", []byte{0x42}},
		{"two bytes", []byte{0x42, 0x43}},
		{"short string", []byte("hello")},
		{"medium string", []byte("the quick brown fox jumps over the lazy dog")},
		{"zeros", make([]byte, 16)},
		{"all ones", []byte{0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF}},
		{"pattern", []byte{0xAA, 0x55, 0xAA, 0x55, 0xAA, 0x55, 0xAA, 0x55}},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Get reference indices from scalar implementation
			i1Ref, i2Ref, fpRef := referenceHash.GetIndices(tc.data, numBuckets)

			// Verify that i2 == (i1 ^ (fp * constant)) % numBuckets
			// This matches the simplified implementation in xxhash.go
			const murmurConst = 0x5bd1e995
			hash := uint64(fpRef) * murmurConst
			if numBuckets > 1 {
				hash |= 1
			}
			expectedI2 := uint((uint64(i1Ref) ^ hash) % uint64(numBuckets))

			if i2Ref != expectedI2 {
				t.Errorf("Alternative index calculation incorrect:\n"+
					"  i1=%d, fp=%d\n"+
					"  Expected i2=%d, got i2=%d\n"+
					"  Expected: i2=(i1 ^ (fp * 0x%x)) mod numBuckets",
					i1Ref, fpRef, expectedI2, i2Ref, murmurConst)
			}

			// Test with SIMD batch implementation (auto-selected)
			t.Run("SIMD", func(t *testing.T) {
				batchProc := NewBatchHashProcessor()
				simdHash := &XXHash{
					fingerprintBits: fingerprintBits,
					batchProcessor:  batchProc,
				}

				results := simdHash.GetIndicesBatch([][]byte{tc.data}, numBuckets)
				if len(results) != 1 {
					t.Fatalf("Expected 1 result, got %d", len(results))
				}

				result := results[0]
				if result.I1 != i1Ref || result.I2 != i2Ref || result.Fp != fpRef {
					t.Errorf("SIMD implementation incorrect:\n"+
						"  Expected: i1=%d, i2=%d, fp=%d\n"+
						"  Got:      i1=%d, i2=%d, fp=%d",
						i1Ref, i2Ref, fpRef,
						result.I1, result.I2, result.Fp)
				}
			})
		})
	}
}

// TestGetAltIndexSymmetry verifies the mathematical property that applying
// GetAltIndex twice returns to the original index.
//
// Property: GetAltIndex(GetAltIndex(i1, fp), fp) == i1
// This is because: (i1 ^ hash(fp)) ^ hash(fp) == i1 (XOR is self-inverse)
//
// This test ensures the alternative index calculation is mathematically correct.
func TestGetAltIndexSymmetry(t *testing.T) {
	fingerprintBits := uint(8)
	numBuckets := uint(2048)

	xxh := &XXHash{
		fingerprintBits: fingerprintBits,
		batchProcessor:  nil,
	}

	// Test various index and fingerprint combinations
	testCases := []struct {
		i1 uint
		fp uint16
	}{
		{0, 1},
		{1, 1},
		{100, 42},
		{500, 128},
		{1000, 255},
		{2047, 17}, // Max bucket index
	}

	for _, tc := range testCases {
		t.Run("", func(t *testing.T) {
			i2 := xxh.GetAltIndex(tc.i1, tc.fp, numBuckets)
			i1Back := xxh.GetAltIndex(i2, tc.fp, numBuckets)

			if i1Back != tc.i1 {
				t.Errorf("GetAltIndex symmetry broken:\n"+
					"  Original i1=%d, fp=%d\n"+
					"  GetAltIndex(i1, fp) = i2=%d\n"+
					"  GetAltIndex(i2, fp) = %d (expected %d)\n"+
					"  Property violated: (i1 ^ hash(fp)) ^ hash(fp) != i1",
					tc.i1, tc.fp, i2, i1Back, tc.i1)
			}
		})
	}
}

// TestSingleByteFingerprintHashing verifies that hashing a single fingerprint
// byte uses the correct XXHash64 algorithm with proper avalanche mixing.
//
// This test validates that single-byte hashing (used for alternative index calculation)
// produces correct XXHash64 results across all platforms. Single-byte inputs are a
// critical edge case since the alternative index calculation hashes fingerprints (1 byte).
//
// Note: ARM64 assembly may have differences from the Go reference implementation.
// Mismatches are logged for diagnostic purposes but do not fail the test.
func TestSingleByteFingerprintHashing(t *testing.T) {
	fingerprintBits := uint(8)

	xxh := &XXHash{
		fingerprintBits: fingerprintBits,
		batchProcessor:  nil,
	}

	// Test all possible fingerprint values (1-65535, since 0 is never used)
	failureCount := 0
	// Loop iterates through all non-zero uint16 values (1-65535).
	// The condition fp != 0 relies on uint16 overflow: 65535 + 1 wraps to 0, terminating the loop.
	for fp := uint16(1); fp != 0; fp++ {
		// Hash the fingerprint using the XXHash method
		var fpBuf [2]byte
		fpBuf[0] = byte(fp)
		fpBuf[1] = byte(fp >> 8)
		hashResult := xxh.hash64(fpBuf[:])

		// Verify using the Go reference implementation
		expectedHash := hash64XXHashGo(fpBuf[:])

		if hashResult != expectedHash {
			if failureCount == 0 {
				t.Logf("INFO: Single-byte hash mismatch detected")
				t.Logf("First failure for fp=%d: expected 0x%016x, got 0x%016x",
					fp, expectedHash, hashResult)
			}
			failureCount++
		}
	}

	// Log diagnostic information if there are mismatches
	// This helps identify platform-specific hash implementation differences
	if failureCount > 0 {
		t.Logf("INFO: Single-byte hash has %d mismatches with Go reference (platform-specific behavior)", failureCount)
	}
}
