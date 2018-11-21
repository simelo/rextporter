package config

import (
	"errors"
	"fmt"
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

// LabelNames return a slice with all the labels name
func (m Metric) LabelNames() []string {
	labelNames := make([]string, len(m.Options.Labels))
	for idxLabel, label := range m.Options.Labels {
		labelNames[idxLabel] = label.Name
	}
	return labelNames
}

func (metric Metric) isHistogram() bool {
	hasBuckets := len(metric.HistogramOptions.ExponentialBuckets) != 0 || len(metric.HistogramOptions.Buckets) != 0
	return hasBuckets || metric.Options.Type == "Histogram"
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
	if metric.HistogramOptions.inferType() == "Histogram" && metric.Options.Type != "Histogram" {
		errs = append(errs, errors.New("the buckets, only apply for metrics of type histogram"))
	}
	errs = append(errs, metric.Options.validate()...)
	if metric.isHistogram() {
		errs = append(errs, metric.HistogramOptions.validate()...)
	}
	return errs
}

// Label to create metrics grouping by json path value, for example:
// {Name: "color", "/properties/color"}
type Label struct {
	// Name the name of the label, different values can be assigned to it
	Name string
	// Path the json path from where you need to extract the label
	Path string
}

func (l *Label) validate() (errs []error) {
	if len(l.Name) == 0 {
		errs = append(errs, errors.New("Name is required in metric"))
	}
	if len(l.Path) == 0 {
		errs = append(errs, errors.New("Path is required in metric"))
	}
	return errs
}

// MetricOptions keep information you about the metric, mostly the type(Counter, Gauge, Summary, and Histogram)
type MetricOptions struct {
	Type        string `json:"type"`
	ItemPath    string `json:"item_path"`
	Description string `json:"description"`
	Labels      []Label
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
	if len(mo.ItemPath) == 0 && len(mo.Labels) != 0 {
		errs = append(errs, errors.New("if you define labels this is a vector and itemPath config is required"))
	}
	if len(mo.ItemPath) != 0 && len(mo.Labels) == 0 {
		errs = append(errs, errors.New("if you define itemPath this is a vector and labels config is required"))
	}
	for _, label := range mo.Labels {
		errs = append(errs, label.validate()...)
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
