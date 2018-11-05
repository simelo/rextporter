package exporter

import (
	"fmt"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/simelo/rextporter/src/common"
	log "github.com/sirupsen/logrus"
)

// SkycoinCollector has the metrics to be exposed
type SkycoinCollector struct {
	Counters []CounterMetric
	Gauges   []GaugeMetric
}

func newSkycoinCollector() (collector *SkycoinCollector, err error) {
	const generalScopeErr = "error creating collector"
	collector = &SkycoinCollector{}
	if collector.Counters, err = createCounters(); err != nil {
		errCause := fmt.Sprintln("error creating counters: ", err.Error())
		return nil, common.ErrorFromThisScope(errCause, generalScopeErr)
	}
	if collector.Gauges, err = createGauges(); err != nil {
		errCause := fmt.Sprintln("error creating gauges: ", err.Error())
		return nil, common.ErrorFromThisScope(errCause, generalScopeErr)
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
	onCollectFail := func(counter CounterMetric, fch chan<- prometheus.Metric) {
		fch <- prometheus.MustNewConstMetric(counter.StatusDesc, prometheus.GaugeValue, 1)
		fch <- prometheus.MustNewConstMetric(counter.MetricDesc, prometheus.CounterValue, counter.lastSuccessValue)
	}
	onCollectSuccess := func(counter CounterMetric, fch chan<- prometheus.Metric, val float64) {
		fch <- prometheus.MustNewConstMetric(counter.StatusDesc, prometheus.GaugeValue, 0)
		fch <- prometheus.MustNewConstMetric(counter.MetricDesc, prometheus.CounterValue, val)
	}
	for _, counter := range collector.Counters {
		if val, err := counter.Client.GetMetric(); err != nil {
			log.WithError(err).Errorln("can not get the data")
			onCollectFail(counter, ch)
		} else {
			typedVal, ok := val.(float64) // FIXME(denisacostaq@gmail.com): make more assertion on this, negative panic
			if ok {
				onCollectSuccess(counter, ch, typedVal)
			} else {
				log.WithField("val", val).Errorln("unable to get value as float64")
				onCollectFail(counter, ch)
			}
		}
	}
}

func (collector *SkycoinCollector) collectGauges(ch chan<- prometheus.Metric) {
	onCollectFail := func(gauge GaugeMetric, fch chan<- prometheus.Metric) {
		fch <- prometheus.MustNewConstMetric(gauge.StatusDesc, prometheus.GaugeValue, 1)
		fch <- prometheus.MustNewConstMetric(gauge.MetricDesc, prometheus.GaugeValue, gauge.lastSuccessValue)
	}
	onCollectSuccess := func(gauge GaugeMetric, fch chan<- prometheus.Metric, val float64) {
		fch <- prometheus.MustNewConstMetric(gauge.StatusDesc, prometheus.GaugeValue, 0)
		fch <- prometheus.MustNewConstMetric(gauge.MetricDesc, prometheus.GaugeValue, val)
	}
	for _, gauge := range collector.Gauges {
		if val, err := gauge.Client.GetMetric(); err != nil {
			log.WithError(err).Errorln("can not get the data")
			onCollectFail(gauge, ch)
		} else {
			typedVal, ok := val.(float64)
			if ok {
				onCollectSuccess(gauge, ch, typedVal)
			} else {
				log.WithField("val", val).Errorln("unable to get value as float64")
				onCollectFail(gauge, ch)
			}
		}
	}
}

//Collect update all the descriptors is values
func (collector *SkycoinCollector) Collect(ch chan<- prometheus.Metric) {
	collector.collectCounters(ch)
	collector.collectGauges(ch)
}
