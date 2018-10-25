package config

import (
	"github.com/spf13/viper"
	"strconv"
	"strings"
	"fmt"
	"errors"
	"bytes"
	"github.com/denisacostaq/rextporter/src/common"
)

type Host struct {
	Ref string
	Location          string `json:"location"`
	Port             int    `json:"port"`
	AuthType         string `json:"auth_type"`
	TokenHeaderKey         string `json:"token_header_key"`
	GenTokenEndpoint string `json:"gen_token_endpoint"`
	TokenKeyFromEndpoint string `json:"token_key_from_endpoint"`
}

// isValidUrl tests a string to determine if it is a url or not.
func isValidUrl(toTest string) bool {
	if _, err := url.ParseRequestURI(toTest); err != nil {
		return false
	}
	return true
}

func (host Host) validate() (errs []error) {
	if len(host.Ref) == 0 {
		errs = append(errs, errors.New("ref is required in host"))
	}
	if len(host.Location) == 0 {
		errs = append(errs, errors.New("location is required in host"))
	}
	if !isValidUrl(host.Location) {
		errs = append(errs, errors.New("location is not a valid url in host"))
	}
	if !isValidUrl(host.UriToGetToken()) {
		errs = append(errs, errors.New("location + port can not form a valid uri in host"))
	}
	if host.Port < 0 || host.Port > 65535 {
		errs = append(errs, errors.New("port number should be between 0 and 65535 in host"))
	}
	if strings.Compare(host.AuthType, "CSRF") == 0 && len(host.TokenHeaderKey) == 0 {
		errs = append(errs, errors.New("TokenHeaderKey is required if you are using CSRF"))
	}
	if strings.Compare(host.AuthType, "CSRF") == 0 && len(host.GenTokenEndpoint) == 0 {
		errs = append(errs, errors.New("GenTokenEndpoint is required if you are using CSRF"))
	}
	if strings.Compare(host.AuthType, "CSRF") == 0 && len(host.TokenKeyFromEndpoint) == 0 {
		errs = append(errs, errors.New("TokenKeyFromEndpoint is required if you are using CSRF"))
	}
	return errs
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
	MetricRef string `json:"metric_ref"`
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

func NewConfig(strConf string) (error) {
	const generalScopeErr = "error creating a config instance"
	viper.SetConfigType("toml")
	buff := bytes.NewBuffer([]byte(strConf))
	if err := viper.ReadConfig(buff); err != nil {
		errCause := fmt.Sprintln("can not read the buffer", err.Error())
		return common.ErrorFromThisScope(errCause, generalScopeErr)
	}
	rootConfig = RootConfig{}
	if err := viper.Unmarshal(&rootConfig); err != nil {
		errCause := fmt.Sprintln("can not parse config data", err.Error())
		return common.ErrorFromThisScope(errCause, generalScopeErr)
	}
	return nil
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
	//// FIXME(denisacostaq@gmail.com): not portable
	//viper.SetConfigFile(os.Getenv("GOPATH") + "/src/github.com/denisacostaq/rextporter/examples/simple2.toml")
	//if err := viper.ReadInConfig(); err != nil {
	//	log.Fatalln("Error loading config file:", err)
	//}
	//if err:= viper.Unmarshal(&rootConfig); err != nil {
	//	log.Fatalln("Error unmarshalling:", err)
	//}
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