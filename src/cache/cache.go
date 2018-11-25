package cache

// Cache mechanism for caching strings
type Cache interface {
	Get(key string) ([]byte, error)
	Set(key string, content []byte)
	Reset()
}

// NewCache create a new cache mechanism
func NewCache() Cache {
	return newMemCache()
}
