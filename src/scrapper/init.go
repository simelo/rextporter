package scrapper

import (
	"github.com/simelo/rextporter/src/client"
)

var workPool *client.Pool

func init() {
	workPool = client.NewPool(6)
	workPool.StartDispatcher()
}
