package metrics

import "github.com/prometheus/client_golang/prometheus"

// DefaultFordwaderMetrics default metrics for metrics fordwader
type DefaultFordwaderMetrics struct {
	FordwaderResponseDuration      *prometheus.GaugeVec
	FordwaderScrapeDurationSeconds *prometheus.GaugeVec
}

// NewDefaultFordwaderMetrics create a new DefaultFordwaderMetrics
func NewDefaultFordwaderMetrics() (fordwaderMetrics *DefaultFordwaderMetrics) {
	fordwaderMetrics = &DefaultFordwaderMetrics{
		FordwaderResponseDuration: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "fordwader_response_duration_seconds",
				Help: "Elapse time(in seconds) to get a response from a fordwader",
			},
			[]string{"job", "instance"},
		),
		FordwaderScrapeDurationSeconds: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "fordwader_scrape_duration_seconds",
				Help: "Elapse time(in seconds) to get a response from a fordwader scrapper",
			},
			[]string{"job", "instance"},
		),
	}
	return fordwaderMetrics
}

// MustRegister register default metrics for metrics fordwader in prometheus
func (fordwaderMetrics DefaultFordwaderMetrics) MustRegister() {
	prometheus.MustRegister(fordwaderMetrics.FordwaderResponseDuration)
	prometheus.MustRegister(fordwaderMetrics.FordwaderScrapeDurationSeconds)
}
