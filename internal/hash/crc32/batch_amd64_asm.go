//go:build amd64
// +build amd64

package crc32hash

// batchCRC32SIMD is implemented in batch_simd_amd64.s
// Processes CRC32C checksums using SSE4.2 CRC32 instruction
//
//go:noescape
func batchCRC32SIMD(items [][]byte, results []uint32)
