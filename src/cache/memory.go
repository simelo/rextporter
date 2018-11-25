package cache

import (
	"fmt"
	"sync"
)

// MemCache implements a cache mechanism in memory bay using a map
type MemCache struct {
	vals map[string][]byte
	mu   *sync.RWMutex
}

func newMemCache() *MemCache {
	return &MemCache{
		vals: make(map[string][]byte),
		mu:   &sync.RWMutex{},
	}
}

// Get return the cached value by a giving key, error if this val not found
func (c MemCache) Get(key string) (val []byte, err error) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	if val, ok := c.vals[key]; ok {
		return val, nil
	}
	return nil, fmt.Errorf("value not found for key %s", key)
}

// Set save data with a giving key in the cache sistem
func (c MemCache) Set(key string, data []byte) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.vals[key] = data
}

// Reset clear al the cached data
func (c MemCache) Reset() {
	c.mu.Lock()
	defer c.mu.Unlock()
	for k := range c.vals {
		delete(c.vals, k)
	}
}
