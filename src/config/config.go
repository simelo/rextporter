package config

import (
	"bytes"
	"container/list"
	"fmt"
	"net/url"
	"strings"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/simelo/rextporter/src/common"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

// RootConfig is the top level node for the config tree, it has a list of metrics and a
// service from which get this metrics.
type RootConfig struct {
	Service Service  `json:"service"`
	Metrics []Metric `json:"metrics"`
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

// MetricName returns a promehteus style name for the giving metric name.
func (conf RootConfig) MetricName(metricName string) string {
	return prometheus.BuildFQName("skycoin", conf.Service.Name, metricName)
}

// NewConfigFromFilePath desserialize from the configured values in the toml file in to the config data structure
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

// FilterMetricsByType will return all the metrics who match whit the 't' parameter.
func (conf RootConfig) FilterMetricsByType(t string) (metrics []Metric) {
	tmpMetrics := list.New()
	for _, metric := range conf.Metrics {
		if strings.Compare(metric.Options.Type, t) == 0 {
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

func (conf RootConfig) validate() {
	var errs []error
	errs = append(errs, conf.Service.validate()...)
	for _, metric := range conf.Metrics {
		errs = append(errs, metric.validate()...)
	}
	for _, metric := range conf.Metrics {
		errs = append(errs, metric.validate()...)
	}
	if len(errs) != 0 {
		defer log.Panicln("some errors found")
		for _, err := range errs {
			log.WithError(err).Errorln("Error")
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
