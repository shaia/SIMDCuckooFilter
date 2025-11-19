//go:build amd64 || arm64

package filter

const maxPowerOf2 = 1 << 63 // Largest power of 2 that fits in uint64

// nextPowerOf2 rounds up to the next power of 2
// Note: This implementation assumes 64-bit platforms (amd64/arm64)
// where uint is 64 bits. The final shift by 32 is necessary for
// values larger than 2^32.
//
// For values >= 2^63, returns 2^63 (the maximum representable power of 2).
// This prevents overflow while maintaining correct behavior for all valid
// cuckoo filter capacities.
func nextPowerOf2(n uint) uint {
	if n == 0 {
		return 1
	}
	// Cap at maximum power of 2 to prevent overflow
	if n > maxPowerOf2 {
		return maxPowerOf2
	}
	n--
	n |= n >> 1
	n |= n >> 2
	n |= n >> 4
	n |= n >> 8
	n |= n >> 16
	n |= n >> 32 // Required for 64-bit uint (values > 2^32)
	n++
	return n
}
