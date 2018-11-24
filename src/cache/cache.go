package cache

// Storage mechanism for caching strings
type Storage interface {
	Get(key string) ([]byte, error)
	Set(key string, content []byte)
	Reset()
}

func NewStorage() Storage {
	return newMemCache()
}
