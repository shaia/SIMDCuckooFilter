package xxhash

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/shaia/simdcuckoofilter/internal/hash/types"
)

// TestAVX2BoundaryConditions explicitly tests the boundary conditions
// for the AVX2 implementation (batch size < 4, = 4, > 4).
func TestAVX2BoundaryConditions(t *testing.T) {
	fingerprintBits := uint(8)
	numBuckets := uint(1024)

	// Generate some test data
	items := make([][]byte, 16)
	for i := range items {
		items[i] = []byte(fmt.Sprintf("item-%d", i))
	}

	// Helper to get expected results using scalar implementation
	getExpected := func(batch [][]byte) []types.HashResult {
		h := &XXHash{fingerprintBits: fingerprintBits}
		results := make([]types.HashResult, len(batch))
		for i, item := range batch {
			i1, i2, fp := h.GetIndices(item, numBuckets)
			results[i] = types.HashResult{I1: i1, I2: i2, Fp: fp}
		}
		return results
	}

	// Helper to verify results
	verify := func(t *testing.T, name string, batch [][]byte) {
		t.Helper()
		t.Run(name, func(t *testing.T) {
			// Use the batch processor (which uses AVX2 if available)
			batchProc := NewBatchHashProcessor()
			h := &XXHash{
				fingerprintBits: fingerprintBits,
				batchProcessor:  batchProc,
			}
			
			got := h.GetIndicesBatch(batch, numBuckets)
			want := getExpected(batch)

			if len(got) != len(want) {
				t.Fatalf("Length mismatch: got %d, want %d", len(got), len(want))
			}

			for i := range got {
				if got[i] != want[i] {
					t.Errorf("Mismatch at index %d:\nGot:  %+v\nWant: %+v", i, got[i], want[i])
				}
			}
		})
	}

	// Case 1: Scalar Path (items < 4)
	// The assembly code checks: SUBQ $4, R15; JL scalar_loop
	verify(t, "ScalarPath_1Item", items[:1])
	verify(t, "ScalarPath_3Items", items[:3])

	// Case 2: Exact AVX2 Path (items = 4)
	// The assembly code should enter simd_loop and process one batch
	verify(t, "AVX2Path_Exact4Items", items[:4])

	// Case 3: AVX2 + Scalar Remainder (items = 5)
	// The assembly code should process 4 items in simd_loop, then 1 in scalar_loop
	verify(t, "AVX2Path_5Items_Remainder", items[:5])

	// Case 4: Multiple AVX2 Batches (items = 8)
	// The assembly code should process two iterations of simd_loop
	verify(t, "AVX2Path_MultipleBatches_8Items", items[:8])

	// Case 5: Multiple AVX2 Batches + Remainder (items = 9)
	verify(t, "AVX2Path_MultipleBatches_Remainder_9Items", items[:9])
}

// BenchmarkAVX2Boundaries benchmarks the performance around the threshold of 4 items.
// This helps confirm that the AVX2 path is actually being used (should see a jump in throughput).
func BenchmarkAVX2Boundaries(b *testing.B) {
	fingerprintBits := uint(8)
	numBuckets := uint(1024)
	
	// Setup processor
	batchProc := NewBatchHashProcessor()
	h := &XXHash{
		fingerprintBits: fingerprintBits,
		batchProcessor:  batchProc,
	}

	// Create a large pool of items to avoid branch prediction optimizing too much on the same data
	poolSize := 1000
	pool := make([][]byte, poolSize)
	for i := 0; i < poolSize; i++ {
		pool[i] = bytes.Repeat([]byte{byte(i)}, 32) // 32-byte items
	}

	sizes := []int{1, 2, 3, 4, 5, 8, 16}

	for _, size := range sizes {
		b.Run(fmt.Sprintf("BatchSize_%d", size), func(b *testing.B) {
			batch := make([][]byte, size)
			b.ResetTimer()
			
			for i := 0; i < b.N; i++ {
				// Rotate through pool to simulate real workload
				start := (i * size) % (poolSize - size)
				for j := 0; j < size; j++ {
					batch[j] = pool[start+j]
				}
				
				_ = h.GetIndicesBatch(batch, numBuckets)
			}
		})
	}
}
