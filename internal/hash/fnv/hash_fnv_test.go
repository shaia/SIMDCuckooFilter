package fnv

import (
	"fmt"
	"hash/fnv"
	"testing"
)

// TestFNVHashConsistency verifies that hash function produces consistent results
func TestFNVHashConsistency(t *testing.T) {
	testCases := [][]byte{
		[]byte(""),
		[]byte("a"),
		[]byte("hello"),
		[]byte("world"),
		[]byte("the quick brown fox jumps over the lazy dog"),
		[]byte("0123456789abcdef"),
		make([]byte, 100),
		make([]byte, 1000),
	}

	h := NewFNVHash(8, nil)
	numBuckets := uint(1024)

	for i, tc := range testCases {
		// Hash the same input twice
		i1_1, i2_1, fp1 := h.GetIndices(tc, numBuckets)
		i1_2, i2_2, fp2 := h.GetIndices(tc, numBuckets)

		if i1_1 != i1_2 || i2_1 != i2_2 || fp1 != fp2 {
			t.Errorf("Test case %d: inconsistent results: (%d,%d,%x) vs (%d,%d,%x)",
				i, i1_1, i2_1, fp1, i1_2, i2_2, fp2)
		}
	}
}

// TestFNVHashCorrectness verifies FNV hash against stdlib implementation
func TestFNVHashCorrectness(t *testing.T) {
	testCases := []struct {
		name string
		data []byte
	}{
		{"empty", []byte("")},
		{"single byte", []byte("a")},
		{"two bytes", []byte("ab")},
		{"short", []byte("hello")},
		{"medium", []byte("the quick brown fox jumps over the lazy dog")},
		{"numbers", []byte("0123456789")},
		{"special chars", []byte("!@#$%^&*()")},
		{"unicode", []byte("héllo wörld")},
		{"long", make([]byte, 1000)},
	}

	h := NewFNVHash(8, nil)

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Compute expected hash using stdlib
			hasher := fnv.New64a()
			hasher.Write(tc.data)
			expectedHash := hasher.Sum64()

			// Compute using our implementation
			i1, _, fp := h.GetIndices(tc.data, 1000000)

			// Verify the indices are derived from the same hash
			// (We can't directly compare internal hash values, but we can verify consistency)
			if fp == 0 {
				t.Errorf("%s: fingerprint is zero (should never happen)", tc.name)
			}

			// Verify index is in valid range
			if i1 >= 1000000 {
				t.Errorf("%s: index out of range: %d", tc.name, i1)
			}

			// Verify the hash value matches stdlib (indirect verification)
			if expectedHash%1000000 != uint64(i1) {
				t.Errorf("%s: index doesn't match expected hash: got %d, expected %d",
					tc.name, i1, expectedHash%1000000)
			}
		})
	}
}

// TestFNVFingerprintBits tests different fingerprint bit sizes
func TestFNVFingerprintBits(t *testing.T) {
	// Fingerprints are stored as bytes, so only 1-8 bits are supported
	testCases := []uint{4, 8}
	testData := []byte("test data")

	for _, bits := range testCases {
		t.Run(fmt.Sprintf("%dbits", bits), func(t *testing.T) {
			h := NewFNVHash(bits, nil)
			_, _, fp := h.GetIndices(testData, 1024)

			// Verify fingerprint is never zero
			if fp == 0 {
				t.Errorf("%d-bit fingerprint is zero", bits)
			}

			// Verify fingerprint is within valid range
			maxFp := byte((1 << bits) - 1)
			if fp > maxFp {
				t.Errorf("%d-bit fingerprint %d exceeds max %d", bits, fp, maxFp)
			}
		})
	}
}

// TestFNVZeroFingerprint ensures fingerprint is never zero
func TestFNVZeroFingerprint(t *testing.T) {
	h := NewFNVHash(8, nil)

	// Test with various inputs to ensure fingerprint is never 0
	testInputs := [][]byte{
		[]byte(""),
		[]byte("a"),
		[]byte("test"),
		[]byte("hello world"),
		make([]byte, 100),
		// Add some inputs that might hash to values with low bits zero
		[]byte("\x00"),
		[]byte("\x00\x00"),
		[]byte("\x00\x00\x00\x00"),
	}

	for i, input := range testInputs {
		_, _, fp := h.GetIndices(input, 1024)
		if fp == 0 {
			t.Errorf("Test case %d: fingerprint is zero for input %v", i, input)
		}
	}
}

// TestFNVGetAltIndex verifies alternative index calculation
func TestFNVGetAltIndex(t *testing.T) {
	h := NewFNVHash(8, nil)
	numBuckets := uint(1024)

	testCases := []struct {
		index uint
		fp    byte
	}{
		{0, 1},
		{100, 255},
		{512, 128},
		{1023, 64},
	}

	for _, tc := range testCases {
		altIdx := h.GetAltIndex(tc.index, tc.fp, numBuckets)

		// Verify alternative index is in valid range
		if altIdx >= numBuckets {
			t.Errorf("Alternative index out of range: %d >= %d", altIdx, numBuckets)
		}

		// Verify GetAltIndex is deterministic
		altIdx2 := h.GetAltIndex(tc.index, tc.fp, numBuckets)
		if altIdx != altIdx2 {
			t.Errorf("GetAltIndex is not deterministic: %d vs %d", altIdx, altIdx2)
		}

		// Verify reversibility: GetAltIndex(GetAltIndex(i, fp)) should equal i
		originalIdx := h.GetAltIndex(altIdx, tc.fp, numBuckets)
		if originalIdx != tc.index {
			t.Errorf("GetAltIndex is not reversible: %d -> %d -> %d", tc.index, altIdx, originalIdx)
		}
	}
}

// TestFNVBatchProcessing tests batch processing functionality
func TestFNVBatchProcessing(t *testing.T) {
	items := [][]byte{
		[]byte("item1"),
		[]byte("item2"),
		[]byte("item3"),
		[]byte("item4"),
		[]byte("item5"),
	}

	h := NewFNVHash(8, nil)
	numBuckets := uint(1024)

	// Get individual results
	individualResults := make([]struct{ i1, i2 uint; fp byte }, len(items))
	for i, item := range items {
		i1, i2, fp := h.GetIndices(item, numBuckets)
		individualResults[i] = struct{ i1, i2 uint; fp byte }{i1, i2, fp}
	}

	// Get batch results
	batchResults := h.GetIndicesBatch(items, numBuckets)

	// Compare results
	if len(batchResults) != len(items) {
		t.Fatalf("Batch result count mismatch: got %d, expected %d", len(batchResults), len(items))
	}

	for i := range items {
		if batchResults[i].I1 != individualResults[i].i1 ||
			batchResults[i].I2 != individualResults[i].i2 ||
			batchResults[i].Fp != individualResults[i].fp {
			t.Errorf("Item %d mismatch: batch=%+v, individual=(%d,%d,%x)",
				i, batchResults[i], individualResults[i].i1, individualResults[i].i2, individualResults[i].fp)
		}
	}
}

// TestFNVEmptyInput tests behavior with empty input
func TestFNVEmptyInput(t *testing.T) {
	h := NewFNVHash(8, nil)
	numBuckets := uint(1024)

	i1, i2, fp := h.GetIndices([]byte(""), numBuckets)

	// Verify fingerprint is not zero
	if fp == 0 {
		t.Error("Empty input produced zero fingerprint")
	}

	// Verify indices are in valid range
	if i1 >= numBuckets {
		t.Errorf("i1 out of range: %d", i1)
	}
	if i2 >= numBuckets {
		t.Errorf("i2 out of range: %d", i2)
	}
}

// TestFNVLargeInput tests behavior with large input
func TestFNVLargeInput(t *testing.T) {
	h := NewFNVHash(8, nil)
	numBuckets := uint(1024)

	// Create a large input (10KB)
	largeInput := make([]byte, 10*1024)
	for i := range largeInput {
		largeInput[i] = byte(i % 256)
	}

	i1, i2, fp := h.GetIndices(largeInput, numBuckets)

	// Verify fingerprint is not zero
	if fp == 0 {
		t.Error("Large input produced zero fingerprint")
	}

	// Verify indices are in valid range
	if i1 >= numBuckets {
		t.Errorf("i1 out of range: %d", i1)
	}
	if i2 >= numBuckets {
		t.Errorf("i2 out of range: %d", i2)
	}

	// Verify consistency with multiple calls
	i1_2, i2_2, fp2 := h.GetIndices(largeInput, numBuckets)
	if i1 != i1_2 || i2 != i2_2 || fp != fp2 {
		t.Error("Large input hashing is not consistent")
	}
}

// TestFNVVariousBucketSizes tests with different bucket counts
func TestFNVVariousBucketSizes(t *testing.T) {
	h := NewFNVHash(8, nil)
	testData := []byte("test data")

	bucketSizes := []uint{1, 2, 4, 8, 16, 32, 64, 128, 256, 512, 1024, 2048, 4096}

	for _, numBuckets := range bucketSizes {
		i1, i2, fp := h.GetIndices(testData, numBuckets)

		if i1 >= numBuckets {
			t.Errorf("i1=%d exceeds numBuckets=%d", i1, numBuckets)
		}
		if i2 >= numBuckets {
			t.Errorf("i2=%d exceeds numBuckets=%d", i2, numBuckets)
		}
		if fp == 0 {
			t.Errorf("fingerprint is zero for numBuckets=%d", numBuckets)
		}
	}
}

// TestFNVDifferentInputs verifies that different inputs produce different hashes
func TestFNVDifferentInputs(t *testing.T) {
	h := NewFNVHash(8, nil)
	numBuckets := uint(1024)

	inputs := [][]byte{
		[]byte("test1"),
		[]byte("test2"),
		[]byte("test3"),
		[]byte("different"),
		[]byte("unique"),
	}

	results := make(map[uint]bool)
	fpResults := make(map[byte]bool)

	for _, input := range inputs {
		i1, _, fp := h.GetIndices(input, numBuckets)
		results[i1] = true
		fpResults[fp] = true
	}

	// While hash collisions are possible, with 1024 buckets and 5 different inputs,
	// we expect at least some variation
	if len(results) < 2 {
		t.Error("All inputs hashed to similar indices - possible hash function issue")
	}
}

// BenchmarkFNVHash benchmarks single hash operation
func BenchmarkFNVHash(b *testing.B) {
	h := NewFNVHash(8, nil)
	data := []byte("benchmark test data with reasonable length")
	numBuckets := uint(1024)

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _, _ = h.GetIndices(data, numBuckets)
	}
}

// BenchmarkFNVBatchHash benchmarks batch hash operations
func BenchmarkFNVBatchHash(b *testing.B) {
	items := make([][]byte, 32)
	for i := range items {
		items[i] = []byte("benchmark test data")
	}

	h := NewFNVHash(8, nil)
	numBuckets := uint(1024)

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = h.GetIndicesBatch(items, numBuckets)
	}
}
