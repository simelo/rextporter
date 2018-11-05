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
	flag.Parse()
	exporter.ExportMetrics(*configFile, 8080)
	waitForEver := make(chan bool)
	<-waitForEver
}
