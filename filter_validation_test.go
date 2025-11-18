package cuckoofilter

import (
	"fmt"
	"testing"

	"github.com/shaia/simdcuckoofilter/internal/hash"
)

// TestNew validates basic filter creation
func TestNew(t *testing.T) {
	cf, err := New(1000)
	if err != nil {
		t.Fatalf("New failed: %v", err)
	}

	if cf == nil {
		t.Fatal("Filter is nil")
	}

	if cf.Count() != 0 {
		t.Errorf("Expected count 0, got %d", cf.Count())
	}

	if cf.Capacity() == 0 {
		t.Error("Capacity should be greater than 0")
	}
}

// TestInsertAndLookup validates basic insert and lookup operations
func TestInsertAndLookup(t *testing.T) {
	cf, _ := New(1000)

	item := []byte("test-item")

	// Insert
	if !cf.Insert(item) {
		t.Fatal("Insert failed")
	}

	// Lookup should find it
	if !cf.Lookup(item) {
		t.Error("Lookup failed to find inserted item")
	}

	// Lookup non-existent item
	if cf.Lookup([]byte("non-existent")) {
		t.Log("False positive (expected occasionally)")
	}
}

// TestDelete validates delete operations
func TestDelete(t *testing.T) {
	cf, _ := New(1000)

	item := []byte("test-item-for-deletion")

	// Insert
	cf.Insert(item)

	// Verify it exists
	if !cf.Lookup(item) {
		t.Error("Item not found after insert")
	}

	// Delete
	if !cf.Delete(item) {
		t.Error("Delete failed")
	}

	// Should not be found after delete
	// Note: Due to hash collisions and fingerprint false positives,
	// there's a small chance the item might still appear to exist
	if cf.Lookup(item) {
		t.Log("Warning: Item still found after delete (false positive)")
	}
}

// TestCount validates item counting
func TestCount(t *testing.T) {
	cf, _ := New(1000)

	if cf.Count() != 0 {
		t.Errorf("Initial count should be 0, got %d", cf.Count())
	}

	// Insert items
	for i := 0; i < 10; i++ {
		cf.Insert([]byte(fmt.Sprintf("item-%d", i)))
	}

	if cf.Count() != 10 {
		t.Errorf("Expected count 10, got %d", cf.Count())
	}
}

// TestReset validates reset functionality
func TestReset(t *testing.T) {
	cf, _ := New(1000)

	// Insert items
	for i := 0; i < 10; i++ {
		cf.Insert([]byte(fmt.Sprintf("item-%d", i)))
	}

	if cf.Count() == 0 {
		t.Error("Count should not be 0 after inserts")
	}

	// Reset
	cf.Reset()

	if cf.Count() != 0 {
		t.Errorf("Count should be 0 after reset, got %d", cf.Count())
	}

	if cf.LoadFactor() != 0 {
		t.Errorf("Load factor should be 0 after reset, got %f", cf.LoadFactor())
	}
}

// TestWithOptions validates filter creation with custom options
func TestWithOptions(t *testing.T) {
	cf, err := New(1000,
		WithFingerprintSize(8), // Maximum supported fingerprint size
		WithHashStrategy(hash.HashStrategyCRC32),
		WithMaxKicks(500),
	)

	if err != nil {
		t.Fatalf("New with options failed: %v", err)
	}

	item := []byte("test")
	if !cf.Insert(item) {
		t.Error("Insert failed")
	}

	if !cf.Lookup(item) {
		t.Error("Lookup failed")
	}
}

// TestInvalidCapacity validates capacity validation
func TestInvalidCapacity(t *testing.T) {
	_, err := New(0)
	if err != ErrInvalidCapacity {
		t.Errorf("Expected ErrInvalidCapacity, got %v", err)
	}
}

// TestInvalidOptions validates option validation
func TestInvalidOptions(t *testing.T) {
	_, err := New(1000, WithBucketSize(3)) // Invalid bucket size
	if err != ErrInvalidBucketSize {
		t.Errorf("Expected ErrInvalidBucketSize, got %v", err)
	}

	_, err = New(1000, WithFingerprintSize(0)) // Invalid: too small
	if err != ErrInvalidFingerprintSize {
		t.Errorf("Expected ErrInvalidFingerprintSize for 0 bits, got %v", err)
	}

	_, err = New(1000, WithFingerprintSize(9)) // Invalid: too large (max is 8)
	if err != ErrInvalidFingerprintSize {
		t.Errorf("Expected ErrInvalidFingerprintSize for 9 bits, got %v", err)
	}
}

// TestLoadFactor validates load factor calculation
func TestLoadFactor(t *testing.T) {
	cf, _ := New(1000)

	// Should start at 0
	if cf.LoadFactor() != 0 {
		t.Errorf("Initial load factor should be 0, got %f", cf.LoadFactor())
	}

	// Insert some items
	numItems := 100
	for i := 0; i < numItems; i++ {
		cf.Insert([]byte(fmt.Sprintf("load-%d", i)))
	}

	lf := cf.LoadFactor()
	if lf <= 0 || lf > 1 {
		t.Errorf("Load factor should be between 0 and 1, got %f", lf)
	}

	// Load factor should increase with more items
	expectedLF := float64(numItems) / float64(cf.Capacity())
	if lf < expectedLF*0.5 || lf > expectedLF*1.5 {
		t.Errorf("Load factor %f seems incorrect for %d items (expected ~%f)", lf, numItems, expectedLF)
	}
}

// TestMultipleInsertsSameItem validates handling of duplicate insertions
func TestMultipleInsertsSameItem(t *testing.T) {
	cf, _ := New(1000)

	item := []byte("duplicate-item")

	// First insert should succeed
	if !cf.Insert(item) {
		t.Error("First insert failed")
	}

	// Second insert of same item - behavior may vary
	// (some implementations allow duplicates, some don't)
	cf.Insert(item)

	// Should still be found
	if !cf.Lookup(item) {
		t.Error("Item not found after duplicate insert")
	}
}

// TestEmptyItem validates handling of empty items
func TestEmptyItem(t *testing.T) {
	cf, _ := New(1000)

	emptyItem := []byte("")

	// Insert empty item
	result := cf.Insert(emptyItem)
	t.Logf("Insert empty item returned: %v", result)

	// Lookup empty item
	found := cf.Lookup(emptyItem)
	t.Logf("Lookup empty item returned: %v", found)
}

// TestLargeItems validates handling of large items
func TestLargeItems(t *testing.T) {
	cf, _ := New(1000)

	// Create a large item (1KB)
	largeItem := make([]byte, 1024)
	for i := range largeItem {
		largeItem[i] = byte(i % 256)
	}

	if !cf.Insert(largeItem) {
		t.Error("Insert large item failed")
	}

	if !cf.Lookup(largeItem) {
		t.Error("Lookup large item failed")
	}
}

// TestAllHashStrategies validates all supported hash strategies
func TestAllHashStrategies(t *testing.T) {
	strategies := []hash.HashStrategy{
		hash.HashStrategyXXHash,
		hash.HashStrategyCRC32,
		hash.HashStrategyFNV,
	}

	for _, strategy := range strategies {
		t.Run(strategy.String(), func(t *testing.T) {
			cf, err := New(1000, WithHashStrategy(strategy))
			if err != nil {
				t.Fatalf("New failed with %s: %v", strategy, err)
			}

			item := []byte(fmt.Sprintf("test-%s", strategy))

			if !cf.Insert(item) {
				t.Errorf("Insert failed with %s", strategy)
			}

			if !cf.Lookup(item) {
				t.Errorf("Lookup failed with %s", strategy)
			}
		})
	}
}

// TestAllFingerprintSizes validates all supported fingerprint sizes
func TestAllFingerprintSizes(t *testing.T) {
	// Fingerprints are stored as bytes, so only 1-8 bits are supported
	sizes := []uint{4, 8}

	for _, size := range sizes {
		t.Run(fmt.Sprintf("%d-bit", size), func(t *testing.T) {
			cf, err := New(1000, WithFingerprintSize(size))
			if err != nil {
				t.Fatalf("New failed with %d-bit fingerprint: %v", size, err)
			}

			item := []byte(fmt.Sprintf("test-%d-bit", size))

			if !cf.Insert(item) {
				t.Errorf("Insert failed with %d-bit fingerprint", size)
			}

			if !cf.Lookup(item) {
				t.Errorf("Lookup failed with %d-bit fingerprint", size)
			}
		})
	}
}

// TestAllBucketSizes validates all supported bucket sizes
func TestAllBucketSizes(t *testing.T) {
	sizes := []uint{4, 8, 16, 32, 64}

	for _, size := range sizes {
		t.Run(fmt.Sprintf("bucket-%d", size), func(t *testing.T) {
			cf, err := New(1000, WithBucketSize(size))
			if err != nil {
				t.Fatalf("New failed with bucket size %d: %v", size, err)
			}

			// Insert items up to bucket size
			for i := uint(0); i < size; i++ {
				item := []byte(fmt.Sprintf("bucket-%d-item-%d", size, i))
				if !cf.Insert(item) {
					t.Errorf("Insert failed for item %d with bucket size %d", i, size)
				}
			}
		})
	}
}
