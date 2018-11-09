package exporter

import (
	"fmt"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/simelo/rextporter/src/config"
	log "github.com/sirupsen/logrus"
)

func exposedMetricsMidleware(metricsMidleware []MetricMidleware, promHandler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		for _, cl := range metricsMidleware {
			log.Println("working with this client", cl)
		}
		promHandler.ServeHTTP(w, r)
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
	http.Handle(handlerEndpint, exposedMetricsMidleware(metricsMidleware, promhttp.Handler()))
	go func() {
		log.Infoln(fmt.Sprintf("Starting server in port %d, path %s ...", listenPort, handlerEndpint))
		log.WithError(srv.ListenAndServe()).Errorln("unable to start the server")
	}()
	return srv
}

// TODO(denisacostaq@gmail.com): you can use a NewProcessCollector, NewGoProcessCollector, make a blockchain collector sense?
