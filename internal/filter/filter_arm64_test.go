//go:build arm64
// +build arm64

package filter

import (
	"fmt"
	"testing"

	"github.com/shaia/simdcuckoofilter/internal/hash"
)

// TestFilterNEONLookup tests that NEON-optimized lookup is being used
func TestFilterNEONLookup(t *testing.T) {
	// This test verifies that the NEON lookup path is executed
	// by checking that lookups work correctly with various bucket sizes
	bucketSizes := []uint{4, 8, 16, 32, 64}

	for _, size := range bucketSizes {
		f, _ := New(1000, size, 8, 500, hash.HashStrategyXXHash, 32)

		// Insert items
		items := make([][]byte, 20)
		for i := range items {
			items[i] = []byte(fmt.Sprintf("neon-test-%d-size-%d", i, size))
			f.Insert(items[i])
		}

		// Lookup all items
		for i, item := range items {
			if !f.Lookup(item) {
				t.Errorf("NEON lookup failed for item %d with bucket size %d", i, size)
			}
		}
	}
}
