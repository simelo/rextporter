package scrapper

import (
	"fmt"

	"github.com/simelo/rextporter/src/client"
	"github.com/simelo/rextporter/src/util"
)

// Numeric scrapper can get numeric(gauges or counters) metrics
type Numeric struct {
	baseScrapper
}

func newNumeric(cl client.Client, p BodyParser, path string) Scrapper {
	return Numeric{baseScrapper: baseScrapper{client: cl, parser: p, jsonPath: path}}
}

// GetMetric returns a single number with the metric value, is a counter or a gauge
func (n Numeric) GetMetric() (val interface{}, err error) {
	const generalScopeErr = "error scrapping numeric(gauge|counter) metric"
	var iBody interface{}
	if iBody, err = getData(n.client, n.parser); err != nil {
		errCause := "client can not decode the body"
		return val, util.ErrorFromThisScope(errCause, generalScopeErr)
	}
	if val, err = n.parser.pathLookup(n.jsonPath, iBody); err != nil {
		errCause := fmt.Sprintln("can not get node: ", err.Error())
		return nil, util.ErrorFromThisScope(errCause, generalScopeErr)
	}
	return val, err
}
