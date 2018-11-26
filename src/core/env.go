package core

// RextEnv encapsulates an implementation of the rextporter top-level environment
type RextEnv interface {
	NewServiceScraper() RextServiceScraper
	NewAuthStrategy(authtype string, options RextKeyValueStore)
	NewMetricsExtractor(scraperType string, options RextKeyValueStore, metrics []RextMetricDef) error
	NewMetricsDatasource(srcType string) RextDataSource
	RegisterScraperForServices(s RextServiceScraper, services ...string) error
	RegisterScraperForStacks(s RextServiceScraper, stacks ...string) error
	GetOptions() RextKeyValueStore
}
