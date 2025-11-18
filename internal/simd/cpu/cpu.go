package cpu

// SIMDType represents the SIMD instruction set being used
type SIMDType int

const (
	SIMDNone SIMDType = iota
	SIMDAVX2
	SIMDNEON
)

func (s SIMDType) String() string {
	switch s {
	case SIMDAVX2:
		return "AVX2"
	case SIMDNEON:
		return "NEON"
	default:
		return "None"
	}
}
