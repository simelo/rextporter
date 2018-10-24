package main

import (
	"github.com/denisacostaq/rextporter/config"
	"log"
	"encoding/json"
	"github.com/denisacostaq/rextporter/client"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	conf := config.Config()
	if /*b*/_, err := json.MarshalIndent(conf, "", " "); err != nil {
		log.Println("Error marshalling:", err)
	} else {
		//os.Stdout.Write(b)
	}

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
