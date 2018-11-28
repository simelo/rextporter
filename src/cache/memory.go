package cache

import (
	"fmt"
	"sync"
)

// MemCache implements a cache mechanism in memory bay using a map
type MemCache struct {
	vals map[string][]byte
	sync.RWMutex
}

func newMemCache() *MemCache {
	return &MemCache{
		vals: make(map[string][]byte),
	}
}

// Get return the cached value by a giving key, error if this val not found
func (c *MemCache) Get(key string) (val []byte, err error) {
	c.RLock()
	defer c.RUnlock()
	if val, ok := c.vals[key]; ok {
		return val, nil
	}
	return nil, fmt.Errorf("value not found for key %s", key)
}

// Set save data with a giving key in the cache sistem
func (c *MemCache) Set(key string, data []byte) {
	c.Lock()
	defer c.Unlock()
	c.vals[key] = data
}

// Reset clear al the cached data
func (c *MemCache) Reset() {
	c.Lock()
	defer c.Unlock()
	for k := range c.vals {
		delete(c.vals, k)
	}
}
