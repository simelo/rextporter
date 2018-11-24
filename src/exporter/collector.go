package exporter

import (
	"errors"
	"fmt"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/simelo/rextporter/src/config"
	"github.com/simelo/rextporter/src/scrapper"
	"github.com/simelo/rextporter/src/util"
	log "github.com/sirupsen/logrus"
)

// SkycoinCollector has the metrics to be exposed
type SkycoinCollector struct {
	Counters   []CounterMetric
	Gauges     []GaugeMetric
	Histograms []HistogramMetric
}

func newSkycoinCollector(conf config.RootConfig) (collector *SkycoinCollector, err error) {
	const generalScopeErr = "error creating collector"
	collector = &SkycoinCollector{}
	if collector.Counters, err = createCounters(conf); err != nil {
		errCause := fmt.Sprintln("error creating counters: ", err.Error())
		return nil, util.ErrorFromThisScope(errCause, generalScopeErr)
	}
	if collector.Gauges, err = createGauges(conf); err != nil {
		errCause := fmt.Sprintln("error creating gauges: ", err.Error())
		return nil, util.ErrorFromThisScope(errCause, generalScopeErr)
	}
	if collector.Histograms, err = createHistograms(conf); err != nil {
		errCause := fmt.Sprintln("error creating histograms: ", err.Error())
		return nil, util.ErrorFromThisScope(errCause, generalScopeErr)
	}
	return collector, err
}

// Describe writes all the descriptors to the prometheus desc channel.
func (collector *SkycoinCollector) Describe(ch chan<- *prometheus.Desc) {
	for _, counter := range collector.Counters {
		ch <- counter.MetricDesc
		ch <- counter.StatusDesc
	}
	for _, gauge := range collector.Gauges {
		ch <- gauge.MetricDesc
		ch <- gauge.StatusDesc
	}
	for _, histogram := range collector.Histograms {
		ch <- histogram.MetricDesc
		ch <- histogram.StatusDesc
	}
}

func (collector *SkycoinCollector) collectCounters(ch chan<- prometheus.Metric) {
	onCollectFail := func(counter CounterMetric, fch chan<- prometheus.Metric) {
		if metric, err := prometheus.NewConstMetric(counter.StatusDesc, prometheus.GaugeValue, 1); err == nil {
			fch <- metric
		} else {
			log.WithError(err).Errorln("collectCounters -> onCollectFail can not set up flag")
		}
		switch counter.lastSuccessValue.(type) {
		case float64:
			val, okVal := counter.lastSuccessValue.(float64)
			if !okVal {
				log.WithField("val", val).Errorln("can not get las success val")
			}
			if metric, err := prometheus.NewConstMetric(counter.MetricDesc, prometheus.CounterValue, val); err == nil {
				fch <- metric
			} else {
				log.WithError(err).Errorln("collectCounters -> onCollectFail can not set the last success value")
			}
		case scrapper.NumericVecVals:
			vals, okVals := counter.lastSuccessValue.(scrapper.NumericVecVals)
			if !okVals {
				log.WithField("vals", vals).Errorln("can not get las success vals")
			}
			for _, val := range vals {
				if metric, err := prometheus.NewConstMetric(counter.MetricDesc, prometheus.CounterValue, val.Val, val.Labels...); err == nil {
					fch <- metric
				} else {
					log.WithError(err).Errorln("collectCounters -> onCollectFail can not set the last success vec value")
				}
			}
		}
	}
	recoverNegativeCounter := func(counter CounterMetric, fch chan<- prometheus.Metric) {
		if r := recover(); r != nil {
			switch val := r.(type) {
			case string:
				log.WithError(errors.New(val)).Errorln("recovered with msg string")
			case error:
				log.WithError(val).Errorln("recovered with error")
			default:
				log.WithField("val", val).Errorln("recovered with un-know type")
			}
			onCollectFail(counter, fch)
		}
	}
	onCollectSuccess := func(counter *CounterMetric, fch chan<- prometheus.Metric, val float64) {
		defer recoverNegativeCounter(*counter, fch)
		if metric, err := prometheus.NewConstMetric(counter.StatusDesc, prometheus.GaugeValue, 0); err == nil {
			fch <- metric
		} else {
			log.WithError(err).Errorln("collectCounters -> onCollectSuccess can not set down flag")
			onCollectFail(*counter, fch)
			return
		}
		if metric, err := prometheus.NewConstMetric(counter.MetricDesc, prometheus.CounterValue, val); err == nil {
			fch <- metric
		} else {
			log.WithError(err).Errorln("collectCounters -> onCollectSuccess can not set the value")
			onCollectFail(*counter, fch)
			return
		}
		counter.lastSuccessValue = val
	}
	onCollectVecSuccess := func(counter *CounterMetric, fch chan<- prometheus.Metric, vals scrapper.NumericVecVals) {
		defer recoverNegativeCounter(*counter, fch)
		if metric, err := prometheus.NewConstMetric(counter.StatusDesc, prometheus.GaugeValue, 0); err == nil {
			fch <- metric
		} else {
			log.WithError(err).Errorln("collectCounters -> onCollectVecSuccess can not set down flag")
			onCollectFail(*counter, fch)
			return
		}
		for _, val := range vals {
			if metric, err := prometheus.NewConstMetric(counter.MetricDesc, prometheus.CounterValue, val.Val, val.Labels...); err == nil {
				fch <- metric
			} else {
				log.WithError(err).Errorln("collectCounters -> onCollectVecSuccess can not set the value")
				onCollectFail(*counter, fch)
				return
			}
		}
		counter.lastSuccessValue = vals
	}
	for idxCounter := range collector.Counters {
		if val, err := collector.Counters[idxCounter].scrapper.GetMetric(); err != nil {
			log.WithError(err).Errorln("can not get the data")
			onCollectFail(collector.Counters[idxCounter], ch)
		} else {
			switch val.(type) {
			case float64:
				counterVal, okCounterVal := val.(float64)
				if okCounterVal {
					onCollectSuccess(&(collector.Counters[idxCounter]), ch, counterVal)
				} else {
					log.WithField("val", val).Errorln(fmt.Sprintf("unable to get value %+v as float64", val))
					onCollectFail(collector.Counters[idxCounter], ch)
				}
			case scrapper.NumericVecVals:
				counterVecVal, okCounterVecVal := val.(scrapper.NumericVecVals)
				if okCounterVecVal {
					onCollectVecSuccess(&(collector.Counters[idxCounter]), ch, counterVecVal)
				} else {
					log.WithField("val", val).Errorln(fmt.Sprintf("unable to get value %+v as float64", val))
					onCollectFail(collector.Counters[idxCounter], ch)
				}
			default:
				log.WithField("val", val).Errorln(fmt.Sprintf("unable to determine value %+v type", val))
				onCollectFail(collector.Counters[idxCounter], ch)
			}
		}
	}
}

func (collector *SkycoinCollector) collectGauges(ch chan<- prometheus.Metric) {
	onCollectFail := func(gauge GaugeMetric, fch chan<- prometheus.Metric) {
		if metric, err := prometheus.NewConstMetric(gauge.StatusDesc, prometheus.GaugeValue, 1); err == nil {
			fch <- metric
		} else {
			log.WithError(err).Errorln("collectGauges -> onCollectFail can not set up flag")
		}
		switch gauge.lastSuccessValue.(type) {
		case float64:
			val := gauge.lastSuccessValue.(float64)
			if metric, err := prometheus.NewConstMetric(gauge.MetricDesc, prometheus.GaugeValue, val); err == nil {
				fch <- metric
			} else {
				log.WithError(err).Errorln("collectGauges -> onCollectFail can not set the last success value")
			}
		case scrapper.NumericVecVals:
			vals := gauge.lastSuccessValue.(scrapper.NumericVecVals)
			for _, val := range vals {
				if metric, err := prometheus.NewConstMetric(gauge.MetricDesc, prometheus.GaugeValue, val.Val, val.Labels...); err == nil {
					fch <- metric
				} else {
					log.WithError(err).Errorln("collectGauges -> onCollectFail can not set the last success vec value")
				}
			}
		}
	}
	onCollectSuccess := func(gauge *GaugeMetric, fch chan<- prometheus.Metric, val float64) {
		if metric, err := prometheus.NewConstMetric(gauge.StatusDesc, prometheus.GaugeValue, 0); err == nil {
			fch <- metric
		} else {
			log.WithError(err).Errorln("collectGauges -> onCollectSuccess can not set down flag")
			onCollectFail(*gauge, fch)
			return
		}
		if metric, err := prometheus.NewConstMetric(gauge.MetricDesc, prometheus.GaugeValue, val); err == nil {
			fch <- metric
		} else {
			log.WithError(err).Errorln("collectGauges -> onCollectSuccess can not set the value")
			onCollectFail(*gauge, fch)
			return
		}
		gauge.lastSuccessValue = val
	}
	onCollectVecSuccess := func(gauge *GaugeMetric, fch chan<- prometheus.Metric, vals scrapper.NumericVecVals) {
		if metric, err := prometheus.NewConstMetric(gauge.StatusDesc, prometheus.GaugeValue, 0); err == nil {
			fch <- metric
		} else {
			log.WithError(err).Errorln("collectGauges -> onCollectVecSuccess can not set down flag")
			onCollectFail(*gauge, fch)
			return
		}
		for _, val := range vals {
			if metric, err := prometheus.NewConstMetric(gauge.MetricDesc, prometheus.GaugeValue, val.Val, val.Labels...); err == nil {
				fch <- metric
			} else {
				log.WithError(err).Errorln("collectGauges -> onCollectVecSuccess can not set the value")
				onCollectFail(*gauge, fch)
				return
			}
		}
		gauge.lastSuccessValue = vals
	}
	for idxGauge := range collector.Gauges {
		if val, err := collector.Gauges[idxGauge].scrapper.GetMetric(); err != nil {
			log.WithError(err).Errorln("can not get the data")
			onCollectFail(collector.Gauges[idxGauge], ch)
		} else {
			switch val.(type) {
			case float64:
				gaugeVal, okGaugeVal := val.(float64)
				if okGaugeVal {
					onCollectSuccess(&(collector.Gauges[idxGauge]), ch, gaugeVal)
				} else {
					log.WithField("val", val).Errorln(fmt.Sprintf("unable to get value %+v as float64", val))
					onCollectFail(collector.Gauges[idxGauge], ch)
				}
			case scrapper.NumericVecVals:
				gaugeVecVal, okGaugeVecVal := val.(scrapper.NumericVecVals)
				if okGaugeVecVal {
					onCollectVecSuccess(&(collector.Gauges[idxGauge]), ch, gaugeVecVal)
				} else {
					log.WithField("val", val).Errorln(fmt.Sprintf("unable to get value %+v as float64", val))
					onCollectFail(collector.Gauges[idxGauge], ch)
				}
			default:
				log.WithField("val", val).Errorln(fmt.Sprintf("unable to determine value %+v type", val))
				onCollectFail(collector.Gauges[idxGauge], ch)
			}
		}
	}
}

func (collector *SkycoinCollector) collectHistograms(ch chan<- prometheus.Metric) {
	onCollectFail := func(histogram HistogramMetric, fch chan<- prometheus.Metric) {
		if metric, err := prometheus.NewConstMetric(histogram.StatusDesc, prometheus.GaugeValue, 1); err == nil {
			fch <- metric
		} else {
			log.WithError(err).Errorln("collectHistograms -> onCollectFail can not set up flag")
		}
		if metric, err := prometheus.NewConstHistogram(
			histogram.MetricDesc,
			histogram.lastSuccessValue.Count,
			histogram.lastSuccessValue.Sum,
			histogram.lastSuccessValue.Buckets,
		); err == nil {
			fch <- metric
		} else {
			log.WithError(err).Errorln("collectHistograms -> onCollectFail can not set the last success value")
		}
	}
	onCollectSuccess := func(histogram *HistogramMetric, fch chan<- prometheus.Metric, val scrapper.HistogramValue) {
		if metric, err := prometheus.NewConstMetric(histogram.StatusDesc, prometheus.GaugeValue, 0); err == nil {
			fch <- metric
		} else {
			log.WithError(err).Errorln("collectHistograms -> onCollectSuccess can not set down flag")
		}
		if metric, err := prometheus.NewConstHistogram(
			histogram.MetricDesc,
			val.Count,
			val.Sum,
			val.Buckets,
		); err == nil {
			fch <- metric
		} else {
			log.WithError(err).Errorln("collectHistograms -> onCollectSuccess can not set the value")
		}
		histogram.lastSuccessValue = val
	}
	for idxHistogram := range collector.Histograms {
		if val, err := collector.Histograms[idxHistogram].scrapper.GetMetric(); err != nil {
			log.WithError(err).Errorln("can not get the data")
			onCollectFail(collector.Histograms[idxHistogram], ch)
		} else {
			metricVal, okMetricVal := val.(scrapper.HistogramValue)
			if okMetricVal {
				onCollectSuccess(&(collector.Histograms[idxHistogram]), ch, metricVal)
			} else {
				log.WithError(err).Errorln("can not assert the metric value to type histogram")
			}
		}
	}
}

//Collect update all the descriptors is values
// TODO(denisacostaq@gmail.com): Make a research about race conditions here, "lastSuccessValue"
func (collector *SkycoinCollector) Collect(ch chan<- prometheus.Metric) {
	collector.collectCounters(ch)
	collector.collectGauges(ch)
	collector.collectHistograms(ch)
}
