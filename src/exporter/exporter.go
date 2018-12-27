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
	"github.com/simelo/rextporter/src/cache"
	"github.com/simelo/rextporter/src/core"
	"github.com/simelo/rextporter/src/scrapper"
	"github.com/simelo/rextporter/src/util"
	"github.com/simelo/rextporter/src/util/metrics"
	log "github.com/sirupsen/logrus"
)

func exposedMetricsMiddleware(scrappers []scrapper.FordwaderScrapper, promHandler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if len(listenAddr) == 0 {
			listenAddr = r.Host
		}
		getDefaultData := func() (data []byte, err error) {
			generalScopeErr := "error reding default data"
			recorder := httptest.NewRecorder()
			promHandler.ServeHTTP(recorder, r)
			var reader io.ReadCloser
			switch recorder.Header().Get("Content-Encoding") {
			case "gzip":
				reader, err = gzip.NewReader(recorder.Body)
				if err != nil {
					errCause := fmt.Sprintln("can not create gzip reader.", err.Error())
					return nil, util.ErrorFromThisScope(errCause, generalScopeErr)
				}
			default:
				reader = ioutil.NopCloser(bytes.NewReader(recorder.Body.Bytes()))
			}
			defer reader.Close()
			if data, err = ioutil.ReadAll(reader); err != nil {
				log.WithError(err).Errorln("can not read recorded default data")
				return nil, err
			}
			return data, nil
		}
		var allData []byte
		if defaultData, err := getDefaultData(); err != nil {
			log.WithError(err).Errorln("error getting default data")
		} else {
			allData = append(allData, defaultData...)
		}
		for _, fs := range scrappers {
			var iMetrics interface{}
			var err error
			if iMetrics, err = fs.GetMetric(); err != nil {
				log.WithError(err).Errorln("error scrapping fordwader metrics")
			} else {
				customData, okCustomData := iMetrics.([]byte)
				if okCustomData {
					allData = append(allData, customData...)
				} else {
					log.WithError(err).Errorln("error asserting fordwader metrics data as []byte")
				}
			}
		}
		w.Header().Set("Content-Type", "text/plain")
		if count, err := w.Write(allData); err != nil || count != len(allData) {
			if err != nil {
				log.WithError(err).Errorln("error writing data")
			}
			if count != len(allData) {
				log.WithFields(log.Fields{
					"wrote":    count,
					"required": len(allData),
				}).Errorln("no enough content wrote")
			}
		}
	})
}

var fDefMetrics *metrics.DefaultFordwaderMetrics

// MustExportMetrics will read the config from mainConfigFile if any or use a default one.
func MustExportMetrics(handlerEndpoint string, listenPort uint16, conf core.RextRoot) (srv *http.Server) {
	c := cache.NewCache()
	if collector, err := newMetricsCollector(c, conf); err != nil {
		log.WithError(err).Panicln("Can not create metrics")
	} else {
		prometheus.MustRegister(collector)
		fDefMetrics = metrics.NewDefaultFordwaderMetrics()
		fDefMetrics.MustRegister()
	}
	metricsForwaders, err := createMetricsForwaders(conf, fDefMetrics)
	if err != nil {
		log.WithError(err).Panicln("Can not create forward_metrics metrics")
	}
	var listenAddrPort string
	if len(listenAddr) == 0 {
		listenAddrPort = fmt.Sprintf(":%d", listenPort)
	} else {
		listenAddrPort = fmt.Sprintf("%s:%d", listenAddr, listenPort)
		listenAddr = listenAddrPort
	}
	srv = &http.Server{Addr: listenAddrPort}
	http.Handle(
		handlerEndpoint,
		gziphandler.GzipHandler(exposedMetricsMiddleware(listenAddr, metricsForwaders, promhttp.Handler())))
	go func() {
		log.Infoln(fmt.Sprintf("Starting server in port %d, path %s ...", listenPort, handlerEndpoint))
		log.WithError(srv.ListenAndServe()).Errorln("unable to start the server")
	}()
	return srv
}

// TODO(denisacostaq@gmail.com): you can use a NewProcessCollector, NewGoProcessCollector, make a blockchain collector sense?
