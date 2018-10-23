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
		if tk, err := client.GetToken(host); err != nil {
			log.Println("Can not get the token:", err)
		} else {
			log.Println("tk:", tk)
			links := filterLinksByHost(host, conf)
			for _, link := range links {
				client.GetMetric(host, link, tk)
				log.Println(client.GetMetric(host, link, tk))
			}
		}
	}
}
