package exporter

import (
	"fmt"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/simelo/rextporter/src/client"
	"github.com/simelo/rextporter/src/config"
	"github.com/simelo/rextporter/src/util"
)

// CounterMetric has the necessary http client to get and updated value for the counter metric
type CounterMetric struct {
	Client           *client.MetricClient
	lastSuccessValue float64
	MetricDesc       *prometheus.Desc
	StatusDesc       *prometheus.Desc
}

func createCounter(metricConf config.Metric, conf config.RootConfig) (metric CounterMetric, err error) {
	generalScopeErr := "can not create metric " + metricConf.Name
	var metricClient *client.MetricClient
	if metricClient, err = client.NewMetricClient(metricConf, conf); err != nil {
		errCause := fmt.Sprintln("error creating metric client: ", err.Error())
		return metric, util.ErrorFromThisScope(errCause, generalScopeErr)
	}
	metric = CounterMetric{
		Client:     metricClient,
		MetricDesc: prometheus.NewDesc(conf.MetricName(metricConf.Name), metricConf.Options.Description, nil, nil),
		StatusDesc: prometheus.NewDesc(conf.MetricName(metricConf.Name)+"_up", "Says if the same name metric("+conf.MetricName(metricConf.Name)+") was success updated, 1 for ok, 0 for failed.", nil, nil),
	}
	return metric, err
}

func createCounters() ([]CounterMetric, error) {
	generalScopeErr := "can not create counters"
	conf := config.Config() // TODO(denisacostaq@gmail.com): recive conf as parameter
	metrics := conf.FilterMetricsByType(config.KeyTypeCounter)
	counters := make([]CounterMetric, len(metrics))
	for idx, metric := range metrics {
		if counter, err := createCounter(metric, conf); err == nil {
			counters[idx] = counter
		} else {
			errCause := "error creating counter: " + err.Error()
			return []CounterMetric{}, util.ErrorFromThisScope(errCause, generalScopeErr)
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

func createGauge(metricConf config.Metric, conf config.RootConfig) (metric GaugeMetric, err error) {
	generalScopeErr := "can not create metric " + metricConf.Name
	var metricClient *client.MetricClient
	if metricClient, err = client.NewMetricClient(metricConf, conf); err != nil {
		errCause := fmt.Sprintln("error creating metric client: ", err.Error())
		return metric, util.ErrorFromThisScope(errCause, generalScopeErr)
	}
	metric = GaugeMetric{
		Client:     metricClient,
		MetricDesc: prometheus.NewDesc(conf.MetricName(metricConf.Name), metricConf.Options.Description, nil, nil),
		StatusDesc: prometheus.NewDesc(conf.MetricName(metricConf.Name)+"_up", "Says if the same name metric("+conf.MetricName(metricConf.Name)+") was success updated, 1 for ok, 0 for failed.", nil, nil),
	}
	return metric, err
}

func createGauges() ([]GaugeMetric, error) {
	generalScopeErr := "can not create gauges"
	conf := config.Config() // TODO(denisacostaq@gmail.com): recive conf as parameter
	metrics := conf.FilterMetricsByType(config.KeyTypeGauge)
	gauges := make([]GaugeMetric, len(metrics))
	for idx, metric := range metrics {
		gauge, err := createGauge(metric, conf)
		if err != nil {
			errCause := fmt.Sprintln("error creating gauge: ", err.Error())
			return []GaugeMetric{}, util.ErrorFromThisScope(errCause, generalScopeErr)
		}
		gauges[idx] = gauge
	}
	return gauges, nil
}
