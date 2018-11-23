package exporter

import (
	"fmt"

	"github.com/denisacostaq/rextporter/src/scrapper"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/simelo/rextporter/src/client"
	"github.com/simelo/rextporter/src/config"
	"github.com/simelo/rextporter/src/util"
)

// MetricMiddleware has the necessary http client to get exposed metric from a service
type MetricMiddleware struct {
	client *client.ProxyMetricClient
}

func createMetricsMiddleware(conf config.RootConfig) (metricsMiddleware []MetricMiddleware, err error) {
	generalScopeErr := "can not create metrics Middleware"
	services := conf.FilterServicesByType(config.ServiceTypeProxy)
	for _, service := range services {
		var cl *client.ProxyMetricClient
		if cl, err = client.NewProxyMetricClient(service); err != nil {
			errCause := fmt.Sprintln("error creating metric client: ", err.Error())
			return metricsMiddleware, util.ErrorFromThisScope(errCause, generalScopeErr)
		}
		metricsMiddleware = append(metricsMiddleware, MetricMiddleware{client: cl})
	}
	return metricsMiddleware, err
}

// CounterMetric has the necessary http client to get and updated value for the counter metric
type CounterMetric struct {
	Client           client.Client
	lastSuccessValue interface{}
	MetricDesc       *prometheus.Desc
	StatusDesc       *prometheus.Desc
}

func createCounter(metricConf config.Metric, srvConf config.Service) (metric CounterMetric, err error) {
	generalScopeErr := "can not create metric " + metricConf.Name
	var metricClient client.Client
	if metricClient, err = client.NewClient(metricConf, srvConf); err != nil {
		errCause := fmt.Sprintln("error creating metric client: ", err.Error())
		return metric, util.ErrorFromThisScope(errCause, generalScopeErr)
	}
	labels := metricConf.LabelNames()
	metric = CounterMetric{
		// FIXME(denisacostaq@gmail.com): if you use a duplicated name can panic?
		Client:     metricClient,
		MetricDesc: prometheus.NewDesc(srvConf.MetricName(metricConf.Name), metricConf.Options.Description, labels, nil),
		StatusDesc: prometheus.NewDesc(srvConf.MetricName(metricConf.Name)+"_up", "Says if the same name metric("+srvConf.MetricName(metricConf.Name)+") was success updated, 1 for ok, 0 for failed.", labels, nil),
	}
	return metric, err
}

func createCounters(conf config.RootConfig) ([]CounterMetric, error) {
	generalScopeErr := "can not create counters"
	services := conf.FilterServicesByType(config.ServiceTypeAPIRest)
	var counterMetricsAmount = 0
	for _, service := range services {
		counterMetricsAmount += service.CountMetricsByType(config.KeyTypeCounter)
	}
	counters := make([]CounterMetric, counterMetricsAmount)
	var idxMetric = 0
	for _, service := range services {
		metricsForService := service.FilterMetricsByType(config.KeyTypeCounter)
		for _, metric := range metricsForService {
			if counter, err := createCounter(metric, service); err == nil {
				counters[idxMetric] = counter
				idxMetric++
			} else {
				errCause := "error creating counter: " + err.Error()
				return []CounterMetric{}, util.ErrorFromThisScope(errCause, generalScopeErr)
			}
		}
	}
	return counters, nil
}

// GaugeMetric has the necessary http client to get and updated value for the counter metric
type GaugeMetric struct {
	scrapper         scrapper.Scrapper
	lastSuccessValue interface{}
	MetricDesc       *prometheus.Desc
	StatusDesc       *prometheus.Desc
}

func createGauge(metricConf config.Metric, srvConf config.Service) (metric GaugeMetric, err error) {
	generalScopeErr := "can not create metric " + metricConf.Name
	var metricClient client.Client
	if metricClient, err = client.NewClient(metricConf, srvConf); err != nil {
		errCause := fmt.Sprintln("error creating metric client: ", err.Error())
		return metric, util.ErrorFromThisScope(errCause, generalScopeErr)
	}
	numScrapper := scrapper.NewNumeric(metricClient, scrapper.JsonParser{}, metricConf.Path)
	labels := metricConf.LabelNames()
	metric = GaugeMetric{
		scrapper:   numScrapper,
		MetricDesc: prometheus.NewDesc(srvConf.MetricName(metricConf.Name), metricConf.Options.Description, labels, nil),
		StatusDesc: prometheus.NewDesc(srvConf.MetricName(metricConf.Name)+"_up", "Says if the same name metric("+srvConf.MetricName(metricConf.Name)+") was success updated, 1 for ok, 0 for failed.", labels, nil),
	}
	return metric, err
}

func createGauges(conf config.RootConfig) ([]GaugeMetric, error) {
	generalScopeErr := "can not create gauges"
	services := conf.FilterServicesByType(config.ServiceTypeAPIRest)
	var gaugeMetricsAmount = 0
	for _, service := range services {
		gaugeMetricsAmount += service.CountMetricsByType(config.KeyTypeGauge)
	}
	gauges := make([]GaugeMetric, gaugeMetricsAmount)
	var idxMetric = 0
	for _, service := range services {
		metricsForService := service.FilterMetricsByType(config.KeyTypeGauge)
		for _, metric := range metricsForService {
			if gauge, err := createGauge(metric, service); err == nil {
				gauges[idxMetric] = gauge
				idxMetric++
			} else {
				errCause := "error creating gauge: " + err.Error()
				return gauges, util.ErrorFromThisScope(errCause, generalScopeErr)
			}
		}
	}
	return gauges, nil
}

// HistogramMetric has the necessary http client to get and updated value for the histogram metric
type HistogramMetric struct {
	Client           client.Client
	lastSuccessValue client.HistogramValue
	MetricDesc       *prometheus.Desc
	StatusDesc       *prometheus.Desc
}

func createHistogram(metricConf config.Metric, service config.Service) (metric HistogramMetric, err error) {
	generalScopeErr := "can not create metric " + metricConf.Name
	var metricClient client.Client
	if metricClient, err = client.NewClient(metricConf, service); err != nil {
		errCause := fmt.Sprintln("error creating metric client: ", err.Error())
		return metric, util.ErrorFromThisScope(errCause, generalScopeErr)
	}
	labels := metricConf.LabelNames()
	metric = HistogramMetric{
		Client:     metricClient,
		MetricDesc: prometheus.NewDesc(service.MetricName(metricConf.Name), metricConf.Options.Description, labels, nil),
		StatusDesc: prometheus.NewDesc(service.MetricName(metricConf.Name)+"_up", "Says if the same name metric("+service.MetricName(metricConf.Name)+") was success updated, 1 for ok, 0 for failed.", labels, nil),
	}
	return metric, err
}

func createHistograms(conf config.RootConfig) ([]HistogramMetric, error) {
	generalScopeErr := "can not create histograms"
	services := conf.FilterServicesByType(config.ServiceTypeAPIRest)
	var histogramMetricsAmount = 0
	for _, service := range services {
		histogramMetricsAmount += service.CountMetricsByType(config.KeyTypeHistogram)
	}
	histograms := make([]HistogramMetric, histogramMetricsAmount)
	var idxMetric = 0
	for _, service := range services {
		metricsForService := service.FilterMetricsByType(config.KeyTypeHistogram)
		for _, metric := range metricsForService {
			if histogram, err := createHistogram(metric, service); err == nil {
				histograms[idxMetric] = histogram
				idxMetric++
			} else {
				errCause := "error creating histogram: " + err.Error()
				return []HistogramMetric{}, util.ErrorFromThisScope(errCause, generalScopeErr)
			}
		}
	}
	return histograms, nil
}
