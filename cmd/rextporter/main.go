package main

import (
	"flag"

	"github.com/simelo/rextporter/src/config"
	"github.com/simelo/rextporter/src/exporter"
)

func main() {
	mainConfigFile := flag.String("config", "", "Metrics main config file path.")
	defaultListenPort := 8080
	listenPort := flag.Uint("port", uint(defaultListenPort), "Listen port.")
	defaultHandlerEndpint := "/metrics"
	handlerEndpint := flag.String("handler", defaultHandlerEndpint, "Handler endpoint.")
	flag.Parse()
	conf := config.NewConfigFromFileSystem(*mainConfigFile)
	exporter.ExportMetrics(handlerEndpint, uint16(*listenPort), conf)
	waitForEver := make(chan bool)
	<-waitForEver
}
