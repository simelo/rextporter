package memconfig

import (
	"github.com/simelo/rextporter/src/core"
)

// MetricDef implements the interface core.RextMetricDef
type MetricDef struct {
	name        string
	mType       string
	description string
	labels      []core.RextKeyValueStore
	options     core.RextKeyValueStore
}

// NewMetricFromDef create a MetricDef from a core.RextMetricDef
func NewMetricFromDef(m core.RextMetricDef) *MetricDef {
	ops, err := m.GetOptions().Clone()
	if err != nil {
		// FIXME(denisacostaq@gmail.com): Handle this error
	}
	return &MetricDef{
		name:        m.GetMetricName(),
		mType:       m.GetMetricType(),
		description: m.GetMetricDescription(),
		// Labels:      m.GetMetricLabels(),
		options: ops,
	}
}

// GetMetricName return the metric name
func (m MetricDef) GetMetricName() string {
	return m.name
}

// GetMetricType return the metric type, one of: gauge, counter, histogram or summary
func (m MetricDef) GetMetricType() string {
	return m.mType
}

// GetMetricDescription return the metric description
func (m MetricDef) GetMetricDescription() string {
	return m.description
}

// SetMetricName can set the metric name
func (m *MetricDef) SetMetricName(name string) {
	m.name = name
}

// SetMetricType can set the metric type
func (m *MetricDef) SetMetricType(tp string) {
	m.mType = tp
}

// SetMetricDescription can set the metric description
func (m *MetricDef) SetMetricDescription(description string) {
	m.description = description
}

// GetMetricLabels return some key/value pair with label name and path
func (m MetricDef) GetMetricLabels() []core.RextKeyValueStore {
	return m.labels
}

// SetMetricLabels receive some key/value pair with label name and path
func (m *MetricDef) SetMetricLabels(labels []core.RextKeyValueStore) {
	m.labels = labels
}

// GetOptions return key/value pairs for extra options
func (m *MetricDef) GetOptions() core.RextKeyValueStore {
	return m.options
}

// NewMetricDef create a new metric definition
func NewMetricDef(name, mType, description string, options core.RextKeyValueStore, labels []core.RextKeyValueStore) *MetricDef {
	return &MetricDef{
		name:        name,
		mType:       mType,
		description: description,
		labels:      labels,
		options:     options,
	}
}
