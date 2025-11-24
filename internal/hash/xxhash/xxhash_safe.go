package xxhash

// hash64XXHashGo is the Go fallback implementation
func hash64XXHashGo(data []byte) uint64 {
	var hash uint64

	if len(data) >= 8 {
		hash = prime64_5 + uint64(len(data))
		for len(data) >= 8 {
			k := uint64(data[0]) | uint64(data[1])<<8 | uint64(data[2])<<16 | uint64(data[3])<<24 |
				uint64(data[4])<<32 | uint64(data[5])<<40 | uint64(data[6])<<48 | uint64(data[7])<<56
			k *= prime64_2
			k = (k << 31) | (k >> 33)
			k *= prime64_1
			hash ^= k
			hash = ((hash << 27) | (hash >> 37)) * prime64_1
			hash += prime64_4
			data = data[8:]
		}
	} else {
		hash = prime64_5 + uint64(len(data))
	}

	for len(data) > 0 {
		hash ^= uint64(data[0]) * prime64_5
		hash = ((hash << 11) | (hash >> 53)) * prime64_1
		data = data[1:]
	}

	hash ^= hash >> 33
	hash *= prime64_2
	hash ^= hash >> 29
	hash *= prime64_3
	hash ^= hash >> 32

	return hash
}
