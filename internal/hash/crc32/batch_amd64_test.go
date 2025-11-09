//go:build amd64
// +build amd64

package crc32hash

import (
	"fmt"
	"hash/crc32"
	"testing"

	"github.com/shaia/cuckoofilter/internal/simd/cpu"
)

// TestSIMDvsNonSIMD compares SIMD and non-SIMD CRC32 batch processing
func TestSIMDvsNonSIMD(t *testing.T) {
	table := crc32.MakeTable(crc32.Castagnoli)
	simdType := cpu.GetBestSIMD(true)

	// Create SIMD and non-SIMD processors
	simdProcessor := NewBatchProcessor(table, simdType)
	nonSIMDProcessor := NewBatchProcessorNoSIMD(table, simdType)

	testCases := []struct {
		name  string
		items [][]byte
	}{
		{
			name: "4 items",
			items: [][]byte{
				[]byte("test1"),
				[]byte("test2"),
				[]byte("test3"),
				[]byte("test4"),
			},
		},
		{
			name: "8 items",
			items: [][]byte{
				[]byte("item1"),
				[]byte("item2"),
				[]byte("item3"),
				[]byte("item4"),
				[]byte("item5"),
				[]byte("item6"),
				[]byte("item7"),
				[]byte("item8"),
			},
		},
		{
			name: "16 items with varying lengths",
			items: func() [][]byte {
				items := make([][]byte, 16)
				for i := range items {
					items[i] = []byte(fmt.Sprintf("variable-length-item-%d-with-more-data", i))
				}
				return items
			}(),
		},
		{
			name: "Empty items",
			items: [][]byte{
				[]byte(""),
				[]byte("a"),
				[]byte(""),
				[]byte("bc"),
			},
		},
		{
			name: "Large items",
			items: func() [][]byte {
				items := make([][]byte, 4)
				for i := range items {
					data := make([]byte, 1024)
					for j := range data {
						data[j] = byte((i * j) % 256)
					}
					items[i] = data
				}
				return items
			}(),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Process with SIMD
			simdResults := simdProcessor.ProcessBatch(tc.items, 8, 1000)

			// Process without SIMD
			nonSIMDResults := nonSIMDProcessor.ProcessBatch(tc.items, 8, 1000)

			// Compare results
			if len(simdResults) != len(nonSIMDResults) {
				t.Fatalf("Result length mismatch: SIMD=%d, non-SIMD=%d", len(simdResults), len(nonSIMDResults))
			}

			for i := range simdResults {
				if simdResults[i].I1 != nonSIMDResults[i].I1 {
					t.Errorf("Item %d: I1 mismatch: SIMD=%d, non-SIMD=%d", i, simdResults[i].I1, nonSIMDResults[i].I1)
				}
				if simdResults[i].I2 != nonSIMDResults[i].I2 {
					t.Errorf("Item %d: I2 mismatch: SIMD=%d, non-SIMD=%d", i, simdResults[i].I2, nonSIMDResults[i].I2)
				}
				if simdResults[i].Fp != nonSIMDResults[i].Fp {
					t.Errorf("Item %d: Fp mismatch: SIMD=%d, non-SIMD=%d", i, simdResults[i].Fp, nonSIMDResults[i].Fp)
				}
			}
		})
	}
}

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
	simdProcessor := NewBatchProcessor(table, cpu.GetBestSIMD(true))
	nonSIMDProcessor := NewBatchProcessorNoSIMD(table, cpu.GetBestSIMD(true))

	batchSizes := []int{1, 2, 3, 4, 5, 8, 12, 16, 32, 64}

	for _, size := range batchSizes {
		t.Run(fmt.Sprintf("size-%d", size), func(t *testing.T) {
			// Generate test data
			items := make([][]byte, size)
			for i := range items {
				items[i] = []byte(fmt.Sprintf("item-%d", i))
			}

			// Process with both methods
			simdResults := simdProcessor.ProcessBatch(items, 8, 1000)
			nonSIMDResults := nonSIMDProcessor.ProcessBatch(items, 8, 1000)

			// Verify they match
			for i := range simdResults {
				if simdResults[i] != nonSIMDResults[i] {
					t.Errorf("Item %d mismatch: SIMD=%+v, non-SIMD=%+v", i, simdResults[i], nonSIMDResults[i])
				}
			}
		})
	}
}

// TestSIMDEdgeCases tests edge cases
func TestSIMDEdgeCases(t *testing.T) {
	table := crc32.MakeTable(crc32.Castagnoli)
	simdProcessor := NewBatchProcessor(table, cpu.GetBestSIMD(true))

	t.Run("all empty", func(t *testing.T) {
		items := [][]byte{[]byte(""), []byte(""), []byte(""), []byte("")}
		results := simdProcessor.ProcessBatch(items, 8, 1000)
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
		results := simdProcessor.ProcessBatch(items, 8, 1000)
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
		results := simdProcessor.ProcessBatch(items, 8, 1000)
		if len(results) != 8 {
			t.Errorf("Expected 8 results, got %d", len(results))
		}
	})
}
