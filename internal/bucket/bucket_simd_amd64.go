//go:build amd64
// +build amd64

package bucket

// containsSIMD checks if a fingerprint exists in the bucket
// AMD64 implementation uses inline scalar code for all bucket sizes
// (AVX2 bucket operations removed as buckets are typically small: 4-64 bytes)
func containsSIMD(data []byte, fp byte) bool {
	return inlineContains(data, fp)
}

// isFullSIMD checks if bucket is full (no zeros)
// AMD64 implementation uses inline scalar code for all bucket sizes
func isFullSIMD(data []byte) bool {
	return inlineIsFull(data)
}

// countSIMD counts non-zero entries
// AMD64 implementation uses inline scalar code for all bucket sizes
func countSIMD(data []byte) uint {
	return inlineCount(data)
}

// findFirstZeroSIMD finds the first zero slot
// AMD64 implementation uses inline scalar code for all bucket sizes
func findFirstZeroSIMD(data []byte) uint {
	return inlineFindFirstZero(data)
}
