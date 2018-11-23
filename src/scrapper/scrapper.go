package scrapper

import (
	"github.com/simelo/rextporter/src/client"
	"github.com/simelo/rextporter/src/util"
)

// Scrapper get metrics from raw data
type Scrapper interface {
	// GetMetric receive some data as input and should return the metric val
	GetMetric() (val interface{}, err error)
}

// BodyParser decode body from different formats, an get some data node
type BodyParser interface {
	decodeBody(body []byte) (val interface{}, err error)
	pathLookup(path string, val interface{}) (node interface{}, err error)
}

func getData(cl client.Client, p BodyParser) (data interface{}, err error) {
	const generalScopeErr = "error getting data"
	var body []byte
	if body, err = cl.GetData(); err != nil {
		errCause := "client can not get data"
		return data, util.ErrorFromThisScope(errCause, generalScopeErr)
	}
	if data, err = p.decodeBody(body); err != nil {
		errCause := "client can not decode the body"
		return data, util.ErrorFromThisScope(errCause, generalScopeErr)
	}
	return data, err
}
