package hash

import (
	"fmt"
	"hash/crc32"
	"testing"

	crc32hash "github.com/shaia/cuckoofilter/internal/hash/crc32"
	"github.com/shaia/cuckoofilter/internal/hash/fnv"
	"github.com/shaia/cuckoofilter/internal/hash/types"
	"github.com/shaia/cuckoofilter/internal/hash/xxhash"
)

func BenchmarkXXHash(b *testing.B) {
	data := make([]byte, 64)
	for i := range data {
		data[i] = byte(i)
	}

	xxh := xxhash.NewXXHash(8, nil)

	b.Run("Optimized", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			_, _, _ = xxh.GetIndices(data, 1024)
		}
	})
}

func BenchmarkBatchXXHash(b *testing.B) {
	// Create test data
	numItems := []int{4, 8, 16, 32, 64, 128}

	for _, n := range numItems {
		items := make([][]byte, n)
		for i := range items {
			items[i] = make([]byte, 32)
			for j := range items[i] {
				items[i][j] = byte(i*n + j)
			}
		}

		b.Run(fmt.Sprintf("Scalar/%ditems", n), func(b *testing.B) {
			xxh := xxhash.NewXXHash(8, nil)
			numBuckets := uint(1024)
			results := make([]types.HashResult, len(items))
			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				for j, item := range items {
					i1, i2, fp := xxh.GetIndices(item, numBuckets)
					results[j] = types.HashResult{I1: i1, I2: i2, Fp: fp}
				}
			}
		})

		// Benchmark SIMD (automatically uses best available for platform)
		b.Run(fmt.Sprintf("SIMD/%ditems", n), func(b *testing.B) {
			proc := xxhash.NewBatchHashProcessor()
			xxh := xxhash.NewXXHash(8, proc)
			numBuckets := uint(1024)
			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				_ = xxh.GetIndicesBatch(items, numBuckets)
			}
		})
	}
}

func TestXXHashConsistency(t *testing.T) {
	// Verify that hash function produces consistent results
	testCases := [][]byte{
		[]byte("hello"),
		[]byte("world"),
		[]byte(""),
		[]byte("a"),
		[]byte("0123456789abcdef"),
		make([]byte, 100),
	}

	xxh := xxhash.NewXXHash(8, nil)
	numBuckets := uint(1024)

	for i, tc := range testCases {
		i1_1, i2_1, fp1 := xxh.GetIndices(tc, numBuckets)
		i1_2, i2_2, fp2 := xxh.GetIndices(tc, numBuckets)

		if i1_1 != i1_2 || i2_1 != i2_2 || fp1 != fp2 {
			t.Errorf("Test case %d: inconsistent results: (%d,%d,%x) vs (%d,%d,%x)",
				i, i1_1, i2_1, fp1, i1_2, i2_2, fp2)
		}
	}
}

func TestBatchXXHashConsistency(t *testing.T) {
	items := [][]byte{
		[]byte("hello"),
		[]byte("world"),
		[]byte("test"),
		[]byte("data"),
		[]byte("0123456789"),
	}

	numBuckets := uint(1024)

	// Get scalar results
	xxhScalar := xxhash.NewXXHash(8, nil)
	scalarResults := make([]types.HashResult, len(items))
	for i, item := range items {
		i1, i2, fp := xxhScalar.GetIndices(item, numBuckets)
		scalarResults[i] = types.HashResult{I1: i1, I2: i2, Fp: fp}
	}

	// Test SIMD batch processing (automatically uses best available)
	proc := xxhash.NewBatchHashProcessor()
	xxhSIMD := xxhash.NewXXHash(8, proc)
	simdResults := xxhSIMD.GetIndicesBatch(items, numBuckets)

	for i := range scalarResults {
		if scalarResults[i] != simdResults[i] {
			t.Errorf("SIMD item %d: scalar=%+v, simd=%+v", i, scalarResults[i], simdResults[i])
		}
	}
	t.Log("SIMD consistency test passed")
}

// BenchmarkCRC32Hash benchmarks CRC32C hash performance
func BenchmarkCRC32Hash(b *testing.B) {
	dataSizes := []int{8, 16, 32, 64, 128, 256}

	for _, size := range dataSizes {
		data := make([]byte, size)
		for i := range data {
			data[i] = byte(i)
		}

		b.Run(fmt.Sprintf("%dbytes", size), func(b *testing.B) {
			table := crc32.MakeTable(crc32.Castagnoli)
			crc := crc32hash.NewCRC32Hash(table, 8, nil)
			numBuckets := uint(1024)
			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				_, _, _ = crc.GetIndices(data, numBuckets)
			}
		})
	}
}

// BenchmarkFNVHash benchmarks FNV-1a hash performance
func BenchmarkFNVHash(b *testing.B) {
	dataSizes := []int{8, 16, 32, 64, 128, 256}

	for _, size := range dataSizes {
		data := make([]byte, size)
		for i := range data {
			data[i] = byte(i)
		}

		b.Run(fmt.Sprintf("%dbytes", size), func(b *testing.B) {
			fnvHash := fnv.NewFNVHash(8, nil)
			numBuckets := uint(1024)
			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				_, _, _ = fnvHash.GetIndices(data, numBuckets)
			}
		})
	}
}

// BenchmarkHashComparison compares all three hash implementations
func BenchmarkHashComparison(b *testing.B) {
	data := make([]byte, 64)
	for i := range data {
		data[i] = byte(i)
	}
	numBuckets := uint(1024)

	b.Run("XXHash", func(b *testing.B) {
		xxh := xxhash.NewXXHash(8, nil)
		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_, _, _ = xxh.GetIndices(data, numBuckets)
		}
	})

	b.Run("CRC32", func(b *testing.B) {
		table := crc32.MakeTable(crc32.Castagnoli)
		crc := crc32hash.NewCRC32Hash(table, 8, nil)
		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_, _, _ = crc.GetIndices(data, numBuckets)
		}
	})

	b.Run("FNV", func(b *testing.B) {
		fnvHash := fnv.NewFNVHash(8, nil)
		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_, _, _ = fnvHash.GetIndices(data, numBuckets)
		}
	})
}

// BenchmarkBatchFNVHash benchmarks FNV batch processing
func BenchmarkBatchFNVHash(b *testing.B) {
	numItems := []int{4, 8, 16, 32, 64}

	for _, n := range numItems {
		items := make([][]byte, n)
		for i := range items {
			items[i] = make([]byte, 32)
			for j := range items[i] {
				items[i][j] = byte(i*n + j)
			}
		}

		b.Run(fmt.Sprintf("Sequential/%ditems", n), func(b *testing.B) {
			fnvHash := fnv.NewFNVHash(8, nil)
			numBuckets := uint(1024)
			results := make([]types.HashResult, len(items))
			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				for j, item := range items {
					i1, i2, fp := fnvHash.GetIndices(item, numBuckets)
					results[j] = types.HashResult{I1: i1, I2: i2, Fp: fp}
				}
			}
		})

		// Benchmark FNV batch processor (AMD64 has parallel implementation)
		b.Run(fmt.Sprintf("Batch/%ditems", n), func(b *testing.B) {
			proc := fnv.NewBatchProcessor()
			fnvHash := fnv.NewFNVHash(8, proc)
			numBuckets := uint(1024)
			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				_ = fnvHash.GetIndicesBatch(items, numBuckets)
			}
		})
	}
}

// BenchmarkGetAltIndex benchmarks alternative index calculation
func BenchmarkGetAltIndex(b *testing.B) {
	numBuckets := uint(1024)
	index := uint(42)
	fp := byte(123)

	b.Run("XXHash", func(b *testing.B) {
		xxh := xxhash.NewXXHash(8, nil)
		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_ = xxh.GetAltIndex(index, fp, numBuckets)
		}
	})

	b.Run("CRC32", func(b *testing.B) {
		table := crc32.MakeTable(crc32.Castagnoli)
		crc := crc32hash.NewCRC32Hash(table, 8, nil)
		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_ = crc.GetAltIndex(index, fp, numBuckets)
		}
	})

	b.Run("FNV", func(b *testing.B) {
		fnvHash := fnv.NewFNVHash(8, nil)
		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_ = fnvHash.GetAltIndex(index, fp, numBuckets)
		}
	})
}
