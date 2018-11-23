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

// NumericVec implements the Client interface(is able to get numeric metrics through `GetMetric` like Gauge and Counter)
type NumericVec struct {
	BaseClient
	labels     []config.Label
	labelsName []string
	itemPath   string
}

// NumericVecItemVal can instances a numeric(Gauge or Counter) vec item with the required labels values
type NumericVecItemVal struct {
	Val    float64
	Labels []string
}

// NumericVecVals can instances a numeric(Gauge or Counter) vec with values and corresponding labels
type NumericVecVals []NumericVecItemVal

func createNumericVec(metric config.Metric, service config.Service) (client NumericVec, err error) {
	const generalScopeErr = "error creating metric vec client"
	client = NumericVec{}
	client.BaseClient.service = service
	client.metricJPath = metric.Path
	client.BaseClient.req, err = http.NewRequest(metric.HTTPMethod, client.service.URIToGetMetric(metric), nil)
	if err != nil {
		errCause := fmt.Sprintln("can not create the request: ", err.Error())
		return client, util.ErrorFromThisScope(errCause, generalScopeErr)
	}
	client.itemPath = metric.Options.ItemPath
	client.labels = metric.Options.Labels
	client.labelsName = metric.LabelNames()
	return client, nil
}

// GetMetric returns a numeric(Gauge or Counter) vector metric by using remote data.
func (client NumericVec) GetMetric() (interface{}, error) {
	const generalScopeErr = "error getting metric vec data"
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
	metricsVal := make(NumericVecVals, len(metricCollection))
	for idxMetric, metric := range metricCollection {
		mJPath := "$" + strings.Replace(client.itemPath, "/", ".", -1)
		var iMetricVal interface{}
		if iMetricVal, err = jsonpath.JsonPathLookup(metric, mJPath); err != nil {
			errCause := fmt.Sprintln("can not locate the metric val: ", err.Error())
			return nil, util.ErrorFromThisScope(errCause, generalScopeErr)
		}
		metricVal, okMetricVal := iMetricVal.(float64)
		if !okMetricVal {
			errCause := fmt.Sprintf("can not assert metric val %+v as float64", iMetricVal)
			return nil, util.ErrorFromThisScope(errCause, generalScopeErr)
		}
		metricsVal[idxMetric].Val = metricVal
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
