package main

import (
	"log"
	"os"

	"github.com/simelo/rextporter/src/client"
	"github.com/simelo/rextporter/src/config"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	// FIXME(denisacostaq@gmail.com): not portable
	if err := config.NewConfigFromFilePath(os.Getenv("GOPATH") + "/src/github.com/denisacostaq/rextporter/examples/simple.toml"); err != nil {
		log.Panicln(err)
	}
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
