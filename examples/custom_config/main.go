package main

import (
	"fmt"
	"log"

	"github.com/shaia/simdcuckoofilter"
)

func main() {
	// Create filter with custom options
	// Note: Fingerprint size is limited to 8 bits (1 byte) in this implementation
	cf, err := cuckoofilter.New(10000,
		cuckoofilter.WithFingerprintSize(8),   // 8-bit fingerprints (maximum supported)
		cuckoofilter.WithBucketSize(32),       // 32 fingerprints per bucket (optimal for AVX2/NEON SIMD)
		cuckoofilter.WithMaxKicks(500),        // Standard relocation limit
	)
	if err != nil {
		log.Fatalf("Failed to create filter: %v", err)
	}

	fmt.Println("Created Cuckoo filter with custom options:")
	fmt.Printf("  Fingerprint size: 8 bits (1 byte)\n")
	fmt.Printf("  Bucket size: 32 fingerprints (optimal for AMD64 AVX2 and ARM64 NEON)\n")
	fmt.Printf("  Max kicks: 500 (relocation attempts)\n")
	fmt.Printf("  Capacity: %d items\n\n", cf.Capacity())

	// Insert many items
	numItems := 5000
	fmt.Printf("Inserting %d items...\n", numItems)
	insertCount := 0
	for i := 0; i < numItems; i++ {
		item := fmt.Sprintf("item-%d", i)
		if cf.Insert([]byte(item)) {
			insertCount++
		}
	}
	fmt.Printf("  Successfully inserted: %d/%d items\n", insertCount, numItems)
	fmt.Printf("  Load factor: %.2f%%\n\n", cf.LoadFactor()*100)

	// Test for false positives
	fmt.Println("Testing for false positives...")
	numTests := 10000
	falsePositives := 0
	for i := numItems; i < numItems+numTests; i++ {
		item := fmt.Sprintf("item-%d", i)
		if cf.Lookup([]byte(item)) {
			falsePositives++
		}
	}

	fpr := float64(falsePositives) / float64(numTests) * 100
	fmt.Printf("  Tested %d non-inserted items\n", numTests)
	fmt.Printf("  False positives: %d\n", falsePositives)
	fmt.Printf("  Measured FPR: %.4f%%\n", fpr)
	fmt.Printf("  Expected FPR: ~0.39%% (1/2^8 for 8-bit fingerprints)\n")

	if fpr < 1.0 {
		fmt.Println("  ✓ False positive rate is within expected range!")
	} else if fpr < 2.0 {
		fmt.Println("  ✓ False positive rate is acceptable")
	} else {
		fmt.Println("  ⚠ False positive rate is higher than expected")
	}
}
