package scrapper

import (
	"bufio"
	"bytes"
	"errors"
	"io/ioutil"
	"net/http/httptest"
	"time"

	io_prometheus_client "github.com/prometheus/client_model/go"
	"github.com/prometheus/common/expfmt"
	"github.com/simelo/rextporter/src/client"
	"github.com/simelo/rextporter/src/util"
	"github.com/simelo/rextporter/src/util/metrics"
	log "github.com/sirupsen/logrus"
)

// MetricsForwader is a scrapper kind capable to forward a metrics endpoint with job and instance labels at least
type MetricsForwader struct {
	baseFordwaderScrapper
	defFordwaderMetrics *metrics.DefaultFordwaderMetrics
}

// GetJobName return the name of the job(service)
func (scrapper MetricsForwader) GetJobName() string {
	return scrapper.jobName
}

// GetInstanceName return the name of the instance(ip:port)
func (scrapper MetricsForwader) GetInstanceName() string {
	return scrapper.instanceName
}

// NewMetricsForwader create a scrapper that handle the forwaded metrics
func NewMetricsForwader(pmcls client.ProxyMetricClientCreator, fDefMetrics *metrics.DefaultFordwaderMetrics) FordwaderScrapper {
	return MetricsForwader{
		baseFordwaderScrapper: baseFordwaderScrapper{
			baseScrapper: baseScrapper{
				jobName:      pmcls.JobName,
				instanceName: pmcls.InstanceName,
			},
			clientFactory: pmcls,
		},
		defFordwaderMetrics: fDefMetrics,
	}
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
func (scrapper MetricsForwader) GetMetric() (val interface{}, err error) {
	getCustomData := func() (data []byte, err error) {
		successResponse := false
		defer func(startTime time.Time) {
			duration := time.Since(startTime).Seconds()
			labels := []string{scrapper.GetJobName(), scrapper.GetInstanceName()}
			if successResponse {
				scrapper.defFordwaderMetrics.FordwaderScrapeDurationSeconds.WithLabelValues(labels...).Set(duration)
			}
		}(time.Now().UTC())
		generalScopeErr := "Error getting custom data for metrics fordwader"
		recorder := httptest.NewRecorder()
		var cl client.FordwaderClient
		if cl, err = scrapper.clientFactory.CreateClient(); err != nil {
			errCause := "can not create client"
			return data, util.ErrorFromThisScope(errCause, generalScopeErr)
		}
		var exposedMetricsData []byte
		if exposedMetricsData, err = cl.GetData(); err != nil {
			log.WithError(err).Error("error getting metrics from service " + scrapper.GetJobName())
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
					Value: &scrapper.jobName,
				},
				&io_prometheus_client.LabelPair{
					Name:  &instance,
					Value: &scrapper.instanceName,
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
		if data, err = ioutil.ReadAll(recorder.Body); err != nil {
			log.WithError(err).Errorln("can not read recorded custom data")
			return nil, err
		}
		successResponse = true
		return data, nil
	}
	if customData, err := getCustomData(); err == nil {
		val = customData
	} else {
		log.WithError(err).Errorln("error getting custom data")
	}
	return val, nil
}
