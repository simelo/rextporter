package scrapper

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http/httptest"
	"regexp"
	"strings"

	"github.com/simelo/rextporter/src/client"
	"github.com/simelo/rextporter/src/util"
	log "github.com/sirupsen/logrus"
)

type metricsForwader struct {
	clientFactory client.ClientFactory
	serviceName   string
}

// MetricsForwaders have a slice of metricsForwader, that are capable of forward a metrics endpoint with service name as prefix
type MetricsForwaders struct {
	servicesMetricsForwader []metricsForwader
}

func newMetricsForwader(clc client.ProxyMetricClientCreator) metricsForwader {
	return metricsForwader{clientFactory: clc, serviceName: clc.ServiceName}
}

// NewMetricsForwaders create a scrapper that handle all the forwaded services
func NewMetricsForwaders(pmclsc []client.ProxyMetricClientCreator) Scrapper {
	scrapper := MetricsForwaders{servicesMetricsForwader: make([]metricsForwader, len(pmclsc))}
	for idxScrapper := range scrapper.servicesMetricsForwader {
		scrapper.servicesMetricsForwader[idxScrapper] = newMetricsForwader(pmclsc[idxScrapper])
	}
	return scrapper
}

func findMetricsName(metricsData string) (metricsNames []string) {
	rex := regexp.MustCompile(`# TYPE [a-zA-Z_:][a-zA-Z0-9_:]*`)
	metricsNameLines := rex.FindAllString(metricsData, -1)
	metricsNames = make([]string, len(metricsNameLines))
	for idx, metricsNameLine := range metricsNameLines {
		metricsNameLineColumns := strings.Split(metricsNameLine, " ")
		// FIXME(denissacostaq@gmail.com): be careful indexing here
		metricsNames[idx] = metricsNameLineColumns[2]
	}
	return metricsNames
}

func appendPrefixForMetrics(prefix string, metricsData string) ([]byte, error) {
	metricsName := findMetricsName(metricsData)
	for _, metricName := range metricsName {
		repl := strings.NewReplacer(
			"# HELP "+metricName+" ", "# HELP "+prefix+"_"+metricName+" ",
			"# TYPE "+metricName+" ", "# TYPE "+prefix+"_"+metricName+" ",
		)
		metricsData = repl.Replace(metricsData)
		metricsData = strings.Replace(metricsData, "\n"+metricName, "\n"+prefix+"_"+metricName, -1)
	}
	if len(metricsName) == 0 {
		err := fmt.Errorf("data from %s not appear to be from a metrics(trough prometheus instrumentation) endpoint", string(prefix))
		log.WithError(err).Errorln("append prefix error, content ignored")
	}
	return []byte(metricsData), nil
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
			if exposedMetricsData, err := cl.GetData(); err != nil {
				log.WithError(err).Error("error getting metrics from service " + mf.serviceName)
				errCause := "can not get the data"
				return data, util.ErrorFromThisScope(errCause, generalScopeErr)
			} else {
				var prefixed []byte
				if prefixed, err = appendPrefixForMetrics(mf.serviceName, string(exposedMetricsData)); err != nil {
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
