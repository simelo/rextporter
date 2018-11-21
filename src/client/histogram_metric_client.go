package client

import (
	"fmt"
	"net/http"

	"github.com/simelo/rextporter/src/config"
	"github.com/simelo/rextporter/src/util"
)

// HistogramClientOptions hold the necessary reference bucket to create an histogram
type HistogramClientOptions struct {
	Buckets []float64
}

// HistogramMetricClient implements the Client interface(is able to get histogram metrics through `GetMetric`)
type HistogramMetricClient struct {
	NumericClient
	histogramClientOptions HistogramClientOptions
}

func createHistogramClient(metric config.Metric, service config.Service) (client HistogramMetricClient, err error) {
	const generalScopeErr = "error creating histogram client"
	client = HistogramMetricClient{}
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

// getHistogramValue will return a histogram data structure from the remote endpoint.
func (client HistogramMetricClient) getHistogramValue() (val HistogramValue, err error) {
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

// GetMetric returns a histogram metric by using remote data.
func (client HistogramMetricClient) GetMetric() (val interface{}, err error) {
	val, err = client.getHistogramValue()
	return val, err
}
