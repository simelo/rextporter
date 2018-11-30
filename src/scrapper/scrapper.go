package scrapper

import (
	"errors"

	"github.com/simelo/rextporter/src/client"
	"github.com/simelo/rextporter/src/config"
	"github.com/simelo/rextporter/src/util"
)

// Scrapper get metrics from raw data
type Scrapper interface {
	// GetMetric receive some data as input and should return the metric val
	GetMetric() (val interface{}, err error)
}

type baseScrapper struct {
	clientFactory client.ClientFactory
	parser        BodyParser
	jsonPath      string
}

// BodyParser decode body from different formats, an get some data node
type BodyParser interface {
	decodeBody(body []byte) (val interface{}, err error)
	pathLookup(path string, val interface{}) (node interface{}, err error)
}

// NewScrapper will put all the required info to scrap metrics from the body returned by the client.
func NewScrapper(cf client.ClientFactory, parser BodyParser, metric config.Metric) (Scrapper, error) {
	if len(metric.LabelNames()) > 0 {
		return createVecScrapper(cf, parser, metric)
	}
	return createAtomicScrapper(cf, parser, metric)
}

func createVecScrapper(cf client.ClientFactory, parser BodyParser, metric config.Metric) (Scrapper, error) {
	if metric.Options.Type == config.KeyTypeCounter || metric.Options.Type == config.KeyTypeGauge {
		return newNumericVec(cf, parser, metric), nil
	}
	return NumericVec{}, errors.New("histogram vec and summary vec are not supported yet")
}

func createAtomicScrapper(cf client.ClientFactory, parser BodyParser, metric config.Metric) (Scrapper, error) {
	if metric.Options.Type == config.KeyTypeSummary {
		return Histogram{}, errors.New("summary scrapper is not supported yet")
	}
	if metric.Options.Type == config.KeyTypeHistogram {
		return newHistogram(cf, parser, metric), nil
	}
	return newNumeric(cf, parser, metric.Path), nil
}

func getData(cf client.ClientFactory, p BodyParser) (data interface{}, err error) {
	const generalScopeErr = "error getting data"
	var cl client.Client
	if cl, err = cf.CreateClient(); err != nil {
		errCause := "can ot create client"
		return data, util.ErrorFromThisScope(errCause, generalScopeErr)
	}
	respC := make(chan []byte)
	defer close(respC)
	errC := make(chan error)
	defer close(errC)
	workPool.Apply(client.RequestInfo{Client: cl, Res: respC, Err: errC})
	select {
	case body := <-respC:
		if data, err = p.decodeBody(body); err != nil {
			errCause := "scrapper can not decode the body"
			return data, util.ErrorFromThisScope(errCause, generalScopeErr)
		}
		return data, err
	case err = <-errC:
		errCause := "client can not get data"
		return data, util.ErrorFromThisScope(errCause, generalScopeErr)
	}
}
