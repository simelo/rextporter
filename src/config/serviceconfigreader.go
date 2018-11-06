package config

// ServiceConfigReader is an interface to get a service config from for example:
// a file, a REST API, a stream and so on...
type ServiceConfigReader interface {
	// GetConfig return a service config or an error if any
	GetConfig() (Service, error)
}
