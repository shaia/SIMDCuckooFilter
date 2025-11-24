//go:build !amd64 && !arm64
// +build !amd64,!arm64

package xxhash

// hash64XXHashInternal calls the Go fallback implementation for generic architectures
func hash64XXHashInternal(data []byte) uint64 {
	return hash64XXHashGo(data)
}
