package exporter

import (
	"fmt"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/simelo/rextporter/src/client"
	"github.com/simelo/rextporter/src/common"
	"github.com/simelo/rextporter/src/config"
)

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
