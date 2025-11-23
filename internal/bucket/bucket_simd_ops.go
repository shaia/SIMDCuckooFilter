package bucket

// SIMDBucket provides SIMD-optimized bucket operations
type SIMDBucket struct {
	*Bucket
}

// NewSIMDBucket creates a bucket with SIMD optimizations
func NewSIMDBucket(size uint) *SIMDBucket {
	return &SIMDBucket{
		Bucket: NewBucket(size),
	}
}

// ContainsSIMD checks if a fingerprint exists using SIMD
func (b *SIMDBucket) ContainsSIMD(fp uint16) bool {
	return containsSIMD(b.fingerprints[:b.size], fp)
}

// IsFullSIMD checks if bucket is full using SIMD
func (b *SIMDBucket) IsFullSIMD() bool {
	return isFullSIMD(b.fingerprints[:b.size])
}

// CountSIMD counts non-zero entries using SIMD
func (b *SIMDBucket) CountSIMD() uint {
	return countSIMD(b.fingerprints[:b.size])
}

// FindFirstZeroSIMD finds the index of the first zero slot using SIMD
// Returns size if no zero found (bucket is full)
func (b *SIMDBucket) FindFirstZeroSIMD() uint {
	return findFirstZeroSIMD(b.fingerprints[:b.size])
}

// InsertSIMD adds a fingerprint using SIMD-accelerated search
func (b *SIMDBucket) InsertSIMD(fp uint16) bool {
	idx := b.FindFirstZeroSIMD()
	if idx < b.size {
		b.fingerprints[idx] = fp
		return true
	}
	return false
}
