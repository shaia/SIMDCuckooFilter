//go:build amd64
// +build amd64

package bucket

//go:noescape
func containsAVX2(data []uint16, fp uint16) bool

// containsSIMD checks if a fingerprint exists in the bucket
// AMD64 implementation uses AVX2 for optimized parallel comparison
func containsSIMD(data []uint16, fp uint16) bool {
	// For very small buckets, the overhead of the assembly call dominates.
	// Use inline scalar implementation for size 4 and below.
	if len(data) <= 4 {
		return inlineContains(data, fp)
	}
	return containsAVX2(data, fp)
}

// isFullSIMD checks if bucket is full (no zeros)
// AMD64 implementation uses inline scalar code for all bucket sizes
func isFullSIMD(data []uint16) bool {
	return inlineIsFull(data)
}

// countSIMD counts non-zero entries
// AMD64 implementation uses inline scalar code for all bucket sizes
func countSIMD(data []uint16) uint {
	return inlineCount(data)
}

// findFirstZeroSIMD finds the first zero slot
// AMD64 implementation uses inline scalar code for all bucket sizes
func findFirstZeroSIMD(data []uint16) uint {
	return inlineFindFirstZero(data)
}
