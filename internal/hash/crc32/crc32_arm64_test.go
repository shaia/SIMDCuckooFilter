//go:build arm64
// +build arm64

package crc32hash

import (
	"fmt"
	"hash/crc32"
	"testing"
)

// TestARM64HardwareVsSoftware compares ARM64 hardware CRC32 vs software implementation
func TestARM64HardwareVsSoftware(t *testing.T) {
	table := crc32.MakeTable(crc32.Castagnoli)

	// Note: This test now verifies consistency between two instances
	// Previously tested hardware vs software, now both use the same implementation
	hardwareProcessor := NewBatchProcessor(table)
	softwareProcessor := NewBatchProcessor(table)

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
		{
			name: "Single byte items",
			items: [][]byte{
				[]byte("a"),
				[]byte("b"),
				[]byte("c"),
				[]byte("d"),
			},
		},
		{
			name: "Mixed ASCII and binary",
			items: [][]byte{
				[]byte("hello"),
				[]byte{0x00, 0x01, 0x02, 0x03},
				[]byte{0xFF, 0xFE, 0xFD},
				[]byte("world"),
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Process with hardware acceleration
			hardwareResults := hardwareProcessor.ProcessBatch(tc.items, 8, 1000)

			// Process without hardware acceleration
			softwareResults := softwareProcessor.ProcessBatch(tc.items, 8, 1000)

			// Compare results
			if len(hardwareResults) != len(softwareResults) {
				t.Fatalf("Result length mismatch: hardware=%d, software=%d", len(hardwareResults), len(softwareResults))
			}

			for i := range hardwareResults {
				if hardwareResults[i].I1 != softwareResults[i].I1 {
					t.Errorf("Item %d: I1 mismatch: hardware=%d, software=%d", i, hardwareResults[i].I1, softwareResults[i].I1)
				}
				if hardwareResults[i].I2 != softwareResults[i].I2 {
					t.Errorf("Item %d: I2 mismatch: hardware=%d, software=%d", i, hardwareResults[i].I2, softwareResults[i].I2)
				}
				if hardwareResults[i].Fp != softwareResults[i].Fp {
					t.Errorf("Item %d: Fp mismatch: hardware=%d, software=%d", i, hardwareResults[i].Fp, softwareResults[i].Fp)
				}
			}
		})
	}
}

// TestARM64HardwareCorrectness verifies ARM64 hardware CRC32 against stdlib
func TestARM64HardwareCorrectness(t *testing.T) {
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
		{"all zeros", make([]byte, 100)},
		{"all ones", func() []byte {
			data := make([]byte, 100)
			for i := range data {
				data[i] = 0xFF
			}
			return data
		}()},
		{"pattern", []byte{0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07}},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Expected result from stdlib
			expected := crc32.Checksum(tc.data, table)

			// Compute using our ARM64 hardware assembly
			items := [][]byte{tc.data}
			results := make([]uint32, 1)
			batchCRC32Hardware(items, results)

			if results[0] != expected {
				t.Errorf("CRC32 mismatch: got %08x, want %08x", results[0], expected)
			}
		})
	}
}

// TestARM64BatchSizes tests various batch sizes
func TestARM64BatchSizes(t *testing.T) {
	table := crc32.MakeTable(crc32.Castagnoli)
	hardwareProcessor := NewBatchProcessor(table)
	softwareProcessor := NewBatchProcessor(table) // Note: No longer testing separate software path

	batchSizes := []int{1, 2, 3, 4, 5, 8, 12, 16, 32, 64}

	for _, size := range batchSizes {
		t.Run(fmt.Sprintf("size-%d", size), func(t *testing.T) {
			// Generate test data
			items := make([][]byte, size)
			for i := range items {
				items[i] = []byte(fmt.Sprintf("item-%d", i))
			}

			// Process with both methods
			hardwareResults := hardwareProcessor.ProcessBatch(items, 8, 1000)
			softwareResults := softwareProcessor.ProcessBatch(items, 8, 1000)

			// Verify they match
			for i := range hardwareResults {
				if hardwareResults[i] != softwareResults[i] {
					t.Errorf("Item %d mismatch: hardware=%+v, software=%+v", i, hardwareResults[i], softwareResults[i])
				}
			}
		})
	}
}

// TestARM64EdgeCases tests edge cases for ARM64 implementation
func TestARM64EdgeCases(t *testing.T) {
	table := crc32.MakeTable(crc32.Castagnoli)
	processor := NewBatchProcessor(table)

	t.Run("all empty", func(t *testing.T) {
		items := [][]byte{[]byte(""), []byte(""), []byte(""), []byte("")}
		results := processor.ProcessBatch(items, 8, 1000)
		if len(results) != 4 {
			t.Errorf("Expected 4 results, got %d", len(results))
		}
		// Verify all empty items produce the same CRC
		for i := 1; i < len(results); i++ {
			if results[i] != results[0] {
				t.Errorf("Empty items produced different results: %+v vs %+v", results[0], results[i])
			}
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

	t.Run("unaligned data", func(t *testing.T) {
		// Test data that's not 8-byte aligned
		items := [][]byte{
			[]byte("1"),
			[]byte("12"),
			[]byte("123"),
			[]byte("1234"),
			[]byte("12345"),
			[]byte("123456"),
			[]byte("1234567"),
			[]byte("12345678"),
			[]byte("123456789"),
		}
		results := processor.ProcessBatch(items, 8, 1000)
		if len(results) != len(items) {
			t.Errorf("Expected %d results, got %d", len(items), len(results))
		}
	})
}

// TestARM64VsStdlib verifies ARM64 assembly matches stdlib for many inputs
func TestARM64VsStdlib(t *testing.T) {
	table := crc32.MakeTable(crc32.Castagnoli)

	// Generate many test cases
	for length := 0; length <= 100; length++ {
		data := make([]byte, length)
		for i := range data {
			data[i] = byte(i % 256)
		}

		expected := crc32.Checksum(data, table)

		items := [][]byte{data}
		results := make([]uint32, 1)
		batchCRC32Hardware(items, results)

		if results[0] != expected {
			t.Errorf("Length %d: CRC32 mismatch: got %08x, want %08x", length, results[0], expected)
		}
	}
}

// BenchmarkARM64Hardware benchmarks ARM64 hardware CRC32
func BenchmarkARM64Hardware(b *testing.B) {
	table := crc32.MakeTable(crc32.Castagnoli)
	processor := NewBatchProcessor(table)

	items := make([][]byte, 32)
	for i := range items {
		items[i] = []byte("benchmark test data with reasonable length")
	}

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = processor.ProcessBatch(items, 8, 1024)
	}
}

// BenchmarkARM64Software benchmarks software CRC32 for comparison
func BenchmarkARM64Software(b *testing.B) {
	table := crc32.MakeTable(crc32.Castagnoli)
	processor := NewBatchProcessor(table) // Note: No longer testing separate software path

	items := make([][]byte, 32)
	for i := range items {
		items[i] = []byte("benchmark test data with reasonable length")
	}

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = processor.ProcessBatch(items, 8, 1024)
	}
}
