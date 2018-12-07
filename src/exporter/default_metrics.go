package exporter

import (
	"github.com/prometheus/client_golang/prometheus"
	log "github.com/sirupsen/logrus"
)

type scrapDurationInInstance map[string]float64
type scrapDurationInJob map[string]scrapDurationInInstance

func newScrapDuration() scrapDurationInJob {
	return make(scrapDurationInJob)
}

func (sd scrapDurationInJob) addSeconds(amount float64, jobName, instanceName string) {
	instances, okInstances := sd[jobName]
	if okInstances {
		instances[instanceName] += amount
	} else {
		instances := make(scrapDurationInInstance)
		instances[instanceName] = amount
		sd[jobName] = instances
	}
}

type defaultMetrics struct {
	scrapedDurations          scrapDurationInJob
	scrapeDurationSecondsDesc *prometheus.Desc
	scrapedSamples            scrapDurationInJob
	scrapeSamplesScrapedDesc  *prometheus.Desc
	datasourceResponseDurationDesc *prometheus.Desc
}

func (defMetrics defaultMetrics) describe(ch chan<- *prometheus.Desc) {
	ch <- defMetrics.scrapeDurationSecondsDesc
	ch <- defMetrics.scrapeSamplesScrapedDesc
	ch <- defMetrics.datasourceResponseDurationDesc
}

func (defMetrics *defaultMetrics) reset() {
	defMetrics.scrapedDurations = newScrapDuration()
	defMetrics.scrapedSamples = newScrapDuration()
}

func (defMetrics defaultMetrics) collectDefaultMetrics(ch chan<- prometheus.Metric) {
	for jobName, job := range defMetrics.scrapedDurations {
		for instanceName, val := range job {
			labels := []string{jobName, instanceName}
			if metric, err := prometheus.NewConstMetric(defMetrics.scrapeDurationSecondsDesc, prometheus.GaugeValue, val, labels...); err == nil {
				ch <- metric
			} else {
				log.WithError(err).Errorln("collectDefaultMetrics -> scrapeDurationSeconds")
			}
		}
	}
	for jobName, job := range defMetrics.scrapedSamples {
		for instanceName, val := range job {
			labels := []string{jobName, instanceName}
			if metric, err := prometheus.NewConstMetric(defMetrics.scrapeSamplesScrapedDesc, prometheus.GaugeValue, val, labels...); err == nil {
				ch <- metric
			} else {
				log.WithError(err).Errorln("collectDefaultMetrics -> scrapeSamplesScrapedDesc")
			}
		}
	}
}
