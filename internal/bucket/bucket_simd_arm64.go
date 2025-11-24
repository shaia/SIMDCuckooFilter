//go:build arm64
// +build arm64

package bucket

// containsSIMD uses NEON to check if a fingerprint exists in the bucket
// For 16-bit fingerprints, we currently fallback to scalar inline code
func containsSIMD(data []uint16, fp uint16) bool {
	return inlineContains(data, fp)
}

// isFullSIMD uses NEON to check if bucket is full (no zeros)
// For 16-bit fingerprints, we currently fallback to scalar inline code
func isFullSIMD(data []uint16) bool {
	return inlineIsFull(data)
}

// countSIMD uses NEON to count non-zero entries
// For 16-bit fingerprints, we currently fallback to scalar inline code
func countSIMD(data []uint16) uint {
	return inlineCount(data)
}

// findFirstZeroSIMD uses NEON to find the first zero slot
// For 16-bit fingerprints, we currently fallback to scalar inline code
func findFirstZeroSIMD(data []uint16) uint {
	return inlineFindFirstZero(data)
}

// Assembly function declarations
// These are implemented in bucket_simd_arm64.s
// Note: Currently unused as we fallback to scalar for 16-bit support
/*
//go:noescape
func containsNEON(data []byte, fp byte) bool

//go:noescape
func isFullNEON(data []byte) bool

//go:noescape
func countNEON(data []byte) uint

//go:noescape
func findFirstZeroNEON(data []byte) uint
*/
