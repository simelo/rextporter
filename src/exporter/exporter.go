package exporter

import (
	"fmt"
	"log"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/simelo/rextporter/src/client"
	"github.com/simelo/rextporter/src/common"
	"github.com/simelo/rextporter/src/config"
)

// SkycoinCollector has the metrics to be exposed
type SkycoinCollector struct {
	Counters []CounterMetric
	Gauges   []GaugeMetric
}

func newSkycoinCollector() (collector *SkycoinCollector, err error) {
	const generalScopeErr = "error creating collector"
	var counters []CounterMetric
	counters, err = createCounters()
	if err != nil {
		errCause := fmt.Sprintln("error creating counters: ", err.Error())
		return nil, common.ErrorFromThisScope(errCause, generalScopeErr)
	}
	var gauges []GaugeMetric
	gauges, err = createGauges()
	if err != nil {
		errCause := fmt.Sprintln("error creating gauges: ", err.Error())
		return nil, common.ErrorFromThisScope(errCause, generalScopeErr)
	}
	collector = &SkycoinCollector{
		Counters: counters,
		Gauges:   gauges,
	}
	return collector, err
}

// Describe writes all the descriptors to the prometheus desc channel.
func (collector *SkycoinCollector) Describe(ch chan<- *prometheus.Desc) {
	for _, counter := range collector.Counters {
		ch <- counter.Desc
	}
	for _, gauge := range collector.Gauges {
		ch <- gauge.Desc
	}
}

//Collect update all the descriptors is values
func (collector *SkycoinCollector) Collect(ch chan<- prometheus.Metric) {
	for _, counter := range collector.Counters {
		val, err := counter.Client.GetMetric()
		if err != nil {
			log.Fatal("can not get the data", err)
		}
		typedVal := val.(float64) // FIXME(denisacostaq@gmail.com): make more assertion on this
		log.Println("getting typedVal:", typedVal)
		ch <- prometheus.MustNewConstMetric(counter.Desc, prometheus.CounterValue, typedVal)
	}
	for _, gauge := range collector.Gauges {
		val, err := gauge.Client.GetMetric()
		if err != nil {
			log.Fatal("can not get the data", err)
		}
		typedVal := val.(float64) // FIXME(denisacostaq@gmail.com): make more assertion on this
		log.Println("getting typedVal:", typedVal)
		ch <- prometheus.MustNewConstMetric(gauge.Desc, prometheus.GaugeValue, typedVal)
	}
}

func createMetricCommonStages(link config.Link) (metricClient *client.MetricClient, description string, err error) {
	const generalScopeErr = "error initializing metric creation scope"
	if metricClient, err = client.NewMetricClient(link); err != nil {
		errCause := fmt.Sprintln("error creating metric client: ", err.Error())
		return metricClient, description, common.ErrorFromThisScope(errCause, generalScopeErr)
	}
	if description, err = link.MetricDescription(); err != nil {
		errCause := fmt.Sprintln("can not build the description: ", err.Error())
		return metricClient, description, common.ErrorFromThisScope(errCause, generalScopeErr)
	}
	return metricClient, description, err
}

func createCounter(link config.Link) (metric ExportableCounterMetric, err error) {
	const generalScopeErr = "error creating a gauge"
	var metricClient *client.MetricClient
	var description string
	if metricClient, description, err = createMetricCommonStages(link); err != nil {
		errCause := fmt.Sprintln("can not get parameters for gauge: ", err.Error())
		return metric, common.ErrorFromThisScope(errCause, generalScopeErr)
	}
	metric = ExportableCounterMetric{
		Client: metricClient,
		Counter: prometheus.NewCounter(
			prometheus.CounterOpts{
				Name: link.MetricName(),
				Help: description,
			},
		),
	}
	return metric, err
}

func createGauge(link config.Link) (metric ExportableGaugeMetric, err error) {
	const generalScopeErr = "error creating a gauge"
	var metricClient *client.MetricClient
	var description string
	if metricClient, description, err = createMetricCommonStages(link); err != nil {
		errCause := fmt.Sprintln("can not build the parameters for counter: ", err.Error())
		return metric, common.ErrorFromThisScope(errCause, generalScopeErr)
	}
	metric = ExportableGaugeMetric{
		Client: metricClient,
		Gauge: prometheus.NewGauge(
			prometheus.GaugeOpts{
				Name: link.MetricName(),
				Help: description,
			},
		),
	}
	return metric, err
}

func getPrometheusMetrics(metrics []Metric) []prometheus.Collector {
	var collectors = make([]prometheus.Collector, len(metrics))
	for metricIdx, metric := range metrics {
		collectors[metricIdx] = metric.prometheusMetric()
	}
	return collectors
}

func createMetric(t string, link config.Link) (metric Metric, err error) {
	// TODO(denisacostaq@gmail.com): refactor the code bellow. var creator func(config.Link) (Metric, error)
	var generalScopeErr = "Error creating metric"
	switch t {
	case "Counter":
		metric, err = createCounter(link)
		if err != nil {
			errCause := fmt.Sprintln("can not crate a counter: ", err.Error())
			return metric, common.ErrorFromThisScope(errCause, generalScopeErr)
		}
	case "Gauge":
		metric, err = createGauge(link)
		if err != nil {
			errCause := fmt.Sprintln("can not crate a counter: ", err.Error())
			return metric, common.ErrorFromThisScope(errCause, generalScopeErr)
		}
	default:
		errCause := fmt.Sprintln("No switch handler for: ", t)
		return metric, common.ErrorFromThisScope(errCause, generalScopeErr)
	}
	return metric, err
}

func updateMetrics(metrics []Metric) {
	log.Println("Updating metrics")
	for _, metric := range metrics {
		metric.update()
	}
}

func createMetrics() (metrics []Metric, err error) {
	var generalScopeErr = "Error creating metrics"
	conf := config.Config()
	metrics = make([]Metric, len(conf.MetricsForHost))
	for linkIdx, link := range conf.MetricsForHost {
		var metricType string
		if metricType, err = link.FindMetricType(); err != nil {
			errCause := fmt.Sprintln("can not find the metric type: ", err.Error())
			return metrics, common.ErrorFromThisScope(errCause, generalScopeErr)
		}
		var metric Metric
		if metric, err = createMetric(metricType, link); err != nil {
			errCause := fmt.Sprintln("can not create the metric: ", err.Error())
			return metrics, common.ErrorFromThisScope(errCause, generalScopeErr)
		}
		metrics[linkIdx] = metric
	}
	prometheus.MustRegister(getPrometheusMetrics(metrics)...)
	return metrics, err
}

func onDemandMetricsUpdateHandler(orgHandler http.Handler) (newHandler http.Handler) {
	metrics, err := createMetrics()
	if err != nil {
		log.Fatalln(err)
	}
	hf := func(w http.ResponseWriter, r *http.Request) {
		log.Println("Before")
		updateMetrics(metrics)
		// http.Error(w, "Unable to update metrics", http.StatusInternalServerError)
		orgHandler.ServeHTTP(w, r) // call original
		log.Println("After")
	}
	return http.HandlerFunc(hf)
}

// ExportMetrics will read the config file from the CLI parammeter `-config` if any
// or use a default one.
func ExportMetrics(configFile string, listenPort uint16) (srv *http.Server) {
	if err := config.NewConfigFromFilePath(configFile); err != nil {
		log.Fatalln("can not open the config file", err.Error())
	}
	port := fmt.Sprintf(":%d", listenPort)
	srv = &http.Server{Addr: port}
	http.Handle("/metrics", onDemandMetricsUpdateHandler(promhttp.Handler()))
	go func() {
		log.Panicln(srv.ListenAndServe())
	}()
	return srv
}

// TODO(denisacostaq@gmail.com): you can use a NewProcessCollector, NewGoProcessCollector, make a blockchain collector sense?
