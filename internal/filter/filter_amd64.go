//go:build amd64
// +build amd64

package filter

import "github.com/shaia/simdcuckoofilter/internal/lookup"

// Lookup uses AVX2-optimized bucket lookup
func (f *simdFilter) Lookup(item []byte) bool {
	f.mu.RLock()
	defer f.mu.RUnlock()

	i1, i2, fp := f.hash.GetIndices(item, f.numBuckets)

	// Use AVX2-optimized lookup
	return lookup.BucketLookup(f.buckets[i1].GetFingerprints(), fp) ||
		lookup.BucketLookup(f.buckets[i2].GetFingerprints(), fp)
}

// LookupBatch uses AVX2-optimized batch processing
func (f *simdFilter) LookupBatch(items [][]byte) []bool {
	results := make([]bool, len(items))

	f.mu.RLock()
	defer f.mu.RUnlock()

	// Use batch hashing for better performance
	hashResults := f.hash.GetIndicesBatch(items, f.numBuckets)

	// Batch process lookups with AVX2
	for i, hr := range hashResults {
		results[i] = lookup.BucketLookup(f.buckets[hr.I1].GetFingerprints(), hr.Fp) ||
			lookup.BucketLookup(f.buckets[hr.I2].GetFingerprints(), hr.Fp)
	}

	return results
}
