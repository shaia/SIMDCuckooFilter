//go:build !amd64 && !arm64
// +build !amd64,!arm64

package xxhash

// hash64XXHashInternal is the Go implementation for non-optimized platforms
func hash64XXHashInternal(data []byte) uint64 {
	return hash64XXHashGo(data)
}
