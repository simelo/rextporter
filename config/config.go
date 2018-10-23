package config

import (
	"github.com/spf13/viper"
	"log"
	"os"
	"strconv"
	"strings"
	"fmt"
	"errors"
)

type Host struct {
	Ref string
	Location          string `json:"location"`
	Port             int    `json:"port"`
	AuthType         string `json:"auth_type"`
	TokenKey         string `json:"token_key"`
	GenTokenEndpoint string `json:"gen_token_endpoint"`
	GenTokenKey string `json:"gen_token_key"`
}

type MetricOptions struct {
	Type string `json:"type"`
	Description string `json:"description"`
}

type Metric struct {
	Name string `json:"name"`
	Options MetricOptions `json:"options"`
}

type Link struct {
	HostRef string `json:"host_ref"`
	Metric string `json:"metric"`
	URL string `json:"url"`
	HttpMethod string `json:"http_method"`
	Path string `json:"path,omitempty"`
}

type RootConfig struct {
	Hosts []Host `json:"hosts"`
	Metrics []Metric `json:"metrics"`
	MetricsForHost []Link `json:"metrics_for_host"`
}

var rootConfig RootConfig

func Config() RootConfig {
	// TODO(denisacostaq@gmail.com): Make it a singleton
	return rootConfig
}

func (conf RootConfig) FindHostByRef(ref string) (host Host, err error) {
	found := false
	for _, host = range conf.Hosts {
		found = strings.Compare(host.Ref, ref) == 0
		if found {
			return
		}
	}
	if !found {
		errCause := fmt.Sprintln("can not find a host for Ref:", ref)
		err = errors.New(errCause)
	}
	return Host{}, err
}

func (host Host) UriToGetMetric(metricInHost Link) string {
	return host.Location + ":" + strconv.Itoa(host.Port) + metricInHost.URL
}

func (host Host) UriToGetToken() string {
	return host.Location + ":" + strconv.Itoa(host.Port) + host.GenTokenEndpoint
}

// TODO(denisacostaq@gmail.com): Fill some data structures for efficient lookup from ref to host for example
func init() {
	// FIXME(denisacostaq@gmail.com): not portable
	viper.SetConfigFile(os.Getenv("GOPATH") + "/src/github.com/denisacostaq/rextporter/examples/simple2.toml")
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalln("Error loading config file:", err)
	}
	if err:= viper.Unmarshal(&rootConfig); err != nil {
		log.Fatalln("Error unmarshalling:", err)
	}
}

func (conf RootConfig) FilterLinksByHost(host Host) []Link {
	var links []Link
	for _,link := range conf.MetricsForHost {
		if strings.Compare(host.Ref, link.HostRef) == 0 {
			links = append(links, link)
		}
	}
	return links
}