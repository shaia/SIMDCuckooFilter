//go:build arm64
// +build arm64

package crc32hash

// batchCRC32Hardware is implemented in batch_simd_arm64.s
// Processes CRC32C checksums using ARMv8 hardware CRC32 instructions
//
//go:noescape
func batchCRC32Hardware(items [][]byte, results []uint32)
