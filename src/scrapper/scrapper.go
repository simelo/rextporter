package scrapper

// Scrapper receive some data as input and should return the metric val
type Scrapper interface {
	GetMetric(data interface{}) interface{}
}
