//go:build amd64 || arm64

package filter

// nextPowerOf2 rounds up to the next power of 2
// Note: This implementation assumes 64-bit platforms (amd64/arm64)
// where uint is 64 bits. The final shift by 32 is necessary for
// values larger than 2^32.
func nextPowerOf2(n uint) uint {
	if n == 0 {
		return 1
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
