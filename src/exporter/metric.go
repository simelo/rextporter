package exporter

import (
	"fmt"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/simelo/rextporter/src/client"
	"github.com/simelo/rextporter/src/config"
	"github.com/simelo/rextporter/src/util"
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

func createCounter(metricConf config.Metric, srvConf config.Service) (metric CounterMetric, err error) {
	generalScopeErr := "can not create metric " + metricConf.Name
	var metricClient *client.MetricClient
	if metricClient, err = client.NewMetricClient(metricConf, srvConf); err != nil {
		errCause := fmt.Sprintln("error creating metric client: ", err.Error())
		return metric, util.ErrorFromThisScope(errCause, generalScopeErr)
	}
	metric = CounterMetric{
		// FIXME(denisacostaq@gmail.com): if you use a duplicated name can panic?
		Client:     metricClient,
		MetricDesc: prometheus.NewDesc(srvConf.MetricName(metricConf.Name), metricConf.Options.Description, nil, nil),
		StatusDesc: prometheus.NewDesc(srvConf.MetricName(metricConf.Name)+"_up", "Says if the same name metric("+srvConf.MetricName(metricConf.Name)+") was success updated, 1 for ok, 0 for failed.", nil, nil),
	}
	return metric, err
}

func createCounters() ([]CounterMetric, error) {
	generalScopeErr := "can not create counters"
	conf := config.Config() // TODO(denisacostaq@gmail.com): recive conf as parameter
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
	Client           *client.MetricClient
	lastSuccessValue float64
	MetricDesc       *prometheus.Desc
	StatusDesc       *prometheus.Desc
}

func createGauge(metricConf config.Metric, srvConf config.Service) (metric GaugeMetric, err error) {
	generalScopeErr := "can not create metric " + metricConf.Name
	var metricClient *client.MetricClient
	if metricClient, err = client.NewMetricClient(metricConf, srvConf); err != nil {
		errCause := fmt.Sprintln("error creating metric client: ", err.Error())
		return metric, util.ErrorFromThisScope(errCause, generalScopeErr)
	}
	metric = GaugeMetric{
		Client:     metricClient,
		MetricDesc: prometheus.NewDesc(srvConf.MetricName(metricConf.Name), metricConf.Options.Description, nil, nil),
		StatusDesc: prometheus.NewDesc(srvConf.MetricName(metricConf.Name)+"_up", "Says if the same name metric("+srvConf.MetricName(metricConf.Name)+") was success updated, 1 for ok, 0 for failed.", nil, nil),
	}
	return metric, err
}

func createGauges() ([]GaugeMetric, error) {
	generalScopeErr := "can not create gauges"
	conf := config.Config() // TODO(denisacostaq@gmail.com): recive conf as parameter
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
			}
		}
	}
	return gauges, nil
}
