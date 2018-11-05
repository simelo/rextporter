package config

import "errors"

// Server the server where is running the service
type Server struct {
	// Location should have the ip or URL.
	Location string `json:"location"`
}

func (server Server) validate() (errs []error) {
	if len(server.Location) == 0 {
		errs = append(errs, errors.New("location is required in server"))
	}
	if !isValidURL(server.Location) {
		errs = append(errs, errors.New("location is not a valid url in server"))
	}
	return errs
}
