package exporter

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"

	"github.com/NYTimes/gziphandler"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/simelo/rextporter/src/common"
	"github.com/simelo/rextporter/src/config"
	log "github.com/sirupsen/logrus"
)

func findMetricsName(metricsData []byte) (metricsName [][]byte) {
	lines := bytes.Split(metricsData, []byte("\n"))
	var typeLines []byte
	for _, line := range lines {
		if bytes.HasPrefix(line, []byte("# TYPE ")) {
			typeLines = append(typeLines, line...)
		}
	}
	splittedTypeLines := bytes.Split(typeLines, []byte("# TYPE "))
	// NOTE(denisacostaq@gmail.com): remove the first empty val a the left of "# TYPE " with the 1:
	for _, splittedType := range splittedTypeLines[1:] {
		splittedTypeColumns := bytes.Split(splittedType, []byte(" "))
		metricsName = append(metricsName, splittedTypeColumns[0])
	}
	metricsName = metricsName[1:]
	return metricsName
}

func appendPrefixForMetrics(prefix []byte, metricsData []byte) (prefixedMetricsData []byte) {
	metricsName := findMetricsName(metricsData)
	prefixedMetricsData = make([]byte, len(metricsData))
	copy(prefixedMetricsData, metricsData)
	for _, metricName := range metricsName {
		newName := append(prefix, []byte("_")...)
		newName = append(newName, metricName...)
		prefixedMetricsData = bytes.Replace(prefixedMetricsData, append(metricName, []byte(" ")...), append(newName, []byte(" ")...), -1)
	}
	return prefixedMetricsData
}

func exposedMetricsMidleware(metricsMidleware []MetricMidleware, promHandler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		getCustomData := func() ([]byte, error) {
			recorder := httptest.NewRecorder()
			for _, cl := range metricsMidleware {
				if exposedMetricsData, err := cl.client.GetExposedMetrics(); err != nil {
					log.WithError(err).Error("error getting metrics from service " + cl.client.Name)
				} else {
					prefixed := appendPrefixForMetrics([]byte(cl.client.Name), exposedMetricsData)
					var count int
					if count, err = recorder.Write(prefixed); err != nil || count != len(prefixed) {
						if err != nil {
							log.WithError(err).Errorln("error writing prefixed content")
						}
						if count != len(prefixed) {
							log.WithFields(log.Fields{
								"wrote":    count,
								"required": len(prefixed),
							}).Errorln("no enough content wrote")
						}
					}
				}
			}
			if data, err := ioutil.ReadAll(recorder.Body); err != nil {
				log.WithError(err).Errorln("can not read recorded custom data")
				return nil, err
			} else {
				return data, nil
			}
		}
		getDefaultData := func() (data []byte, err error) {
			generalScopeErr := "error reding default data"
			recorder := httptest.NewRecorder()
			promHandler.ServeHTTP(recorder, r)
			var reader io.ReadCloser
			// BUG(denisacostaq@gmail.com): close this reader.
			// defer reader.Close()
			switch recorder.Header().Get("Content-Encoding") {
			case "gzip":
				reader, err = gzip.NewReader(recorder.Body)
				if err != nil {
					errCause := fmt.Sprintln("can not create gzip reader.", err.Error())
					return nil, common.ErrorFromThisScope(errCause, generalScopeErr)
				}
			default:
				reader = ioutil.NopCloser(bytes.NewReader(recorder.Body.Bytes()))
			}
			if data, err = ioutil.ReadAll(reader); err != nil {
				log.WithError(err).Errorln("can not read recorded default data")
				return nil, err
			} else {
				return data, nil
			}
		}
		var allData []byte
		if defaultData, err := getDefaultData(); err != nil {
			log.WithError(err).Errorln("error getting default data")
		} else {
			allData = append(allData, defaultData...)
		}
		if customData, err := getCustomData(); err != nil {
			log.WithError(err).Errorln("error getting custom data")
		} else {
			allData = append(allData, customData...)
		}
		w.Header().Set("Content-Type", "text/plain")
		if allData == nil {
			allData = []byte("a")
		}
		w.Write(allData)
	})
}

// ExportMetrics will read the config from mainConfigFile if any or use a default one.
func ExportMetrics(mainConfigFile, handlerEndpint string, listenPort uint16) (srv *http.Server) {
	config.NewConfigFromFileSystem(mainConfigFile)
	if collector, err := newSkycoinCollector(); err != nil {
		log.WithError(err).Panicln("Can not create metrics")
	} else {
		prometheus.MustRegister(collector)
	}
	metricsMidleware, err := createMetricsMidleware()
	if err != nil {
		log.WithError(err).Panicln("Can not create proxy metrics")
	}
	port := fmt.Sprintf(":%d", listenPort)
	srv = &http.Server{Addr: port}
	hdl := promhttp.InstrumentMetricHandler(
		prometheus.DefaultRegisterer,
		promhttp.HandlerFor(prometheus.DefaultGatherer, promhttp.HandlerOpts{DisableCompression: false}),
	)
	http.Handle(handlerEndpint, gziphandler.GzipHandler(exposedMetricsMidleware(metricsMidleware, hdl)))
	go func() {
		log.Infoln(fmt.Sprintf("Starting server in port %d, path %s ...", listenPort, handlerEndpint))
		log.WithError(srv.ListenAndServe()).Errorln("unable to start the server")
	}()
	return srv
}

// TODO(denisacostaq@gmail.com): you can use a NewProcessCollector, NewGoProcessCollector, make a blockchain collector sense?
