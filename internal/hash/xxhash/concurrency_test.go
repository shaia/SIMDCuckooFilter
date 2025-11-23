package xxhash

import (
	"sync"
	"testing"
)

// TestConcurrentGetAltIndex verifies that GetAltIndex is safe for concurrent use.
// This test prevents regression of the data race bug where fpBuf was shared state.
func TestConcurrentGetAltIndex(t *testing.T) {
	xxh := NewXXHash(8, nil)
	numBuckets := uint(1024)

	// Run many goroutines concurrently calling GetAltIndex
	const numGoroutines = 100
	const iterations = 1000

	var wg sync.WaitGroup
	wg.Add(numGoroutines)

	// Channel to collect any errors
	errors := make(chan error, numGoroutines)

	for g := 0; g < numGoroutines; g++ {
		go func(id int) {
			defer wg.Done()

			for i := 0; i < iterations; i++ {
				index := uint(i % 500)
				fp := uint16((id + i) % 65535)
				if fp == 0 {
					fp = 1
				}

				// Call GetAltIndex - this would race on fpBuf if not fixed
				i2 := xxh.GetAltIndex(index, fp, numBuckets)

				// Verify result is within bounds
				if i2 >= numBuckets {
					errors <- nil // Signal error without details
					return
				}

				// Verify symmetry property
				i1Back := xxh.GetAltIndex(i2, fp, numBuckets)
				if i1Back != index {
					errors <- nil // Signal error
					return
				}
			}
		}(g)
	}

	wg.Wait()
	close(errors)

	// Check if any goroutine reported an error
	if len(errors) > 0 {
		t.Error("Concurrent GetAltIndex produced incorrect results")
	}
}

// TestConcurrentGetIndices verifies that GetIndices is safe for concurrent use.
func TestConcurrentGetIndices(t *testing.T) {
	xxh := NewXXHash(8, nil)
	numBuckets := uint(2048)

	const numGoroutines = 100
	const iterations = 1000

	var wg sync.WaitGroup
	wg.Add(numGoroutines)

	errors := make(chan error, numGoroutines)

	for g := 0; g < numGoroutines; g++ {
		go func(id int) {
			defer wg.Done()

			for i := 0; i < iterations; i++ {
				// Create unique data for this iteration
				data := []byte{byte(id), byte(i >> 8), byte(i)}

				// Call GetIndices - internally calls GetAltIndex
				i1, i2, fp := xxh.GetIndices(data, numBuckets)

				// Verify results are valid
				if i1 >= numBuckets || i2 >= numBuckets {
					errors <- nil
					return
				}

				if fp == 0 {
					errors <- nil
					return
				}

				// Verify we can get back to i1 from i2
				i1Back := xxh.GetAltIndex(i2, fp, numBuckets)
				if i1Back != i1 {
					errors <- nil
					return
				}
			}
		}(g)
	}

	wg.Wait()
	close(errors)

	if len(errors) > 0 {
		t.Error("Concurrent GetIndices produced incorrect results")
	}
}

// TestConcurrentMixedOperations tests concurrent calls to different methods.
func TestConcurrentMixedOperations(t *testing.T) {
	xxh := NewXXHash(8, nil)
	numBuckets := uint(1024)

	const numGoroutines = 50
	const iterations = 500

	var wg sync.WaitGroup
	wg.Add(numGoroutines * 3) // 3 types of operations

	errors := make(chan error, numGoroutines*3)

	// Goroutines calling GetIndices
	for g := 0; g < numGoroutines; g++ {
		go func(id int) {
			defer wg.Done()
			for i := 0; i < iterations; i++ {
				data := []byte{byte(id), byte(i)}
				i1, i2, fp := xxh.GetIndices(data, numBuckets)
				if i1 >= numBuckets || i2 >= numBuckets || fp == 0 {
					errors <- nil
					return
				}
			}
		}(g)
	}

	// Goroutines calling GetAltIndex
	for g := 0; g < numGoroutines; g++ {
		go func(id int) {
			defer wg.Done()
			for i := 0; i < iterations; i++ {
				index := uint(i % 500)
				fp := uint16((id+i)%65534 + 1)
				i2 := xxh.GetAltIndex(index, fp, numBuckets)
				if i2 >= numBuckets {
					errors <- nil
					return
				}
			}
		}(g)
	}

	// Goroutines calling hash64 directly
	for g := 0; g < numGoroutines; g++ {
		go func(id int) {
			defer wg.Done()
			for i := 0; i < iterations; i++ {
				data := []byte{byte(id), byte(i >> 8), byte(i)}
				hash := xxh.hash64(data)
				if hash == 0 {
					errors <- nil
					return
				}
			}
		}(g)
	}

	wg.Wait()
	close(errors)

	if len(errors) > 0 {
		t.Error("Concurrent mixed operations produced incorrect results")
	}
}
