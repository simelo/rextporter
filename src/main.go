package main

import (
	"github.com/denisacostaq/rextporter/src/config"
	"log"
	"encoding/json"
	"github.com/denisacostaq/rextporter/src/client"
	"os"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	// FIXME(denisacostaq@gmail.com): not portable
	config.NewConfigFromFilePath(os.Getenv("GOPATH") + "/src/github.com/denisacostaq/rextporter/examples/simple.toml")
	conf := config.Config()

	for _, host := range conf.Hosts {
		// cl, err := client.NewTokenClient(host)
		// log.Println("tk:", tk)
		links := conf.FilterLinksByHost(host)
		for _, link := range links {
			if cl, err := client.NewMetricClient(link); err != nil {
				log.Println(err.Error())
			} else {
				log.Println(cl.GetMetric())
			}
		}
	}
}
