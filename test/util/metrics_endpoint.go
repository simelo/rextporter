package util

import (
	"bytes"
	"strings"

	"github.com/denisacostaq/rextporter/src/core"
	"github.com/prometheus/common/expfmt"
	log "github.com/sirupsen/logrus"
)

// FoundMetric return true if found a metric with mtrName inside a metrics endpoint response as plain text.
func FoundMetric(metrics []byte, mtrName string) (bool, error) {
	var parser expfmt.TextParser
	in := bytes.NewReader(metrics)
	metricFamilies, err := parser.TextToMetricFamilies(in)
	if err != nil {
		log.WithError(err).Errorln("error, reading text format failed")
		return false, core.ErrKeyDecodingFile
	}
	for _, mf := range metricFamilies {
		if mtrName == *mf.Name {
			return true, nil
		}
	}
	return false, err
}

// GetGaugeValue return a number(the gauge value) if found a metric with mtrName and if a Gauge kind metric
// from a metrics endpoint response as plain text.
func GetGaugeValue(metrics []byte, mtrName string) (float64, error) {
	var parser expfmt.TextParser
	in := bytes.NewReader(metrics)
	metricFamilies, err := parser.TextToMetricFamilies(in)
	if err != nil {
		log.WithError(err).Errorln("error, reading text format failed")
		return -1, core.ErrKeyDecodingFile
	}
	for _, mf := range metricFamilies {
		if mtrName == *mf.Name {
			if (*mf.Type).String() != strings.ToUpper(core.KeyMetricTypeGauge) {
				return -1, core.ErrKeyInvalidType
			}
			return *mf.Metric[0].Gauge.Value, nil
		}
	}
	return -1, core.ErrKeyNotFound
}

// GetCounterValue return a number(the counter value) if found a metric with mtrName and if a Counter kind metric
// from a metrics endpoint response as plain text.
func GetCounterValue(metrics []byte, mtrName string) (float64, error) {
	var parser expfmt.TextParser
	in := bytes.NewReader(metrics)
	metricFamilies, err := parser.TextToMetricFamilies(in)
	if err != nil {
		log.WithError(err).Errorln("error, reading text format failed")
		return -1, core.ErrKeyDecodingFile
	}
	for _, mf := range metricFamilies {
		if mtrName == *mf.Name {
			if (*mf.Type).String() != strings.ToUpper(core.KeyMetricTypeCounter) {
				return -1, core.ErrKeyInvalidType
			}
			return *mf.Metric[0].Counter.Value, nil
		}
	}
	return -1, core.ErrKeyNotFound
}
