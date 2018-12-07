package exporter

import (
	"fmt"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/simelo/rextporter/src/cache"
	"github.com/simelo/rextporter/src/client"
	"github.com/simelo/rextporter/src/config"
	"github.com/simelo/rextporter/src/scrapper"
	"github.com/simelo/rextporter/src/util"
)

func createMetricsForwaders(conf config.RootConfig) (fordwaderScrappers []scrapper.Scrapper, err error) {
	generalScopeErr := "can not create metrics Middleware"
	services := conf.FilterServicesByType(config.ServiceTypeProxy)
	fordwaderScrappers = make([]scrapper.Scrapper, len(services))
	for idxService := range services {
		var metricFordwaderCreator client.ProxyMetricClientCreator
		if metricFordwaderCreator, err = client.CreateProxyMetricClientCreator(services[idxService]); err != nil {
			errCause := fmt.Sprintln("error creating metric client: ", err.Error())
			return nil, util.ErrorFromThisScope(errCause, generalScopeErr)
		}
		fordwaderScrappers[idxService] = scrapper.NewMetricsForwader(metricFordwaderCreator)
	}
	return fordwaderScrappers, nil
}

// constMetric has a scrapper to get remote data, can be a previously cached content
type constMetric struct {
	kind       string
	scrapper   scrapper.Scrapper
	metricDesc *prometheus.Desc
}

type endpointData2MetricsConsumer map[string][]constMetric

func createMetrics(cache cache.Cache, srvsConf []config.Service, datasourceResponseDurationDesc *prometheus.Desc) (metrics endpointData2MetricsConsumer, err error) {
	generalScopeErr := "can not create metrics"
	metrics = make(endpointData2MetricsConsumer)
	for _, srvConf := range srvsConf {
		for _, mConf := range srvConf.Metrics {
			k := srvConf.URIToGetMetric(mConf)
			var m constMetric
			if m, err = createConstMetric(cache, mConf, srvConf, datasourceResponseDurationDesc); err != nil {
				errCause := fmt.Sprintln(fmt.Sprintf("error creating metric client for %s metric of kind %s. ", mConf.Name, mConf.Options.Type), err.Error())
				return metrics, util.ErrorFromThisScope(errCause, generalScopeErr)
			}
			metrics[k] = append(metrics[k], m)
		}
	}
	return metrics, err
}

var instance4JobLabels = []string{"job", "instance"}

func newDefaultMetrics() *defaultMetrics {
	scrapeDurationSecondsDesc := prometheus.NewDesc(
		"scrape_duration_seconds",
		"Scrape duration in seconds",
		instance4JobLabels,
		nil,
	)
	scrapeSamplesScrapedDesc := prometheus.NewDesc(
		"scrape_samples_scraped",
		"The number of samples the target exposed",
		instance4JobLabels,
		nil,
	)
	datasourceResponseDurationDesc := prometheus.NewDesc(
		"datasource_response_duration",
		"Elapse time to get a response from a datasource",
		append(instance4JobLabels, "datasource"),
		nil,
	)
	return &defaultMetrics{
		scrapedDurations:               newScrapDuration(),
		scrapeDurationSecondsDesc:      scrapeDurationSecondsDesc,
		scrapedSamples:                 newScrapDuration(),
		scrapeSamplesScrapedDesc:       scrapeSamplesScrapedDesc,
		datasourceResponseDurationDesc: datasourceResponseDurationDesc,
	}
}

func createConstMetric(cache cache.Cache, metricConf config.Metric, srvConf config.Service, datasourceResponseDurationDesc *prometheus.Desc) (metric constMetric, err error) {
	generalScopeErr := "can not create metric " + metricConf.Name
	var ccf client.CacheableFactory
	if ccf, err = client.CreateAPIRestCreator(metricConf, srvConf, datasourceResponseDurationDesc); err != nil {
		errCause := fmt.Sprintln("error creating metric client: ", err.Error())
		return metric, util.ErrorFromThisScope(errCause, generalScopeErr)
	}
	cc := client.CatcherCreator{Cache: cache, ClientFactory: ccf}
	var numScrapper scrapper.Scrapper
	if numScrapper, err = scrapper.NewScrapper(cc, scrapper.JSONParser{}, metricConf, srvConf); err != nil {
		errCause := fmt.Sprintln("error creating metric client: ", err.Error())
		return metric, util.ErrorFromThisScope(errCause, generalScopeErr)
	}
	labels := append(metricConf.LabelNames(), instance4JobLabels...)
	metric = constMetric{
		kind:     metricConf.Options.Type,
		scrapper: numScrapper,
		// FIXME(denisacostaq@gmail.com): if you use a duplicated name can panic?
		metricDesc: prometheus.NewDesc(metricConf.Name, metricConf.Options.Description, labels, nil),
	}
	return metric, err
}
