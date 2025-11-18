//go:build amd64
// +build amd64

package crc32hash

import (
	"fmt"
	"hash/crc32"
	"testing"
)

// BenchmarkSIMDBatchSizes benchmarks different batch sizes
func BenchmarkSIMDBatchSizes(b *testing.B) {
	table := crc32.MakeTable(crc32.Castagnoli)

	batchSizes := []int{4, 8, 16, 32, 64}

	for _, size := range batchSizes {
		// Generate test data
		items := make([][]byte, size)
		for i := range items {
			items[i] = []byte(fmt.Sprintf("benchmark-item-%d-with-some-data", i))
		}

		b.Run(fmt.Sprintf("size-%d", size), func(b *testing.B) {
			processor := NewBatchProcessor(table)
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				processor.ProcessBatch(items, 8, 1000)
			}
		})
	}
}

// BenchmarkSIMDRawCRC32 benchmarks raw CRC32 computation
func BenchmarkSIMDRawCRC32(b *testing.B) {
	items := make([][]byte, 16)
	for i := range items {
		items[i] = []byte(fmt.Sprintf("benchmark-item-%d", i))
	}
	results := make([]uint32, 16)

	b.Run("SIMD-Assembly", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			batchCRC32SIMD(items, results)
		}
	})

	b.Run("Stdlib-Sequential", func(b *testing.B) {
		table := crc32.MakeTable(crc32.Castagnoli)
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			for j, item := range items {
				results[j] = crc32.Checksum(item, table)
			}
		}
	})
}

// BenchmarkSIMDItemSizes benchmarks different item sizes
func BenchmarkSIMDItemSizes(b *testing.B) {
	table := crc32.MakeTable(crc32.Castagnoli)
	processor := NewBatchProcessor(table)

	itemSizes := []int{8, 16, 32, 64, 128, 256, 512, 1024}

	for _, itemSize := range itemSizes {
		items := make([][]byte, 16)
		for i := range items {
			data := make([]byte, itemSize)
			for j := range data {
				data[j] = byte(j % 256)
			}
			items[i] = data
		}

		b.Run(fmt.Sprintf("itemsize-%d", itemSize), func(b *testing.B) {
			b.SetBytes(int64(itemSize * 16)) // Total bytes processed per iteration
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				processor.ProcessBatch(items, 8, 1000)
			}
		})
	}
}

// BenchmarkSIMDThroughput measures throughput
func BenchmarkSIMDThroughput(b *testing.B) {
	table := crc32.MakeTable(crc32.Castagnoli)

	// Create large batch with 64-byte items
	const batchSize = 1000
	const itemSize = 64
	items := make([][]byte, batchSize)
	for i := range items {
		data := make([]byte, itemSize)
		for j := range data {
			data[j] = byte((i + j) % 256)
		}
		items[i] = data
	}

	b.Run("BatchProcessor", func(b *testing.B) {
		processor := NewBatchProcessor(table)
		b.SetBytes(batchSize * itemSize)
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			processor.ProcessBatch(items, 8, 10000)
		}
	})
}

// BenchmarkCompareSingleVsBatch compares single vs batch processing
func BenchmarkCompareSingleVsBatch(b *testing.B) {
	table := crc32.MakeTable(crc32.Castagnoli)
	processor := NewBatchProcessor(table)

	items := make([][]byte, 16)
	for i := range items {
		items[i] = []byte(fmt.Sprintf("test-item-%d", i))
	}

	b.Run("Single", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			for _, item := range items {
				_ = crc32.Checksum(item, table)
			}
		}
	})

	b.Run("BatchSIMD", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			processor.ProcessBatch(items, 8, 1000)
		}
	})
}
