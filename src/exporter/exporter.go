package exporter

import (
	"fmt"
	"log"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/simelo/rextporter/src/config"
)

// ExportMetrics will read the config file from the CLI parammeter `-config` if any
// or use a default one.
func ExportMetrics(configFile string, listenPort uint16) (srv *http.Server) {
	if err := config.NewConfigFromFilePath(configFile); err != nil {
		log.Fatalln("can not open the config file", err.Error())
	}
	if collector, err := newSkycoinCollector(); err != nil {
		log.Panicln("Can not create metrics:", err)
	} else {
		prometheus.MustRegister(collector)
	}
	port := fmt.Sprintf(":%d", listenPort)
	srv = &http.Server{Addr: port}
	http.Handle("/metrics", promhttp.Handler())
	go func() {
		log.Panicln(srv.ListenAndServe())
	}()
	return srv
}

// TODO(denisacostaq@gmail.com): you can use a NewProcessCollector, NewGoProcessCollector, make a blockchain collector sense?
