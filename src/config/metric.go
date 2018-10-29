package config

import "errors"

// Metric keep the metric name as an instance of MetricOptions
type Metric struct {
	Name    string        `json:"name"`
	Options MetricOptions `json:"options"`
}

// MetricOptions keep information you about the metric, mostly the type(Counter, Gauge, Summary, and Histogram)
type MetricOptions struct {
	Type        string `json:"type"`
	Description string `json:"description"`
}

func (metric Metric) validate() (errs []error) {
	if len(metric.Name) == 0 {
		errs = append(errs, errors.New("name is required in metric"))
	}
	errs = append(errs, metric.Options.validate()...)
	return errs
}

func (mo MetricOptions) validate() (errs []error) {
	if len(mo.Type) == 0 {
		errs = append(errs, errors.New("type is required in metric"))
	}
	return errs
}
