package exporter

import (
	"bytes"
	"fmt"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
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
		prefixedMetricsData = bytes.Replace(prefixedMetricsData, metricName, newName, -1)
	}
	return prefixedMetricsData
}

func exposedMetricsMidleware(metricsMidleware []MetricMidleware, promHandler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// TODO(denisacostaq@gmail.com): Track all the data and write a header with Content-Length compressed
		promHandler.ServeHTTP(w, r)
		for _, cl := range metricsMidleware {
			if exposedMetricsData, err := cl.client.GetExposedMetrics(); err != nil {
				log.WithError(err).Error("error getting metrics from service " + cl.client.Name)
			} else {
				prefixed := appendPrefixForMetrics([]byte(cl.client.Name), exposedMetricsData)
				w.Write(prefixed)
			}
		}
  // TODO(denisacostaq@gmail.com): compre all the content and use the promhttp.Handler() who wrte compressed content
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
		promhttp.HandlerFor(prometheus.DefaultGatherer, promhttp.HandlerOpts{DisableCompression: true}),
	)
	http.Handle(handlerEndpint, exposedMetricsMidleware(metricsMidleware, hdl))
	go func() {
		log.Infoln(fmt.Sprintf("Starting server in port %d, path %s ...", listenPort, handlerEndpint))
		log.WithError(srv.ListenAndServe()).Errorln("unable to start the server")
	}()
	return srv
}

// TODO(denisacostaq@gmail.com): you can use a NewProcessCollector, NewGoProcessCollector, make a blockchain collector sense?
