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

func createMetricsForwaders(conf config.RootConfig) (scrapper.Scrapper, error) {
	generalScopeErr := "can not create metrics Middleware"
	services := conf.FilterServicesByType(config.ServiceTypeProxy)
	proxyMetricClientCreators := make([]client.ProxyMetricClientCreator, len(services))
	for idxService := range services {
		var err error
		if proxyMetricClientCreators[idxService], err = client.CreateProxyMetricClientCreator(services[idxService]); err != nil {
			errCause := fmt.Sprintln("error creating metric client: ", err.Error())
			return nil, util.ErrorFromThisScope(errCause, generalScopeErr)
		}
	}
	return scrapper.NewMetricsForwaders(proxyMetricClientCreators), nil
}

// constMetric has a scrapper to get remote data, can be a previously cached content
type constMetric struct {
	kind             string
	scrapper         scrapper.APIRestScrapper
	lastSuccessValue interface{}
	metricDesc       *prometheus.Desc
	statusDesc       *prometheus.Desc
}

type endpointData2MetricsConsumer map[string][]constMetric

func createMetrics(cache cache.Cache, srvsConf []config.Service) (metrics endpointData2MetricsConsumer, err error) {
	generalScopeErr := "can not create metrics"
	metrics = make(endpointData2MetricsConsumer)
	for _, srvConf := range srvsConf {
		for _, mConf := range srvConf.Metrics {
			k := srvConf.URIToGetMetric(mConf)
			var m constMetric
			if m, err = createConstMetric(cache, mConf, srvConf); err != nil {
				errCause := fmt.Sprintln(fmt.Sprintf("error creating metric client for %s metric of kind %s. ", mConf.Name, mConf.Options.Type), err.Error())
				return metrics, util.ErrorFromThisScope(errCause, generalScopeErr)
			}
			metrics[k] = append(metrics[k], m)
		}
	}
	return metrics, err
}

func newDefaultMetrics() *defaultMetrics {
	labels := []string{"job", "instance"}
	scrapeDurationSecondsDesc := prometheus.NewDesc("scrape_duration_seconds", "Scrape duration in seconds", labels, nil)
	scrapeSamplesScrapedDesc := prometheus.NewDesc("scrape_samples_scraped", "The number of samples the target exposed", labels, nil)
	return &defaultMetrics{
		scrapesDuration:           newScrapDuration(),
		scrapeDurationSecondsDesc: scrapeDurationSecondsDesc,
		scrapeSamplesScraped:      newScrapDuration(),
		scrapeSamplesScrapedDesc:  scrapeSamplesScrapedDesc,
	}
}

func createConstMetric(cache cache.Cache, metricConf config.Metric, srvConf config.Service) (metric constMetric, err error) {
	generalScopeErr := "can not create metric " + metricConf.Name
	var ccf client.CacheableFactory
	if ccf, err = client.CreateAPIRestCreator(metricConf, srvConf); err != nil {
		errCause := fmt.Sprintln("error creating metric client: ", err.Error())
		return metric, util.ErrorFromThisScope(errCause, generalScopeErr)
	}
	cc := client.CatcherCreator{Cache: cache, ClientFactory: ccf}
	var numScrapper scrapper.APIRestScrapper
	if numScrapper, err = scrapper.NewScrapper(cc, scrapper.JSONParser{}, metricConf, srvConf); err != nil {
		errCause := fmt.Sprintln("error creating metric client: ", err.Error())
		return metric, util.ErrorFromThisScope(errCause, generalScopeErr)
	}
	labels := metricConf.LabelNames()
	metric = constMetric{
		kind:     metricConf.Options.Type,
		scrapper: numScrapper,
		// FIXME(denisacostaq@gmail.com): if you use a duplicated name can panic?
		metricDesc: prometheus.NewDesc(srvConf.MetricName(metricConf.Name), metricConf.Options.Description, labels, nil),
		statusDesc: prometheus.NewDesc(srvConf.MetricName(metricConf.Name)+"_up", "Says if the same name metric("+srvConf.MetricName(metricConf.Name)+") was success updated, 1 for ok, 0 for failed.", nil, nil),
	}
	return metric, err
}
