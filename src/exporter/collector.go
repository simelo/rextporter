package exporter

import (
	"errors"
	"fmt"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/simelo/rextporter/src/cache"
	"github.com/simelo/rextporter/src/config"
	"github.com/simelo/rextporter/src/scrapper"
	"github.com/simelo/rextporter/src/util"
	log "github.com/sirupsen/logrus"
)

// MetricsCollector has the metrics to be exposed
type MetricsCollector struct {
	metrics    endpointData2MetricsConsumer
	cache      cache.Cache
	defMetrics *defaultMetrics
}

func newMetricsCollector(c cache.Cache, conf config.RootConfig) (collector *MetricsCollector, err error) {
	const generalScopeErr = "error creating collector"
	var metrics endpointData2MetricsConsumer
	if metrics, err = createMetrics(c, conf.Services); err != nil {
		errCause := fmt.Sprintln("error creating metrics: ", err.Error())
		return nil, util.ErrorFromThisScope(errCause, generalScopeErr)
	}
	collector = &MetricsCollector{metrics: metrics, cache: c, defMetrics: newDefaultMetrics()}
	return collector, err
}

// Describe writes all the descriptors to the prometheus desc channel.
func (collector *MetricsCollector) Describe(ch chan<- *prometheus.Desc) {
	for k := range collector.metrics {
		for idxMColl := range collector.metrics[k] {
			ch <- collector.metrics[k][idxMColl].metricDesc
		}
	}
}

func collectCounters(metricsColl []constMetric, defMetrics *defaultMetrics, ch chan<- prometheus.Metric) {
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
		}
	}
	onCollectSuccess := func(counter *constMetric, jobName, instanceName string, fch chan<- prometheus.Metric, val float64) {
		defer recoverNegativeCounter(*counter, fch)
		if metric, err := prometheus.NewConstMetric(counter.metricDesc, prometheus.CounterValue, val, jobName, instanceName); err == nil {
			fch <- metric
		} else {
			log.WithError(err).Errorln("collectCounter -> onCollectSuccess can not set the value")
			// FIXME(denisacostaq@gmail.com): onCollectFail(*counter, jobName, instanceName, fch)
			return
		}
	}
	onCollectVecSuccess := func(counter *constMetric, jobName, instanceName string, fch chan<- prometheus.Metric, vals scrapper.NumericVecVals) {
		defer recoverNegativeCounter(*counter, fch)
		for _, val := range vals {
			labels := append(val.Labels, jobName, instanceName)
			if metric, err := prometheus.NewConstMetric(counter.metricDesc, prometheus.CounterValue, val.Val, labels...); err == nil {
				fch <- metric
			} else {
				log.WithError(err).Errorln("collectCounter -> onCollectVecSuccess can not set the value")
			}
		}
	}
	resC := make(chan scrapper.ScrapResult)
	defer close(resC)
	errC := make(chan scrapper.ScrapErrResult)
	defer close(errC)
	var metricsNum = 0
	startScrappingInPool := time.Now().UTC()
	for idxM, mColl := range metricsColl {
		scrapper.WorkPool.Apply(
			scrapper.ScrapRequest{
				Scrap:            mColl.scrapper,
				Res:              resC,
				ConstMetricIdxIn: idxM,
				JobName:          mColl.scrapper.GetJobName(),
				InstanceName:     mColl.scrapper.GetInstanceName(),
				Err:              errC,
			},
		)
		metricsNum++
	}
	for i := 0; i < metricsNum; i++ {
		select {
		case res := <-resC:
			switch res.Val.(type) {
			case float64:
				counterVal, okCounterVal := res.Val.(float64)
				defMetrics.scrapedSamples.addSeconds(1, res.JobName, res.InstanceName)
				if okCounterVal {
					onCollectSuccess(&(metricsColl[res.ConstMetricIdxOut]), res.JobName, res.InstanceName, ch, counterVal)
				} else {
					log.WithField("val", res.Val).Errorln(fmt.Sprintf("unable to get value %+v as float64", res.Val))
				}
			case scrapper.NumericVecVals:
				counterVecVal, okCounterVecVal := res.Val.(scrapper.NumericVecVals)
				defMetrics.scrapedSamples.addSeconds(float64(len(counterVecVal)), res.JobName, res.InstanceName)
				if okCounterVecVal {
					onCollectVecSuccess(&(metricsColl[res.ConstMetricIdxOut]), res.JobName, res.InstanceName, ch, counterVecVal)
				} else {
					log.WithField("val", res.Val).Errorln(fmt.Sprintf("unable to get value %+v as float64", res.Val))
				}
			default:
				log.WithField("val", res.Val).Errorln(fmt.Sprintf("unable to determine value %+v type", res.Val))
			}
			defMetrics.scrapedDurations.addSeconds(time.Since(startScrappingInPool).Seconds(), res.JobName, res.InstanceName)
		case err := <-errC:
			log.WithError(err.Err).Errorln("can not get the data")
			// FIXME(denisacostaq@gmail.com): onCollectFail(metricsColl[err.ConstMetricIdxOut], err.JobName, err.InstanceName, ch)
			defMetrics.scrapedDurations.addSeconds(time.Since(startScrappingInPool).Seconds(), err.JobName, err.InstanceName)
		}
	}
}

func collectGauges(metricsColl []constMetric, defMetrics *defaultMetrics, ch chan<- prometheus.Metric) {
	onCollectSuccess := func(gauge *constMetric, jobName, instanceName string, fch chan<- prometheus.Metric, val float64) {
		if metric, err := prometheus.NewConstMetric(gauge.metricDesc, prometheus.GaugeValue, val, jobName, instanceName); err == nil {
			fch <- metric
		} else {
			log.WithError(err).Errorln("collectGauge -> onCollectSuccess can not set the value")
			// FIXME(denisacostaq@gmail.com): onCollectFail(*gauge, jobName, instanceName, fch)
			return
		}
	}
	onCollectVecSuccess := func(gauge *constMetric, jobName, instanceName string, fch chan<- prometheus.Metric, vals scrapper.NumericVecVals) {
		for _, val := range vals {
			labels := append(val.Labels, jobName, instanceName)
			if metric, err := prometheus.NewConstMetric(gauge.metricDesc, prometheus.GaugeValue, val.Val, labels...); err == nil {
				fch <- metric
			} else {
				log.WithError(err).Errorln("collectGauge -> onCollectVecSuccess can not set the value")
			}
		}
	}
	resC := make(chan scrapper.ScrapResult)
	defer close(resC)
	errC := make(chan scrapper.ScrapErrResult)
	defer close(errC)
	var metricsNum = 0
	startScrappingInPool := time.Now().UTC()
	for idxM, mColl := range metricsColl {
		scrapper.WorkPool.Apply(
			scrapper.ScrapRequest{
				Scrap:            mColl.scrapper,
				Res:              resC,
				ConstMetricIdxIn: idxM,
				JobName:          mColl.scrapper.GetJobName(),
				InstanceName:     mColl.scrapper.GetInstanceName(),
				Err:              errC,
			},
		)
		metricsNum++
	}
	for i := 0; i < metricsNum; i++ {
		select {
		case res := <-resC:
			switch res.Val.(type) {
			case float64:
				defMetrics.scrapedSamples.addSeconds(1, res.JobName, res.InstanceName)
				gaugeVal, okGaugeVal := res.Val.(float64)
				if okGaugeVal {
					onCollectSuccess(&(metricsColl[res.ConstMetricIdxOut]), res.JobName, res.InstanceName, ch, gaugeVal)
				} else {
					log.WithField("val", res.Val).Errorln(fmt.Sprintf("unable to get value %+v as float64", res.Val))
					// FIXME(denisacostaq@gmail.com): onCollectFail(metricsColl[res.ConstMetricIdxOut], res.JobName, res.InstanceName, ch)
				}
			case scrapper.NumericVecVals:
				gaugeVecVal, okGaugeVecVal := res.Val.(scrapper.NumericVecVals)
				defMetrics.scrapedSamples.addSeconds(float64(len(gaugeVecVal)), res.JobName, res.InstanceName)
				if okGaugeVecVal {
					onCollectVecSuccess(&(metricsColl[res.ConstMetricIdxOut]), res.JobName, res.InstanceName, ch, gaugeVecVal)
				} else {
					log.WithField("val", res.Val).Errorln(fmt.Sprintf("unable to get value %+v as float64", res.Val))
					// FIXME(denisacostaq@gmail.com): onCollectFail(metricsColl[res.ConstMetricIdxOut], res.JobName, res.InstanceName, ch)
				}
			default:
				log.WithField("val", res.Val).Errorln(fmt.Sprintf("unable to determine value %+v type", res.Val))
				// FIXME(denisacostaq@gmail.com): onCollectFail(metricsColl[res.ConstMetricIdxOut], res.JobName, res.InstanceName, ch)
			}
			defMetrics.scrapedDurations.addSeconds(time.Since(startScrappingInPool).Seconds(), res.JobName, res.InstanceName)
		case err := <-errC:
			log.WithError(err.Err).Errorln("can not get the data")
			// FIXME(denisacostaq@gmail.com): onCollectFail(metricsColl[err.ConstMetricIdxOut], err.JobName, err.InstanceName, ch)
			defMetrics.scrapedDurations.addSeconds(time.Since(startScrappingInPool).Seconds(), err.JobName, err.InstanceName)
		}
	}
}

func collectHistograms(metricsColl []constMetric, defMetrics *defaultMetrics, ch chan<- prometheus.Metric) {
	onCollectSuccess := func(histogram *constMetric, jobName, instanceName string, fch chan<- prometheus.Metric, val scrapper.HistogramValue) {
		if metric, err := prometheus.NewConstHistogram(
			histogram.metricDesc,
			val.Count,
			val.Sum,
			val.Buckets,
			jobName,
			instanceName,
		); err == nil {
			fch <- metric
		} else {
			log.WithError(err).Errorln("collectHistogram -> onCollectSuccess can not set the value")
		}
	}
	resC := make(chan scrapper.ScrapResult)
	defer close(resC)
	errC := make(chan scrapper.ScrapErrResult)
	defer close(errC)
	var metricsNum = 0
	startScrappingInPool := time.Now().UTC()
	for idxM, mColl := range metricsColl {
		scrapper.WorkPool.Apply(
			scrapper.ScrapRequest{
				Scrap:            mColl.scrapper,
				Res:              resC,
				ConstMetricIdxIn: idxM,
				JobName:          mColl.scrapper.GetJobName(),
				InstanceName:     mColl.scrapper.GetInstanceName(),
				Err:              errC,
			},
		)
		metricsNum++
	}
	for i := 0; i < metricsNum; i++ {
		select {
		case res := <-resC:
			metricVal, okMetricVal := res.Val.(scrapper.HistogramValue)
			defMetrics.scrapedSamples.addSeconds(float64(len(metricVal.Buckets)+2), res.JobName, res.InstanceName)
			if okMetricVal {
				onCollectSuccess(&(metricsColl[res.ConstMetricIdxOut]), res.JobName, res.InstanceName, ch, metricVal)
			} else {
				log.WithField("val", res.Val).Errorln("can not assert the metric value to type histogram")
			}
			defMetrics.scrapedDurations.addSeconds(time.Since(startScrappingInPool).Seconds(), res.JobName, res.InstanceName)
		case err := <-errC:
			log.WithError(err.Err).Errorln("can not get the data")
			// FIXME(denisacostaq@gmail.com): onCollectFail(metricsColl[err.ConstMetricIdxOut], err.JobName, err.InstanceName, ch)
			defMetrics.scrapedDurations.addSeconds(time.Since(startScrappingInPool).Seconds(), err.JobName, err.InstanceName)
		}
	}
}

//Collect update all the descriptors is values
func (collector *MetricsCollector) Collect(ch chan<- prometheus.Metric) {
	filterMetricsByKind := func(kind string, orgMetrics []constMetric) (filteredMetrics []constMetric) {
		for _, metric := range orgMetrics {
			if metric.kind == kind {
				filteredMetrics = append(filteredMetrics, metric)
			}
		}
		return filteredMetrics
	}
	collector.defMetrics.reset()
	for k := range collector.metrics {
		counters := filterMetricsByKind(config.KeyTypeCounter, collector.metrics[k])
		gauges := filterMetricsByKind(config.KeyTypeGauge, collector.metrics[k])
		histograms := filterMetricsByKind(config.KeyTypeHistogram, collector.metrics[k])
		collectCounters(counters, collector.defMetrics, ch)
		collectGauges(gauges, collector.defMetrics, ch)
		collectHistograms(histograms, collector.defMetrics, ch)
		collector.cache.Reset()
	}
	collector.defMetrics.collectDefaultMetrics(ch)
}
