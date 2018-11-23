package scrapper

import (
	"fmt"
	"strings"

	"github.com/oliveagle/jsonpath"
	"github.com/simelo/rextporter/src/client"
	"github.com/simelo/rextporter/src/util"
)

// Numeric scrapper can get numeric(gauges or counters) metrics
type Numeric struct {
	client   client.Client
	parser   BodyParser
	jsonPath string
}

func NewNumeric(cl, client.Client, p client.BodyParser, path string) Scrapper {
	return Numeric{client: cl, parser: p, jsonPath: jPath}
}

func (n Numeric) GetMetric() (val interface{}, err error) {
	const generalScopeErr = "error scrapping numeric(gauge|counter) metric"
	var iBody interface{}
	if iBody, err = getData(client, parser); err != nil {
		errCause := "client can not decode the body"
		return val, util.ErrorFromThisScope(errCause, generalScopeErr)
	}
	if val, err = parser.pathLookup(n.jsonPath, iBody); err != nil {
		errCause := fmt.Sprintln("can not get node: ", err.Error())
		return nil, util.ErrorFromThisScope(errCause, generalScopeErr)
	}

	return val, err
}
