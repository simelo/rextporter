package memconfig

import "github.com/simelo/rextporter/src/core"

// MetricsExtractor implements the core.RextMetricsExtractor interface
type MetricsExtractor struct {
	metrics []*MetricDef
	options core.RextKeyValueStore
}

// Apply a rule for the extractor
func (me MetricsExtractor) Apply(rule core.RextMetricDef) (core.RextMetricDef, error) {
	var m *MetricDef
	if metricRule, isMetric := rule.(*MetricDef); isMetric {
		m = metricRule
	} else {
		m = NewMetricFromDef(rule)
	}
	me.metrics = append(me.metrics, m)
	return m, nil
}

// ApplyMany multiple rules for the extractor
func (me MetricsExtractor) ApplyMany(rules []core.RextMetricDef) ([]core.RextMetricDef, error) {
	newRules := make([]core.RextMetricDef, len(rules))
	var err error
	for idxRule, rule := range rules {
		if newRules[idxRule], err = me.Apply(rule); err != nil {
			return newRules, err
		}
	}
	return rules, err
}

// ExtractMetrics from a target
func (me MetricsExtractor) ExtractMetrics(target interface{}) ([]core.RextMetric, error) {
	return nil, nil
}

// GetOptions return key/value pairs for extra options
func (me MetricsExtractor) GetOptions() core.RextKeyValueStore {
	return me.options
}

// NewMetricsExtractor create a new MetricsExtractor
func NewMetricsExtractor(metrics []*MetricDef, options core.RextKeyValueStore) core.RextMetricsExtractor {
	return &MetricsExtractor{
		metrics: metrics,
		options: options,
	}
}
