package main

import (
	"flag"

	"github.com/simelo/rextporter/src/exporter"
)

func main() {
	mainConfigFile := flag.String("config-path", "", "Metrics main config file path.")
	defaultListenPort := 8080
	listenPort := flag.Uint("port", uint(defaultListenPort), "Listen port.")
	flag.Parse()

	exporter.ExportMetrics(*mainConfigFile, uint16(*listenPort))
	waitForEver := make(chan bool)
	<-waitForEver
}
