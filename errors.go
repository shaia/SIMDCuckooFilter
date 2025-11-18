package cuckoofilter

import "errors"

var (
	// ErrInvalidCapacity is returned when capacity is zero or invalid
	ErrInvalidCapacity = errors.New("capacity must be greater than zero")

	// ErrInvalidBucketSize is returned when bucket size is not valid
	ErrInvalidBucketSize = errors.New("bucket size must be 2, 4, 8, 16, 32, or 64")

	// ErrInvalidFingerprintSize is returned when fingerprint size is invalid
	ErrInvalidFingerprintSize = errors.New("fingerprint size must be between 1 and 8 bits (stored as bytes)")

	// ErrInvalidHashStrategy is returned when hash strategy is unknown
	ErrInvalidHashStrategy = errors.New("invalid hash strategy")
)
