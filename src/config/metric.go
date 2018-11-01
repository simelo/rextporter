package config

import "errors"

const (
	// KeyTypeCounter is the key you should define in the config file for counters.
	KeyTypeCounter = "Counter"
	// KeyTypeGauge is the key you should define in the config file for gauges.
	KeyTypeGauge = "Gauge"
	// KeyTypeHistogram is the key you should define in the config file for histograms.
	KeyTypeHistogram = "Histogram"
	// KeyTypeSummary is the key you should define in the config file for summaries.
	KeyTypeSummary = "Summary"
)

// Metric keep the metric name as an instance of MetricOptions
type Metric struct {
	Name    string        `json:"name"`
	Options MetricOptions `json:"options"`
}

func (metric Metric) validate() (errs []error) {
	if len(metric.Name) == 0 {
		errs = append(errs, errors.New("name is required in metric"))
	}
	errs = append(errs, metric.Options.validate()...)
	return errs
}

// MetricOptions keep information you about the metric, mostly the type(Counter, Gauge, Summary, and Histogram)
type MetricOptions struct {
	Type        string `json:"type"`
	Description string `json:"description"`
}

func (mo MetricOptions) validate() (errs []error) {
	if len(mo.Type) == 0 {
		errs = append(errs, errors.New("type is required in metric"))
	}
	return errs
}
