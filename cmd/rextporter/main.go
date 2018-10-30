package main

import (
	"flag"
	"os"

	"github.com/simelo/rextporter/src/exporter"
)

func main() {
	gopath := os.Getenv("GOPATH")
	defaultConfigFilePath := gopath + "/src/github.com/simelo/rextporter/examples/simple.toml"
	configFile := flag.String("config", defaultConfigFilePath, "Config file path.")
	flag.Parse()
	exporter.ExportMetrics(*configFile, 8080)
}
