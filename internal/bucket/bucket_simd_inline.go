package bucket

// Shared scalar implementations used across bucket operations.
// These eliminate code duplication and provide consistent behavior.
//
// For SIMD operations on very small data (< 4 bytes), these are used
// as inline fallbacks since SIMD overhead exceeds any performance benefit.
//
// For standard bucket operations, these provide the core loop logic.

// inlineContains checks if a fingerprint exists in data using scalar loop
func inlineContains(data []uint16, fp uint16) bool {
	for _, b := range data {
		if b == fp {
			return true
		}
	}
	return false
}

// inlineIsFull checks if data is full (no zeros) using scalar loop
func inlineIsFull(data []uint16) bool {
	for _, b := range data {
		if b == 0 {
			return false
		}
	}
	return true
}

// inlineCount counts non-zero entries in data using scalar loop
func inlineCount(data []uint16) uint {
	count := uint(0)
	for _, b := range data {
		if b != 0 {
			count++
		}
	}
	return count
}

// inlineFindFirstZero finds the first zero slot in data using scalar loop
// Returns len(data) if no zero found
func inlineFindFirstZero(data []uint16) uint {
	for i, b := range data {
		if b == 0 {
			return uint(i)
		}
	}
	return uint(len(data))
}

// inlineRemove finds and removes (zeros) the first occurrence of a fingerprint
// Returns true if found and removed, false otherwise
func inlineRemove(data []uint16, fp uint16) bool {
	for i, b := range data {
		if b == fp {
			data[i] = 0
			return true
		}
	}
	return false
}
