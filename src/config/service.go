package config

import (
	"errors"
	"fmt"
	"strings"

	"github.com/prometheus/client_golang/prometheus"
)

// Service is a concept to grab information about a running server, for example:
// where is it http://localhost:1234 (Location + : + Port), what auth kind you need to use?
// what is the header key you in which you need to send the token, and so on.
type Service struct {
	Name string `json:"name"`
	Mode string `json:"mode"`
	// Scheme is http or https
	Scheme               string `json:"scheme"`
	Port                 uint16 `json:"port"`
	BasePath             string `json:"basePath"`
	AuthType             string `json:"authType"`
	TokenHeaderKey       string `json:"tokenHeaderKey"`
	GenTokenEndpoint     string `json:"genTokenEndpoint"`
	TokenKeyFromEndpoint string `json:"tokenKeyFromEndpoint"`
	Location             Server `json:"location"`
}

// MetricName returns a promehteus style name for the giving metric name.
func (srv Service) MetricName(metricName string) string {
	return prometheus.BuildFQName("skycoin", srv.Name, metricName)
}

// URIToGetMetric build the URI from where you will to get metric information.
func (srv Service) URIToGetMetric(metric Metric) string {
	return fmt.Sprintf("%s://%s:%d%s%s", srv.Scheme, srv.Location.Location, srv.Port, srv.BasePath, metric.URL)
}

// URIToGetMetric build the URI from where you will to get the exposed metrics.
func (srv Service) URIToGetExposedMetric() string {
	return fmt.Sprintf("%s://%s:%d%s", srv.Scheme, srv.Location.Location, srv.Port, srv.BasePath)
}

// URIToGetToken build the URI from where you will to get the token
func (srv Service) URIToGetToken() string {
	return fmt.Sprintf("%s://%s:%d%s%s", srv.Scheme, srv.Location.Location, srv.Port, srv.BasePath, srv.GenTokenEndpoint)
}

func (srv Service) validateProxy() (errs []error) {
	if !isValidURL(srv.URIToGetExposedMetric()) {
		errs = append(errs, errors.New("can not create a valid url to get the exposed metric."))
	}
	return errs
}

func (srv Service) validateApiRest() (errs []error) {
	if !isValidURL(srv.URIToGetToken()) {
		errs = append(errs, errors.New("can not create a valid url to get token: "+srv.URIToGetToken()))
	}
	for _, metric := range Config().Metrics {
		if !isValidURL(srv.URIToGetMetric(metric)) {
			errs = append(errs, errors.New("can not create a valid url to get metric: "+srv.URIToGetMetric(metric)))
		}
	}
	if strings.Compare(srv.AuthType, "CSRF") == 0 && len(srv.TokenHeaderKey) == 0 {
		errs = append(errs, errors.New("TokenHeaderKey is required if you are using CSRF"))
	}
	if strings.Compare(srv.AuthType, "CSRF") == 0 && len(srv.TokenKeyFromEndpoint) == 0 {
		errs = append(errs, errors.New("TokenKeyFromEndpoint is required if you are using CSRF"))
	}
	if strings.Compare(srv.AuthType, "CSRF") == 0 && len(srv.GenTokenEndpoint) == 0 {
		errs = append(errs, errors.New("GenTokenEndpoint is required if you are using CSRF"))
	}
	return errs
}

func (srv Service) validate() (errs []error) {
	if len(srv.Name) == 0 {
		errs = append(errs, errors.New("name is required in service"))
	}
	if len(srv.Scheme) == 0 {
		errs = append(errs, errors.New("scheme is required in service"))
	}
	if srv.Port < 1 || srv.Port > 65535 {
		errs = append(errs, errors.New("port must be betwen 1 and 65535"))
	}
	// if len(srv.BasePath) == 0 {
	// 	// TODO(denisacosta): What make sense in this?
	// }
	switch srv.Mode {
	case "proxy":
		errs = append(errs, srv.validateProxy()...)
	case "apiRest":
		errs = append(errs, srv.validateApiRest()...)
	default:
		if len(srv.Mode) == 0 {
			errs = append(errs, errors.New("mode is required in service"))
		} else if strings.Compare(srv.Mode, "proxy") != 0 && strings.Compare(srv.Mode, "apiRest") != 0 {
			errs = append(errs, errors.New("mode should be proxy or apiRest"))
		}
	}
	errs = append(errs, srv.Location.validate()...)
	return errs
}
