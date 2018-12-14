package memconfig

import (
	"github.com/simelo/rextporter/src/core"
)

// RootConfig implements core.RextEnv
type RootConfig struct {
	options core.RextKeyValueStore
}

// NewServiceScraper create a new service scrapper
func (conf RootConfig) NewServiceScraper() (core.RextServiceScraper, error) {
	return &Scrapper{
		// SupportedServiceNames: nil,
		// SupportedStackNames:   nil,
		// Definitions:           make(map[string]interface{}),
		// Sources:               nil,
		options: NewOptionsMap(),
	}, nil
}

// NewAuthStrategy create a new auth strategy
func (conf RootConfig) NewAuthStrategy(authtype string, options core.RextKeyValueStore) (core.RextAuth, error) {
	return &HTTPAuth{}, nil
}

// NewMetricsExtractor create a new metrics extractor
func (conf RootConfig) NewMetricsExtractor(scraperType string, options core.RextKeyValueStore, metrics []core.RextMetricDef) (core.RextMetricsExtractor, error) {
	return &MetricsExtractor{}, nil
}

// NewMetricsDatasource create a new datasource of type srcType
func (conf RootConfig) NewMetricsDatasource(srcType string) (core.RextDataSource, error) {
	return &Service{}, nil
}

// RegisterScraperForServices register service scrapper "s" for all services in parameter
func (conf RootConfig) RegisterScraperForServices(s core.RextServiceScraper, services ...string) error {
	if ds, isScrapper := s.(*Scrapper); isScrapper {
		ds.supportedServiceNames = services
		return nil
	}
	return core.ErrInvalidType
}

// RegisterScraperForStacks register service scrapper "s" for all stacks in parameter
func (conf RootConfig) RegisterScraperForStacks(s core.RextServiceScraper, stacks ...string) error {
	if ds, isScrapper := s.(*Scrapper); isScrapper {
		ds.supportedStackNames = stacks
		return nil
	}
	return core.ErrInvalidType
}

// GetOptions return key/value pairs for extra options
func (conf RootConfig) GetOptions() core.RextKeyValueStore {
	return conf.options
}

// NewRootConfig create a new root config
func NewRootConfig(options core.RextKeyValueStore) core.RextEnv {
	return &RootConfig{
		options: options,
	}
}
