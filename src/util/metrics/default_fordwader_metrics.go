package metrics

import "github.com/prometheus/client_golang/prometheus"

type DefaultFordwaderMetrics struct {
	FordwaderDatasourceResponseDuration *prometheus.GaugeVec
}

var fDefMetrics DefaultFordwaderMetrics

func NewDefaultFordwaderMetrics() (fordwaderMetrics *DefaultFordwaderMetrics) {
	fordwaderMetrics = &DefaultFordwaderMetrics{
		FordwaderDatasourceResponseDuration: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "fordwader_datasource_response_duration",
				Help: "Elapse time(in seconds) to get a response from a fordwader datasource",
			},
			[]string{"job", "instance", "datasource"},
		),
	}
	return fordwaderMetrics
}

func (fordwaderMetrics DefaultFordwaderMetrics) MustRegister() {
	prometheus.MustRegister(fordwaderMetrics.FordwaderDatasourceResponseDuration)
}
