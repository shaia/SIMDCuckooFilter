package cuckoofilter

import "github.com/shaia/simdcuckoofilter/internal/hash"

// hashStrategy is the internal type for hash function selection.
// Users should use the With*Hash() functions instead of using this type directly.
type hashStrategy int

const (
	hashStrategyFNV    hashStrategy = hashStrategy(hash.HashStrategyFNV)
	hashStrategyCRC32  hashStrategy = hashStrategy(hash.HashStrategyCRC32)
	hashStrategyXXHash hashStrategy = hashStrategy(hash.HashStrategyXXHash)
)

// String returns the string representation of the hash strategy
func (s hashStrategy) String() string {
	switch s {
	case hashStrategyFNV:
		return "FNV-1a"
	case hashStrategyCRC32:
		return "CRC32C"
	case hashStrategyXXHash:
		return "XXHash64"
	default:
		return "Unknown"
	}
}

// WithFNVHash configures the filter to use FNV-1a hash function.
// FNV-1a provides good distribution and compatibility with moderate speed.
// This is the default hash function.
func WithFNVHash() Option {
	return func(o *Options) {
		o.hashStrategy = hashStrategyFNV
	}
}

// WithCRC32Hash configures the filter to use CRC32C hash function.
// CRC32C is hardware-accelerated on modern CPUs (SSE4.2) and offers the fastest performance.
// Best for high-throughput scenarios where speed is critical.
func WithCRC32Hash() Option {
	return func(o *Options) {
		o.hashStrategy = hashStrategyCRC32
	}
}

// WithXXHash configures the filter to use XXHash64 hash function.
// XXHash64 provides excellent distribution and fast performance.
// Recommended for general-purpose use when you need better hash quality than FNV.
func WithXXHash() Option {
	return func(o *Options) {
		o.hashStrategy = hashStrategyXXHash
	}
}
