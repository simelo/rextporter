package scrapper

import (
	"errors"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/simelo/rextporter/src/client"
	"github.com/simelo/rextporter/src/config"
	"github.com/simelo/rextporter/src/util"
)

// Scrapper get metrics from raw data
type Scrapper interface {
	// GetMetric recive the metrics collector channel and should return the metric val
	GetMetric(metricsCollector chan<- prometheus.Metric) (val interface{}, err error)
	GetJobName() string
	GetInstanceName() string
	GetDatasource() string
}

// FordwaderScrapper get metrics from an already metrics endpoint
type FordwaderScrapper interface {
	// GetMetric should return the metrics vals as raw string
	GetMetric() (val interface{}, err error)
	GetJobName() string
	GetInstanceName() string
}

type baseScrapper struct {
	jobName      string
	instanceName string
}

type baseAPIScrapper struct {
	baseScrapper
	datasource    string
	clientFactory client.Factory
	parser        BodyParser
	jsonPath      string
}

type baseFordwaderScrapper struct {
	baseScrapper
	clientFactory client.FordwaderFactory
}

func (s baseScrapper) GetJobName() string {
	return s.jobName
}

func (s baseScrapper) GetInstanceName() string {
	return s.instanceName
}

func (s baseAPIScrapper) GetDatasource() string {
	return s.datasource
}

// BodyParser decode body from different formats, an get some data node
type BodyParser interface {
	decodeBody(body []byte) (val interface{}, err error)
	pathLookup(path string, val interface{}) (node interface{}, err error)
}

// NewScrapper will put all the required info to scrap metrics from the body returned by the client.
func NewScrapper(cf client.Factory, parser BodyParser, metric config.Metric, srvConf config.Service) (Scrapper, error) {
	jobName := srvConf.JobName()
	instanceName := srvConf.InstanceName()
	datasource := metric.URL
	if len(metric.LabelNames()) > 0 {
		return createVecScrapper(cf, parser, metric, jobName, instanceName, datasource)
	}
	return createAtomicScrapper(cf, parser, metric, jobName, instanceName, datasource)
}

func createVecScrapper(cf client.Factory, parser BodyParser, metric config.Metric, jobName, instanceName, datasource string) (Scrapper, error) {
	if metric.Options.Type == config.KeyTypeCounter || metric.Options.Type == config.KeyTypeGauge {
		return newNumericVec(cf, parser, metric, jobName, instanceName, datasource), nil
	}
	return NumericVec{}, errors.New("histogram vec and summary vec are not supported yet")
}

func createAtomicScrapper(cf client.Factory, parser BodyParser, metric config.Metric, jobName, instanceName, datasource string) (Scrapper, error) {
	if metric.Options.Type == config.KeyTypeSummary {
		return Histogram{}, errors.New("summary scrapper is not supported yet")
	}
	if metric.Options.Type == config.KeyTypeHistogram {
		return newHistogram(cf, parser, metric, jobName, instanceName, datasource), nil
	}
	return newNumeric(cf, parser, metric.Path, jobName, instanceName, datasource), nil
}

func getData(cf client.Factory, p BodyParser, metricsCollector chan<- prometheus.Metric) (data interface{}, err error) {
	const generalScopeErr = "error getting data"
	var cl client.Client
	if cl, err = cf.CreateClient(); err != nil {
		errCause := "can ot create client"
		return data, util.ErrorFromThisScope(errCause, generalScopeErr)
	}
	var body []byte
	if body, err = cl.GetData(metricsCollector); err != nil {
		errCause := "client can not get data"
		return data, util.ErrorFromThisScope(errCause, generalScopeErr)
	}
	if data, err = p.decodeBody(body); err != nil {
		errCause := "scrapper can not decode the body"
		return data, util.ErrorFromThisScope(errCause, generalScopeErr)
	}
	return data, err
}
