// Package types provides common types shared across hash implementations.
package types

// HashResult holds the result of hashing an item for cuckoo filter insertion.
// It contains the two bucket indices (I1, I2) for cuckoo hashing and the
// fingerprint (Fp) value that will be stored in the bucket.
type HashResult struct {
	I1, I2 uint // Two bucket indices for cuckoo hashing
	Fp     byte // Fingerprint value (never zero, as 0 indicates empty slot)
}
