//go:build arm64
// +build arm64

package xxhash

//go:noescape
func hash64XXHashInternal(data []byte) uint64
