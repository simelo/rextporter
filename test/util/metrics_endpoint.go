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
				log.WithFields(log.Fields{
					"current_type": (*mf.Type).String(),
					"looking_for":  strings.ToUpper(core.KeyMetricTypeHistogram),
				}).Errorln("metric is not one of the spected type")
				return -1, core.ErrKeyInvalidType
			}
			return *mf.Metric[0].Gauge.Value, nil
		}
	}
	log.WithField("wanted_name", mtrName).Errorln("metric name not found")
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
				log.WithFields(log.Fields{
					"current_type": (*mf.Type).String(),
					"looking_for":  strings.ToUpper(core.KeyMetricTypeHistogram),
				}).Errorln("metric is not one of the spected type")
				return -1, core.ErrKeyInvalidType
			}
			return *mf.Metric[0].Counter.Value, nil
		}
	}
	log.WithField("wanted_name", mtrName).Errorln("metric name not found")
	return -1, core.ErrKeyNotFound
}

type HistogramValue struct {
	SampleCount uint64
	SampleSum   float64
	Buckets     map[float64]uint64
}

func GetHistogramValue(metrics []byte, mtrName string) (HistogramValue, error) {
	var parser expfmt.TextParser
	in := bytes.NewReader(metrics)
	metricFamilies, err := parser.TextToMetricFamilies(in)
	if err != nil {
		log.WithError(err).Errorln("error, reading text format failed")
		return HistogramValue{}, core.ErrKeyDecodingFile
	}
	for _, mf := range metricFamilies {
		if mtrName == *mf.Name {
			if (*mf.Type).String() != strings.ToUpper(core.KeyMetricTypeHistogram) {
				log.WithFields(log.Fields{
					"current_type": (*mf.Type).String(),
					"looking_for":  strings.ToUpper(core.KeyMetricTypeHistogram),
				}).Errorln("metric is not one of the spected type")
				return HistogramValue{}, core.ErrKeyInvalidType
			}
			mtr := *mf.Metric[0]
			hv := HistogramValue{
				SampleCount: *mtr.Histogram.SampleCount,
				SampleSum:   *mtr.Histogram.SampleSum,
				Buckets:     make(map[float64]uint64),
			}
			for _, b := range mtr.Histogram.Bucket {
				hv.Buckets[*b.UpperBound] = *b.CumulativeCount
			}
			return hv, nil
		}
	}
	log.WithField("wanted_name", mtrName).Errorln("metric name not found")
	return HistogramValue{}, core.ErrKeyNotFound
}

type LabelValue struct {
	Name  string
	Value string
}

type MetricVal struct {
	Labels []LabelValue
	Number float64
}

type NumericVec struct {
	Values []MetricVal
}

func GetNumericVecValues(metrics []byte, mtrName string) (NumericVec, error) {
	var parser expfmt.TextParser
	in := bytes.NewReader(metrics)
	metricFamilies, err := parser.TextToMetricFamilies(in)
	if err != nil {
		log.WithError(err).Errorln("error, reading text format failed")
		return NumericVec{}, core.ErrKeyDecodingFile
	}
	for _, mf := range metricFamilies {
		if mtrName == *mf.Name {
			if (*mf.Type).String() != strings.ToUpper(core.KeyMetricTypeGauge) && (*mf.Type).String() != strings.ToUpper(core.KeyMetricTypeCounter) {
				log.WithFields(log.Fields{
					"current_type": (*mf.Type).String(),
					"looking_for":  strings.ToUpper(core.KeyMetricTypeHistogram),
				}).Errorln("metric is not one of the spected type")
				return NumericVec{}, core.ErrKeyInvalidType
			}
			nv := NumericVec{}
			for _, mtrVal := range mf.Metric {
				var val float64
				if (*mf.Type).String() == strings.ToUpper(core.KeyMetricTypeGauge) {
					val = *((*mtrVal).Gauge.Value)
				} else {
					val = *((*mtrVal).Counter.Value)
				}
				mv := MetricVal{Number: val}
				for _, label := range (*mtrVal).Label {
					lv := LabelValue{Name: *label.Name, Value: *label.Value}
					mv.Labels = append(mv.Labels, lv)
				}
				nv.Values = append(nv.Values, mv)
			}
			return nv, nil
		}
	}
	log.WithField("wanted_name", mtrName).Errorln("metric name not found")
	return NumericVec{}, core.ErrKeyNotFound
}
