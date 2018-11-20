package core

type RextDataSource interface {
	GetLocation() string
	SetLocation(url string) error
	ActivateScraper(scraper RextScraper)
}
