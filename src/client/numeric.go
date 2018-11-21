package client

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/oliveagle/jsonpath"
	"github.com/simelo/rextporter/src/config"
	"github.com/simelo/rextporter/src/util"
)

// Numeric implements the GetMetric method from `client.Client` interface by using some `.toml` config parameters
// like for example: where is the host? it should be a GET, a POST or some other? ...
// sa NewMetricClient method.
type Numeric struct {
	BaseClient
}

func createNumeric(metric config.Metric, service config.Service) (client Numeric, err error) {
	const generalScopeErr = "error creating number(gauge | counter) client"
	client = Numeric{}
	client.BaseClient.service = service
	client.metricJPath = metric.Path
	client.BaseClient.req, err = http.NewRequest(metric.HTTPMethod, client.service.URIToGetMetric(metric), nil)
	if err != nil {
		errCause := fmt.Sprintln("can not create the request: ", err.Error())
		return client, util.ErrorFromThisScope(errCause, generalScopeErr)
	}
	return client, nil
}

// GetMetric returns a numeric(Gauge or Counter) metric by using remote data.
func (client Numeric) GetMetric() (val interface{}, err error) {
	const generalScopeErr = "error getting metric data"
	var data []byte
	if data, err = client.getRemoteInfo(); err != nil {
		errCause := fmt.Sprintln("can not get remote info: ", err.Error())
		return nil, util.ErrorFromThisScope(errCause, generalScopeErr)
	}
	var jsonData interface{}
	if err = json.Unmarshal(data, &jsonData); err != nil {
		errCause := fmt.Sprintf("can not decode the body: %s. Err: %s", string(data), err.Error())
		return nil, util.ErrorFromThisScope(errCause, generalScopeErr)
	}
	jPath := "$" + strings.Replace(client.BaseClient.metricJPath, "/", ".", -1)
	if val, err = jsonpath.JsonPathLookup(jsonData, jPath); err != nil {
		errCause := fmt.Sprintln("can not locate the path: ", err.Error())
		return nil, util.ErrorFromThisScope(errCause, generalScopeErr)
	}
	return val, nil
}
