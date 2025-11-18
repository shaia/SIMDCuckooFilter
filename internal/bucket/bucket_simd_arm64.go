//go:build arm64
// +build arm64

package bucket

// containsSIMD uses NEON to check if a fingerprint exists in the bucket
// For very small data (< 4 bytes), uses inline scalar code since SIMD overhead exceeds benefit
func containsSIMD(data []byte, fp byte) bool {
	if len(data) < 4 {
		return inlineContains(data, fp)
	}
	return containsNEON(data, fp)
}

// isFullSIMD uses NEON to check if bucket is full (no zeros)
// For very small data (< 4 bytes), uses inline scalar code since SIMD overhead exceeds benefit
func isFullSIMD(data []byte) bool {
	if len(data) < 4 {
		return inlineIsFull(data)
	}
	return isFullNEON(data)
}

// countSIMD uses NEON to count non-zero entries
// For very small data (< 4 bytes), uses inline scalar code since SIMD overhead exceeds benefit
func countSIMD(data []byte) uint {
	if len(data) < 4 {
		return inlineCount(data)
	}
	return countNEON(data)
}

// findFirstZeroSIMD uses NEON to find the first zero slot
// For very small data (< 4 bytes), uses inline scalar code since SIMD overhead exceeds benefit
func findFirstZeroSIMD(data []byte) uint {
	if len(data) < 4 {
		return inlineFindFirstZero(data)
	}
	return findFirstZeroNEON(data)
}

// Assembly function declarations
// These are implemented in bucket_simd_arm64.s

//go:noescape
func containsNEON(data []byte, fp byte) bool

//go:noescape
func isFullNEON(data []byte) bool

//go:noescape
func countNEON(data []byte) uint

//go:noescape
func findFirstZeroNEON(data []byte) uint
