//go:build arm64
// +build arm64

package cpu

// DefaultSIMD is the default SIMD type for ARM64.
// NEON is mandatory in ARMv8, so it's always available.
const DefaultSIMD = SIMDNEON

// FallbackSIMD is the same as DefaultSIMD for ARM64.
const FallbackSIMD = SIMDNEON

// GetBestSIMD returns the best SIMD type for ARM64.
// NEON is always available on ARM64, so this always returns SIMDNEON.
// The parameter is ignored but kept for API compatibility with AMD64.
func GetBestSIMD(_ bool) SIMDType {
	return SIMDNEON
}
