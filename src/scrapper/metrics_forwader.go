package scrapper

import (
	"bufio"
	"bytes"
	"errors"
	"io/ioutil"
	"net/http/httptest"

	io_prometheus_client "github.com/prometheus/client_model/go"
	"github.com/prometheus/common/expfmt"
	"github.com/simelo/rextporter/src/client"
	"github.com/simelo/rextporter/src/util"
	log "github.com/sirupsen/logrus"
)

type metricsForwader struct {
	clientFactory client.Factory
	serviceName   string
	jobName       string
	instanceName  string
}

// MetricsForwaders have a slice of metricsForwader, that are capable of forward a metrics endpoint with service name as prefix
type MetricsForwaders struct {
	servicesMetricsForwader []metricsForwader
}

func newMetricsForwader(clc client.ProxyMetricClientCreator) metricsForwader {
	return metricsForwader{
		clientFactory: clc,
		serviceName:   clc.ServiceName,
		jobName:       clc.JobName,
		instanceName:  clc.InstanceName,
	}
}

// NewMetricsForwaders create a scrapper that handle all the forwaded services
func NewMetricsForwaders(pmclsc []client.ProxyMetricClientCreator) Scrapper {
	scrapper := MetricsForwaders{
		servicesMetricsForwader: make([]metricsForwader, len(pmclsc)),
	}
	for idxScrapper := range scrapper.servicesMetricsForwader {
		scrapper.servicesMetricsForwader[idxScrapper] = newMetricsForwader(pmclsc[idxScrapper])
	}
	return scrapper
}

func appendLables(metrics []byte, labels []*io_prometheus_client.LabelPair) ([]byte, error) {
	var parser expfmt.TextParser
	in := bytes.NewReader(metrics)
	metricFamilies, err := parser.TextToMetricFamilies(in)
	if err != nil {
		log.WithError(err).Errorln("error, reading text format failed")
		return metrics, err
	}
	var buff bytes.Buffer
	writer := bufio.NewWriter(&buff)
	encoder := expfmt.NewEncoder(writer, expfmt.FmtText)
	for _, mf := range metricFamilies {
		for idxMetrics := range mf.Metric {
			mf.Metric[idxMetrics].Label = append(mf.Metric[idxMetrics].Label, labels...)
		}
		err := encoder.Encode(mf)
		if err != nil {
			log.WithFields(log.Fields{"err": err, "metric family": mf}).Errorln("can not encode metric family")
			return metrics, err
		}
	}
	writer.Flush()
	return buff.Bytes(), nil
}

// GetMetric return the original metrics but with a service name as prefix in his names
func (mfs MetricsForwaders) GetMetric() (val interface{}, err error) {
	getCustomData := func() (data []byte, err error) {
		generalScopeErr := "Error getting custom data for metrics fordwader"
		recorder := httptest.NewRecorder()
		for _, mf := range mfs.servicesMetricsForwader {
			var cl client.Client
			if cl, err = mf.clientFactory.CreateClient(); err != nil {
				errCause := "can not create client"
				return data, util.ErrorFromThisScope(errCause, generalScopeErr)
			}
			var exposedMetricsData []byte
			if exposedMetricsData, err = cl.GetData(); err != nil {
				log.WithError(err).Error("error getting metrics from service " + mf.serviceName)
				errCause := "can not get the data"
				return data, util.ErrorFromThisScope(errCause, generalScopeErr)
			}
			job := "job"
			instance := "instance"
			prefixed, err := appendLables(
				exposedMetricsData,
				[]*io_prometheus_client.LabelPair{
					&io_prometheus_client.LabelPair{
						Name:  &job,
						Value: &mf.jobName,
					},
					&io_prometheus_client.LabelPair{
						Name:  &instance,
						Value: &mf.instanceName,
					},
				},
			)
			if err != nil {
				return nil, err
			}
			if count, err := recorder.Write(prefixed); err != nil || count != len(prefixed) {
				if err != nil {
					log.WithError(err).Errorln("error writing prefixed content")
				}
				if count != len(prefixed) {
					log.WithFields(log.Fields{
						"wrote":    count,
						"required": len(prefixed),
					}).Errorln("no enough content wrote")
					return nil, errors.New("no enough content wrote")
				}
			}
		}
		if data, err = ioutil.ReadAll(recorder.Body); err != nil {
			log.WithError(err).Errorln("can not read recorded custom data")
			return nil, err
		}
		return data, nil
	}
	if len(mfs.servicesMetricsForwader) == 0 {
		return nil, nil
	}
	if customData, err := getCustomData(); err == nil {
		val = customData
	} else {
		log.WithError(err).Errorln("error getting custom data")
	}
	return val, nil
}
