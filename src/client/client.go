package client

import (
	"errors"
	"log"

	"github.com/simelo/rextporter/src/config"
	"github.com/simelo/rextporter/src/util"
)

// Client is an http wrapper(implement the GetMetric).
type Client interface {
	// GetMetric return a metric val by using some `.toml` config parameters
	// like for example: where is the host? it should be a GET, a POST or some other? ...
	// sa NewClient method.
	GetMetric() (val interface{}, err error)
}

// NewClient will put all the required info to be able to do http requests to get the remote data.
func NewClient(metric config.Metric, service config.Service) (Client, error) {
	const generalScopeErr = "error creating a client to get a metric from remote endpoint"
	if util.StrSliceContains(service.Modes, config.ServiceTypeAPIRest) {
		errCause := "can not create an api rest metric client from a service witout mode " + config.ServiceTypeAPIRest
		return nil, util.ErrorFromThisScope(errCause, generalScopeErr)
	}
	if len(metric.LabelNames()) > 0 {
		return createVecClient(metric, service)
	}
	return createAtomicClient(metric, service)
}

func createVecClient(metric config.Metric, service config.Service) (Client, error) {
	if metric.Options.Type == config.KeyTypeHistogram {
		// FIXME(denisacostaq@gmail.com): work on this feacture
		var v HistogramVec
		var err error
		if v, err = createHistogramVec(metric, service); err != nil {
			log.Println(v, err)
		}
		return v, errors.New("create histogram vec is not supported yet, see return createHistogramVec(metric, service)")
	} else if metric.Options.Type == config.KeyTypeSummary {
		return HistogramVec{}, errors.New("create summary vec is not supported yet")
	}
	return createNumericVec(metric, service)
}

func createAtomicClient(metric config.Metric, service config.Service) (Client, error) {
	if metric.Options.Type == config.KeyTypeHistogram {
		return createHistogram(metric, service)
	}
	return createNumeric(metric, service)
}

// TODO(denisacostaq@gmail.com): check out http://localhost:6060/pkg/github.com/prometheus/client_golang/api/#NewClient
