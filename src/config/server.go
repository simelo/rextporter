package config

// Server the server where is running the service
type Server struct {
	// Location should have the ip or URL.
	Location string `json:"location"`
}
