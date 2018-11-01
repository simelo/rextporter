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

// CounterMetric has the necessary http client to get and updated value for the counter metric
type CounterMetric struct {
	Client *client.MetricClient
	Desc   *prometheus.Desc
}

func createCounter(link config.Link) (metric CounterMetric, err error) {
	generalScopeErr := "can not create metric " + link.MetricRef
	var metricClient *client.MetricClient
	var description string
	if metricClient, description, err = createMetricCommonStages(link); err != nil {
		errCause := fmt.Sprintln("can not get parameters for counter creation stage: ", err.Error())
		return metric, common.ErrorFromThisScope(errCause, generalScopeErr)
	}
	metric = CounterMetric{
		Client: metricClient,
		Desc:   prometheus.NewDesc(link.MetricRef, description, nil, nil),
	}
	return metric, err
}

func createCounters() ([]CounterMetric, error) {
	generalScopeErr := "can not create counters"
	conf := config.Config()
	links, err := config.FilterLinksByMetricType(conf.MetricsForHost, config.KeyTypeCounter)
	if err != nil {
		errCause := fmt.Sprintln("can not filter links by type type: ", config.KeyTypeCounter, " ", err.Error())
		return []CounterMetric{}, common.ErrorFromThisScope(errCause, generalScopeErr)
	}
	counters := make([]CounterMetric, len(links))
	for idx, link := range links {
		counter, err := createCounter(link)
		if err != nil {
			errCause := fmt.Sprintln("error creating counter: ", err.Error())
			return []CounterMetric{}, common.ErrorFromThisScope(errCause, generalScopeErr)
		}
		counters[idx] = counter
	}
	return counters, nil
}

// GaugeMetric has the necessary http client to get and updated value for the counter metric
type GaugeMetric struct {
	Client *client.MetricClient
	Desc   *prometheus.Desc
}

func createGauge(link config.Link) (metric GaugeMetric, err error) {
	generalScopeErr := "can not create metric " + link.MetricRef
	var metricClient *client.MetricClient
	var description string
	if metricClient, description, err = createMetricCommonStages(link); err != nil {
		errCause := fmt.Sprintln("can not get parameters for gauge creation stage: ", err.Error())
		return metric, common.ErrorFromThisScope(errCause, generalScopeErr)
	}
	metric = GaugeMetric{
		Client: metricClient,
		Desc:   prometheus.NewDesc(link.MetricRef, description, nil, nil),
	}
	return metric, err
}

func createGauges() ([]GaugeMetric, error) {
	generalScopeErr := "can not create gauges"
	conf := config.Config()
	links, err := config.FilterLinksByMetricType(conf.MetricsForHost, config.KeyTypeGauge)
	if err != nil {
		errCause := fmt.Sprintln("can not filter links by type type: ", config.KeyTypeGauge, " ", err.Error())
		return []GaugeMetric{}, common.ErrorFromThisScope(errCause, generalScopeErr)
	}
	gauges := make([]GaugeMetric, len(links))
	for idx, link := range links {
		gauge, err := createGauge(link)
		if err != nil {
			errCause := fmt.Sprintln("error creating gauge: ", err.Error())
			return []GaugeMetric{}, common.ErrorFromThisScope(errCause, generalScopeErr)
		}
		gauges[idx] = gauge
	}
	return gauges, nil
}

// ExportMetrics will read the config file from the CLI parammeter `-config` if any
// or use a default one.
func ExportMetrics(configFile string, listenPort uint16) (srv *http.Server) {
	if err := config.NewConfigFromFilePath(configFile); err != nil {
		log.Fatalln("can not open the config file", err.Error())
	}
	if collector, err := newSkycoinCollector(); err != nil {
		log.Panicln("Can not create metrics:", err)
	} else {
		prometheus.MustRegister(collector)
	}
	port := fmt.Sprintf(":%d", listenPort)
	srv = &http.Server{Addr: port}
	http.Handle("/metrics", promhttp.Handler())
	log.Panicln(srv.ListenAndServe())
	return srv
}

// TODO(denisacostaq@gmail.com): you can use a NewProcessCollector, NewGoProcessCollector, make a blockchain collector sense?
