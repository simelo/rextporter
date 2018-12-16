package core

// RextEnv encapsulates an implementation of the rextporter top-level environment
type RextEnv interface {
	NewAuthStrategy(authtype string, options RextKeyValueStore) (RextAuthDef, error)
	// NewMetricsDatasource(srcType string) (RextDataSource, error)
	// service name -> template(metrics)
	// RegisterScraperForServices(s RextServiceScraperDef, services ...string) error
	// stack name -> template(metrics)
	// RegisterScraperForStacks(s RextServiceScraperDef, stacks ...string) error
	GetOptions() RextKeyValueStore
}
