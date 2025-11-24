package cuckoofilter

// Options configures a Cuckoo filter
type Options struct {
	bucketSize      uint
	fingerprintBits uint
	maxKicks        uint
	hashStrategy    hashStrategy
	preferSIMD      bool
	preferAVX2      bool
	batchSize       uint
}

// Option is a function that configures Options
type Option func(*Options)

// defaultOptions returns the default configuration
func defaultOptions() Options {
	return Options{
		bucketSize:      4,
		fingerprintBits: 8,
		maxKicks:        500,
		hashStrategy:    hashStrategyFNV,
		preferSIMD:      true,
		preferAVX2:      true,
		batchSize:       32,
	}
}

// Validate checks if the options are valid
func (o *Options) Validate() error {
	if o.bucketSize != 2 && o.bucketSize != 4 && o.bucketSize != 8 &&
		o.bucketSize != 16 && o.bucketSize != 32 && o.bucketSize != 64 {
		return ErrInvalidBucketSize
	}
	// Fingerprints are stored as uint16, so 1-16 bits are supported
	if o.fingerprintBits < 1 || o.fingerprintBits > 16 {
		return ErrInvalidFingerprintSize
	}
	return nil
}

// WithBucketSize sets the number of fingerprints per bucket (2, 4, 8, 16, 32, or 64)
// Larger sizes provide better load factors and benefit more from SIMD optimizations.
// Recommended: 8 for balanced performance, 32 for maximum load factor, 64 for AVX2 and cache line alignment.
func WithBucketSize(size uint) Option {
	return func(o *Options) {
		o.bucketSize = size
	}
}

// WithFingerprintSize sets the fingerprint size in bits (1-16)
// Fingerprints are stored as uint16, so the maximum is 16 bits.
// Common values: 8 (standard), 12 (low false positive), 16 (very low false positive)
func WithFingerprintSize(bits uint) Option {
	return func(o *Options) {
		o.fingerprintBits = bits
	}
}

// WithMaxKicks sets the maximum number of relocation attempts
func WithMaxKicks(kicks uint) Option {
	return func(o *Options) {
		o.maxKicks = kicks
	}
}

// WithSIMD enables or disables SIMD optimizations
func WithSIMD(enabled bool) Option {
	return func(o *Options) {
		o.preferSIMD = enabled
	}
}

// WithAVX2 enables or disables AVX2 optimizations
func WithAVX2(prefer bool) Option {
	return func(o *Options) {
		o.preferAVX2 = prefer
	}
}

// WithBatchSize sets the batch operation size
func WithBatchSize(size uint) Option {
	return func(o *Options) {
		o.batchSize = size
	}
}
