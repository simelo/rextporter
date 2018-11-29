package scrapper

import (
	"fmt"

	"github.com/simelo/rextporter/src/client"
	"github.com/simelo/rextporter/src/config"
	"github.com/simelo/rextporter/src/util"
)

// NumericVec implements the Client interface(is able to get numeric metrics through `GetMetric` like Gauge and Counter)
type NumericVec struct {
	baseScrapper
	labels     []config.Label
	labelsName []string
	itemPath   string
}

func newNumericVec(cf client.ClientFactory, p BodyParser, metric config.Metric) Scrapper {
	return NumericVec{
		baseScrapper: baseScrapper{clientFactory: cf, parser: p, jsonPath: metric.Path},
		labels:       metric.Options.Labels,
		labelsName:   metric.LabelNames(),
		itemPath:     metric.Options.ItemPath}
}

// NumericVecItemVal can instances a numeric(Gauge or Counter) vec item with the required labels values
type NumericVecItemVal struct {
	Val    float64
	Labels []string
}

// NumericVecVals can instances a numeric(Gauge or Counter) vec with values and corresponding labels
type NumericVecVals []NumericVecItemVal

// GetMetric returns a numeric(Gauge or Counter) vector metric by using remote data.
func (nv NumericVec) GetMetric() (val interface{}, err error) {
	const generalScopeErr = "error scrapping numeric vec(gauge|counter) metric vec"
	var iBody interface{}
	if iBody, err = getData(nv.clientFactory, nv.parser); err != nil {
		errCause := "numeric vec client can not decode the body"
		return val, util.ErrorFromThisScope(errCause, generalScopeErr)
	}
	var iValColl interface{}
	if iValColl, err = nv.parser.pathLookup(nv.jsonPath, iBody); err != nil {
		errCause := fmt.Sprintln("can not get collection node: ", err.Error())
		return nil, util.ErrorFromThisScope(errCause, generalScopeErr)
	}
	metricCollection, okMetricCollection := iValColl.([]interface{})
	if !okMetricCollection {
		errCause := fmt.Sprintln("can not assert the metric type as slice")
		return nil, util.ErrorFromThisScope(errCause, generalScopeErr)
	}
	metricsVal := make(NumericVecVals, len(metricCollection))
	for idxMetric, metricItem := range metricCollection {
		var iMetricVal interface{}
		if iMetricVal, err = nv.parser.pathLookup(nv.itemPath, metricItem); err != nil {
			errCause := fmt.Sprintln("can not locate the metric val: ", err.Error())
			return nil, util.ErrorFromThisScope(errCause, generalScopeErr)
		}
		metricVal, okMetricVal := iMetricVal.(float64)
		if !okMetricVal {
			errCause := fmt.Sprintf("can not assert metric val %+v as float64", iMetricVal)
			return nil, util.ErrorFromThisScope(errCause, generalScopeErr)
		}
		metricsVal[idxMetric].Val = metricVal
		metricsVal[idxMetric].Labels = make([]string, len(nv.labels))
		for idxLabel, label := range nv.labels {
			var iLabelVal interface{}
			if iLabelVal, err = nv.parser.pathLookup(label.Path, metricItem); err != nil {
				errCause := fmt.Sprintln("can not locate the path for label: ", err.Error())
				return nil, util.ErrorFromThisScope(errCause, generalScopeErr)
			}
			labelVal, okLabelVal := iLabelVal.(string)
			if !okLabelVal {
				errCause := fmt.Sprintf("can not assert metric label %s %+v as string", label.Name, iLabelVal)
				return nil, util.ErrorFromThisScope(errCause, generalScopeErr)
			}
			metricsVal[idxMetric].Labels[idxLabel] = labelVal
		}
	}
	return metricsVal, nil
}
