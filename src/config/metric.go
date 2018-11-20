package config

import (
	"errors"
	"fmt"
	"strings"
)

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
	Name             string           `json:"name"`
	URL              string           `json:"url"`
	HTTPMethod       string           `json:"http_method"`
	Path             string           `json:"path,omitempty"`
	Options          MetricOptions    `json:"options"`
	HistogramOptions HistogramOptions `json:"histogram_options"`
}

func (metric Metric) isHistogram() bool {
	hasBuckets := len(metric.HistogramOptions.ExponentialBuckets) != 0 || len(metric.HistogramOptions.Buckets) != 0
	return hasBuckets || strings.Compare(metric.Options.Type, "Histogram") == 0
}

func (metric Metric) validate() (errs []error) {
	if len(metric.Name) == 0 {
		errs = append(errs, errors.New("name is required in metric"))
	}
	if len(metric.URL) == 0 {
		errs = append(errs, errors.New("url is required in metric"))
	}
	if len(metric.HTTPMethod) == 0 {
		errs = append(errs, errors.New("HttpMethod is required in metric"))
	}
	if len(metric.Path) == 0 {
		errs = append(errs, errors.New("path is required in metric"))
	}
	if strings.Compare(metric.HistogramOptions.inferType(), "Histogram") == 0 && strings.Compare(metric.Options.Type, "Histogram") != 0 {
		errs = append(errs, errors.New("the buckets, only apply for metrics of type histogram"))
	}
	errs = append(errs, metric.Options.validate()...)
	if metric.isHistogram() {
		errs = append(errs, metric.HistogramOptions.validate()...)
	}
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
	switch mo.Type {
	case KeyTypeCounter, KeyTypeGauge, KeyTypeHistogram:
	case KeyTypeSummary:
		errs = append(errs, fmt.Errorf("type %s is not supported yet", KeyTypeSummary))
	default:
		errs = append(errs, fmt.Errorf("type should be one of %s, %s, %s or %s", KeyTypeCounter, KeyTypeGauge, KeyTypeSummary, KeyTypeHistogram))
	}
	return errs
}

// HistogramOptions allows you to define the histogram is buckets.
type HistogramOptions struct {
	Buckets []float64 `json:"buckets"`

	// ExponentialBuckets is a len three array where:
	// - The first value is the low bound start bucket.
	// - The second vale is the growing factor.
	// - The three one is the buckets amount.
	ExponentialBuckets []float64 `json:"exponential_buckets"`
}

func (ho HistogramOptions) validate() (errs []error) {
	if len(ho.Buckets) > 0 && len(ho.ExponentialBuckets) > 0 {
		errs = append(errs, errors.New("You should define only one betwen the 'buckets' and the 'exponentialBuckets'"))
	}
	if len(ho.Buckets) == 0 && len(ho.ExponentialBuckets) == 0 {
		errs = append(errs, errors.New("At least one should be defined the 'buckets' or the 'exponentialBuckets'"))
	}
	if len(ho.ExponentialBuckets) != 0 && len(ho.ExponentialBuckets) != 3 {
		errs = append(errs, errors.New("'exponentialBuckets' should have an exact length of 3(start, factor, amount)"))
	}
	return errs
}

func (ho HistogramOptions) inferType() (t string) {
	if len(ho.Buckets) != 0 || len(ho.ExponentialBuckets) != 0 {
		t = "Histogram"
	}
	return t
}
