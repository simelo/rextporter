package rxt

import "github.com/simelo/rextporter/src/core"

type OptionsMap map[string]string

func newOptionsMap() OptionsMap {
	return make(OptionsMap)
}

type RxtDefScraperDataset struct {
	SupportedServiceNames []string
	SupportedStackNames   []string
	Definitions           map[string]interface{}
	Sources               []RxtDefDataSource
	Options               OptionsMap
}

type RxtDefAuth struct {
	options OptionsMap
}

type RxtDefSource struct {
	Method   string
	Type     string
	Location string
	Scrapers RxtDefScraper
	Options  OptionsMap
}

type RxtDefScraper struct {
	Type    string
	Metrics []RxtDefMetric
	Options OptionsMap
}

type RxtDefMetric struct {
	Type        string
	Name        string
	Description string
	Labels      []string
	Options     OptionsMap
}
