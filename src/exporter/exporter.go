package exporter

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/simelo/rextporter/src/client"
	"github.com/simelo/rextporter/src/common"
	"github.com/simelo/rextporter/src/config"
)

// Metric interface allow you to put all the metrics in a container an update all of them.
type Metric interface {
	update()
	prometheusMetric() prometheus.Collector
}

// ExportableCounterMetric has the necessary http client to get and updated value for the counter metric
type ExportableCounterMetric struct {
	Client  *client.MetricClient
	Counter prometheus.Counter
}

func (metric ExportableCounterMetric) update() {
	log.Println("Update")
	val, err := metric.Client.GetMetric()
	if err != nil {
		log.Fatal("can not get the data", err)
	}
	metric.Counter.Add(val.(float64))
}

func (metric ExportableCounterMetric) prometheusMetric() prometheus.Collector {
	return metric.Counter
}

// ExportableGaugeMetric has the necessary http client to get and updated value for the counter metric
type ExportableGaugeMetric struct {
	Client *client.MetricClient
	Gauge  prometheus.Gauge
}

func (metric ExportableGaugeMetric) update() {
	val, err := metric.Client.GetMetric()
	if err != nil {
		log.Fatal("can not get the data", err)
	}
	metric.Gauge.Add(val.(float64))
}

func (metric ExportableGaugeMetric) prometheusMetric() prometheus.Collector {
	return metric.Gauge
}

func createMetricCommonStages(link config.Link) (metricClient *client.MetricClient, description string, err error) {
	const generalScopeErr = "error initializing metric creation scope"
	if metricClient, err = client.NewMetricClient(link); err != nil {
		errCause := fmt.Sprintln("error creating metric client", err.Error())
		return metricClient, description, common.ErrorFromThisScope(errCause, generalScopeErr)
	}
	if description, err = link.MetricDescription(); err != nil {
		errCause := fmt.Sprintln("can not build the description", err.Error())
		return metricClient, description, common.ErrorFromThisScope(errCause, generalScopeErr)
	}
	return metricClient, description, err
}

func createCounter(link config.Link) (metric ExportableCounterMetric, err error) {
	const generalScopeErr = "error creating a gauge"
	var metricClient *client.MetricClient
	var description string
	if metricClient, description, err = createMetricCommonStages(link); err != nil {
		errCause := fmt.Sprintln("can not get parameters for gauge", err.Error())
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
		errCause := fmt.Sprintln("can not build the parameters for counter", err.Error())
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
	var generalScopeErr string
	switch t {
	case "Counter":
		metric, err = createCounter(link)
		if err != nil {
			errCause := fmt.Sprintln("can not crate a counter", err.Error())
			return metric, common.ErrorFromThisScope(errCause, generalScopeErr)
		}
	case "Gauge":
		metric, err = createGauge(link)
		if err != nil {
			errCause := fmt.Sprintln("can not crate a counter", err.Error())
			return metric, common.ErrorFromThisScope(errCause, generalScopeErr)
		}
	default:
		errCause := fmt.Sprintln("No switch handler for", t)
		return metric, common.ErrorFromThisScope(errCause, generalScopeErr)
	}
	return metric, err
}

// ExportMetrics will read the config value and created all the specified metrics from the config file.
func ExportMetrics() {
	gopath := os.Getenv("GOPATH")
	filePath := gopath + "/src/github.com/denisacostaq/rextporter/examples/simple.toml"
	if err := config.NewConfigFromFilePath(filePath); err != nil {
		log.Fatalln("can not open file", err)
	}
	conf := config.Config()
	metrics := make([]Metric, len(conf.MetricsForHost))
	for linkIdx, link := range conf.MetricsForHost {
		var metricType string
		var err error
		if metricType, err = link.FindMetricType(); err != nil {
			log.Panicln(err)
		}
		var metric Metric
		if metric, err = createMetric(metricType, link); err != nil {
			log.Panicln(err)
		}
		metrics[linkIdx] = metric
	}
	prometheus.MustRegister(getPrometheusMetrics(metrics)...)
	http.Handle("/metric", promhttp.Handler())
	// NOTE(denisacostaq@gmail.com): This is a fate test, it should be removed
	// check https://github.com/skycoin/skycoin/blob/develop/src/api/http.go
	// make a wrapper in the http handler. TODO
	go func() {
		t := time.NewTimer(time.Second * 5)
		<-t.C
		for _, metric := range metrics {
			metric.update()
		}
	}()
	log.Panicln(http.ListenAndServe(":8000", nil))
}

// TODO(denisacostaq@gmail.com): you can use a NewProcessCollector, NewGoProcessCollector, make a blockchain collector sense?
