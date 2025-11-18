package bucket

import (
	"fmt"
	"testing"
)

func TestSIMDBucketContains(t *testing.T) {
	sizes := []uint{2, 4, 8, 16, 32, 64}

	for _, size := range sizes {
		t.Run(fmt.Sprintf("Size%d", size), func(t *testing.T) {
			b := NewSIMDBucket(size)

			// Insert some fingerprints
			testFPs := []byte{10, 20, 30, 40}
			for i := uint(0); i < size && i < uint(len(testFPs)); i++ {
				b.Insert(testFPs[i])
			}

			// Test Contains for inserted values
			for i := uint(0); i < size && i < uint(len(testFPs)); i++ {
				if !b.ContainsSIMD(testFPs[i]) {
					t.Errorf("ContainsSIMD(%d) = false, want true", testFPs[i])
				}
			}

			// Test Contains for non-existent value
			if b.ContainsSIMD(99) {
				t.Error("ContainsSIMD(99) = true, want false")
			}

			// Verify SIMD and scalar give same results
			for fp := byte(0); fp < 255; fp++ {
				simdResult := b.ContainsSIMD(fp)
				scalarResult := b.Contains(fp)
				if simdResult != scalarResult {
					t.Errorf("ContainsSIMD(%d) = %v, Contains(%d) = %v", fp, simdResult, fp, scalarResult)
				}
			}
		})
	}
}

func TestSIMDBucketIsFull(t *testing.T) {
	sizes := []uint{2, 4, 8, 16, 32, 64}

	for _, size := range sizes {
		t.Run(fmt.Sprintf("Size%d", size), func(t *testing.T) {
			b := NewSIMDBucket(size)

			// Initially not full
			if b.IsFullSIMD() {
				t.Error("IsFullSIMD() = true for empty bucket, want false")
			}

			// Fill bucket
			for i := uint(0); i < size; i++ {
				b.Insert(byte(i + 1))
			}

			// Now should be full
			if !b.IsFullSIMD() {
				t.Error("IsFullSIMD() = false for full bucket, want true")
			}

			// Verify SIMD and scalar match
			if b.IsFullSIMD() != b.IsFull() {
				t.Error("IsFullSIMD() != IsFull()")
			}
		})
	}
}

func TestSIMDBucketCount(t *testing.T) {
	sizes := []uint{2, 4, 8, 16, 32, 64}

	for _, size := range sizes {
		t.Run(fmt.Sprintf("Size%d", size), func(t *testing.T) {
			b := NewSIMDBucket(size)

			// Test count as we add items
			for i := uint(0); i < size; i++ {
				if b.CountSIMD() != i {
					t.Errorf("CountSIMD() = %d, want %d", b.CountSIMD(), i)
				}
				b.Insert(byte(i + 1))
			}

			// Verify final count
			if b.CountSIMD() != size {
				t.Errorf("CountSIMD() = %d, want %d", b.CountSIMD(), size)
			}

			// Verify SIMD and scalar match
			if b.CountSIMD() != b.Count() {
				t.Errorf("CountSIMD() = %d, Count() = %d", b.CountSIMD(), b.Count())
			}
		})
	}
}

func TestSIMDBucketFindFirstZero(t *testing.T) {
	sizes := []uint{2, 4, 8, 16, 32, 64}

	for _, size := range sizes {
		t.Run(fmt.Sprintf("Size%d", size), func(t *testing.T) {
			b := NewSIMDBucket(size)

			// Initially, first zero should be at index 0
			if idx := b.FindFirstZeroSIMD(); idx != 0 {
				t.Errorf("FindFirstZeroSIMD() = %d, want 0", idx)
			}

			// Fill bucket gradually and check
			for i := uint(0); i < size; i++ {
				if idx := b.FindFirstZeroSIMD(); idx != i {
					t.Errorf("FindFirstZeroSIMD() = %d, want %d", idx, i)
				}
				b.Insert(byte(i + 1))
			}

			// When full, should return size
			if idx := b.FindFirstZeroSIMD(); idx != size {
				t.Errorf("FindFirstZeroSIMD() = %d, want %d (size)", idx, size)
			}
		})
	}
}

func TestSIMDBucketInsertSIMD(t *testing.T) {
	sizes := []uint{2, 4, 8, 16, 32, 64}

	for _, size := range sizes {
		t.Run(fmt.Sprintf("Size%d", size), func(t *testing.T) {
			b := NewSIMDBucket(size)

			// Insert until full
			for i := uint(0); i < size; i++ {
				if !b.InsertSIMD(byte(i + 1)) {
					t.Errorf("InsertSIMD(%d) = false, want true", i+1)
				}
			}

			// Next insert should fail
			if b.InsertSIMD(99) {
				t.Error("InsertSIMD(99) = true on full bucket, want false")
			}

			// Verify all items are present
			for i := uint(0); i < size; i++ {
				if !b.Contains(byte(i + 1)) {
					t.Errorf("Contains(%d) = false after InsertSIMD, want true", i+1)
				}
			}
		})
	}
}

// Benchmarks

func BenchmarkBucketContains(b *testing.B) {
	sizes := []uint{2, 4, 8, 16, 32, 64}

	for _, size := range sizes {
		// Setup bucket
		bucket := NewSIMDBucket(size)
		for i := uint(0); i < size; i++ {
			bucket.Insert(byte(i + 1))
		}
		fp := byte(size / 2) // Middle fingerprint

		b.Run(fmt.Sprintf("Scalar/Size%d", size), func(b *testing.B) {
			b.ReportAllocs()
			for i := 0; i < b.N; i++ {
				_ = bucket.Contains(fp)
			}
		})

		// Explicit SIMD version
		b.Run(fmt.Sprintf("SIMD/Size%d", size), func(b *testing.B) {
			b.ReportAllocs()
			for i := 0; i < b.N; i++ {
				_ = bucket.ContainsSIMD(fp)
			}
		})
	}
}

func BenchmarkBucketIsFull(b *testing.B) {
	sizes := []uint{2, 4, 8, 16, 32, 64}

	for _, size := range sizes {
		bucket := NewSIMDBucket(size)
		for i := uint(0); i < size; i++ {
			bucket.Insert(byte(i + 1))
		}

		b.Run(fmt.Sprintf("Scalar/Size%d", size), func(b *testing.B) {
			b.ReportAllocs()
			for i := 0; i < b.N; i++ {
				_ = bucket.IsFull()
			}
		})

		// Explicit SIMD version
		b.Run(fmt.Sprintf("SIMD/Size%d", size), func(b *testing.B) {
			b.ReportAllocs()
			for i := 0; i < b.N; i++ {
				_ = bucket.IsFullSIMD()
			}
		})
	}
}

func BenchmarkBucketCount(b *testing.B) {
	sizes := []uint{2, 4, 8, 16, 32, 64}

	for _, size := range sizes {
		bucket := NewSIMDBucket(size)
		for i := uint(0); i < size-1; i++ { // Leave one empty
			bucket.Insert(byte(i + 1))
		}

		b.Run(fmt.Sprintf("Scalar/Size%d", size), func(b *testing.B) {
			b.ReportAllocs()
			for i := 0; i < b.N; i++ {
				_ = bucket.Count()
			}
		})

		// Explicit SIMD version
		b.Run(fmt.Sprintf("SIMD/Size%d", size), func(b *testing.B) {
			b.ReportAllocs()
			for i := 0; i < b.N; i++ {
				_ = bucket.CountSIMD()
			}
		})
	}
}

func BenchmarkBucketFindFirstZero(b *testing.B) {
	sizes := []uint{2, 4, 8, 16, 32, 64}

	for _, size := range sizes {
		bucket := NewSIMDBucket(size)
		// Fill half the bucket
		for i := uint(0); i < size/2; i++ {
			bucket.Insert(byte(i + 1))
		}

		b.Run(fmt.Sprintf("FindFirstZero/Size%d", size), func(b *testing.B) {
			b.ReportAllocs()
			for i := 0; i < b.N; i++ {
				_ = bucket.FindFirstZeroSIMD()
			}
		})
	}
}
