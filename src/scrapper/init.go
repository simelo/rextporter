package scrapper

var WorkPool *Pool

func init() {
	WorkPool = NewPool(6)
	WorkPool.StartDispatcher()
}
