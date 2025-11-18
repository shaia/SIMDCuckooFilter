//go:build amd64
// +build amd64

package cpu

// DefaultSIMD is the default/preferred SIMD type for AMD64.
// AVX2 provides best performance with 32-byte parallel processing.
const DefaultSIMD = SIMDAVX2

// FallbackSIMD is the same as DefaultSIMD for AMD64 (AVX2-only).
const FallbackSIMD = SIMDAVX2

// GetBestSIMD returns the SIMD type for AMD64.
// Always returns AVX2 as SSE2 support has been removed.
// The parameter is kept for API compatibility.
func GetBestSIMD(_ bool) SIMDType {
	return SIMDAVX2
}
