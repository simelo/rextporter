package client

import (
	"github.com/simelo/rextporter/src/cache"
)

// Catcher have a client and a cache to save time loading data
type Catcher struct {
	cache cache.Cache
	cl    CacheableClient
}

// NewCatcher create a Client compatible  catcher
func NewCatcher(cl CacheableClient, cache cache.Cache) Client {
	return Catcher{cache: cache, cl: cl}
}

// GetData return the data, can be from local cache or making the original request
func (cl Catcher) GetData() (body []byte, err error) {
	dataKey := cl.cl.DataPath()
	if body, err = cl.cache.Get(dataKey); err == nil {
		return body, err
	}
	if body, err = cl.cl.GetData(); err == nil {
		cl.cache.Set(dataKey, body)
	}
	return body, err
	// return cl.cl.GetData()
}
