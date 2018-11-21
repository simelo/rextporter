package client

import (
	"errors"

	"github.com/simelo/rextporter/src/config"
)

// Client is an http wrapper(implement the GetMetric).
type Client interface {
	// GetMetric return a metric val by using some `.toml` config parameters
	// like for example: where is the host? it should be a GET, a POST or some other? ...
	// sa NewMetricClient method.
	GetMetric() (val interface{}, err error)
}

// NewMetricClient will put all the required info to be able to do http requests to get the remote data.
func NewMetricClient(metric config.Metric, service config.Service) (Client, error) {
	const generalScopeErr = "error creating a client to get a metric from remote endpoint"
	if service.Mode != config.ServiceTypeAPIRest {
		return nil, errors.New("can not create an api rest metric client from a service of type " + service.Mode)
	}
	// BUG(denisacostaq@gmail.com): type can collide with labels, for example type histogram with values
	if metric.Options.Type == config.KeyTypeHistogram {
		return createHistogramClient(metric, service)
	}
	if len(metric.LabelNames()) > 0 {
		return createVecMetricClient(metric, service)
	}
	return createNumberClient(metric, service)
}

// TODO(denisacostaq@gmail.com): check out http://localhost:6060/pkg/github.com/prometheus/client_golang/api/#NewClient
