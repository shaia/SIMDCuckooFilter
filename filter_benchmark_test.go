package cuckoofilter

import (
	"fmt"
	"testing"

	"github.com/shaia/simdcuckoofilter/internal/hash"
)

// BenchmarkInsert benchmarks insert operations
func BenchmarkInsert(b *testing.B) {
	cf, _ := New(100000)
	item := []byte("benchmark-item")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		cf.Insert(item)
	}
}

// BenchmarkLookup benchmarks lookup operations
func BenchmarkLookup(b *testing.B) {
	cf, _ := New(100000)
	item := []byte("benchmark-item")
	cf.Insert(item)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		cf.Lookup(item)
	}
}

// BenchmarkDelete benchmarks delete operations
func BenchmarkDelete(b *testing.B) {
	cf, _ := New(100000)
	item := []byte("benchmark-item")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		cf.Insert(item)
		cf.Delete(item)
	}
}

// BenchmarkInsertBatch benchmarks batch insert operations
func BenchmarkInsertBatch(b *testing.B) {
	cf, _ := New(100000)
	bf, ok := cf.(BatchFilter)
	if !ok {
		b.Skip("Filter does not implement BatchFilter")
	}

	items := make([][]byte, 16)
	for i := range items {
		items[i] = []byte(fmt.Sprintf("batch-item-%d", i))
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		bf.InsertBatch(items)
	}
}

// BenchmarkLookupBatch benchmarks batch lookup operations
func BenchmarkLookupBatch(b *testing.B) {
	cf, _ := New(100000)
	bf, ok := cf.(BatchFilter)
	if !ok {
		b.Skip("Filter does not implement BatchFilter")
	}

	items := make([][]byte, 16)
	for i := range items {
		items[i] = []byte(fmt.Sprintf("batch-item-%d", i))
	}
	bf.InsertBatch(items)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		bf.LookupBatch(items)
	}
}

// BenchmarkHashStrategies benchmarks different hash strategies
func BenchmarkHashStrategies(b *testing.B) {
	strategies := []hash.HashStrategy{
		hash.HashStrategyXXHash,
		hash.HashStrategyCRC32,
		hash.HashStrategyFNV,
	}

	for _, strategy := range strategies {
		b.Run(strategy.String(), func(b *testing.B) {
			cf, _ := New(100000, WithHashStrategy(strategy))
			item := []byte("benchmark-item")

			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				cf.Insert(item)
			}
		})
	}
}

// BenchmarkFingerprintSizes benchmarks different fingerprint sizes
func BenchmarkFingerprintSizes(b *testing.B) {
	sizes := []uint{4, 6, 8} // Only 1-8 bits supported (stored as bytes)

	for _, size := range sizes {
		b.Run(fmt.Sprintf("%d-bit", size), func(b *testing.B) {
			cf, err := New(100000, WithFingerprintSize(size))
			if err != nil {
				b.Fatalf("Failed to create filter: %v", err)
			}
			item := []byte("benchmark-item")

			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				cf.Insert(item)
			}
		})
	}
}

// BenchmarkBucketSizes benchmarks different bucket sizes
func BenchmarkBucketSizes(b *testing.B) {
	sizes := []uint{4, 8, 16, 32, 64}

	for _, size := range sizes {
		b.Run(fmt.Sprintf("bucket-%d", size), func(b *testing.B) {
			cf, _ := New(100000, WithBucketSize(size))
			item := []byte("benchmark-item")

			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				cf.Insert(item)
			}
		})
	}
}

// BenchmarkBatchSizes benchmarks different batch sizes
func BenchmarkBatchSizes(b *testing.B) {
	cf, _ := New(100000)
	bf, ok := cf.(BatchFilter)
	if !ok {
		b.Skip("Filter does not implement BatchFilter")
	}

	batchSizes := []int{4, 8, 16, 32, 64}

	for _, size := range batchSizes {
		b.Run(fmt.Sprintf("batch-%d", size), func(b *testing.B) {
			items := make([][]byte, size)
			for i := range items {
				items[i] = []byte(fmt.Sprintf("item-%d", i))
			}

			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				bf.InsertBatch(items)
			}
		})
	}
}

// BenchmarkInsertWithCapacities benchmarks insert with different capacities
func BenchmarkInsertWithCapacities(b *testing.B) {
	capacities := []uint{1000, 10000, 100000, 1000000}

	for _, capacity := range capacities {
		b.Run(fmt.Sprintf("capacity-%d", capacity), func(b *testing.B) {
			cf, _ := New(capacity)
			item := []byte("benchmark-item")

			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				cf.Insert(item)
			}
		})
	}
}

// BenchmarkMixedOperations benchmarks a mix of operations
func BenchmarkMixedOperations(b *testing.B) {
	cf, _ := New(100000)
	items := make([][]byte, 100)
	for i := range items {
		items[i] = []byte(fmt.Sprintf("mixed-item-%d", i))
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		idx := i % len(items)
		switch i % 3 {
		case 0:
			cf.Insert(items[idx])
		case 1:
			cf.Lookup(items[idx])
		case 2:
			cf.Delete(items[idx])
		}
	}
}

// BenchmarkConcurrentInsert benchmarks concurrent insert operations
func BenchmarkConcurrentInsert(b *testing.B) {
	cf, _ := New(100000)

	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			item := []byte(fmt.Sprintf("concurrent-%d", i))
			cf.Insert(item)
			i++
		}
	})
}

// BenchmarkConcurrentLookup benchmarks concurrent lookup operations
func BenchmarkConcurrentLookup(b *testing.B) {
	cf, _ := New(100000)

	// Pre-populate
	for i := 0; i < 1000; i++ {
		cf.Insert([]byte(fmt.Sprintf("concurrent-%d", i)))
	}

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			item := []byte(fmt.Sprintf("concurrent-%d", i%1000))
			cf.Lookup(item)
			i++
		}
	})
}

// BenchmarkConcurrentMixed benchmarks concurrent mixed operations
func BenchmarkConcurrentMixed(b *testing.B) {
	cf, _ := New(100000)

	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			item := []byte(fmt.Sprintf("concurrent-%d", i))
			switch i % 3 {
			case 0:
				cf.Insert(item)
			case 1:
				cf.Lookup(item)
			case 2:
				cf.Delete(item)
			}
			i++
		}
	})
}
