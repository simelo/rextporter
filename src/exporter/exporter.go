package exporter

import (
	"fmt"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/simelo/rextporter/src/config"
	log "github.com/sirupsen/logrus"
)

// ExportMetrics will read the config file from the CLI parammeter `-config` if any
// or use a default one.
func ExportMetrics(configFile string, listenPort uint16) (srv *http.Server) {
	if err := config.NewConfigFromFilePath(configFile); err != nil {
		log.WithError(err).Fatalln("can not open the config file")
	}
	if collector, err := newSkycoinCollector(); err != nil {
		log.WithError(err).Panicln("Can not create metrics")
	} else {
		prometheus.MustRegister(collector)
	}
	port := fmt.Sprintf(":%d", listenPort)
	srv = &http.Server{Addr: port}
	http.Handle("/metrics", promhttp.Handler())
	go func() {
		log.Infoln(fmt.Sprintf("Starting server in port %d, path /metrics ...", listenPort))
		log.WithError(srv.ListenAndServe()).Panicln("unable to start the server")
	}()
	return srv
}

// TODO(denisacostaq@gmail.com): you can use a NewProcessCollector, NewGoProcessCollector, make a blockchain collector sense?
