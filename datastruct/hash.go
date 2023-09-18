package datastruct

const prim32 = uint32(16777619)

func fvn32(key string) uint32 {
	hash := uint32(2166136261)
	for i := 0; i < len(key); i++ {
		hash *= prim32
		hash ^= uint32(key[i])
	}
	return hash
}
