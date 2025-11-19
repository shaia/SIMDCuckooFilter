//go:build amd64 || arm64

package filter

import (
	"fmt"
	"sync"
	"testing"

	"github.com/shaia/simdcuckoofilter/internal/hash"
)

// TestNewFilter tests filter creation
func TestNewFilter(t *testing.T) {
	f, err := New(1000, 4, 8, 500, hash.HashStrategyXXHash, 32)
	if err != nil {
		t.Fatalf("New failed: %v", err)
	}

	if f == nil {
		t.Fatal("Filter is nil")
	}

	if f.numBuckets == 0 {
		t.Error("numBuckets should not be 0")
	}

	if f.numItems != 0 {
		t.Errorf("numItems should be 0, got %d", f.numItems)
	}

	if f.bucketSize != 4 {
		t.Errorf("bucketSize should be 4, got %d", f.bucketSize)
	}

	if f.maxKicks != 500 {
		t.Errorf("maxKicks should be 500, got %d", f.maxKicks)
	}

	if f.batchSize != 32 {
		t.Errorf("batchSize should be 32, got %d", f.batchSize)
	}
}

// TestFilterInsertBasic tests basic insert operation
func TestFilterInsertBasic(t *testing.T) {
	f, _ := New(1000, 4, 8, 500, hash.HashStrategyXXHash, 32)

	item := []byte("test-item")
	if !f.Insert(item) {
		t.Error("Insert failed")
	}

	if f.numItems != 1 {
		t.Errorf("Expected numItems=1, got %d", f.numItems)
	}
}

// TestFilterLookup tests lookup operation
func TestFilterLookup(t *testing.T) {
	f, _ := New(1000, 4, 8, 500, hash.HashStrategyXXHash, 32)

	item := []byte("test-item")
	f.Insert(item)

	if !f.Lookup(item) {
		t.Error("Lookup failed to find inserted item")
	}

	// Non-existent item should not be found (usually)
	if f.Lookup([]byte("non-existent-unique-string-12345")) {
		t.Log("False positive (expected occasionally)")
	}
}

// TestFilterDelete tests delete operation
func TestFilterDelete(t *testing.T) {
	f, _ := New(1000, 4, 8, 500, hash.HashStrategyXXHash, 32)

	item := []byte("test-item")
	f.Insert(item)

	if !f.Delete(item) {
		t.Error("Delete failed")
	}

	if f.numItems != 0 {
		t.Errorf("Expected numItems=0 after delete, got %d", f.numItems)
	}
}

// TestFilterCapacity tests capacity calculation
func TestFilterCapacity(t *testing.T) {
	f, _ := New(1000, 4, 8, 500, hash.HashStrategyXXHash, 32)

	capacity := f.Capacity()
	if capacity == 0 {
		t.Error("Capacity should not be 0")
	}

	expectedCapacity := f.numBuckets * f.bucketSize
	if capacity != expectedCapacity {
		t.Errorf("Expected capacity %d, got %d", expectedCapacity, capacity)
	}
}

// TestFilterLoadFactor tests load factor calculation
func TestFilterLoadFactor(t *testing.T) {
	f, _ := New(1000, 4, 8, 500, hash.HashStrategyXXHash, 32)

	// Initially 0
	if f.LoadFactor() != 0 {
		t.Errorf("Initial load factor should be 0, got %f", f.LoadFactor())
	}

	// Insert some items
	for i := 0; i < 100; i++ {
		f.Insert([]byte(fmt.Sprintf("item-%d", i)))
	}

	lf := f.LoadFactor()
	if lf <= 0 || lf > 1 {
		t.Errorf("Load factor should be between 0 and 1, got %f", lf)
	}
}

// TestFilterReset tests reset operation
func TestFilterReset(t *testing.T) {
	f, _ := New(1000, 4, 8, 500, hash.HashStrategyXXHash, 32)

	// Insert items
	for i := 0; i < 10; i++ {
		f.Insert([]byte(fmt.Sprintf("item-%d", i)))
	}

	if f.numItems == 0 {
		t.Error("numItems should not be 0 after inserts")
	}

	f.Reset()

	if f.numItems != 0 {
		t.Errorf("Expected numItems=0 after reset, got %d", f.numItems)
	}

	if f.LoadFactor() != 0 {
		t.Errorf("Expected load factor=0 after reset, got %f", f.LoadFactor())
	}
}

// TestFilterBucketDistribution tests that items are distributed across buckets
func TestFilterBucketDistribution(t *testing.T) {
	f, _ := New(1000, 4, 8, 500, hash.HashStrategyXXHash, 32)

	// Insert many items
	numItems := 100
	for i := 0; i < numItems; i++ {
		f.Insert([]byte(fmt.Sprintf("item-%d", i)))
	}

	// Count non-empty buckets
	nonEmptyBuckets := 0
	for _, b := range f.buckets {
		if b.Count() > 0 {
			nonEmptyBuckets++
		}
	}

	// Should use multiple buckets (not all in one bucket)
	if nonEmptyBuckets <= 1 {
		t.Errorf("Items should be distributed across multiple buckets, got %d", nonEmptyBuckets)
	}

	t.Logf("Items distributed across %d/%d buckets", nonEmptyBuckets, f.numBuckets)
}

// TestFilterRelocate tests the relocation mechanism
func TestFilterRelocate(t *testing.T) {
	// Create small filter to force relocations
	f, _ := New(50, 4, 8, 500, hash.HashStrategyXXHash, 32)

	// Insert items within capacity - should succeed
	numItems := 40 // Less than capacity
	for i := 0; i < numItems; i++ {
		item := []byte(fmt.Sprintf("item-%d", i))
		if !f.Insert(item) {
			t.Errorf("Insert failed for item %d within capacity", i)
		}
	}

	// All items should be findable
	foundCount := 0
	for i := 0; i < numItems; i++ {
		item := []byte(fmt.Sprintf("item-%d", i))
		if f.Lookup(item) {
			foundCount++
		}
	}

	// Should find most items (allowing for some hash collisions)
	minExpected := numItems * 80 / 100 // At least 80%
	if foundCount < minExpected {
		t.Errorf("Expected to find at least %d items, found %d", minExpected, foundCount)
	}

	t.Logf("Successfully inserted and found %d/%d items", foundCount, numItems)
}

// TestFilterConcurrentInsert tests concurrent insert operations
func TestFilterConcurrentInsert(t *testing.T) {
	f, _ := New(10000, 4, 8, 500, hash.HashStrategyXXHash, 32)

	var wg sync.WaitGroup
	numGoroutines := 10
	itemsPerGoroutine := 100

	for g := 0; g < numGoroutines; g++ {
		wg.Add(1)
		go func(goroutineID int) {
			defer wg.Done()
			for i := 0; i < itemsPerGoroutine; i++ {
				item := []byte(fmt.Sprintf("goroutine-%d-item-%d", goroutineID, i))
				f.Insert(item)
			}
		}(g)
	}

	wg.Wait()

	// Count should be close to total items (some may fail due to collisions)
	minExpected := uint(numGoroutines * itemsPerGoroutine * 80 / 100) // At least 80%
	if f.Count() < minExpected {
		t.Errorf("Expected at least %d items, got %d", minExpected, f.Count())
	}
}

// TestFilterConcurrentLookup tests concurrent lookup operations
func TestFilterConcurrentLookup(t *testing.T) {
	f, _ := New(1000, 4, 8, 500, hash.HashStrategyXXHash, 32)

	// Pre-populate
	numItems := 100
	for i := 0; i < numItems; i++ {
		f.Insert([]byte(fmt.Sprintf("item-%d", i)))
	}

	var wg sync.WaitGroup
	numGoroutines := 10

	// Channel to collect failures from goroutines (thread-safe)
	failures := make(chan string, numGoroutines*numItems)

	for g := 0; g < numGoroutines; g++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for i := 0; i < numItems; i++ {
				item := []byte(fmt.Sprintf("item-%d", i))
				if !f.Lookup(item) {
					failures <- fmt.Sprintf("Failed to find item-%d", i)
				}
			}
		}()
	}

	wg.Wait()
	close(failures)

	// Report any failures after goroutines complete (thread-safe)
	for failure := range failures {
		t.Error(failure)
	}
}

// TestFilterBatchInsert tests batch insert operation
func TestFilterBatchInsert(t *testing.T) {
	f, _ := New(1000, 4, 8, 500, hash.HashStrategyXXHash, 32)

	items := [][]byte{
		[]byte("batch-item-1"),
		[]byte("batch-item-2"),
		[]byte("batch-item-3"),
		[]byte("batch-item-4"),
	}

	results := f.InsertBatch(items)

	if len(results) != len(items) {
		t.Errorf("Expected %d results, got %d", len(items), len(results))
	}

	for i, success := range results {
		if !success {
			t.Errorf("Batch insert failed for item %d", i)
		}
	}
}

// TestFilterBatchLookup tests batch lookup operation
func TestFilterBatchLookup(t *testing.T) {
	f, _ := New(1000, 4, 8, 500, hash.HashStrategyXXHash, 32)

	items := [][]byte{
		[]byte("batch-item-1"),
		[]byte("batch-item-2"),
		[]byte("batch-item-3"),
		[]byte("batch-item-4"),
	}

	// Insert items
	f.InsertBatch(items)

	// Batch lookup
	results := f.LookupBatch(items)

	if len(results) != len(items) {
		t.Errorf("Expected %d results, got %d", len(items), len(results))
	}

	for i, found := range results {
		if !found {
			t.Errorf("Batch lookup failed to find item %d", i)
		}
	}
}

// TestFilterBatchDelete tests batch delete operation
func TestFilterBatchDelete(t *testing.T) {
	f, _ := New(1000, 4, 8, 500, hash.HashStrategyXXHash, 32)

	items := [][]byte{
		[]byte("batch-delete-1"),
		[]byte("batch-delete-2"),
		[]byte("batch-delete-3"),
		[]byte("batch-delete-4"),
	}

	// Insert items
	f.InsertBatch(items)

	initialCount := f.Count()

	// Batch delete
	results := f.DeleteBatch(items)

	if len(results) != len(items) {
		t.Errorf("Expected %d results, got %d", len(items), len(results))
	}

	for i, deleted := range results {
		if !deleted {
			t.Errorf("Batch delete failed for item %d", i)
		}
	}

	if f.Count() != initialCount-uint(len(items)) {
		t.Errorf("Expected count=%d after batch delete, got %d", initialCount-uint(len(items)), f.Count())
	}
}

// TestFilterOptimalBatchSize tests optimal batch size reporting
func TestFilterOptimalBatchSize(t *testing.T) {
	f, _ := New(1000, 4, 8, 500, hash.HashStrategyXXHash, 32)

	batchSize := f.OptimalBatchSize()
	if batchSize != 32 {
		t.Errorf("Expected optimal batch size=32, got %d", batchSize)
	}
}

// TestFilterWithDifferentHashStrategies tests filter with different hash strategies
func TestFilterWithDifferentHashStrategies(t *testing.T) {
	strategies := []hash.HashStrategy{
		hash.HashStrategyXXHash,
		hash.HashStrategyCRC32,
		hash.HashStrategyFNV,
	}

	for _, strategy := range strategies {
		t.Run(strategy.String(), func(t *testing.T) {
			f, err := New(1000, 4, 8, 500, strategy, 32)
			if err != nil {
				t.Fatalf("New failed with %s: %v", strategy, err)
			}

			item := []byte(fmt.Sprintf("test-%s", strategy))

			if !f.Insert(item) {
				t.Errorf("Insert failed with %s", strategy)
			}

			if !f.Lookup(item) {
				t.Errorf("Lookup failed with %s", strategy)
			}
		})
	}
}

// TestFilterWithDifferentBucketSizes tests filter with different bucket sizes
func TestFilterWithDifferentBucketSizes(t *testing.T) {
	sizes := []uint{4, 8, 16, 32, 64}

	for _, size := range sizes {
		t.Run(fmt.Sprintf("bucket-%d", size), func(t *testing.T) {
			f, err := New(1000, size, 8, 500, hash.HashStrategyXXHash, 32)
			if err != nil {
				t.Fatalf("New failed with bucket size %d: %v", size, err)
			}

			if f.bucketSize != size {
				t.Errorf("Expected bucket size %d, got %d", size, f.bucketSize)
			}

			// Insert items
			for i := 0; i < 10; i++ {
				item := []byte(fmt.Sprintf("bucket-%d-item-%d", size, i))
				if !f.Insert(item) {
					t.Errorf("Insert failed for bucket size %d", size)
				}
			}
		})
	}
}

// TestFilterNextPowerOf2 tests the nextPowerOf2 helper function
func TestFilterNextPowerOf2(t *testing.T) {
	tests := []struct {
		input    uint
		expected uint
	}{
		{0, 1},
		{1, 1},
		{2, 2},
		{3, 4},
		{4, 4},
		{5, 8},
		{7, 8},
		{8, 8},
		{9, 16},
		{15, 16},
		{16, 16},
		{17, 32},
		{1000, 1024},
	}

	for _, tc := range tests {
		result := nextPowerOf2(tc.input)
		if result != tc.expected {
			t.Errorf("nextPowerOf2(%d) = %d, expected %d", tc.input, result, tc.expected)
		}
	}
}

// TestFilterNumBucketsCalculation tests bucket count calculation
func TestFilterNumBucketsCalculation(t *testing.T) {
	tests := []struct {
		capacity   uint
		bucketSize uint
		minBuckets uint
	}{
		{100, 4, 25},   // 100/4 = 25, next power of 2 = 32
		{1000, 4, 250}, // 1000/4 = 250, next power of 2 = 256
		{1000, 8, 125}, // 1000/8 = 125, next power of 2 = 128
	}

	for _, tc := range tests {
		f, _ := New(tc.capacity, tc.bucketSize, 8, 500, hash.HashStrategyXXHash, 32)
		if f.numBuckets < tc.minBuckets {
			t.Errorf("Expected at least %d buckets, got %d", tc.minBuckets, f.numBuckets)
		}
		// Should be power of 2
		if f.numBuckets&(f.numBuckets-1) != 0 {
			t.Errorf("numBuckets %d is not a power of 2", f.numBuckets)
		}
	}
}
