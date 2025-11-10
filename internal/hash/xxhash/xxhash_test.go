//go:build amd64 || arm64
// +build amd64 arm64

package xxhash

import (
	"testing"
)

// TestXXHashAssemblyVsGo verifies assembly implementation matches Go reference
func TestXXHashAssemblyVsGo(t *testing.T) {
	testCases := []struct {
		name string
		data []byte
	}{
		{"empty", []byte("")},
		{"single byte", []byte("a")},
		{"two bytes", []byte("ab")},
		{"three bytes", []byte("abc")},
		{"four bytes", []byte("abcd")},
		{"five bytes", []byte("abcde")},
		{"six bytes", []byte("abcdef")},
		{"seven bytes", []byte("abcdefg")},
		{"eight bytes", []byte("abcdefgh")},
		{"nine bytes", []byte("abcdefghi")},
		{"ten bytes", []byte("abcdefghij")},
		{"medium", []byte("the quick brown fox jumps over the lazy dog")},
		{"zeros", []byte{0, 0, 0, 0, 0}},
		{"high bytes", []byte{255, 254, 253, 252, 251}},
		{"mixed", []byte{0, 1, 2, 255, 254, 253}},
		{"long", make([]byte, 100)},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Assembly implementation
			asmResult := hash64XXHashInternal(tc.data)

			// Go reference implementation
			goResult := hash64XXHashGo(tc.data)

			if asmResult != goResult {
				t.Errorf("Hash mismatch for %q:\n  Assembly: %016x\n  Go:       %016x",
					tc.name, asmResult, goResult)
			}
		})
	}
}

// TestXXHashFinalByteProcessing specifically tests the bug fix for final byte processing
// The bug was: assembly did `hash ^= byte; hash *= prime64_5` instead of `hash ^= byte * prime64_5`
func TestXXHashFinalByteProcessing(t *testing.T) {
	// These test cases are specifically designed to trigger the final byte loop
	// by having lengths that are not multiples of 8
	testCases := [][]byte{
		[]byte("1"),        // 1 byte - triggers short_byte_loop
		[]byte("12"),       // 2 bytes - triggers short_byte_loop
		[]byte("123"),      // 3 bytes - triggers short_byte_loop
		[]byte("1234567"),  // 7 bytes - triggers short_byte_loop
		[]byte("12345678"), // 8 bytes - uses chunk_loop only
		[]byte("123456789"), // 9 bytes - chunk_loop + 1 final byte
		[]byte("1234567890"), // 10 bytes - chunk_loop + 2 final bytes
		[]byte("12345678901234567"), // 17 bytes - 2 chunks + 1 final byte
		[]byte("123456789012345678901234567"), // 27 bytes - 3 chunks + 3 final bytes
	}

	for _, data := range testCases {
		asmResult := hash64XXHashInternal(data)
		goResult := hash64XXHashGo(data)

		if asmResult != goResult {
			t.Errorf("Hash mismatch for %d-byte input %q:\n  Assembly: %016x\n  Go:       %016x",
				len(data), string(data), asmResult, goResult)
			t.Logf("Input bytes: %v", data)
		}
	}
}

// TestFingerprintZeroHandling verifies that fingerprints are never zero.
// A zero fingerprint indicates an empty slot in the cuckoo filter.
//
// This test ensures that the fingerprint function always returns non-zero values,
// even when the hash value would naturally produce a zero fingerprint after masking.
func TestFingerprintZeroHandling(t *testing.T) {
	fingerprintBits := uint(8)

	xxh := &XXHash{
		fingerprintBits: fingerprintBits,
		batchProcessor:  nil,
	}

	// Test many different inputs to ensure we never get fp=0
	testInputs := make([][]byte, 1000)
	for i := 0; i < 1000; i++ {
		data := make([]byte, 1+i%32)
		for j := range data {
			data[j] = byte(i + j)
		}
		testInputs[i] = data
	}

	for i, input := range testInputs {
		_, _, fp := xxh.GetIndices(input, 1024)
		if fp == 0 {
			t.Errorf("Input %d produced zero fingerprint (forbidden): %v", i, input)
		}
	}

	// Test the edge case where hash would naturally produce fp=0
	// The fingerprint function should convert it to 1
	hashValZeroFp := uint64(0) // This would give fp=0 without the check
	fp := fingerprint(hashValZeroFp, fingerprintBits)
	if fp != 1 {
		t.Errorf("Zero fingerprint not converted to 1: got %d", fp)
	}

	// Test hash values that would produce 0 after masking
	for bits := uint(4); bits <= 8; bits++ {
		mask := (uint64(1) << bits) - 1
		testValues := []uint64{
			0,                    // Direct zero
			^mask,                // All high bits set, low bits zero
			0xFFFFFFFFFFFFFF00,   // Example of non-zero hash with zero fingerprint
			uint64(1) << bits,    // Power of 2 (would be zero after mask)
			uint64(256),          // Would be zero for 8-bit fingerprint
			uint64(16),           // Would be zero for 4-bit fingerprint
		}

		for _, hashVal := range testValues {
			fp := fingerprint(hashVal, bits)
			if fp == 0 {
				t.Errorf("fingerprint(%d, %d bits) returned 0, should return 1",
					hashVal, bits)
			}
		}
	}
}
