//go:build arm64
// +build arm64

package lookup

// BucketLookup performs NEON-optimized lookup in a bucket for ARM64.
// Uses NEON instructions to process 16 bytes in parallel.
//
// Performance: ~2-3x faster than scalar implementation (when NEON assembly is implemented).
//
// TODO: Currently uses scalar fallback. Implement NEON assembly in bucket_lookup_neon_arm64.s
func BucketLookup(fingerprints []byte, target byte) bool {
	if len(fingerprints) == 0 {
		return false
	}
	// TODO: return bucketLookupNEON(fingerprints, target)
	return BucketLookupScalar(fingerprints, target)
}

// BucketLookupScalar provides scalar lookup for benchmarking.
// Simple loop implementation without SIMD.
func BucketLookupScalar(fingerprints []byte, target byte) bool {
	for _, fp := range fingerprints {
		if fp == target {
			return true
		}
	}
	return false
}

// TODO: Implement bucketLookupNEON in assembly
// bucketLookupNEON performs NEON-optimized bucket lookup.
// Implemented in bucket_lookup_neon_arm64.s
//
// //go:noescape
// func bucketLookupNEON(fingerprints []byte, target byte) bool
