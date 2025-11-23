package cuckoofilter

import (
	"fmt"
	"runtime"
	"testing"
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
		WithCRC32Hash(),
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

	_, err = New(1000, WithFingerprintSize(17)) // Invalid: too large (max is 16)
	if err != ErrInvalidFingerprintSize {
		t.Errorf("Expected ErrInvalidFingerprintSize for 17 bits, got %v", err)
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
	tests := []struct {
		name   string
		option Option
	}{
		{"XXHash64", WithXXHash()},
		{"CRC32C", WithCRC32Hash()},
		{"FNV-1a", WithFNVHash()},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cf, err := New(1000, tt.option)
			if err != nil {
				t.Fatalf("New failed with %s: %v", tt.name, err)
			}

			item := []byte(fmt.Sprintf("test-%s", tt.name))

			if !cf.Insert(item) {
				t.Errorf("Insert failed with %s", tt.name)
			}

			if !cf.Lookup(item) {
				t.Errorf("Lookup failed with %s", tt.name)
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

// TestIntegration tests a complete workflow
func TestIntegration(t *testing.T) {
	cf, err := New(1000)
	if err != nil {
		t.Fatalf("New failed: %v", err)
	}

	// Insert multiple items
	items := [][]byte{
		[]byte("integration-1"),
		[]byte("integration-2"),
		[]byte("integration-3"),
		[]byte("integration-4"),
		[]byte("integration-5"),
	}

	for _, item := range items {
		if !cf.Insert(item) {
			t.Errorf("Failed to insert item: %s", item)
		}
	}

	// Verify all items exist
	for _, item := range items {
		if !cf.Lookup(item) {
			t.Errorf("Failed to find item: %s", item)
		}
	}

	// Check count
	if cf.Count() != uint(len(items)) {
		t.Errorf("Expected count %d, got %d", len(items), cf.Count())
	}

	// Delete some items
	if !cf.Delete(items[0]) {
		t.Error("Failed to delete first item")
	}
	if !cf.Delete(items[2]) {
		t.Error("Failed to delete third item")
	}

	// Verify count updated
	expectedCount := uint(len(items) - 2)
	if cf.Count() != expectedCount {
		t.Errorf("Expected count %d after deletes, got %d", expectedCount, cf.Count())
	}

	// Reset and verify empty
	cf.Reset()
	if cf.Count() != 0 {
		t.Errorf("Expected count 0 after reset, got %d", cf.Count())
	}

	// Should not find items after reset
	for _, item := range items {
		if cf.Lookup(item) {
			t.Logf("Warning: Item %s still found after reset (false positive)", item)
		}
	}
}

// TestSIMD tests SIMD-specific filter functionality
func TestSIMD(t *testing.T) {
	// Skip on architectures without SIMD support
	if runtime.GOARCH != "amd64" && runtime.GOARCH != "arm64" {
		t.Skip("SIMD not supported on this architecture")
	}

	cf, err := New(1000)
	if err != nil {
		t.Fatalf("New failed: %v", err)
	}

	item := []byte("test-simd")
	if !cf.Insert(item) {
		t.Error("Insert failed")
	}

	if !cf.Lookup(item) {
		t.Error("Lookup failed")
	}
}

// TestBatchOperations tests SIMD batch processing
func TestBatchOperations(t *testing.T) {
	cf, _ := New(1000)

	// Type assert to BatchFilter
	bf, ok := cf.(BatchFilter)
	if !ok {
		t.Skip("Filter does not implement BatchFilter")
	}

	items := [][]byte{
		[]byte("item1"),
		[]byte("item2"),
		[]byte("item3"),
		[]byte("item4"),
		[]byte("item5"),
		[]byte("item6"),
		[]byte("item7"),
		[]byte("item8"),
	}

	// Batch insert
	results := bf.InsertBatch(items)
	for i, success := range results {
		if !success {
			t.Errorf("Batch insert failed for item %d", i)
		}
	}

	// Batch lookup
	found := bf.LookupBatch(items)
	for i, exists := range found {
		if !exists {
			t.Errorf("Batch lookup failed for item %d", i)
		}
	}
}

// TestBatchDelete tests batch delete operations
func TestBatchDelete(t *testing.T) {
	cf, _ := New(1000)

	bf, ok := cf.(BatchFilter)
	if !ok {
		t.Skip("Filter does not implement BatchFilter")
	}

	items := [][]byte{
		[]byte("delete1"),
		[]byte("delete2"),
		[]byte("delete3"),
		[]byte("delete4"),
	}

	// Insert items
	bf.InsertBatch(items)

	// Verify they exist
	found := bf.LookupBatch(items)
	for i, exists := range found {
		if !exists {
			t.Errorf("Item %d not found after batch insert", i)
		}
	}

	// Batch delete
	results := bf.DeleteBatch(items)
	for i, success := range results {
		if !success {
			t.Logf("Batch delete returned false for item %d (may not exist)", i)
		}
	}
}

// TestBatchWithDifferentHashStrategies tests batch operations with different hash algorithms
func TestBatchWithDifferentHashStrategies(t *testing.T) {
	tests := []struct {
		name   string
		option Option
	}{
		{"XXHash64", WithXXHash()},
		{"CRC32C", WithCRC32Hash()},
		{"FNV-1a", WithFNVHash()},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cf, err := New(1000, tt.option)
			if err != nil {
				t.Fatalf("New failed with %s: %v", tt.name, err)
			}

			bf, ok := cf.(BatchFilter)
			if !ok {
				t.Skip("Filter does not implement BatchFilter")
			}

			items := make([][]byte, 16)
			for i := range items {
				items[i] = []byte(fmt.Sprintf("item-%s-%d", tt.name, i))
			}

			// Batch insert
			results := bf.InsertBatch(items)
			for i, success := range results {
				if !success {
					t.Errorf("Batch insert failed for item %d with %s", i, tt.name)
				}
			}

			// Batch lookup
			found := bf.LookupBatch(items)
			for i, exists := range found {
				if !exists {
					t.Errorf("Batch lookup failed for item %d with %s", i, tt.name)
				}
			}
		})
	}
}
