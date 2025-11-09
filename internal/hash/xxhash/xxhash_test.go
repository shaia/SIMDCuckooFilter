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

// TestXXHashZeroFingerprint ensures fingerprint is never zero
func TestXXHashZeroFingerprint(t *testing.T) {
	// Test with various inputs to ensure fingerprint is never 0
	testInputs := [][]byte{
		[]byte(""),
		[]byte("a"),
		[]byte("test"),
		[]byte("hello world"),
		make([]byte, 100),
	}

	for _, input := range testInputs {
		hash := hash64XXHashInternal(input)
		fp8 := fingerprint(hash, 8)
		fp16 := fingerprint(hash, 16)
		fp32 := fingerprint(hash, 32)

		if fp8 == 0 {
			t.Errorf("8-bit fingerprint is zero for input %q (hash=%016x)", input, hash)
		}
		if fp16 == 0 {
			t.Errorf("16-bit fingerprint is zero for input %q (hash=%016x)", input, hash)
		}
		if fp32 == 0 {
			t.Errorf("32-bit fingerprint is zero for input %q (hash=%016x)", input, hash)
		}
	}
}
