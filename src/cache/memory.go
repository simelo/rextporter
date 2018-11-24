package cache

import (
	"errors"
	"fmt"
	"sync"
)

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

func (c MemCache) Get(key string) (val []byte, err error) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	if val, ok := c.vals[key]; ok {
		return val, nil
	}
	return nil, errors.New(fmt.Sprintf("value not found for key %s", key))
}

func (c MemCache) Set(key string, data []byte) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.vals[key] = data
}

func (c MemCache) Reset() {
	c.mu.Lock()
	defer c.mu.Unlock()
	for k := range c.vals {
		delete(c.vals, k)
	}
}
