package main

import (
	"fmt"
	"log"

	"github.com/shaia/simdcuckoofilter"
)

func main() {
	// Create a new Cuckoo filter with capacity for 10,000 items
	// On ARM64 (Apple Silicon, AWS Graviton), uses optimized assembly for bucket operations
	// On AMD64 and other platforms, uses Go's compiler-optimized code
	cf, err := cuckoofilter.New(10000)
	if err != nil {
		log.Fatalf("Failed to create filter: %v", err)
	}

	fmt.Println("Created Cuckoo filter with default options")
	fmt.Printf("Capacity: %d items\n\n", cf.Capacity())

	// Insert items
	items := []string{"apple", "banana", "cherry", "date", "elderberry"}
	fmt.Println("Inserting items...")
	for _, item := range items {
		if cf.Insert([]byte(item)) {
			fmt.Printf("  ✓ Inserted: %s\n", item)
		} else {
			fmt.Printf("  ✗ Failed to insert: %s\n", item)
		}
	}

	// Lookup items
	fmt.Println("\nLooking up items...")
	for _, item := range items {
		if cf.Lookup([]byte(item)) {
			fmt.Printf("  ✓ Found: %s\n", item)
		} else {
			fmt.Printf("  ✗ Not found: %s\n", item)
		}
	}

	// Test false positives
	fmt.Println("\nTesting items that were not inserted...")
	notInserted := []string{"grape", "fig", "kiwi"}
	for _, item := range notInserted {
		if cf.Lookup([]byte(item)) {
			fmt.Printf("  ⚠ False positive: %s\n", item)
		} else {
			fmt.Printf("  ✓ Correctly not found: %s\n", item)
		}
	}

	// Delete an item
	fmt.Println("\nDeleting 'banana'...")
	if cf.Delete([]byte("banana")) {
		fmt.Println("  ✓ Deleted successfully")
	} else {
		fmt.Println("  ✗ Failed to delete")
	}

	// Verify deletion
	if !cf.Lookup([]byte("banana")) {
		fmt.Println("  ✓ Confirmed: 'banana' is no longer in the filter")
	}

	// Check statistics
	fmt.Println("\nFilter statistics:")
	fmt.Printf("  Count: %d items\n", cf.Count())
	fmt.Printf("  Capacity: %d items\n", cf.Capacity())
	fmt.Printf("  Load factor: %.2f%%\n", cf.LoadFactor()*100)
}
