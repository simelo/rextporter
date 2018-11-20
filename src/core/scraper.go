package core

type RextServiceScraper interface {
	ApplyToServices(services ...string)
	ApplyToStacks(stacks ...string)
	AddAuthStrategy(auth RextAuth, name string)
	AddSource(source RextDataSource)
	AddSources(sources ...RextDataSource)
}

type RextMetricsExtractor interface {
	ExtractMetrics(target interface{}) []RextMetric
	MergeMetrics(overwrite bool, metrics ...[]RextMetric) []RextMetric
}

type RextMetricScraper interface {
	Apply(rule RextMetricDef)
	ApplyMany(rules ...RextMetricDef)
}

type RextMetricDef interface {
	GetMetricName() string
	GetMetricType() string
	GetMetricDescription() string
	GetMetricLabels() []string
}
