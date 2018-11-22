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

// HistogramClientOptions hold the necessary reference bucket to create an histogram
type HistogramClientOptions struct {
	Buckets []float64
}

// Histogram implements the Client interface(is able to get histogram metrics through `GetMetric`)
type Histogram struct {
	Numeric
	histogramClientOptions HistogramClientOptions
}

func createHistogram(metric config.Metric, service config.Service) (client Histogram, err error) {
	const generalScopeErr = "error creating histogram client"
	client = Histogram{}
	client.BaseClient.service = service
	client.metricJPath = metric.Path
	client.BaseClient.req, err = http.NewRequest(metric.HTTPMethod, client.service.URIToGetMetric(metric), nil)
	if err != nil {
		errCause := fmt.Sprintln("can not create the request: ", err.Error())
		return client, util.ErrorFromThisScope(errCause, generalScopeErr)
	}
	client.histogramClientOptions = HistogramClientOptions{
		Buckets: metric.HistogramOptions.Buckets,
	}
	return client, nil
}

// HistogramValue hold the required values to create a histogram metric, the Count, Sum and buckets.
type HistogramValue struct {
	Count   uint64
	Sum     float64
	Buckets map[float64]uint64
}

func newHistogram(buckets []float64) HistogramValue {
	val := HistogramValue{
		Count:   0,
		Sum:     0,
		Buckets: make(map[float64]uint64, len(buckets)),
	}
	for _, bucket := range buckets {
		val.Buckets[bucket] = 0
	}
	return val
}

// GetMetric returns a histogram metric by using remote data.
func (client Histogram) GetMetric() (interface{}, error) {
	generalScopeErr := "error getting histogram values"
	var data []byte
	var err error
	if data, err = client.getRemoteInfo(); err != nil {
		errCause := fmt.Sprintln("can not get remote info: ", err.Error())
		return nil, util.ErrorFromThisScope(errCause, generalScopeErr)
	}
	var jsonData interface{}
	if err = json.Unmarshal(data, &jsonData); err != nil {
		errCause := fmt.Sprintf("can not decode the body: %s. Err: %s", string(data), err.Error())
		return nil, util.ErrorFromThisScope(errCause, generalScopeErr)
	}
	jPath := "$" + strings.Replace(client.metricJPath, "/", ".", -1)
	var val interface{}
	if val, err = jsonpath.JsonPathLookup(jsonData, jPath); err != nil {
		errCause := fmt.Sprintln("can not locate the path: ", err.Error())
		return nil, util.ErrorFromThisScope(errCause, generalScopeErr)
	}
	return createHistogramFromData(client.histogramClientOptions.Buckets, val)
}

func createHistogramFromData(buckets []float64, data interface{}) (val interface{}, err error) {
	generalScopeErr := "creating histogram from data"
	collection, okCollection := data.([]interface{})
	if !okCollection {
		errCause := fmt.Sprintln("can not assert the metric type as slice")
		return val, util.ErrorFromThisScope(errCause, generalScopeErr)
	}
	histogram := newHistogram(buckets)
	for _, item := range collection {
		histogram.Count++
		metricVal, okMetricVal := item.(float64)
		if !okMetricVal {
			errCause := fmt.Sprintf("can not assert the metric value %+v to type float", item)
			return val, util.ErrorFromThisScope(errCause, generalScopeErr)
		}
		histogram.Sum += metricVal
		for _, bucket := range buckets {
			if bucket <= metricVal {
				histogram.Buckets[bucket]++
			}
		}
	}
	val = histogram
	return val, err
}
