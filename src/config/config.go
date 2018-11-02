package config

import (
	"bytes"
	"container/list"
	"errors"
	"fmt"
	"log"
	"net/url"
	"strings"

	"github.com/simelo/rextporter/src/common"
	"github.com/spf13/viper"
)

// RootConfig is the top level node for the config tree, it has a list of hosts, a list of metrics
// and a list of links(MetricsForHost, says how a metric is mapped in a host).
type RootConfig struct {
	Hosts          []Host   `json:"hosts"`
	Metrics        []Metric `json:"metrics"`
	MetricsForHost []Link   `json:"metrics_for_host"`
}

var rootConfig RootConfig

// Config TODO(denisacostaq@gmail.com): make a singleton
func Config() RootConfig {
	//if b, err := json.MarshalIndent(rootConfig, "", " "); err != nil {
	//	log.Println("Error marshaling:", err)
	//} else {
	//	os.Stdout.Write(b)
	//	log.Println("\n\n\n\n\n")
	//}
	// TODO(denisacostaq@gmail.com): Make it a singleton
	return rootConfig
}

// NewConfigFromRawString allow you to define a `.toml` config in the fly, a raw string with the "config content"
func NewConfigFromRawString(strConf string) error {
	const generalScopeErr = "error creating a config instance"
	viper.SetConfigType("toml")
	buff := bytes.NewBuffer([]byte(strConf))
	if err := viper.ReadConfig(buff); err != nil {
		errCause := fmt.Sprintln("can not read the buffer: ", err.Error())
		return common.ErrorFromThisScope(errCause, generalScopeErr)
	}
	rootConfig = RootConfig{}
	if err := viper.Unmarshal(&rootConfig); err != nil {
		errCause := fmt.Sprintln("can not decode the config data: ", err.Error())
		return common.ErrorFromThisScope(errCause, generalScopeErr)
	}
	rootConfig.validate()
	return nil
}

// NewConfigFromFilePath TODO(denisacostaq@gmail.com): Fill some data structures for efficient lookup from ref to host for example
func NewConfigFromFilePath(path string) error {
	const generalScopeErr = "error creating a config instance"
	viper.SetConfigFile(path)
	if err := viper.ReadInConfig(); err != nil {
		errCause := fmt.Sprintln("error reading config file: ", path, err.Error())
		return common.ErrorFromThisScope(errCause, generalScopeErr)
	}
	if err := viper.Unmarshal(&rootConfig); err != nil {
		errCause := fmt.Sprintln("can not decode the config data: ", err.Error())
		return common.ErrorFromThisScope(errCause, generalScopeErr)
	}
	rootConfig.validate()
	return nil
}

// FindHostByRef will return a host where you can match the host.Ref with the ref parameter
// or an error if not found.
func (conf RootConfig) FindHostByRef(ref string) (host Host, err error) {
	found := false
	for _, host = range conf.Hosts {
		found = strings.Compare(host.Ref, ref) == 0
		if found {
			return
		}
	}
	if !found {
		errCause := fmt.Sprintln("can not find a host for Ref: ", ref)
		err = errors.New(errCause)
	}
	return Host{}, err
}

// FilterLinksByHost will return all links where you can match the host.Ref with link.HostRef
func (conf RootConfig) FilterLinksByHost(host Host) []Link {
	var links []Link
	for _, link := range conf.MetricsForHost {
		if strings.Compare(host.Ref, link.HostRef) == 0 {
			links = append(links, link)
		}
	}
	return links
}

// FilterLinksByMetricType will return all the links who match whit the type parameter.
func FilterLinksByMetricType(links []Link, t string) ([]Link, error) {
	const generalScopeErr = "error filtering links by metric type"
	tmpLinks := list.New()
	for _, link := range links {
		mt, err := link.FindMetricType()
		if err != nil {
			errCause := fmt.Sprintln("can not find metric type: ", err.Error())
			return []Link{}, common.ErrorFromThisScope(errCause, generalScopeErr)
		}
		if strings.Compare(t, mt) == 0 {
			tmpLinks.PushBack(link)
		}
	}
	retLinks := make([]Link, tmpLinks.Len())
	idxLink := 0
	for it := tmpLinks.Front(); it != nil; it = it.Next() {
		retLinks[idxLink] = it.Value.(Link)
		idxLink++
	}
	return retLinks, nil
}

// findMetricByRef will return a metric where you can match the metric.Ref with the ref parameter
// or an error if not found.
func (conf RootConfig) findMetricByRef(ref string) (metric Metric, err error) {
	found := false
	for _, metric = range conf.Metrics {
		found = strings.Compare(metric.Name, ref) == 0
		if found {
			return
		}
	}
	if !found {
		errCause := fmt.Sprintln("can not find a metric for Ref: ", ref)
		err = errors.New(errCause)
	}
	return Metric{}, err
}

func (conf RootConfig) validate() {
	var errs []error
	for _, host := range conf.Hosts {
		errs = append(errs, host.validate()...)
	}
	for _, metric := range conf.Metrics {
		errs = append(errs, metric.validate()...)
	}
	for _, mHost := range conf.MetricsForHost {
		errs = append(errs, mHost.validate()...)
	}
	if len(errs) != 0 {
		defer log.Panicln("some errors found")
		for _, err := range errs {
			log.Println(err.Error())
		}
	}
}

// isValidUrl tests a string to determine if it is a valid URL or not.
func isValidURL(toTest string) bool {
	if _, err := url.ParseRequestURI(toTest); err != nil {
		return false
	}
	return true
}
