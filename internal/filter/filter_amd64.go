//go:build amd64
// +build amd64

package filter

// Lookup uses optimized bucket lookup
func (f *simdFilter) Lookup(item []byte) bool {
	f.mu.RLock()
	defer f.mu.RUnlock()

	i1, i2, fp := f.hash.GetIndices(item, f.numBuckets)

	// Use bucket's optimized lookup (handles 16-bit fingerprints)
	return f.buckets[i1].Contains(fp) || f.buckets[i2].Contains(fp)
}

// LookupBatch uses optimized batch processing
func (f *simdFilter) LookupBatch(items [][]byte) []bool {
	results := make([]bool, len(items))

	f.mu.RLock()
	defer f.mu.RUnlock()

	// Use batch hashing for better performance
	hashResults := f.hash.GetIndicesBatch(items, f.numBuckets)

	// Batch process lookups
	for i, hr := range hashResults {
		results[i] = f.buckets[hr.I1].Contains(hr.Fp) ||
			f.buckets[hr.I2].Contains(hr.Fp)
	}

	return results
}
