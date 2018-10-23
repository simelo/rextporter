package main

import (
	"github.com/denisacostaq/rextporter/config"
	"log"
	"encoding/json"
	"github.com/denisacostaq/rextporter/client"
	"strings"
)

func filterLinksByHost(host config.Host, conf config.RootConfig) []config.Link {
	var links []config.Link
	for _,link := range conf.MetricsForHost {
		if strings.Compare(host.Ref, link.HostRef) == 0 {
			links = append(links, link)
		}
	}
	return links
}

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
		links := filterLinksByHost(host, conf)
		for _, link := range links {
			if cl, err := client.NewMetricClient(link); err != nil {
				log.Println(err.Error())
			} else {
				log.Println(cl.GetMetric())
			}
		}
	}
}
