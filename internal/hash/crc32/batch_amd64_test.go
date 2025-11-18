//go:build amd64
// +build amd64

package crc32hash

import (
	"fmt"
	"hash/crc32"
	"testing"
)

// TestSIMDCorrectness verifies SIMD CRC32 against stdlib
func TestSIMDCorrectness(t *testing.T) {
	table := crc32.MakeTable(crc32.Castagnoli)

	testCases := []struct {
		name string
		data []byte
	}{
		{"empty", []byte("")},
		{"single byte", []byte("a")},
		{"short", []byte("hello")},
		{"medium", []byte("the quick brown fox jumps over the lazy dog")},
		{"long", func() []byte {
			data := make([]byte, 1000)
			for i := range data {
				data[i] = byte(i % 256)
			}
			return data
		}()},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Expected result from stdlib
			expected := crc32.Checksum(tc.data, table)

			// Compute using our SIMD assembly
			items := [][]byte{tc.data}
			results := make([]uint32, 1)
			batchCRC32SIMD(items, results)

			if results[0] != expected {
				t.Errorf("CRC32 mismatch: got %08x, want %08x", results[0], expected)
			}
		})
	}
}

// TestSIMDBatchSizes tests various batch sizes
func TestSIMDBatchSizes(t *testing.T) {
	table := crc32.MakeTable(crc32.Castagnoli)
	processor := NewBatchProcessor(table)

	batchSizes := []int{1, 2, 3, 4, 5, 8, 12, 16, 32, 64}

	for _, size := range batchSizes {
		t.Run(fmt.Sprintf("size-%d", size), func(t *testing.T) {
			// Generate test data
			items := make([][]byte, size)
			for i := range items {
				items[i] = []byte(fmt.Sprintf("item-%d", i))
			}

			// Process batch
			results := processor.ProcessBatch(items, 8, 1000)

			// Verify correct number of results
			if len(results) != size {
				t.Errorf("Expected %d results, got %d", size, len(results))
			}

			// Verify all results are valid (non-zero fingerprints for non-empty items)
			for i, result := range results {
				if result.Fp == 0 && len(items[i]) > 0 {
					t.Errorf("Item %d: unexpected zero fingerprint for non-empty item", i)
				}
			}
		})
	}
}

// TestSIMDEdgeCases tests edge cases
func TestSIMDEdgeCases(t *testing.T) {
	table := crc32.MakeTable(crc32.Castagnoli)
	processor := NewBatchProcessor(table)

	t.Run("all empty", func(t *testing.T) {
		items := [][]byte{[]byte(""), []byte(""), []byte(""), []byte("")}
		results := processor.ProcessBatch(items, 8, 1000)
		if len(results) != 4 {
			t.Errorf("Expected 4 results, got %d", len(results))
		}
	})

	t.Run("single large item", func(t *testing.T) {
		largeData := make([]byte, 10000)
		for i := range largeData {
			largeData[i] = byte(i % 256)
		}
		items := [][]byte{largeData}
		results := processor.ProcessBatch(items, 8, 1000)
		if len(results) != 1 {
			t.Errorf("Expected 1 result, got %d", len(results))
		}
	})

	t.Run("mixed sizes", func(t *testing.T) {
		items := [][]byte{
			[]byte(""),
			[]byte("a"),
			[]byte("ab"),
			[]byte("abc"),
			[]byte("abcd"),
			[]byte("abcde"),
			[]byte("abcdef"),
			[]byte("abcdefg"),
		}
		results := processor.ProcessBatch(items, 8, 1000)
		if len(results) != 8 {
			t.Errorf("Expected 8 results, got %d", len(results))
		}
	})
}
