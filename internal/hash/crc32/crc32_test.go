package crc32hash

import (
	"fmt"
	"hash/crc32"
	"testing"
)

// TestCRC32HashConsistency verifies that hash function produces consistent results
func TestCRC32HashConsistency(t *testing.T) {
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

	table := crc32.MakeTable(crc32.Castagnoli)
	h := NewCRC32Hash(table, 8, nil)
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

// TestCRC32HashCorrectness verifies CRC32C hash against stdlib implementation
func TestCRC32HashCorrectness(t *testing.T) {
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
		{"zeros", make([]byte, 100)},
	}

	table := crc32.MakeTable(crc32.Castagnoli)
	h := NewCRC32Hash(table, 8, nil)

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Compute expected hash using stdlib
			expectedHash := crc32.Checksum(tc.data, table)

			// Compute using our implementation
			i1, _, fp := h.GetIndices(tc.data, 1000000)

			// Verify the indices are derived from the same hash
			if fp == 0 {
				t.Errorf("%s: fingerprint is zero (should never happen)", tc.name)
			}

			// Verify index is in valid range
			if i1 >= 1000000 {
				t.Errorf("%s: index out of range: %d", tc.name, i1)
			}

			// Verify the hash value matches stdlib
			if uint32(expectedHash)%uint32(1000000) != uint32(i1) {
				t.Errorf("%s: index doesn't match expected hash: got %d, expected %d",
					tc.name, i1, uint32(expectedHash)%1000000)
			}
		})
	}
}

// TestCRC32FingerprintBits tests different fingerprint bit sizes
func TestCRC32FingerprintBits(t *testing.T) {
	// Fingerprints are stored as bytes, so only 1-8 bits are supported
	testCases := []uint{4, 8}
	testData := []byte("test data")
	table := crc32.MakeTable(crc32.Castagnoli)

	for _, bits := range testCases {
		t.Run(fmt.Sprintf("%dbits", bits), func(t *testing.T) {
			h := NewCRC32Hash(table, bits, nil)
			_, _, fp := h.GetIndices(testData, 1024)

			// Verify fingerprint is never zero
			if fp == 0 {
				t.Errorf("%d-bit fingerprint is zero", bits)
			}

			// Verify fingerprint is within valid range
			maxFp := uint16((1 << bits) - 1)
			if fp > maxFp {
				t.Errorf("%d-bit fingerprint %d exceeds max %d", bits, fp, maxFp)
			}
		})
	}
}

// TestCRC32ZeroFingerprint ensures fingerprint is never zero
func TestCRC32ZeroFingerprint(t *testing.T) {
	table := crc32.MakeTable(crc32.Castagnoli)
	h := NewCRC32Hash(table, 8, nil)

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

// TestCRC32GetAltIndex verifies alternative index calculation
func TestCRC32GetAltIndex(t *testing.T) {
	table := crc32.MakeTable(crc32.Castagnoli)
	h := NewCRC32Hash(table, 8, nil)
	numBuckets := uint(1024)

	testCases := []struct {
		index uint
		fp    uint16
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

// TestCRC32BatchProcessing tests batch processing functionality
func TestCRC32BatchProcessing(t *testing.T) {
	items := [][]byte{
		[]byte("item1"),
		[]byte("item2"),
		[]byte("item3"),
		[]byte("item4"),
		[]byte("item5"),
	}

	table := crc32.MakeTable(crc32.Castagnoli)
	h := NewCRC32Hash(table, 8, nil)
	numBuckets := uint(1024)

	// Get individual results
	individualResults := make([]struct {
		i1, i2 uint
		fp     uint16
	}, len(items))
	for i, item := range items {
		i1, i2, fp := h.GetIndices(item, numBuckets)
		individualResults[i] = struct {
			i1, i2 uint
			fp     uint16
		}{i1, i2, fp}
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

// TestCRC32EmptyInput tests behavior with empty input
func TestCRC32EmptyInput(t *testing.T) {
	table := crc32.MakeTable(crc32.Castagnoli)
	h := NewCRC32Hash(table, 8, nil)
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

// TestCRC32LargeInput tests behavior with large input
func TestCRC32LargeInput(t *testing.T) {
	table := crc32.MakeTable(crc32.Castagnoli)
	h := NewCRC32Hash(table, 8, nil)
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

// TestCRC32VariousBucketSizes tests with different bucket counts
func TestCRC32VariousBucketSizes(t *testing.T) {
	table := crc32.MakeTable(crc32.Castagnoli)
	h := NewCRC32Hash(table, 8, nil)
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

// TestCRC32DifferentInputs verifies that different inputs produce different hashes
func TestCRC32DifferentInputs(t *testing.T) {
	table := crc32.MakeTable(crc32.Castagnoli)
	h := NewCRC32Hash(table, 8, nil)
	numBuckets := uint(1024)

	inputs := [][]byte{
		[]byte("test1"),
		[]byte("test2"),
		[]byte("test3"),
		[]byte("different"),
		[]byte("unique"),
	}

	results := make(map[uint]bool)
	fpResults := make(map[uint16]bool)

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

// TestCRC32Tables tests with different CRC32 polynomial tables
func TestCRC32Tables(t *testing.T) {
	testData := []byte("test data")
	numBuckets := uint(1024)

	tables := []struct {
		name  string
		table *crc32.Table
	}{
		{"Castagnoli", crc32.MakeTable(crc32.Castagnoli)},
		{"IEEE", crc32.MakeTable(crc32.IEEE)},
		{"Koopman", crc32.MakeTable(crc32.Koopman)},
	}

	for _, tc := range tables {
		t.Run(tc.name, func(t *testing.T) {
			h := NewCRC32Hash(tc.table, 8, nil)
			i1, i2, fp := h.GetIndices(testData, numBuckets)

			if i1 >= numBuckets {
				t.Errorf("i1 out of range: %d", i1)
			}
			if i2 >= numBuckets {
				t.Errorf("i2 out of range: %d", i2)
			}
			if fp == 0 {
				t.Error("fingerprint is zero")
			}
		})
	}
}

// TestCRC32FingerprintFunction tests the fingerprint extraction
func TestCRC32FingerprintFunction(t *testing.T) {
	testCases := []struct {
		hashVal uint64
		bits    uint
		expect  uint16
	}{
		{0x00, 8, 1},      // Zero should become 1
		{0xFF, 8, 0xFF},   // All bits set
		{0x80, 8, 0x80},   // High bit set
		{0x01, 8, 0x01},   // Low bit set
		{0x100, 8, 1},     // Overflow, low bits zero -> 1
		{0x1234, 8, 0x34}, // Extract low 8 bits
		{0x00, 4, 1},      // Zero with 4 bits -> 1
		{0x0F, 4, 0x0F},   // Max 4-bit value
	}

	for _, tc := range testCases {
		result := fingerprint(tc.hashVal, tc.bits)
		if result != tc.expect {
			t.Errorf("fingerprint(0x%x, %d) = 0x%x, want 0x%x",
				tc.hashVal, tc.bits, result, tc.expect)
		}
	}
}

// BenchmarkCRC32Hash benchmarks single hash operation
func BenchmarkCRC32Hash(b *testing.B) {
	table := crc32.MakeTable(crc32.Castagnoli)
	h := NewCRC32Hash(table, 8, nil)
	data := []byte("benchmark test data with reasonable length")
	numBuckets := uint(1024)

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _, _ = h.GetIndices(data, numBuckets)
	}
}

// BenchmarkCRC32BatchHash benchmarks batch hash operations
func BenchmarkCRC32BatchHash(b *testing.B) {
	items := make([][]byte, 32)
	for i := range items {
		items[i] = []byte("benchmark test data")
	}

	table := crc32.MakeTable(crc32.Castagnoli)
	h := NewCRC32Hash(table, 8, nil)
	numBuckets := uint(1024)

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = h.GetIndicesBatch(items, numBuckets)
	}
}

// BenchmarkCRC32vsStdlib compares our implementation with stdlib
func BenchmarkCRC32vsStdlib(b *testing.B) {
	data := []byte("benchmark test data")
	table := crc32.MakeTable(crc32.Castagnoli)

	b.Run("Stdlib", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			_ = crc32.Checksum(data, table)
		}
	})

	b.Run("OurImplementation", func(b *testing.B) {
		h := NewCRC32Hash(table, 8, nil)
		numBuckets := uint(1024)
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			_, _, _ = h.GetIndices(data, numBuckets)
		}
	})
}
