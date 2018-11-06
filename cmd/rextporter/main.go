package main

import (
	"flag"
	"os"
	"path/filepath"

	"github.com/simelo/rextporter/src/exporter"
)

func main() {
	gopath := os.Getenv("GOPATH")
	defaultConfigFilePath := filepath.Join(gopath, "", "src", "github.com", "simelo", "rextporter", "conf", "default", "skycoin.toml")
	configFile := flag.String("config", defaultConfigFilePath, "Config file path.")
	defaultListenPort := 8080
	listenPort := flag.Uint("port", uint(defaultListenPort), "Listen port.")
	flag.Parse()

	exporter.ExportMetrics(*configFile, uint16(*listenPort))
	waitForEver := make(chan bool)
	<-waitForEver
}
