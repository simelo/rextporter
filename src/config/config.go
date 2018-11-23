package config

import (
	"bytes"
	"container/list"
	"fmt"
	"net/url"

	"github.com/simelo/rextporter/src/util"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

// RootConfig is the top level node for the config tree, it has a list of metrics and a
// service from which get this metrics.
type RootConfig struct {
	Services []Service `json:"services"`
}

// NewConfigFromRawString allow you to define a `.toml` config in the fly, a raw string with the "config content"
func NewConfigFromRawString(strConf string) (conf RootConfig, err error) {
	const generalScopeErr = "error creating a config instance"
	viper.SetConfigType("toml")
	buff := bytes.NewBuffer([]byte(strConf))
	if err = viper.ReadConfig(buff); err != nil {
		errCause := fmt.Sprintln("can not read the buffer: ", err.Error())
		return conf, util.ErrorFromThisScope(errCause, generalScopeErr)
	}
	if err = viper.Unmarshal(&conf); err != nil {
		errCause := fmt.Sprintln("can not decode the config data: ", err.Error())
		return conf, util.ErrorFromThisScope(errCause, generalScopeErr)
	}
	conf.validate()
	return conf, err
}

// newMetricsConfig desserialize a metrics config from the 'toml' file path
func newMetricsConfig(path string) (metricsConf []Metric, err error) {
	const generalScopeErr = "error reading metrics config"
	if len(path) == 0 {
		errCause := "path should not be null"
		return metricsConf, util.ErrorFromThisScope(errCause, generalScopeErr)
	}
	viper.SetConfigFile(path)
	if err := viper.ReadInConfig(); err != nil {
		errCause := fmt.Sprintln("error reading config file: ", path, err.Error())
		return metricsConf, util.ErrorFromThisScope(errCause, generalScopeErr)
	}
	type metricsForService struct {
		Metrics []Metric
	}
	var root metricsForService
	if err := viper.Unmarshal(&root); err != nil {
		errCause := fmt.Sprintln("can not decode the config data: ", err.Error())
		return metricsConf, util.ErrorFromThisScope(errCause, generalScopeErr)
	}
	metricsConf = root.Metrics
	return metricsConf, nil
}

// newServicesConfigFromFile desserialize a service config from the 'toml' file path
func newServicesConfigFromFile(path string, conf mainConfigData) (servicesConf []Service, err error) {
	const generalScopeErr = "error reading service config"
	servicesConfReader := NewServicesConfigFromFile(path)
	if servicesConf, err = servicesConfReader.GetConfig(); err != nil {
		errCause := "error reading service config"
		return servicesConf, util.ErrorFromThisScope(errCause, generalScopeErr)
	}
	for idxService, service := range servicesConf {
		if util.StrSliceContains(service.Modes, ServiceTypeAPIRest) {
			if servicesConf[idxService].Metrics, err = newMetricsConfig(conf.MetricsConfigPath(service.Name)); err != nil {
				errCause := "error reading metrics config: " + err.Error()
				panic(util.ErrorFromThisScope(errCause, generalScopeErr))
			}
		}
	}
	return servicesConf, err
}

// NewConfigFromFileSystem will read the config from the file system, you should send the
// metric config file path and service config file path into metricsPath, servicePath respectively.
// This function can cause a panic.
func NewConfigFromFileSystem(mainConfigPath string) (rootConf RootConfig) {
	const generalScopeErr = "error getting config values from file system"
	var conf mainConfigData
	var err error
	if conf, err = newMainConfigData(mainConfigPath); err != nil {
		errCause := "error reading metrics config: " + err.Error()
		panic(errCause)
	}
	if rootConf.Services, err = newServicesConfigFromFile(conf.ServicesConfigPath(), conf); err != nil {
		errCause := "root cause: " + err.Error()
		panic(util.ErrorFromThisScope(errCause, generalScopeErr))
	}
	rootConf.validate()
	return rootConf
}

// FilterMetricsByType will return all the metrics who match with the 't' parameter.
func (conf RootConfig) FilterMetricsByType(t string) (metrics []Metric) {
	tmpMetrics := list.New()
	for _, service := range conf.Services {
		metricsForService := service.FilterMetricsByType(t)
		for _, metric := range metricsForService {
			tmpMetrics.PushBack(metric)
		}
	}
	metrics = make([]Metric, tmpMetrics.Len())
	idxLink := 0
	for it := tmpMetrics.Front(); it != nil; it = it.Next() {
		metrics[idxLink] = it.Value.(Metric)
		idxLink++
	}
	return metrics
}

// FilterServicesByType will return all the services who match with the 't' parameter.
func (conf RootConfig) FilterServicesByType(t string) (services []Service) {
	return filterServicesByType(t, conf.Services)
}

func filterServicesByType(t string, services []Service) (filteredService []Service) {
	tmpServices := list.New()
	for _, service := range services {
		if util.StrSliceContains(service.Modes, t) {
			tmpServices.PushBack(service)
		}
	}
	filteredService = make([]Service, tmpServices.Len())
	idxLink := 0
	for it := tmpServices.Front(); it != nil; it = it.Next() {
		filteredService[idxLink] = it.Value.(Service)
		idxLink++
	}
	return filteredService
}

func (conf RootConfig) validate() {
	var errs []error
	for _, service := range conf.Services {
		errs = append(errs, service.validate()...)
	}
	if len(errs) != 0 {
		for _, err := range errs {
			log.WithError(err).Errorln("Error")
		}
		log.Panicln("some errors found")
	}
}

// isValidUrl tests a string to determine if it is a valid URL or not.
func isValidURL(toTest string) bool {
	if _, err := url.ParseRequestURI(toTest); err != nil {
		return false
	}
	return true
}
