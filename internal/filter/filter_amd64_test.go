//go:build amd64
// +build amd64

package filter

import (
	"fmt"
	"testing"

	"github.com/shaia/simdcuckoofilter/internal/hash"
)

// TestFilterAVX2Lookup tests that AVX2-optimized lookup is being used
func TestFilterAVX2Lookup(t *testing.T) {
	// This test verifies that the AVX2 lookup path is executed
	// by checking that lookups work correctly with various bucket sizes
	bucketSizes := []uint{4, 8, 16, 32, 64}

	for _, size := range bucketSizes {
		f, _ := New(1000, size, 8, 500, hash.HashStrategyXXHash, 32)

		// Insert items
		items := make([][]byte, 20)
		for i := range items {
			items[i] = []byte(fmt.Sprintf("avx2-test-%d-size-%d", i, size))
			f.Insert(items[i])
		}

		// Lookup all items
		for i, item := range items {
			if !f.Lookup(item) {
				t.Errorf("AVX2 lookup failed for item %d with bucket size %d", i, size)
			}
		}
	}
}
