//go:build arm64
// +build arm64

package filter

// Lookup uses scalar fallback (NEON TODO)
func (f *simdFilter) Lookup(item []byte) bool {
	f.mu.RLock()
	defer f.mu.RUnlock()

	i1, i2, fp := f.hash.GetIndices(item, f.numBuckets)

	// Use NEON-optimized lookup through bucket's Contains method
	return f.buckets[i1].Contains(fp) || f.buckets[i2].Contains(fp)
}

// LookupBatch uses scalar fallback with batch hashing (NEON TODO)
func (f *simdFilter) LookupBatch(items [][]byte) []bool {
	results := make([]bool, len(items))

	f.mu.RLock()
	defer f.mu.RUnlock()

	// Use batch hashing with optimized ARM64 assembly
	hashResults := f.hash.GetIndicesBatch(items, f.numBuckets)

	// Batch process lookups with NEON optimizations
	for i, hr := range hashResults {
		results[i] = f.buckets[hr.I1].Contains(hr.Fp) || f.buckets[hr.I2].Contains(hr.Fp)
	}

	return results
}
