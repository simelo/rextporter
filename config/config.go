package config

import (
	"github.com/spf13/viper"
	"log"
	"os"
)

type MetricOptions struct {
	Type string `json:"type"`
	Description string `json:"description"`
}

type Metric struct {
	Name string `json:"name"`
	Options MetricOptions `json:"options"`
}

type Link struct {
	Host string `json:"host"`
	Metric string `json:"metric"`
	URL string `json:"url"`
	Path string `json:"path,omitempty"`
}

type RootConfig struct {
	Hosts []string `json:"hosts"`
	Metrics []Metric `json:"metrics"`
	MetricsForHost []Link `json:"metricsforhost"`
}

var rootConfig RootConfig

func Config() RootConfig {
	return rootConfig
}

func init() {
	viper.SetConfigFile(os.Getenv("GOPATH") + "/src/github.com/denisacostaq/rextporter/examples/simple.toml")
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalln("Error loading config file:", err)
	}
	if err:= viper.Unmarshal(&rootConfig); err != nil {
		log.Fatalln("Error unmarshalling:", err)
	}
}
