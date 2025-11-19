package cuckoofilter

import (
	"github.com/shaia/simdcuckoofilter/internal/filter"
	"github.com/shaia/simdcuckoofilter/internal/hash"
)

// CuckooFilter is a probabilistic data structure for set membership testing
type CuckooFilter interface {
	// Insert adds an item to the filter
	// Returns true if successful, false if filter is full
	Insert(item []byte) bool

	// Lookup checks if an item might be in the filter
	// Returns true if item might be present (with false positive probability)
	// Returns false if item is definitely not present
	Lookup(item []byte) bool

	// Delete removes an item from the filter
	// Returns true if item was found and deleted
	Delete(item []byte) bool

	// Count returns the approximate number of items in the filter
	Count() uint

	// LoadFactor returns current load factor (0.0 to 1.0)
	LoadFactor() float64

	// Capacity returns the total capacity of the filter
	Capacity() uint

	// Reset clears all items from the filter
	Reset()
}

// BatchFilter extends CuckooFilter with batch operations (SIMD-optimized)
type BatchFilter interface {
	CuckooFilter

	// InsertBatch inserts multiple items
	InsertBatch(items [][]byte) []bool

	// LookupBatch checks multiple items
	LookupBatch(items [][]byte) []bool

	// DeleteBatch deletes multiple items
	DeleteBatch(items [][]byte) []bool

	// OptimalBatchSize returns recommended batch size for this implementation
	OptimalBatchSize() int
}

// New creates a SIMD-optimized Cuckoo filter with the specified capacity.
// Uses SIMD implementation based on platform:
//   - AMD64: AVX2
//   - ARM64: NEON
//
// Examples:
//   cf, _ := cuckoofilter.New(10000)
//   cf, _ := cuckoofilter.New(10000, cuckoofilter.WithFingerprintSize(8))
//   cf, _ := cuckoofilter.New(10000, cuckoofilter.WithBucketSize(32))
func New(capacity uint, opts ...Option) (CuckooFilter, error) {
	if capacity == 0 {
		return nil, ErrInvalidCapacity
	}

	options := defaultOptions()
	for _, opt := range opts {
		opt(&options)
	}

	if err := options.Validate(); err != nil {
		return nil, err
	}

	return filter.New(
		capacity,
		options.bucketSize,
		options.fingerprintBits,
		options.maxKicks,
		hash.HashStrategy(options.hashStrategy),
		options.batchSize,
	)
}
