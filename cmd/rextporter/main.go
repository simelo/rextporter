package main

import (
	"flag"
	"os"
	"path/filepath"

	"github.com/simelo/rextporter/src/exporter"
)

func main() {
	gopath := os.Getenv("GOPATH")
	defaultMetricsConfFilePath := filepath.Join(gopath, "", "src", "github.com", "simelo", "rextporter", "conf", "default", "skycoin", "metrics.toml")
	metricsConfFile := flag.String("metrics-config-path", defaultMetricsConfFilePath, "Metrics config file path.")
	defaultServiceConfFilePath := filepath.Join(gopath, "", "src", "github.com", "simelo", "rextporter", "conf", "default", "skycoin", "service.toml")
	serviceConfFile := flag.String("service-config-path", defaultServiceConfFilePath, "Service config file path.")
	defaultListenPort := 8080
	listenPort := flag.Uint("port", uint(defaultListenPort), "Listen port.")
	flag.Parse()

	exporter.ExportMetrics(*metricsConfFile, *serviceConfFile, uint16(*listenPort))
	waitForEver := make(chan bool)
	<-waitForEver
}
