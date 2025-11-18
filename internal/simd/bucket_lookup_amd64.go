//go:build amd64
// +build amd64

package simd

// BucketLookup performs AVX2-optimized lookup in a bucket for AMD64.
// Uses AVX2 instructions to process 32 bytes in parallel.
func BucketLookup(fingerprints []byte, target byte) bool {
	if len(fingerprints) == 0 {
		return false
	}
	return bucketLookupAVX2(fingerprints, target)
}

// bucketLookupAVX2 performs AVX2-optimized bucket lookup.
// Implemented in bucket_lookup_avx2_amd64.s
//
//go:noescape
func bucketLookupAVX2(fingerprints []byte, target byte) bool
