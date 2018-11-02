package exporter

import (
	"fmt"
	"log"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/simelo/rextporter/src/common"
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
		ch <- counter.MetricDesc
	}
	for _, gauge := range collector.Gauges {
		ch <- gauge.MetricDesc
	}
}

func (collector *SkycoinCollector) collectCounters(ch chan<- prometheus.Metric) {
	for _, counter := range collector.Counters {
		val, err := counter.Client.GetMetric()
		if err != nil {
			log.Println("can not get the data:", err)
			ch <- prometheus.MustNewConstMetric(counter.StatusDesc, prometheus.GaugeValue, 0)
		} else {
			typedVal := val.(float64) // FIXME(denisacostaq@gmail.com): make more assertion on this
			ch <- prometheus.MustNewConstMetric(counter.MetricDesc, prometheus.CounterValue, typedVal)
			ch <- prometheus.MustNewConstMetric(counter.StatusDesc, prometheus.GaugeValue, 1)
		}
	}
}

func (collector *SkycoinCollector) collectGauges(ch chan<- prometheus.Metric) {
	for _, gauge := range collector.Gauges {
		val, err := gauge.Client.GetMetric()
		if err != nil {
			log.Println("can not get the data", err)
			ch <- prometheus.MustNewConstMetric(gauge.StatusDesc, prometheus.GaugeValue, 0)
		} else {
			typedVal := val.(float64) // FIXME(denisacostaq@gmail.com): make more assertion on this
			ch <- prometheus.MustNewConstMetric(gauge.MetricDesc, prometheus.GaugeValue, typedVal)
			ch <- prometheus.MustNewConstMetric(gauge.StatusDesc, prometheus.GaugeValue, 1)
		}
	}
}

//Collect update all the descriptors is values
func (collector *SkycoinCollector) Collect(ch chan<- prometheus.Metric) {
	collector.collectCounters(ch)
	collector.collectGauges(ch)
}
