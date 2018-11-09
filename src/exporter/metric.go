package exporter

import (
	"fmt"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/simelo/rextporter/src/client"
	"github.com/simelo/rextporter/src/common"
	"github.com/simelo/rextporter/src/config"
)

// MetricMidleware has the necessary http client to get exposed metric from a service
type MetricMidleware struct {
	client *client.ProxyMetricClient
}

func createMetricsMidleware() (metricsMidleware []MetricMidleware, err error) {
	generalScopeErr := "can not create metrics midleware"
	conf := config.Config()
	services := conf.FilterServicesByType(config.ServiceTypeProxy)
	for _, service := range services {
		var cl *client.ProxyMetricClient
		if cl, err = client.NewProxyMetricClient(service); err != nil {
			errCause := fmt.Sprintln("error creating metric client: ", err.Error())
			return metricsMidleware, common.ErrorFromThisScope(errCause, generalScopeErr)
		}
		metricsMidleware = append(metricsMidleware, MetricMidleware{client: cl})
	}
	return metricsMidleware, err
}

// CounterMetric has the necessary http client to get and updated value for the counter metric
type CounterMetric struct {
	Client           *client.MetricClient
	lastSuccessValue float64
	MetricDesc       *prometheus.Desc
	StatusDesc       *prometheus.Desc
}

func createCounter(metricConf config.Metric, service config.Service) (metric CounterMetric, err error) {
	generalScopeErr := "can not create metric " + metricConf.Name
	var metricClient *client.MetricClient
	if metricClient, err = client.NewMetricClient(metricConf, service); err != nil {
		errCause := fmt.Sprintln("error creating metric client: ", err.Error())
		return metric, common.ErrorFromThisScope(errCause, generalScopeErr)
	}
	metric = CounterMetric{
		// FIXME(denisacostaq@gmail.com): if you use a duplicated name can panic?
		Client:     metricClient,
		MetricDesc: prometheus.NewDesc(service.MetricName(metricConf.Name), metricConf.Options.Description, nil, nil),
		StatusDesc: prometheus.NewDesc(service.MetricName(metricConf.Name)+"_up", "Says if the same name metric("+service.MetricName(metricConf.Name)+") was success updated, 1 for ok, 0 for failed.", nil, nil),
	}
	return metric, err
}

func createCounters() ([]CounterMetric, error) {
	generalScopeErr := "can not create counters"
	conf := config.Config() // TODO(denisacostaq@gmail.com): recive conf as parameter
	metrics := conf.FilterMetricsByType(config.KeyTypeCounter)
	services := conf.FilterServicesByType(config.ServiceTypeApiRest)
	counters := make([]CounterMetric, len(metrics)*len(services))
	for idxService, service := range services {
		for idxMetric, metric := range metrics {
			if counter, err := createCounter(metric, service); err == nil {
				counters[idxService*len(conf.Metrics)+idxMetric] = counter
			} else {
				errCause := "error creating counter: " + err.Error()
				return []CounterMetric{}, common.ErrorFromThisScope(errCause, generalScopeErr)
			}
		}
	}
	return counters, nil
}

// GaugeMetric has the necessary http client to get and updated value for the counter metric
type GaugeMetric struct {
	Client           *client.MetricClient
	lastSuccessValue float64
	MetricDesc       *prometheus.Desc
	StatusDesc       *prometheus.Desc
}

func createGauge(metricConf config.Metric, service config.Service) (metric GaugeMetric, err error) {
	generalScopeErr := "can not create metric " + metricConf.Name
	var metricClient *client.MetricClient
	if metricClient, err = client.NewMetricClient(metricConf, service); err != nil {
		errCause := fmt.Sprintln("error creating metric client: ", err.Error())
		return metric, common.ErrorFromThisScope(errCause, generalScopeErr)
	}
	metric = GaugeMetric{
		Client:     metricClient,
		MetricDesc: prometheus.NewDesc(service.MetricName(metricConf.Name), metricConf.Options.Description, nil, nil),
		StatusDesc: prometheus.NewDesc(service.MetricName(metricConf.Name)+"_up", "Says if the same name metric("+service.MetricName(metricConf.Name)+") was success updated, 1 for ok, 0 for failed.", nil, nil),
	}
	return metric, err
}

func createGauges() ([]GaugeMetric, error) {
	generalScopeErr := "can not create gauges"
	conf := config.Config() // TODO(denisacostaq@gmail.com): recive conf as parameter
	metrics := conf.FilterMetricsByType(config.KeyTypeGauge)
	services := conf.FilterServicesByType(config.ServiceTypeApiRest)
	gauges := make([]GaugeMetric, len(metrics)*len(services))
	for idxService, service := range services {
		for idxMetric, metric := range metrics {
			gauge, err := createGauge(metric, service)
			if err != nil {
				errCause := fmt.Sprintln("error creating gauge: ", err.Error())
				return []GaugeMetric{}, common.ErrorFromThisScope(errCause, generalScopeErr)
			}
			gauges[idxService*len(conf.Metrics)+idxMetric] = gauge
		}
	}
	return gauges, nil
}
