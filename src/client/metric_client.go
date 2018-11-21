package client

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/oliveagle/jsonpath"
	"github.com/simelo/rextporter/src/config"
	"github.com/simelo/rextporter/src/util"
)


// HistogramClientOptions hold the necessary reference bucket to create an histogram
type HistogramClientOptions struct {
	Buckets []float64
}

// MetricClient implements the getRemoteInfo method from `client.Client` interface by using some `.toml` config parameters
// like for example: where is the host? it should be a GET, a POST or some other? ...
// sa NewMetricClient method.
type MetricClient struct {
	BaseClient
	metricJPath            string
	histogramClientOptions HistogramClientOptions
}

// NewMetricClient will put all the required info to be able to do http requests to get the remote data.
func NewMetricClient(metric config.Metric, service config.Service) (client *MetricClient, err error) {
	const generalScopeErr = "error creating a client to get a metric from remote endpoint"
	if strings.Compare(service.Mode, config.ServiceTypeAPIRest) != 0 {
		return client, errors.New("can not create an api rest metric client from a service of type " + service.Mode)
	}
	client = new(MetricClient)
	client.BaseClient.service = service
	client.metricJPath = metric.Path
	if strings.Compare(metric.Options.Type, config.KeyTypeHistogram) == 0 {
		client.histogramClientOptions = HistogramClientOptions{
			Buckets: metric.HistogramOptions.Buckets,
		}
	}
	client.BaseClient.req, err = http.NewRequest(metric.HTTPMethod, client.service.URIToGetMetric(metric), nil)
	if err != nil {
		errCause := fmt.Sprintln("can not create the request: ", err.Error())
		return nil, util.ErrorFromThisScope(errCause, generalScopeErr)
	}
	return client, nil
}

// GetMetric returns the metric previously bound through config parameters like:
// url(endpoint), json path, type and so on.
func (client *MetricClient) GetMetric() (val interface{}, err error) {
	const generalScopeErr = "error getting metric data"
	var data []byte
	if data, err = client.getRemoteInfo(); err != nil {
		return nil, util.ErrorFromThisScope(err.Error(), generalScopeErr)
	}
	var jsonData interface{}
	if err = json.Unmarshal(data, &jsonData); err != nil {
		errCause := fmt.Sprintln("can not decode the body: ", string(data), " ", err.Error())
		return nil, util.ErrorFromThisScope(errCause, generalScopeErr)
	}
	jPath := "$" + strings.Replace(client.metricJPath, "/", ".", -1)
	if val, err = jsonpath.JsonPathLookup(jsonData, jPath); err != nil {
		errCause := fmt.Sprintln("can not locate the path: ", err.Error())
		return nil, util.ErrorFromThisScope(errCause, generalScopeErr)
	}
	return val, nil
}

// GetHistogramValue will return a histogram data structure from the remote endpoint.
func (client *MetricClient) GetHistogramValue() (val HistogramValue, err error) {
	generalScopeErr := "error getting histogram values"
	var metric interface{}
	if metric, err = client.GetMetric(); err != nil {
		errCause := fmt.Sprintln("can not get the metric data: ", err.Error())
		return val, util.ErrorFromThisScope(errCause, generalScopeErr)
	}
	metricCollection, okMetricCollection := metric.([]interface{})
	if !okMetricCollection {
		errCause := fmt.Sprintln("can not assert the metric type as slice")
		return val, util.ErrorFromThisScope(errCause, generalScopeErr)
	}
	val = newHistogram(client.histogramClientOptions.Buckets)
	for _, metricItem := range metricCollection {
		val.Count++
		metricVal, okMetricVal := metricItem.(float64)
		if !okMetricVal {
			errCause := fmt.Sprintln("can not assert the metric value to type float")
			return val, util.ErrorFromThisScope(errCause, generalScopeErr)
		}
		val.Sum += metricVal
		for _, bucket := range client.histogramClientOptions.Buckets {
			if bucket <= metricVal {
				val.Buckets[bucket]++
			}
		}
	}
	return val, err
}
