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

// HistogramVec implements the Client interface(is able to get numeric metrics through `GetMetric` like Gauge and Counter)
type HistogramVec struct {
	BaseClient
	histogramClientOptions HistogramClientOptions
	labels                 []config.Label
	labelsName             []string
	itemPath               string
}

// HistogramVecItemVal can instance a histogram vec item with the required labels values
type HistogramVecItemVal struct {
	Val    HistogramValue
	Labels []string
}

// HistogramVecVals is a shortcut for HistogramVecItemVal slice(related to metric vector)
type HistogramVecVals []HistogramVecItemVal

func createHistogramVec(metric config.Metric, service config.Service) (client HistogramVec, err error) {
	const generalScopeErr = "error creating metric vec client"
	client = HistogramVec{}
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
	client.itemPath = metric.Options.ItemPath
	client.labels = metric.Options.Labels
	client.labelsName = metric.LabelNames()
	return client, nil
}

// GetMetric returns a Histogram vector metric by using remote data.
func (client HistogramVec) GetMetric() (interface{}, error) {
	const generalScopeErr = "error getting histogram vec data"
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
	cJPath := "$" + strings.Replace(client.BaseClient.metricJPath, "/", ".", -1)
	var iValColl interface{}
	if iValColl, err = jsonpath.JsonPathLookup(jsonData, cJPath); err != nil {
		errCause := fmt.Sprintln("can not locate the path: ", err.Error())
		return nil, util.ErrorFromThisScope(errCause, generalScopeErr)
	}
	metricCollection, okMetricCollection := iValColl.([]interface{})
	if !okMetricCollection {
		errCause := fmt.Sprintln("can not assert the metric type as slice")
		return nil, util.ErrorFromThisScope(errCause, generalScopeErr)
	}
	metricsVal := make(HistogramVecVals, len(metricCollection))
	for idxMetric, metric := range metricCollection {
		mJPath := "$" + strings.Replace(client.itemPath, "/", ".", -1)
		var iMetricVal interface{}
		if iMetricVal, err = jsonpath.JsonPathLookup(metric, mJPath); err != nil {
			errCause := fmt.Sprintln("can not locate the metric val: ", err.Error())
			return nil, util.ErrorFromThisScope(errCause, generalScopeErr)
		}
		histogram, err := createHistogramFromData(client.histogramClientOptions.Buckets, metric)
		if err != nil {
			errCause := fmt.Sprintf("can not create histogram value from data %+v", iMetricVal)
			return nil, util.ErrorFromThisScope(errCause, generalScopeErr)
		}
		metricsVal[idxMetric].Val = histogram
		metricsVal[idxMetric].Labels = make([]string, len(client.labels))
		for idxLabel, label := range client.labels {
			lJPath := "$" + strings.Replace(label.Path, "/", ".", -1)
			var iLabelVal interface{}
			if iLabelVal, err = jsonpath.JsonPathLookup(metric, lJPath); err != nil {
				errCause := fmt.Sprintln("can not locate the path for label: ", err.Error())
				return nil, util.ErrorFromThisScope(errCause, generalScopeErr)
			}
			labelVal, okLabelVal := iLabelVal.(string)
			if !okLabelVal {
				errCause := fmt.Sprintf("can not assert metric label %s %+v as string", label.Name, iMetricVal)
				return nil, util.ErrorFromThisScope(errCause, generalScopeErr)
			}
			metricsVal[idxMetric].Labels[idxLabel] = labelVal
		}
	}
	return metricsVal, nil
}
