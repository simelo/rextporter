package core

type RextServiceScraper interface {
	ApplyToServices(services ...string)
	ApplyToStacks(stacks ...string)
	AddAuthStrategy(auth RextAuth, name string)
	AddSource(source RextDataSource)
	AddSources(sources ...RextDataSource)
}

type RextMetricScrapper interface {
	Apply(rule RextScrapingRule)
	ApplyMany(rules ...RextScrapingRule)
}

type RextScrapingRule interface {
	GetMetricName() string
	GetMetricType() string
	GetMetricDescription() string
	GetMetricLabels() []string
}
