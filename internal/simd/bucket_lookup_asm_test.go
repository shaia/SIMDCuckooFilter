package simd

import (
	"testing"
)

// TestBucketLookupAssemblyEdgeCases tests critical edge cases in assembly implementations
// These tests target bugs like incorrect return value offsets and register clobbering

// TestAVX2ReturnValueOffset specifically tests the bug where return value was written
// to ret+24(FP) instead of ret+25(FP) in the AVX2 assembly implementation
func TestAVX2ReturnValueOffset(t *testing.T) {
	// Test with various data patterns to ensure return value is correctly placed

	testCases := []struct {
		name     string
		data     []byte
		target   byte
		expected bool
	}{
		{
			name:     "found_first_position",
			data:     []byte{42, 1, 2, 3, 4, 5, 6, 7},
			target:   42,
			expected: true,
		},
		{
			name:     "found_last_position",
			data:     []byte{1, 2, 3, 4, 5, 6, 7, 42},
			target:   42,
			expected: true,
		},
		{
			name:     "not_found_empty",
			data:     []byte{},
			target:   42,
			expected: false,
		},
		{
			name:     "not_found_single",
			data:     []byte{1},
			target:   42,
			expected: false,
		},
		{
			name:     "not_found_multiple",
			data:     []byte{1, 2, 3, 4, 5, 6, 7, 8},
			target:   42,
			expected: false,
		},
		{
			name:     "found_multiple_matches",
			data:     []byte{42, 1, 42, 3, 42, 5, 42, 7},
			target:   42,
			expected: true,
		},
		{
			name:     "zero_target_not_found",
			data:     []byte{0, 0, 0, 0, 0, 0, 0, 0},
			target:   1,
			expected: false,
		},
		{
			name:     "zero_target_found",
			data:     []byte{1, 2, 3, 0, 5, 6, 7, 8},
			target:   0,
			expected: true,
		},
		{
			name:     "max_byte_value",
			data:     []byte{255, 254, 253, 252, 251, 250, 249, 248},
			target:   255,
			expected: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := BucketLookup(tc.data, tc.target)
			if result != tc.expected {
				t.Errorf("BucketLookup(%v, %d) = %v, want %v",
					tc.data, tc.target, result, tc.expected)
			}
		})
	}
}

// TestAVX2LargeBuffers tests with buffers larger than 32 bytes (AVX2 processes 32 at a time)
func TestAVX2LargeBuffers(t *testing.T) {
	sizes := []int{32, 33, 64, 65, 100, 128, 256}

	for _, size := range sizes {
		t.Run("", func(t *testing.T) {
			data := make([]byte, size)

			// Fill with non-matching values
			for i := 0; i < size; i++ {
				data[i] = byte(i % 255)
			}

			// Test not found
			if BucketLookup(data, 255) {
				t.Errorf("Found non-existent value in buffer of size %d", size)
			}

			// Place target at various positions
			positions := []int{0, 1, 31, 32, 33, 63, 64, size - 1}
			for _, pos := range positions {
				if pos >= size {
					continue
				}

				// Reset and place target
				for i := 0; i < size; i++ {
					data[i] = byte(i % 255)
				}
				data[pos] = 255

				if !BucketLookup(data, 255) {
					t.Errorf("Failed to find value at position %d in buffer of size %d", pos, size)
				}
			}
		})
	}
}

// TestAVX2RemainderHandling tests the remainder loop (for buffers not aligned to 32 bytes)
func TestAVX2RemainderHandling(t *testing.T) {
	// Test sizes that exercise the remainder loop: 1-31 bytes
	for size := 1; size <= 31; size++ {
		t.Run("", func(t *testing.T) {
			data := make([]byte, size)

			// Fill with sequential values
			for i := 0; i < size; i++ {
				data[i] = byte(i + 1)
			}

			// Test finding each position
			for target := 1; target <= size; target++ {
				if !BucketLookup(data, byte(target)) {
					t.Errorf("Failed to find byte %d in buffer of size %d", target, size)
				}
			}

			// Test not found
			if BucketLookup(data, 0) {
				t.Errorf("Found non-existent 0 in buffer of size %d", size)
			}
			if BucketLookup(data, byte(size+1)) {
				t.Errorf("Found non-existent %d in buffer of size %d", size+1, size)
			}
		})
	}
}

// TestBucketLookupBoundaryValues tests all possible byte values
func TestBucketLookupBoundaryValues(t *testing.T) {
	data := make([]byte, 64)

	// Test each possible byte value (1-255, skip 0 since zero-filled buffer contains 0)
	for target := 1; target <= 255; target++ {
		// Clear buffer (fill with zeros)
		for i := range data {
			data[i] = 0
		}

		// Should not find target (buffer is all zeros)
		if BucketLookup(data, byte(target)) {
			t.Errorf("Found %d in zero-filled buffer", target)
		}

		// Place target in middle
		data[32] = byte(target)

		// Should find target
		if !BucketLookup(data, byte(target)) {
			t.Errorf("Failed to find %d after placing in buffer", target)
		}
	}

	// Special test for 0: zero-filled buffer should find 0
	for i := range data {
		data[i] = 0
	}
	if !BucketLookup(data, 0) {
		t.Error("Failed to find 0 in zero-filled buffer")
	}
}

// TestBucketLookupEmptyBuffer tests edge case of empty buffer
func TestBucketLookupEmptyBuffer(t *testing.T) {
	data := []byte{}

	for target := 0; target <= 255; target++ {
		if BucketLookup(data, byte(target)) {
			t.Errorf("Found %d in empty buffer", target)
		}
	}
}

// TestBucketLookupSingleByte tests edge case of single byte buffer
func TestBucketLookupSingleByte(t *testing.T) {
	for value := 0; value <= 255; value++ {
		data := []byte{byte(value)}

		// Should find the exact value
		if !BucketLookup(data, byte(value)) {
			t.Errorf("Failed to find %d in single-byte buffer containing %d", value, value)
		}

		// Should not find other values
		otherValue := byte((value + 1) % 256)
		if BucketLookup(data, otherValue) {
			t.Errorf("Found %d in single-byte buffer containing %d", otherValue, value)
		}
	}
}

// TestBucketLookupAlignment tests with various alignment patterns
func TestBucketLookupAlignment(t *testing.T) {
	// Create a larger buffer and test with slices at different alignments
	largeBuffer := make([]byte, 128)
	for i := range largeBuffer {
		largeBuffer[i] = byte(i)
	}

	// Test slices starting at different offsets
	for offset := 0; offset < 32; offset++ {
		for size := 1; size <= 64 && offset+size <= len(largeBuffer); size++ {
			slice := largeBuffer[offset : offset+size]

			// Test finding first element
			if !BucketLookup(slice, slice[0]) {
				t.Errorf("Failed to find first element at offset %d, size %d", offset, size)
			}

			// Test finding last element
			if !BucketLookup(slice, slice[len(slice)-1]) {
				t.Errorf("Failed to find last element at offset %d, size %d", offset, size)
			}

			// Test not found
			if BucketLookup(slice, 255) {
				t.Errorf("Found non-existent value at offset %d, size %d", offset, size)
			}
		}
	}
}

// TestBucketLookupConsecutiveRepeats tests with repeated values
func TestBucketLookupConsecutiveRepeats(t *testing.T) {
	sizes := []int{4, 8, 16, 32, 64, 128}

	for _, size := range sizes {
		data := make([]byte, size)

		// All same value
		for i := range data {
			data[i] = 42
		}

		if !BucketLookup(data, 42) {
			t.Errorf("Failed to find repeated value in buffer of size %d", size)
		}

		if BucketLookup(data, 43) {
			t.Errorf("Found non-existent value in repeated buffer of size %d", size)
		}
	}
}

// TestBucketLookupAlternatingPattern tests with alternating byte patterns
func TestBucketLookupAlternatingPattern(t *testing.T) {
	size := 128
	data := make([]byte, size)

	// Alternating 0xAA and 0x55
	for i := range data {
		if i%2 == 0 {
			data[i] = 0xAA
		} else {
			data[i] = 0x55
		}
	}

	if !BucketLookup(data, 0xAA) {
		t.Error("Failed to find 0xAA in alternating pattern")
	}

	if !BucketLookup(data, 0x55) {
		t.Error("Failed to find 0x55 in alternating pattern")
	}

	if BucketLookup(data, 0x00) {
		t.Error("Found non-existent 0x00 in alternating pattern")
	}
}

// TestBucketLookupStressTest performs many random lookups to catch intermittent bugs
func TestBucketLookupStressTest(t *testing.T) {
	iterations := 1000

	for iter := 0; iter < iterations; iter++ {
		// Create buffer with pseudo-random content
		size := 1 + (iter % 127)
		data := make([]byte, size)

		for i := range data {
			data[i] = byte((iter*7 + i*13) % 256)
		}

		// Test each byte in the buffer
		for i := 0; i < size; i++ {
			target := data[i]
			if !BucketLookup(data, target) {
				t.Fatalf("Stress test iteration %d: failed to find byte %d at position %d in buffer of size %d",
					iter, target, i, size)
			}
		}

		// Test a value that's unlikely to be in the buffer
		unlikely := byte((iter * 11) % 256)
		found := false
		for i := 0; i < size; i++ {
			if data[i] == unlikely {
				found = true
				break
			}
		}

		result := BucketLookup(data, unlikely)
		if result != found {
			t.Fatalf("Stress test iteration %d: lookup(%d) = %v, expected %v",
				iter, unlikely, result, found)
		}
	}
}

// TestBucketLookupStackCorruption tests that assembly doesn't corrupt the stack
// by making many consecutive calls and verifying results
func TestBucketLookupStackCorruption(t *testing.T) {
	data1 := []byte{1, 2, 3, 4, 5}
	data2 := []byte{10, 20, 30, 40, 50}
	data3 := []byte{100, 101, 102, 103, 104}

	// Make many interleaved calls
	for i := 0; i < 100; i++ {
		r1 := BucketLookup(data1, 3)
		r2 := BucketLookup(data2, 30)
		r3 := BucketLookup(data3, 102)
		r4 := BucketLookup(data1, 99)
		r5 := BucketLookup(data2, 99)

		if !r1 {
			t.Fatal("Stack corruption: failed to find 3 in data1")
		}
		if !r2 {
			t.Fatal("Stack corruption: failed to find 30 in data2")
		}
		if !r3 {
			t.Fatal("Stack corruption: failed to find 102 in data3")
		}
		if r4 {
			t.Fatal("Stack corruption: found non-existent 99 in data1")
		}
		if r5 {
			t.Fatal("Stack corruption: found non-existent 99 in data2")
		}
	}
}
