//go:build amd64
// +build amd64

package xxhash

//go:noescape
func hash64XXHashInternal(data []byte) uint64
