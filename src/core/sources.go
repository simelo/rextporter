package core

type RextDataSource interface {
	SetBaseURL(url string) string
	GetResourceLocation() string
	SetResourceLocation(url string) error
	ActivateScraper(scraper RextMetricScraper)
	RetrieveData() (interface{}, error)
}
