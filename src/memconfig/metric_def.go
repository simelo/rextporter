package memconfig

import (
	"github.com/simelo/rextporter/src/core"
)

// MetricDef implements the interface core.RextMetricDef
type MetricDef struct {
	name        string
	mType       string
	nodeSolver  core.RextNodeSolver
	description string
	labels      []core.RextLabelDef
	options     core.RextKeyValueStore
}

// GetMetricName return the metric name
func (m MetricDef) GetMetricName() string {
	return m.name
}

// GetMetricType return the metric type, one of: gauge, counter, histogram or summary
func (m MetricDef) GetMetricType() string {
	return m.mType
}

// GetNodeSolver return solver type
func (dp MetricDef) GetNodeSolver() core.RextNodeSolver {
	return dp.nodeSolver
}

// SetNodeSolver set the node solver
func (dp *MetricDef) SetNodeSolver(nodeSolver core.RextNodeSolver) {
	dp.nodeSolver = nodeSolver
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

// GetLabels return labels
func (m MetricDef) GetLabels() []core.RextLabelDef {
	return m.labels
}

// AddLabel receive label to be append to the current list
func (m *MetricDef) AddLabel(label core.RextLabelDef) {
	m.labels = append(m.labels, label)
}

// GetOptions return key/value pairs for extra options
func (m *MetricDef) GetOptions() core.RextKeyValueStore {
	if m.options == nil {
		m.options = NewOptionsMap()
	}
	return m.options
}

// NewMetricDef create a new metric definition
func NewMetricDef(name, mType, description string, nodeSolver core.RextNodeSolver, options core.RextKeyValueStore, labels []core.RextLabelDef) *MetricDef {
	return &MetricDef{
		name:        name,
		mType:       mType,
		nodeSolver:  nodeSolver,
		description: description,
		labels:      labels,
		options:     options,
	}
}
