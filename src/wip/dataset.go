package wip

// RextMetric provides access to values measured for a given metric
type RextMetric interface {
	GetMetadata() RextMetricDef
	// TODO: Methods to retrieve values measured for a given metric
}
