package config

import (
	"errors"
	"fmt"
	"strings"
)

// Host is a concept to grab information about a running server, for example:
// where is it http://localhost:1234 (Location + : + Port), what auth kind you need to use?
// what is the header key you in which you need to send the token, and so on.
type Host struct {
	Ref                  string
	Location             string `json:"location"`
	Port                 int    `json:"port"`
	AuthType             string `json:"auth_type"`
	TokenHeaderKey       string `json:"token_header_key"`
	GenTokenEndpoint     string `json:"gen_token_endpoint"`
	TokenKeyFromEndpoint string `json:"token_key_from_endpoint"`
}

// URIToGetMetric build the URI from where you will to get metric information
func (host Host) URIToGetMetric(metricInHost Link) string {
	return fmt.Sprintf("%s:%d%s", host.Location, host.Port, metricInHost.URL)
}

// URIToGetToken build the URI from where you will to get the token
func (host Host) URIToGetToken() string {
	return fmt.Sprintf("%s:%d%s", host.Location, host.Port, host.TokenKeyFromEndpoint)
}

func (host Host) validate() (errs []error) {
	if len(host.Ref) == 0 {
		errs = append(errs, errors.New("ref is required in host"))
	}
	if len(host.Location) == 0 {
		errs = append(errs, errors.New("location is required in host"))
	}
	if !isValidURL(host.Location) {
		errs = append(errs, errors.New("location is not a valid url in host"))
	}
	if !isValidURL(host.URIToGetToken()) {
		errs = append(errs, errors.New("location + port can not form a valid uri in host"))
	}
	if host.Port < 0 || host.Port > 65535 {
		errs = append(errs, errors.New("port number should be between 0 and 65535 in host"))
	}
	if strings.Compare(host.AuthType, "CSRF") == 0 && len(host.TokenHeaderKey) == 0 {
		errs = append(errs, errors.New("TokenHeaderKey is required if you are using CSRF"))
	}
	if strings.Compare(host.AuthType, "CSRF") == 0 && len(host.GenTokenEndpoint) == 0 {
		errs = append(errs, errors.New("GenTokenEndpoint is required if you are using CSRF"))
	}
	if strings.Compare(host.AuthType, "CSRF") == 0 && len(host.TokenKeyFromEndpoint) == 0 {
		errs = append(errs, errors.New("TokenKeyFromEndpoint is required if you are using CSRF"))
	}
	return errs
}
