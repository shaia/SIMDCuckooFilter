package xxhash

import (
	"testing"

	"github.com/shaia/cuckoofilter/internal/hash/types"
	"github.com/shaia/cuckoofilter/internal/simd/cpu"
)

// TestBatchConsistencyAllSizes verifies that SIMD batch processing produces
// identical results to scalar processing for all batch sizes.
//
// This test ensures consistency between SIMD and scalar code paths, preventing
// regressions where SIMD implementations diverge from the reference implementation.
func TestBatchConsistencyAllSizes(t *testing.T) {
	fingerprintBits := uint(8)
	numBuckets := uint(2048)

	// Create test data of various sizes
	testData := [][]byte{
		{0x01},
		{0x02, 0x03},
		{0x04, 0x05, 0x06},
		[]byte("test"),
		[]byte("hello world"),
		[]byte("the quick brown fox"),
		make([]byte, 100),
	}

	// Initialize pattern in the 100-byte array
	for i := range testData[len(testData)-1] {
		testData[len(testData)-1][i] = byte(i)
	}

	// Reference implementation (scalar)
	refHash := &XXHash{
		fingerprintBits: fingerprintBits,
		batchProcessor:  nil,
	}

	// Test all batch sizes from 1 to 16
	for batchSize := 1; batchSize <= 16; batchSize++ {
		t.Run(string(rune('0'+batchSize/10))+string(rune('0'+batchSize%10))+"_items", func(t *testing.T) {
			// Create batch
			batch := make([][]byte, batchSize)
			for i := 0; i < batchSize; i++ {
				batch[i] = testData[i%len(testData)]
			}

			// Get reference results
			refResults := make([]types.HashResult, batchSize)
			for i, item := range batch {
				i1, i2, fp := refHash.GetIndices(item, numBuckets)
				refResults[i] = types.HashResult{I1: i1, I2: i2, Fp: fp}
			}

			// Test with best available SIMD implementation
			bestSIMD := cpu.GetBestSIMD(true)
			if bestSIMD != cpu.SIMDNone {
				t.Run(bestSIMD.String(), func(t *testing.T) {
					batchProc := NewBatchHashProcessor()
					simdHash := &XXHash{
						fingerprintBits: fingerprintBits,
						batchProcessor:  batchProc,
					}

					results := simdHash.GetIndicesBatch(batch, numBuckets)

					if len(results) != len(refResults) {
						t.Fatalf("Expected %d results, got %d", len(refResults), len(results))
					}

					for i := 0; i < len(results); i++ {
						if results[i] != refResults[i] {
							t.Errorf("Batch item %d mismatch:\n"+
								"  Input: %v\n"+
								"  Expected: i1=%d, i2=%d, fp=%d\n"+
								"  Got:      i1=%d, i2=%d, fp=%d",
								i, batch[i],
								refResults[i].I1, refResults[i].I2, refResults[i].Fp,
								results[i].I1, results[i].I2, results[i].Fp)
						}
					}
				})
			}
		})
	}
}

// TestSIMDScalarFallback verifies that SIMD implementations correctly fall back
// to scalar processing when batch size is less than the SIMD width.
//
// This ensures that edge cases (batch size < 4 for AVX2) are handled correctly.
func TestSIMDScalarFallback(t *testing.T) {
	fingerprintBits := uint(8)
	numBuckets := uint(1024)

	testData := [][]byte{
		{0x42},
		{0x43, 0x44},
		{0x45, 0x46, 0x47},
	}

	// Reference
	refHash := &XXHash{
		fingerprintBits: fingerprintBits,
		batchProcessor:  nil,
	}

	refResults := make([]types.HashResult, len(testData))
	for i, item := range testData {
		i1, i2, fp := refHash.GetIndices(item, numBuckets)
		refResults[i] = types.HashResult{I1: i1, I2: i2, Fp: fp}
	}

	// Test with best available SIMD
	bestSIMD := cpu.GetBestSIMD(true)
	if bestSIMD == cpu.SIMDNone {
		t.Skip("No SIMD support available")
	}

	// Test various batch sizes
	batchSizes := []int{1, 2, 3}
	for _, batchSize := range batchSizes {
		t.Run(bestSIMD.String()+"_"+string(rune('0'+batchSize))+"_items", func(t *testing.T) {
			batchProc := NewBatchHashProcessor()
			simdHash := &XXHash{
				fingerprintBits: fingerprintBits,
				batchProcessor:  batchProc,
			}

			batch := testData[:batchSize]
			results := simdHash.GetIndicesBatch(batch, numBuckets)

			if len(results) != batchSize {
				t.Fatalf("Expected %d results, got %d", batchSize, len(results))
			}

			for i := 0; i < batchSize; i++ {
				if results[i] != refResults[i] {
					t.Errorf("Item %d mismatch:\n"+
						"  Expected: i1=%d, i2=%d, fp=%d\n"+
						"  Got:      i1=%d, i2=%d, fp=%d",
						i,
						refResults[i].I1, refResults[i].I2, refResults[i].Fp,
						results[i].I1, results[i].I2, results[i].Fp)
				}
			}
		})
	}
}

// TestBatchProcessingMemorySafety tests that batch processing doesn't cause
// memory corruption or stack overflow issues.
//
// This test verifies that the stack frame sizes are correct and that
// processing large batches doesn't corrupt memory.
//
// Regression test for: Stack overflow bug where AVX2 implementation wrote
// to offset 128(SP) with only a 128-byte stack frame.
func TestBatchProcessingMemorySafety(t *testing.T) {
	fingerprintBits := uint(8)
	numBuckets := uint(4096)

	// Create large batch
	batchSize := 1024
	batch := make([][]byte, batchSize)
	for i := 0; i < batchSize; i++ {
		data := make([]byte, 1+(i%100))
		for j := range data {
			data[j] = byte(i + j)
		}
		batch[i] = data
	}

	// Reference results
	refHash := &XXHash{
		fingerprintBits: fingerprintBits,
		batchProcessor:  nil,
	}

	refResults := make([]types.HashResult, batchSize)
	for i, item := range batch {
		i1, i2, fp := refHash.GetIndices(item, numBuckets)
		refResults[i] = types.HashResult{I1: i1, I2: i2, Fp: fp}
	}

	// Test with best available SIMD
	bestSIMD := cpu.GetBestSIMD(true)
	if bestSIMD == cpu.SIMDNone {
		t.Skip("No SIMD support available")
	}

	t.Run(bestSIMD.String(), func(t *testing.T) {
		batchProc := NewBatchHashProcessor()
		simdHash := &XXHash{
			fingerprintBits: fingerprintBits,
			batchProcessor:  batchProc,
		}

		results := simdHash.GetIndicesBatch(batch, numBuckets)

		if len(results) != batchSize {
			t.Fatalf("Expected %d results, got %d", batchSize, len(results))
		}

		// Check random samples to verify correctness
		samples := []int{0, 1, 10, 100, 500, batchSize - 1}
		for _, idx := range samples {
			if results[idx] != refResults[idx] {
				t.Errorf("Sample %d mismatch:\n"+
					"  Expected: i1=%d, i2=%d, fp=%d\n"+
					"  Got:      i1=%d, i2=%d, fp=%d",
					idx,
					refResults[idx].I1, refResults[idx].I2, refResults[idx].Fp,
					results[idx].I1, results[idx].I2, results[idx].Fp)
			}
		}
	})
}
