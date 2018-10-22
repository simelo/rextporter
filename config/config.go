package config

import (
	"github.com/spf13/viper"
	"log"
	"os"
)

type Host struct {
	Ref string
	Location          string `json:"location"`
	Port             int    `json:"port"`
	AuthType         string `json:"auth_type"`
	TokenKey         string `json:"token_key"`
	GenTokenEndpoint string `json:"gen_token_endpoint"`
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
