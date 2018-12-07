package exporter

import (
	"github.com/prometheus/client_golang/prometheus"
	log "github.com/sirupsen/logrus"
)

type scrapDurationInInstance map[string]float64
type scrapDurationInJob map[string]scrapDurationInInstance

type datasource2Value map[string]float64
type instance2D2V map[string]datasource2Value
type job2I map[string]instance2D2V

func newResponseDuration() job2I {
	return make(job2I)
}

var instance4JobLabels = []string{"job", "instance"}

func newDefaultMetrics() *defaultMetrics {
	scrapeDurationDesc := prometheus.NewDesc(
		"scrape_duration",
		"Elapse time(in seconds) to get a response from a scrapper",
		instance4JobLabels,
		nil,
	)
	scrapeSamplesScrapedDesc := prometheus.NewDesc(
		"scrape_samples_scraped",
		"The number of samples the target exposed",
		instance4JobLabels,
		nil,
	)
	dataSourceScrapeDurationDesc := prometheus.NewDesc(
		"datasource_scrape_duration",
		"Elapse time(in seconds) to get a response from a datasource in a scrapper",
		append(instance4JobLabels, "datasource"),
		nil,
	)
	datasourceResponseDurationDesc := prometheus.NewDesc(
		"datasource_response_duration",
		"Elapse time(in seconds) to get a response from a datasource",
		append(instance4JobLabels, "datasource"),
		nil,
	)
	return &defaultMetrics{
		scrapeDurations:                newScrapDuration(),
		scrapeDurationDesc:             scrapeDurationDesc,
		scrapeSamples:                  newScrapDuration(),
		scrapeSamplesScrapedDesc:       scrapeSamplesScrapedDesc,
		datasourceResponseDurationDesc: datasourceResponseDurationDesc,
		dataSourceScrapeDuration:       newResponseDuration(),
		dataSourceScrapeDurationDesc:   dataSourceScrapeDurationDesc,
	}
}

func (instances job2I) add(amount float64, jobName, instanceName, datasourceName string) {
	datasources, okDatasources := instances[jobName]
	if okDatasources {
		values, okValues := datasources[instanceName]
		if okValues {
			values[datasourceName] += amount
		} else {
			values := make(datasource2Value)
			values[datasourceName] = amount
			datasources[instanceName] = values
		}
	} else {
		values := make(datasource2Value)
		values[datasourceName] = amount
		datasources := make(instance2D2V)
		datasources[instanceName] = values
		instances[jobName] = datasources
	}
}

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
	scrapeDurations                scrapDurationInJob
	scrapeDurationDesc             *prometheus.Desc
	scrapeSamples                  scrapDurationInJob
	scrapeSamplesScrapedDesc       *prometheus.Desc
	dataSourceScrapeDuration       job2I
	dataSourceScrapeDurationDesc   *prometheus.Desc
	datasourceResponseDurationDesc *prometheus.Desc
}

func (defMetrics defaultMetrics) describe(ch chan<- *prometheus.Desc) {
	ch <- defMetrics.scrapeDurationDesc
	ch <- defMetrics.scrapeSamplesScrapedDesc
	ch <- defMetrics.dataSourceScrapeDurationDesc
	ch <- defMetrics.datasourceResponseDurationDesc
}

func (defMetrics *defaultMetrics) reset() {
	defMetrics.scrapeDurations = newScrapDuration()
	defMetrics.scrapeSamples = newScrapDuration()
	defMetrics.dataSourceScrapeDuration = newResponseDuration()
}

func (defMetrics defaultMetrics) collectDefaultMetrics(ch chan<- prometheus.Metric) {
	for jobName, job := range defMetrics.scrapeDurations {
		for instanceName, val := range job {
			labels := []string{jobName, instanceName}
			if metric, err := prometheus.NewConstMetric(defMetrics.scrapeDurationDesc, prometheus.GaugeValue, val, labels...); err == nil {
				ch <- metric
			} else {
				log.WithError(err).Errorln("collectDefaultMetrics -> scrapeDurationDesc")
			}
		}
	}
	for jobName, job := range defMetrics.scrapeSamples {
		for instanceName, val := range job {
			labels := []string{jobName, instanceName}
			if metric, err := prometheus.NewConstMetric(defMetrics.scrapeSamplesScrapedDesc, prometheus.GaugeValue, val, labels...); err == nil {
				ch <- metric
			} else {
				log.WithError(err).Errorln("collectDefaultMetrics -> scrapeSamplesScrapedDesc")
			}
		}
	}
	for jobName, instances := range defMetrics.dataSourceScrapeDuration {
		for instanceName, values := range instances {
			for datasourceName, val := range values {
				labels := []string{jobName, instanceName, datasourceName}
				if metric, err := prometheus.NewConstMetric(defMetrics.dataSourceScrapeDurationDesc, prometheus.GaugeValue, val, labels...); err == nil {
					ch <- metric
				} else {
					log.WithError(err).Errorln("collectDefaultMetrics -> dataSourceScrapeDuration")
				}
			}
		}
	}
}
