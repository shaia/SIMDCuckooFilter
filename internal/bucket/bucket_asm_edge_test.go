package bucket

import (
	"testing"
)

// TestAssemblyEdgeCases tests critical edge cases that assembly implementations must handle correctly.
// These tests are designed to catch bugs like:
// - Register clobbering (e.g., overwriting data pointers with return values)
// - Incorrect return value offsets in stack frames
// - Off-by-one errors in boundary conditions
// - Empty array handling
// - Full array handling

// TestIsFullEmptyBucket specifically tests the bug where R0 (data pointer) was overwritten
// with value 1 before being used in the scalar loop (ARM64 assembly bug at line 743)
func TestIsFullEmptyBucket(t *testing.T) {
	// Test with bucket of size 0 (empty bucket created differently)
	// This should not crash even though the assembly code has to handle len=0

	// Test sizes including very small buckets
	sizes := []uint{1, 2, 3, 4, 8, 16, 32, 64}

	for _, size := range sizes {
		t.Run("", func(t *testing.T) {
			b := NewSIMDBucket(size)

			// Empty bucket should not be full
			if b.IsFullSIMD() {
				t.Errorf("Empty bucket of size %d reported as full", size)
			}

			// Verify consistency between SIMD and scalar
			if b.IsFullSIMD() != b.IsFull() {
				t.Errorf("IsFullSIMD() != IsFull() for empty bucket of size %d", size)
			}
		})
	}
}

// TestIsFullSingleElement tests edge case with single element bucket
// This tests the scalar fallback path thoroughly
func TestIsFullSingleElement(t *testing.T) {
	b := NewSIMDBucket(1)

	// Initially empty, should return false
	if b.IsFullSIMD() {
		t.Error("Single-element empty bucket reported as full")
	}

	// Insert one element
	b.Insert(42)

	// Now should be full
	if !b.IsFullSIMD() {
		t.Error("Single-element full bucket reported as not full")
	}
}

// TestContainsFirstPosition tests that assembly correctly checks the first byte
func TestContainsFirstPosition(t *testing.T) {
	sizes := []uint{2, 4, 8, 16, 32, 64}

	for _, size := range sizes {
		t.Run("", func(t *testing.T) {
			b := NewSIMDBucket(size)
			b.fingerprints[0] = 123 // Put target in first position

			if !b.ContainsSIMD(123) {
				t.Errorf("Failed to find fingerprint at position 0 in bucket of size %d", size)
			}
		})
	}
}

// TestContainsLastPosition tests that assembly correctly checks the last byte
func TestContainsLastPosition(t *testing.T) {
	sizes := []uint{2, 4, 8, 16, 32, 64}

	for _, size := range sizes {
		t.Run("", func(t *testing.T) {
			b := NewSIMDBucket(size)
			// Fill with non-matching values except last
			for i := uint(0); i < size-1; i++ {
				b.fingerprints[i] = uint16(i + 1)
			}
			b.fingerprints[size-1] = 123 // Put target in last position

			if !b.ContainsSIMD(123) {
				t.Errorf("Failed to find fingerprint at last position in bucket of size %d", size)
			}
		})
	}
}

// TestCountZeroFingerprints ensures 0 fingerprints are correctly counted as empty
func TestCountZeroFingerprints(t *testing.T) {
	sizes := []uint{2, 4, 8, 16, 32, 64}

	for _, size := range sizes {
		t.Run("", func(t *testing.T) {
			b := NewSIMDBucket(size)

			// All zeros should count as 0
			count := b.CountSIMD()
			if count != 0 {
				t.Errorf("Empty bucket of size %d counted %d items, want 0", size, count)
			}

			// Insert zeros explicitly (should still count as 0 because 0 means empty)
			for i := uint(0); i < size; i++ {
				b.fingerprints[i] = 0
			}

			count = b.CountSIMD()
			if count != 0 {
				t.Errorf("Bucket with explicit zeros counted %d items, want 0", count)
			}
		})
	}
}

// TestCountOneFingerprint tests counting with exactly one non-zero fingerprint
func TestCountOneFingerprint(t *testing.T) {
	sizes := []uint{2, 4, 8, 16, 32, 64}

	for _, size := range sizes {
		t.Run("", func(t *testing.T) {
			for pos := uint(0); pos < size; pos++ {
				b := NewSIMDBucket(size)
				b.fingerprints[pos] = 42

				count := b.CountSIMD()
				if count != 1 {
					t.Errorf("Bucket size %d with one fingerprint at pos %d counted %d items, want 1",
						size, pos, count)
				}
			}
		})
	}
}

// TestFindFirstZeroAtStart tests finding zero at the very first position
func TestFindFirstZeroAtStart(t *testing.T) {
	sizes := []uint{2, 4, 8, 16, 32, 64}

	for _, size := range sizes {
		t.Run("", func(t *testing.T) {
			b := NewSIMDBucket(size)
			// All zeros, so first zero should be at index 0

			idx := b.FindFirstZeroSIMD()
			if idx != 0 {
				t.Errorf("FindFirstZeroSIMD() on empty bucket of size %d returned %d, want 0", size, idx)
			}
		})
	}
}

// TestFindFirstZeroAtEnd tests finding zero at the last position
func TestFindFirstZeroAtEnd(t *testing.T) {
	sizes := []uint{2, 4, 8, 16, 32, 64}

	for _, size := range sizes {
		t.Run("", func(t *testing.T) {
			b := NewSIMDBucket(size)
			// Fill all but last
			for i := uint(0); i < size-1; i++ {
				b.fingerprints[i] = uint16(i + 1)
			}
			// Last position is 0

			idx := b.FindFirstZeroSIMD()
			if idx != size-1 {
				t.Errorf("FindFirstZeroSIMD() with zero at end returned %d, want %d", idx, size-1)
			}
		})
	}
}

// TestFindFirstZeroNone tests that we return size when no zeros exist
func TestFindFirstZeroNone(t *testing.T) {
	sizes := []uint{2, 4, 8, 16, 32, 64}

	for _, size := range sizes {
		t.Run("", func(t *testing.T) {
			b := NewSIMDBucket(size)
			// Fill completely with non-zero values
			for i := uint(0); i < size; i++ {
				b.fingerprints[i] = uint16(i + 1)
			}

			idx := b.FindFirstZeroSIMD()
			if idx != size {
				t.Errorf("FindFirstZeroSIMD() on full bucket returned %d, want %d (size)", idx, size)
			}
		})
	}
}

// TestContainsZeroFingerprint tests searching for 0
// Note: While fingerprints are never 0 in practice (hash functions convert 0 to 1),
// the bucket-level Contains operation can still search for and find 0.
func TestContainsZeroFingerprint(t *testing.T) {
	sizes := []uint{2, 4, 8, 16, 32, 64}

	for _, size := range sizes {
		t.Run("", func(t *testing.T) {
			b := NewSIMDBucket(size)

			// Bucket with all zeros - should find 0
			simdResult := b.ContainsSIMD(0)
			scalarResult := b.Contains(0)
			if simdResult != scalarResult {
				t.Errorf("ContainsSIMD(0) = %v, Contains(0) = %v (mismatch)", simdResult, scalarResult)
			}

			// Explicitly set a position to 0
			for i := uint(0); i < size; i++ {
				b.fingerprints[i] = uint16(i + 1)
			}
			b.fingerprints[0] = 0

			simdResult = b.ContainsSIMD(0)
			scalarResult = b.Contains(0)
			if simdResult != scalarResult {
				t.Errorf("ContainsSIMD(0) = %v, Contains(0) = %v (mismatch with explicit 0)", simdResult, scalarResult)
			}

			// No zeros - should not find 0
			for i := uint(0); i < size; i++ {
				b.fingerprints[i] = uint16(i + 1)
			}

			simdResult = b.ContainsSIMD(0)
			scalarResult = b.Contains(0)
			if simdResult != scalarResult {
				t.Errorf("ContainsSIMD(0) = %v, Contains(0) = %v (mismatch with no zeros)", simdResult, scalarResult)
			}
		})
	}
}

// TestAllBytesValue tests with all 256 possible byte values
func TestAllBytesValue(t *testing.T) {
	size := uint(8) // Use size 8 for reasonable test time

	// Loop iterates through all non-zero uint16 values (1-65535).
	// The condition target != 0 relies on uint16 overflow: 65535 + 1 wraps to 0, terminating the loop.
	for target := uint16(1); target != 0; target++ { // 1-65535 (skip 0 as it means empty)
		b := NewSIMDBucket(size)

		// Test not found case
		if b.ContainsSIMD(target) {
			t.Errorf("ContainsSIMD(%d) returned true for empty bucket", target)
		}

		// Insert the target
		b.fingerprints[size/2] = target

		// Test found case
		if !b.ContainsSIMD(target) {
			t.Errorf("ContainsSIMD(%d) returned false when fingerprint exists", target)
		}
	}
}

// TestIsFullBoundaryValues tests IsFull with specific boundary patterns
func TestIsFullBoundaryValues(t *testing.T) {
	sizes := []uint{2, 4, 8, 16, 32, 64}

	for _, size := range sizes {
		t.Run("", func(t *testing.T) {
			b := NewSIMDBucket(size)

			// Pattern: 0xFFFF (all bits set)
			for i := uint(0); i < size; i++ {
				b.fingerprints[i] = 0xFFFF
			}
			if !b.IsFullSIMD() {
				t.Errorf("Bucket filled with 0xFFFF not detected as full (size %d)", size)
			}

			// Pattern: 0x01 (minimal non-zero)
			for i := uint(0); i < size; i++ {
				b.fingerprints[i] = 0x01
			}
			if !b.IsFullSIMD() {
				t.Errorf("Bucket filled with 0x01 not detected as full (size %d)", size)
			}

			// Pattern: alternating high/low
			for i := uint(0); i < size; i++ {
				if i%2 == 0 {
					b.fingerprints[i] = 0xAAAA
				} else {
					b.fingerprints[i] = 0x5555
				}
			}
			if !b.IsFullSIMD() {
				t.Errorf("Bucket with alternating pattern not detected as full (size %d)", size)
			}
		})
	}
}

// TestCountBoundaryValues tests Count with specific patterns
func TestCountBoundaryValues(t *testing.T) {
	size := uint(16)

	testCases := []struct {
		name     string
		pattern  func(uint) uint16
		expected uint
	}{
		{
			name:     "all_ones",
			pattern:  func(i uint) uint16 { return 1 },
			expected: size,
		},
		{
			name:     "all_max",
			pattern:  func(i uint) uint16 { return 0xFFFF },
			expected: size,
		},
		{
			name: "alternating_zero_one",
			pattern: func(i uint) uint16 {
				if i%2 == 0 {
					return 0
				} else {
					return 1
				}
			},
			expected: size / 2,
		},
		{
			name: "first_half_filled",
			pattern: func(i uint) uint16 {
				if i < size/2 {
					return 1
				} else {
					return 0
				}
			},
			expected: size / 2,
		},
		{
			name: "last_half_filled",
			pattern: func(i uint) uint16 {
				if i >= size/2 {
					return 1
				} else {
					return 0
				}
			},
			expected: size / 2,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			b := NewSIMDBucket(size)
			for i := uint(0); i < size; i++ {
				b.fingerprints[i] = tc.pattern(i)
			}

			count := b.CountSIMD()
			if count != tc.expected {
				t.Errorf("CountSIMD() = %d, want %d for pattern %s", count, tc.expected, tc.name)
			}

			// Verify against scalar implementation
			scalarCount := b.Count()
			if count != scalarCount {
				t.Errorf("CountSIMD() = %d, Count() = %d (mismatch)", count, scalarCount)
			}
		})
	}
}

// TestReturnValueConsistency ensures return values are correctly placed in stack frame
// This would catch the AVX2 bug where ret+24(FP) was used instead of ret+25(FP)
func TestReturnValueConsistency(t *testing.T) {
	b := NewSIMDBucket(8)

	// Test ContainsSIMD return values
	for i := 0; i < 1000; i++ {
		b.fingerprints[0] = uint16(i % 65536)
		target := uint16(i % 65536)

		result := b.ContainsSIMD(target)
		expected := b.Contains(target)

		if result != expected {
			t.Fatalf("Return value inconsistency: ContainsSIMD(%d) = %v, Contains(%d) = %v",
				target, result, target, expected)
		}
	}

	// Test IsFullSIMD return values
	for i := uint(0); i <= b.size; i++ {
		// Clear and fill i positions
		for j := uint(0); j < b.size; j++ {
			b.fingerprints[j] = 0
		}
		for j := uint(0); j < i; j++ {
			b.fingerprints[j] = uint16(j + 1)
		}

		result := b.IsFullSIMD()
		expected := b.IsFull()

		if result != expected {
			t.Fatalf("Return value inconsistency at fill level %d: IsFullSIMD() = %v, IsFull() = %v",
				i, result, expected)
		}
	}
}

// TestRegisterClobbering tests for register corruption bugs
// This would catch bugs like R0 being overwritten in ARM64 assembly
func TestRegisterClobbering(t *testing.T) {
	// Run many iterations with different data patterns
	// to increase chance of catching register corruption

	sizes := []uint{2, 4, 8, 16, 32, 64}
	iterations := 100

	for _, size := range sizes {
		for iter := 0; iter < iterations; iter++ {
			b := NewSIMDBucket(size)

			// Fill with pseudo-random pattern
			for i := uint(0); i < size; i++ {
				b.fingerprints[i] = uint16((iter*7 + int(i)*13) % 65536)
			}

			// Test all operations
			target := uint16((iter * 3) % 65536)

			simdContains := b.ContainsSIMD(target)
			scalarContains := b.Contains(target)
			if simdContains != scalarContains {
				t.Fatalf("Register clobbering detected in Contains: iter=%d, size=%d, target=%d, SIMD=%v, scalar=%v",
					iter, size, target, simdContains, scalarContains)
			}

			simdFull := b.IsFullSIMD()
			scalarFull := b.IsFull()
			if simdFull != scalarFull {
				t.Fatalf("Register clobbering detected in IsFull: iter=%d, size=%d, SIMD=%v, scalar=%v",
					iter, size, simdFull, scalarFull)
			}

			simdCount := b.CountSIMD()
			scalarCount := b.Count()
			if simdCount != scalarCount {
				t.Fatalf("Register clobbering detected in Count: iter=%d, size=%d, SIMD=%d, scalar=%d",
					iter, size, simdCount, scalarCount)
			}

			simdIdx := b.FindFirstZeroSIMD()
			scalarIdx := inlineFindFirstZero(b.fingerprints[:b.size])
			if simdIdx != scalarIdx {
				t.Fatalf("Register clobbering detected in FindFirstZero: iter=%d, size=%d, SIMD=%d, scalar=%d",
					iter, size, simdIdx, scalarIdx)
			}
		}
	}
}
