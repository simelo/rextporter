package exporter

import (
	"errors"
	"fmt"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/simelo/rextporter/src/cache"
	"github.com/simelo/rextporter/src/config"
	"github.com/simelo/rextporter/src/scrapper"
	"github.com/simelo/rextporter/src/util"
	log "github.com/sirupsen/logrus"
)

// SkycoinCollector has the metrics to be exposed
type SkycoinCollector struct {
	metrics endpointData2MetricsConsumer
	cache   cache.Cache
}

func newSkycoinCollector(c cache.Cache, conf config.RootConfig) (collector *SkycoinCollector, err error) {
	const generalScopeErr = "error creating collector"
	var metrics endpointData2MetricsConsumer
	if metrics, err = createMetrics(c, conf.Services); err != nil {
		errCause := fmt.Sprintln("error creating metrics: ", err.Error())
		return nil, util.ErrorFromThisScope(errCause, generalScopeErr)
	}
	collector = &SkycoinCollector{metrics: metrics, cache: c}
	return collector, err
}

// Describe writes all the descriptors to the prometheus desc channel.
func (collector *SkycoinCollector) Describe(ch chan<- *prometheus.Desc) {
	for k := range collector.metrics {
		for idxMColl := range collector.metrics[k] {
			ch <- collector.metrics[k][idxMColl].metricDesc
			ch <- collector.metrics[k][idxMColl].statusDesc
		}
	}
}

func collectCounter(mColl *constMetric, ch chan<- prometheus.Metric) {
	onCollectFail := func(counter constMetric, fch chan<- prometheus.Metric) {
		if metric, err := prometheus.NewConstMetric(counter.statusDesc, prometheus.GaugeValue, 1); err == nil {
			fch <- metric
		} else {
			log.WithError(err).Errorln("collectCounter -> onCollectFail can not set up flag")
		}
		switch counter.lastSuccessValue.(type) {
		case float64:
			val, okVal := counter.lastSuccessValue.(float64)
			if !okVal {
				log.WithField("val", val).Errorln("can not get las success val")
			}
			if metric, err := prometheus.NewConstMetric(counter.metricDesc, prometheus.CounterValue, val); err == nil {
				fch <- metric
			} else {
				log.WithError(err).Errorln("collectCounter -> onCollectFail can not set the last success value")
			}
		case scrapper.NumericVecVals:
			vals, okVals := counter.lastSuccessValue.(scrapper.NumericVecVals)
			if !okVals {
				log.WithField("vals", vals).Errorln("can not get las success vals")
			}
			for _, val := range vals {
				if metric, err := prometheus.NewConstMetric(counter.metricDesc, prometheus.CounterValue, val.Val, val.Labels...); err == nil {
					fch <- metric
				} else {
					log.WithError(err).Errorln("collectCounter -> onCollectFail can not set the last success vec value")
				}
			}
		}
	}
	recoverNegativeCounter := func(counter constMetric, fch chan<- prometheus.Metric) {
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
	onCollectSuccess := func(counter *constMetric, fch chan<- prometheus.Metric, val float64) {
		defer recoverNegativeCounter(*counter, fch)
		if metric, err := prometheus.NewConstMetric(counter.statusDesc, prometheus.GaugeValue, 0); err == nil {
			fch <- metric
		} else {
			log.WithError(err).Errorln("collectCounter -> onCollectSuccess can not set down flag")
			onCollectFail(*counter, fch)
			return
		}
		if metric, err := prometheus.NewConstMetric(counter.metricDesc, prometheus.CounterValue, val); err == nil {
			fch <- metric
		} else {
			log.WithError(err).Errorln("collectCounter -> onCollectSuccess can not set the value")
			onCollectFail(*counter, fch)
			return
		}
		counter.lastSuccessValue = val
	}
	onCollectVecSuccess := func(counter *constMetric, fch chan<- prometheus.Metric, vals scrapper.NumericVecVals) {
		defer recoverNegativeCounter(*counter, fch)
		if metric, err := prometheus.NewConstMetric(counter.statusDesc, prometheus.GaugeValue, 0); err == nil {
			fch <- metric
		} else {
			log.WithError(err).Errorln("collectCounter -> onCollectVecSuccess can not set down flag")
			onCollectFail(*counter, fch)
			return
		}
		for _, val := range vals {
			if metric, err := prometheus.NewConstMetric(counter.metricDesc, prometheus.CounterValue, val.Val, val.Labels...); err == nil {
				fch <- metric
			} else {
				log.WithError(err).Errorln("collectCounter -> onCollectVecSuccess can not set the value")
				onCollectFail(*counter, fch)
				return
			}
		}
		counter.lastSuccessValue = vals
	}
	if val, err := mColl.scrapper.GetMetric(); err != nil {
		log.WithError(err).Errorln("can not get the data")
		onCollectFail(*mColl, ch)
	} else {
		switch val.(type) {
		case float64:
			counterVal, okCounterVal := val.(float64)
			if okCounterVal {
				onCollectSuccess(mColl, ch, counterVal)
			} else {
				log.WithField("val", val).Errorln(fmt.Sprintf("unable to get value %+v as float64", val))
				onCollectFail(*mColl, ch)
			}
		case scrapper.NumericVecVals:
			counterVecVal, okCounterVecVal := val.(scrapper.NumericVecVals)
			if okCounterVecVal {
				onCollectVecSuccess(mColl, ch, counterVecVal)
			} else {
				log.WithField("val", val).Errorln(fmt.Sprintf("unable to get value %+v as float64", val))
				onCollectFail(*mColl, ch)
			}
		default:
			log.WithField("val", val).Errorln(fmt.Sprintf("unable to determine value %+v type", val))
			onCollectFail(*mColl, ch)
		}
	}
}

func collectGauge(mColl *constMetric, ch chan<- prometheus.Metric) {
	onCollectFail := func(gauge constMetric, fch chan<- prometheus.Metric) {
		if metric, err := prometheus.NewConstMetric(gauge.statusDesc, prometheus.GaugeValue, 1); err == nil {
			fch <- metric
		} else {
			log.WithError(err).Errorln("collectGauge -> onCollectFail can not set up flag")
		}
		switch gauge.lastSuccessValue.(type) {
		case float64:
			val := gauge.lastSuccessValue.(float64)
			if metric, err := prometheus.NewConstMetric(gauge.metricDesc, prometheus.GaugeValue, val); err == nil {
				fch <- metric
			} else {
				log.WithError(err).Errorln("collectGauge -> onCollectFail can not set the last success value")
			}
		case scrapper.NumericVecVals:
			vals := gauge.lastSuccessValue.(scrapper.NumericVecVals)
			for _, val := range vals {
				if metric, err := prometheus.NewConstMetric(gauge.metricDesc, prometheus.GaugeValue, val.Val, val.Labels...); err == nil {
					fch <- metric
				} else {
					log.WithError(err).Errorln("collectGauge -> onCollectFail can not set the last success vec value")
				}
			}
		}
	}
	onCollectSuccess := func(gauge *constMetric, fch chan<- prometheus.Metric, val float64) {
		if metric, err := prometheus.NewConstMetric(gauge.statusDesc, prometheus.GaugeValue, 0); err == nil {
			fch <- metric
		} else {
			log.WithError(err).Errorln("collectGauge -> onCollectSuccess can not set down flag")
			onCollectFail(*gauge, fch)
			return
		}
		if metric, err := prometheus.NewConstMetric(gauge.metricDesc, prometheus.GaugeValue, val); err == nil {
			fch <- metric
		} else {
			log.WithError(err).Errorln("collectGauge -> onCollectSuccess can not set the value")
			onCollectFail(*gauge, fch)
			return
		}
		gauge.lastSuccessValue = val
	}
	onCollectVecSuccess := func(gauge *constMetric, fch chan<- prometheus.Metric, vals scrapper.NumericVecVals) {
		if metric, err := prometheus.NewConstMetric(gauge.statusDesc, prometheus.GaugeValue, 0); err == nil {
			fch <- metric
		} else {
			log.WithError(err).Errorln("collectGauge -> onCollectVecSuccess can not set down flag")
			onCollectFail(*gauge, fch)
			return
		}
		for _, val := range vals {
			if metric, err := prometheus.NewConstMetric(gauge.metricDesc, prometheus.GaugeValue, val.Val, val.Labels...); err == nil {
				fch <- metric
			} else {
				log.WithError(err).Errorln("collectGauge -> onCollectVecSuccess can not set the value")
				onCollectFail(*gauge, fch)
				return
			}
		}
		gauge.lastSuccessValue = vals
	}
	if val, err := mColl.scrapper.GetMetric(); err != nil {
		log.WithError(err).Errorln("can not get the data")
		onCollectFail(*mColl, ch)
	} else {
		switch val.(type) {
		case float64:
			gaugeVal, okGaugeVal := val.(float64)
			if okGaugeVal {
				onCollectSuccess(mColl, ch, gaugeVal)
			} else {
				log.WithField("val", val).Errorln(fmt.Sprintf("unable to get value %+v as float64", val))
				onCollectFail(*mColl, ch)
			}
		case scrapper.NumericVecVals:
			gaugeVecVal, okGaugeVecVal := val.(scrapper.NumericVecVals)
			if okGaugeVecVal {
				onCollectVecSuccess(mColl, ch, gaugeVecVal)
			} else {
				log.WithField("val", val).Errorln(fmt.Sprintf("unable to get value %+v as float64", val))
				onCollectFail(*mColl, ch)
			}
		default:
			log.WithField("val", val).Errorln(fmt.Sprintf("unable to determine value %+v type", val))
			onCollectFail(*mColl, ch)
		}
	}
}

func collectHistogram(mColl *constMetric, ch chan<- prometheus.Metric) {
	onCollectFail := func(histogram constMetric, fch chan<- prometheus.Metric) {
		if metric, err := prometheus.NewConstMetric(histogram.statusDesc, prometheus.GaugeValue, 1); err == nil {
			fch <- metric
		} else {
			log.WithError(err).Errorln("collectHistogram -> onCollectFail can not set up flag")
		}
		lastSuccessVal, okLastSuccessVal := histogram.lastSuccessValue.(scrapper.HistogramValue)
		if okLastSuccessVal {
			if metric, err := prometheus.NewConstHistogram(
				histogram.metricDesc,
				lastSuccessVal.Count,
				lastSuccessVal.Sum,
				lastSuccessVal.Buckets,
			); err == nil {
				fch <- metric
			} else {
				log.WithError(err).Errorln("collectHistogram -> onCollectFail can not set the last success value")
			}
		} else {
			log.WithField("val", histogram).Errorln(fmt.Sprintf("unable to get value %+v as histogram val", histogram))
		}
	}
	onCollectSuccess := func(histogram *constMetric, fch chan<- prometheus.Metric, val scrapper.HistogramValue) {
		if metric, err := prometheus.NewConstMetric(histogram.statusDesc, prometheus.GaugeValue, 0); err == nil {
			fch <- metric
		} else {
			log.WithError(err).Errorln("collectHistogram -> onCollectSuccess can not set down flag")
		}
		if metric, err := prometheus.NewConstHistogram(
			histogram.metricDesc,
			val.Count,
			val.Sum,
			val.Buckets,
		); err == nil {
			fch <- metric
		} else {
			log.WithError(err).Errorln("collectHistogram -> onCollectSuccess can not set the value")
		}
		histogram.lastSuccessValue = val
	}
	if val, err := mColl.scrapper.GetMetric(); err != nil {
		log.WithError(err).Errorln("can not get the data")
		onCollectFail(*mColl, ch)
	} else {
		metricVal, okMetricVal := val.(scrapper.HistogramValue)
		if okMetricVal {
			onCollectSuccess(mColl, ch, metricVal)
		} else {
			log.WithError(err).Errorln("can not assert the metric value to type histogram")
		}
	}
}

//Collect update all the descriptors is values
// TODO(denisacostaq@gmail.com): Make a research about race conditions here, "lastSuccessValue"
func (collector *SkycoinCollector) Collect(ch chan<- prometheus.Metric) {
	for k := range collector.metrics {
		for idxMetric := range collector.metrics[k] {
			switch collector.metrics[k][idxMetric].kind {
			case config.KeyTypeCounter:
				collectCounter(&(collector.metrics[k][idxMetric]), ch)
			case config.KeyTypeGauge:
				collectGauge(&(collector.metrics[k][idxMetric]), ch)
			case config.KeyTypeHistogram:
				collectHistogram(&(collector.metrics[k][idxMetric]), ch)
			default:
				log.Println("error")
			}
		}
		collector.cache.Reset()
	}
}
