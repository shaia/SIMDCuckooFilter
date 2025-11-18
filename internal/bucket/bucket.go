package bucket

// Bucket represents a fixed-size array of fingerprints
// Each bucket can hold up to 'size' fingerprints (typically 4)
type Bucket struct {
	fingerprints []byte
	size         uint
}

// NewBucket creates a new bucket with the specified size
func NewBucket(size uint) *Bucket {
	return &Bucket{
		fingerprints: make([]byte, size),
		size:         size,
	}
}

// Insert adds a fingerprint to the bucket if there's space
// Returns true if successful, false if bucket is full
func (b *Bucket) Insert(fp byte) bool {
	idx := inlineFindFirstZero(b.fingerprints[:b.size])
	if idx < b.size {
		b.fingerprints[idx] = fp
		return true
	}
	return false
}

// Remove removes a fingerprint from the bucket
// Returns true if found and removed, false otherwise
func (b *Bucket) Remove(fp byte) bool {
	return inlineRemove(b.fingerprints[:b.size], fp)
}

// Contains checks if a fingerprint exists in the bucket
func (b *Bucket) Contains(fp byte) bool {
	return inlineContains(b.fingerprints[:b.size], fp)
}

// IsFull returns true if the bucket has no empty slots
func (b *Bucket) IsFull() bool {
	return inlineIsFull(b.fingerprints[:b.size])
}

// Count returns the number of non-zero fingerprints in the bucket
func (b *Bucket) Count() uint {
	return inlineCount(b.fingerprints[:b.size])
}

// Swap replaces a fingerprint at the given index and returns the old value
// This is used during cuckoo hashing relocation
func (b *Bucket) Swap(index uint, fp byte) byte {
	if index >= b.size {
		return 0
	}
	old := b.fingerprints[index]
	b.fingerprints[index] = fp
	return old
}

// Reset clears all fingerprints in the bucket
func (b *Bucket) Reset() {
	for i := range b.fingerprints {
		b.fingerprints[i] = 0
	}
}

// GetFingerprints returns the underlying fingerprint slice (for SIMD operations)
func (b *Bucket) GetFingerprints() []byte {
	return b.fingerprints
}
