//go:build amd64 || arm64

package filter

import (
	"math/rand"
	"sync"

	"github.com/shaia/simdcuckoofilter/internal/bucket"
	"github.com/shaia/simdcuckoofilter/internal/hash"
)

// simdFilter is the platform-optimized filter implementation
type simdFilter struct {
	buckets    []*bucket.Bucket
	numBuckets uint
	numItems   uint
	maxKicks   uint
	bucketSize uint
	hash       hash.HashInterface
	batchSize  uint
	mu         sync.RWMutex
}

func New(capacity, bucketSize, fingerprintBits, maxKicks uint, hashStrategy hash.HashStrategy, batchSize uint) (*simdFilter, error) {
	// Calculate number of buckets
	numBuckets := nextPowerOf2((capacity + bucketSize - 1) / bucketSize)
	if numBuckets == 0 {
		numBuckets = 1
	}

	// Create buckets
	buckets := make([]*bucket.Bucket, numBuckets)
	for i := range buckets {
		buckets[i] = bucket.NewBucket(bucketSize)
	}

	return &simdFilter{
		buckets:    buckets,
		numBuckets: numBuckets,
		numItems:   0,
		maxKicks:   maxKicks,
		bucketSize: bucketSize,
		hash:       hash.NewHashFunction(hashStrategy, fingerprintBits),
		batchSize:  batchSize,
	}, nil
}

func (f *simdFilter) Insert(item []byte) bool {
	f.mu.Lock()
	defer f.mu.Unlock()

	i1, i2, fp := f.hash.GetIndices(item, f.numBuckets)

	// Try first bucket
	if f.buckets[i1].Insert(fp) {
		f.numItems++
		return true
	}

	// Try second bucket
	if f.buckets[i2].Insert(fp) {
		f.numItems++
		return true
	}

	// Both full, need to relocate
	return f.relocate(i1, i2, fp)
}

func (f *simdFilter) relocate(i1, i2 uint, fp byte) bool {
	// Start from random bucket
	index := i1
	if rand.Intn(2) == 1 {
		index = i2
	}

	currentFp := fp

	for i := uint(0); i < f.maxKicks; i++ {
		// Randomly select a position in the bucket (standard cuckoo hashing)
		pos := uint(rand.Intn(int(f.bucketSize)))

		// Swap the fingerprint at the random position
		oldFp := f.buckets[index].Swap(pos, currentFp)
		if oldFp == 0 {
			// Found an empty slot
			f.numItems++
			return true
		}

		// Continue with the evicted fingerprint
		currentFp = oldFp

		// Calculate alternative index for the evicted fingerprint
		index = f.hash.GetAltIndex(index, currentFp, f.numBuckets)

		// Try to insert the evicted fingerprint into its alternative bucket
		if f.buckets[index].Insert(currentFp) {
			f.numItems++
			return true
		}
	}

	return false
}

// Lookup is implemented in platform-specific files:
// - filter_amd64.go: Uses AVX2-optimized lookup
// - filter_arm64.go: Uses scalar fallback (NEON TODO)

// LookupBatch is implemented in platform-specific files for optimized batch processing

func (f *simdFilter) Delete(item []byte) bool {
	f.mu.Lock()
	defer f.mu.Unlock()

	i1, i2, fp := f.hash.GetIndices(item, f.numBuckets)

	if f.buckets[i1].Remove(fp) {
		f.numItems--
		return true
	}

	if f.buckets[i2].Remove(fp) {
		f.numItems--
		return true
	}

	return false
}

func (f *simdFilter) Count() uint {
	f.mu.RLock()
	defer f.mu.RUnlock()
	return f.numItems
}

func (f *simdFilter) LoadFactor() float64 {
	f.mu.RLock()
	defer f.mu.RUnlock()

	totalSlots := f.numBuckets * f.bucketSize
	if totalSlots == 0 {
		return 0
	}

	return float64(f.numItems) / float64(totalSlots)
}

func (f *simdFilter) Capacity() uint {
	return f.numBuckets * f.bucketSize
}

func (f *simdFilter) Reset() {
	f.mu.Lock()
	defer f.mu.Unlock()

	for _, b := range f.buckets {
		b.Reset()
	}
	f.numItems = 0
}

// Batch operations
func (f *simdFilter) InsertBatch(items [][]byte) []bool {
	results := make([]bool, len(items))
	for i, item := range items {
		results[i] = f.Insert(item)
	}
	return results
}

func (f *simdFilter) DeleteBatch(items [][]byte) []bool {
	results := make([]bool, len(items))
	for i, item := range items {
		results[i] = f.Delete(item)
	}
	return results
}

func (f *simdFilter) OptimalBatchSize() int {
	return int(f.batchSize)
}
