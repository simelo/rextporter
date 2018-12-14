package memconfig

import "github.com/simelo/rextporter/src/core"

// Scrapper implements core.RextServiceScraper
type Scrapper struct {
	auths                 map[string]core.RextAuth // FIXME(denisacostaq@gmail.com): make new
	rextDataSources       []core.RextDataSource
	supportedServiceNames []string
	supportedStackNames   []string
	options               core.RextKeyValueStore
}

// AddAuthStrategy add a new auth strategy to the scrapper
func (s *Scrapper) AddAuthStrategy(auth core.RextAuth, name string) {
	s.auths[name] = auth
}

// AddSource add a source to the scrapper
func (s *Scrapper) AddSource(source core.RextDataSource) {
	s.rextDataSources = append(s.rextDataSources, source)
}

// AddSources add multiple sources to the scrapper
func (s *Scrapper) AddSources(sources ...core.RextDataSource) {
	for _, source := range sources {
		s.AddSource(source)
	}
}

// GetOptions return key/value pairs for extra options
func (s Scrapper) GetOptions() core.RextKeyValueStore {
	return s.options
}

// NewScrapperConf create a new service scrapper
func NewScrapperConf(options core.RextKeyValueStore) core.RextServiceScraper {
	return &Scrapper{
		// 	auths                 map[string]core.RextAuth // FIXME(denisacostaq@gmail.com): make new
		// rextDataSources       []core.RextDataSource
		// supportedServiceNames []string
		// supportedStackNames   []string
		options: options,
	}
}
