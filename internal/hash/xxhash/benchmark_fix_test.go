package xxhash

import (
	"fmt"
	"testing"
)

// BenchmarkSIMDVsScalar compares the performance of the "SIMD" path (which uses
// our new unrolled scalar implementation) vs the original "Scalar" path (which
// uses the loop-based scalar implementation).
//
// We exploit the fact that:
// - Batch size >= 4 triggers the SIMD loop (our fix)
// - Batch size < 4 triggers the Scalar loop (original fallback)
func BenchmarkSIMDVsScalar(b *testing.B) {
	// Item sizes to test
	itemSizes := []int{32, 128, 1024}

	for _, size := range itemSizes {
		// Create a single item of data
		itemData := make([]byte, size)
		for i := range itemData {
			itemData[i] = byte(i)
		}

		// 1. Benchmark SIMD Path (Batch Size 4)
		// This forces execution of the unrolled code we just wrote.
		b.Run(fmt.Sprintf("SIMD_Fix/Batch4/%dbytes", size), func(b *testing.B) {
			proc := NewBatchHashProcessor()
			xxh := NewXXHash(8, proc)
			numBuckets := uint(1024)

			// Create batch of 4 items
			items := make([][]byte, 4)
			for i := range items {
				items[i] = itemData
			}

			b.SetBytes(int64(4 * size))
			b.ReportAllocs()
			b.ResetTimer()

			for i := 0; i < b.N; i++ {
				_ = xxh.GetIndicesBatch(items, numBuckets)
			}
		})

		// 2. Benchmark Scalar Path (Batch Size 3)
		// This forces execution of the original scalar_loop label in assembly.
		b.Run(fmt.Sprintf("Scalar_Fallback/Batch3/%dbytes", size), func(b *testing.B) {
			proc := NewBatchHashProcessor()
			xxh := NewXXHash(8, proc)
			numBuckets := uint(1024)

			// Create batch of 3 items
			items := make([][]byte, 3)
			for i := range items {
				items[i] = itemData
			}

			b.SetBytes(int64(3 * size))
			b.ReportAllocs()
			b.ResetTimer()

			for i := 0; i < b.N; i++ {
				_ = xxh.GetIndicesBatch(items, numBuckets)
			}
		})
	}
}
