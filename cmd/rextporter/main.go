package main

import (
	"flag"

	"github.com/simelo/rextporter/src/exporter"
)

func main() {
	mainConfigFile := flag.String("config", "", "Metrics main config file path.")
	defaultListenPort := 8080
	listenPort := flag.Uint("port", uint(defaultListenPort), "Listen port.")
	defaultHandlerEndpint := "/metrics"
	handlerEndpint := flag.String("handler", defaultHandlerEndpint, "Handler endpoint.")
	flag.Parse()
	exporter.ExportMetrics(*mainConfigFile, *handlerEndpint, uint16(*listenPort))
	waitForEver := make(chan bool)
	<-waitForEver
}
