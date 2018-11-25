package client

import (
	"github.com/simelo/rextporter/src/cache"
)

type ClientCatcher struct {
	cache cache.Cache
	cl    CacheableClient
}

func NewClientCatcher(cl CacheableClient, cache cache.Cache) Client {
	return ClientCatcher{cache: cache, cl: cl}
}

func (cl ClientCatcher) GetData() (body []byte, err error) {
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
