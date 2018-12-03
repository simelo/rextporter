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

// MetricsCollector has the metrics to be exposed
type MetricsCollector struct {
	metrics endpointData2MetricsConsumer
	cache   cache.Cache
}

func newMetricsCollector(c cache.Cache, conf config.RootConfig) (collector *MetricsCollector, err error) {
	const generalScopeErr = "error creating collector"
	var metrics endpointData2MetricsConsumer
	if metrics, err = createMetrics(c, conf.Services); err != nil {
		errCause := fmt.Sprintln("error creating metrics: ", err.Error())
		return nil, util.ErrorFromThisScope(errCause, generalScopeErr)
	}
	collector = &MetricsCollector{metrics: metrics, cache: c}
	return collector, err
}

// Describe writes all the descriptors to the prometheus desc channel.
func (collector *MetricsCollector) Describe(ch chan<- *prometheus.Desc) {
	for k := range collector.metrics {
		for idxMColl := range collector.metrics[k] {
			ch <- collector.metrics[k][idxMColl].metricDesc
			ch <- collector.metrics[k][idxMColl].statusDesc
		}
	}
}

func collectCounters(metricsColl []constMetric, ch chan<- prometheus.Metric) {
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
	resC := make(chan scrapper.ScrapResult)
	defer close(resC)
	errC := make(chan scrapper.ScrapErrResult)
	defer close(errC)
	var metrcisNum = 0
	for idxM, mColl := range metricsColl {
		scrapper.WorkPool.Apply(
			scrapper.ScrapRequest{Scrap: mColl.scrapper, Res: resC, ConstMetricIdxIn: idxM, Err: errC})
		metrcisNum++
	}
	for i := 0; i < metrcisNum; i++ {
		select {
		case res := <-resC:
			switch res.Val.(type) {
			case float64:
				counterVal, okCounterVal := res.Val.(float64)
				if okCounterVal {
					onCollectSuccess(&(metricsColl[res.ConstMetricIdxOut]), ch, counterVal)
				} else {
					log.WithField("val", res.Val).Errorln(fmt.Sprintf("unable to get value %+v as float64", res.Val))
					onCollectFail(metricsColl[res.ConstMetricIdxOut], ch)
				}
			case scrapper.NumericVecVals:
				counterVecVal, okCounterVecVal := res.Val.(scrapper.NumericVecVals)
				if okCounterVecVal {
					onCollectVecSuccess(&(metricsColl[res.ConstMetricIdxOut]), ch, counterVecVal)
				} else {
					log.WithField("val", res.Val).Errorln(fmt.Sprintf("unable to get value %+v as float64", res.Val))
					onCollectFail(metricsColl[res.ConstMetricIdxOut], ch)
				}
			default:
				log.WithField("val", res.Val).Errorln(fmt.Sprintf("unable to determine value %+v type", res.Val))
				onCollectFail(metricsColl[res.ConstMetricIdxOut], ch)
			}
		case err := <-errC:
			log.WithError(err.Err).Errorln("can not get the data")
			onCollectFail(metricsColl[err.ConstMetricIdxOut], ch)
		}
	}
}

func collectGauges(metricsColl []constMetric, ch chan<- prometheus.Metric) {
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
	resC := make(chan scrapper.ScrapResult)
	defer close(resC)
	errC := make(chan scrapper.ScrapErrResult)
	defer close(errC)
	var metrcisNum = 0
	for idxM, mColl := range metricsColl {
		scrapper.WorkPool.Apply(
			scrapper.ScrapRequest{Scrap: mColl.scrapper, Res: resC, ConstMetricIdxIn: idxM, Err: errC})
		metrcisNum++
	}
	for i := 0; i < metrcisNum; i++ {
		select {
		case res := <-resC:
			switch res.Val.(type) {
			case float64:
				gaugeVal, okGaugeVal := res.Val.(float64)
				if okGaugeVal {
					onCollectSuccess(&(metricsColl[res.ConstMetricIdxOut]), ch, gaugeVal)
				} else {
					log.WithField("val", res.Val).Errorln(fmt.Sprintf("unable to get value %+v as float64", res.Val))
					onCollectFail(metricsColl[res.ConstMetricIdxOut], ch)
				}
			case scrapper.NumericVecVals:
				gaugeVecVal, okGaugeVecVal := res.Val.(scrapper.NumericVecVals)
				if okGaugeVecVal {
					onCollectVecSuccess(&(metricsColl[res.ConstMetricIdxOut]), ch, gaugeVecVal)
				} else {
					log.WithField("val", res.Val).Errorln(fmt.Sprintf("unable to get value %+v as float64", res.Val))
					onCollectFail(metricsColl[res.ConstMetricIdxOut], ch)
				}
			default:
				log.WithField("val", res.Val).Errorln(fmt.Sprintf("unable to determine value %+v type", res.Val))
				onCollectFail(metricsColl[res.ConstMetricIdxOut], ch)
			}
		case err := <-errC:
			log.WithError(err.Err).Errorln("can not get the data")
			onCollectFail(metricsColl[err.ConstMetricIdxOut], ch)
		}
	}
}

func collectHistograms(metricsColl []constMetric, ch chan<- prometheus.Metric) {
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
	resC := make(chan scrapper.ScrapResult)
	defer close(resC)
	errC := make(chan scrapper.ScrapErrResult)
	defer close(errC)
	var metrcisNum = 0
	for idxM, mColl := range metricsColl {
		scrapper.WorkPool.Apply(
			scrapper.ScrapRequest{Scrap: mColl.scrapper, Res: resC, ConstMetricIdxIn: idxM, Err: errC})
		metrcisNum++
	}
	for i := 0; i < metrcisNum; i++ {
		select {
		case res := <-resC:
			metricVal, okMetricVal := res.Val.(scrapper.HistogramValue)
			if okMetricVal {
				onCollectSuccess(&(metricsColl[res.ConstMetricIdxOut]), ch, metricVal)
			} else {
				log.WithField("val", res.Val).Errorln("can not assert the metric value to type histogram")
			}
		case err := <-errC:
			log.WithError(err.Err).Errorln("can not get the data")
			onCollectFail(metricsColl[err.ConstMetricIdxOut], ch)
		}
	}
}

//Collect update all the descriptors is values
// TODO(denisacostaq@gmail.com): Make a research about race conditions here, "lastSuccessValue"
func (collector *MetricsCollector) Collect(ch chan<- prometheus.Metric) {
	filterMetricsByKind := func(kind string, orgMetrics []constMetric) (filteredMetrics []constMetric) {
		for _, metric := range orgMetrics {
			if metric.kind == kind {
				filteredMetrics = append(filteredMetrics, metric)
			}
		}
		return filteredMetrics
	}
	for k := range collector.metrics {
		counters := filterMetricsByKind(config.KeyTypeCounter, collector.metrics[k])
		gauges := filterMetricsByKind(config.KeyTypeGauge, collector.metrics[k])
		histograms := filterMetricsByKind(config.KeyTypeHistogram, collector.metrics[k])
		collectCounters(counters, ch)
		collectGauges(gauges, ch)
		collectHistograms(histograms, ch)
		collector.cache.Reset()
	}
}
