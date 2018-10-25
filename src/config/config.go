package config

import (
	"github.com/spf13/viper"
	"strconv"
	"strings"
	"fmt"
	"errors"
	"bytes"
	"github.com/denisacostaq/rextporter/src/common"
	"log"
	"net/url"
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

func (mo MetricOptions) validate() (errs []error) {
	if len(mo.Type) == 0 {
		errs = append(errs, errors.New("type is required in metric"))
	}
	return errs
}

type Metric struct {
	Name string `json:"name"`
	Options MetricOptions `json:"options"`
}

func (metric Metric) validate() (errs []error) {
	if len(metric.Name) == 0 {
		errs = append(errs, errors.New("name is required in metric"))
	}
	errs = append(errs, metric.Options.validate()...)
	return errs
}

type Link struct {
	HostRef string `json:"host_ref"`
	MetricRef string `json:"metric_ref"`
	URL string `json:"url"`
	HttpMethod string `json:"http_method"`
	Path string `json:"path,omitempty"`
}

func (link Link) validate() (errs []error) {
	if len(link.HostRef) == 0 {
		errs = append(errs, errors.New("HostRef is required in Link(metric fo host)"))
	}
	if len(link.MetricRef) == 0 {
		errs = append(errs, errors.New("HostRef is required in Link(metric fo host)"))
	}
	if len(link.URL) == 0 {
		errs = append(errs, errors.New("url is required"))
	}
	if len(link.HttpMethod) == 0 {
		errs = append(errs, errors.New("HttpMethod is required in Link(metric fo host)"))
	}
	if len(link.Path) == 0 {
		errs = append(errs, errors.New("path is required in Link(metric fo host)"))
	}
	host, hostNotFound := Config().FindHostByRef(link.HostRef)
	if hostNotFound != nil {
		errs = append(errs, hostNotFound)
	} else {
		if !isValidUrl(host.UriToGetMetric(link)) {
			errs = append(errs, errors.New("can not create a valid uri under link"))
		}
		errs = append(errs, host.validate()...)
	}
	metric, metricNotFound := Config().FindMetricByRef(link.MetricRef)
	if metricNotFound != nil {
		errs = append(errs, metricNotFound)
	} else {
		errs = append(errs, metric.validate()...)
	}
	return errs
}

type RootConfig struct {
	Hosts []Host `json:"hosts"`
	Metrics []Metric `json:"metrics"`
	MetricsForHost []Link `json:"metrics_for_host"`
}

var rootConfig RootConfig

func (conf RootConfig) validate() {
	var errs []error
	for _, host := range conf.Hosts {
		errs = append(errs, host.validate()...)
	}
	for _, metric := range conf.Metrics {
		errs = append(errs, metric.validate()...)
	}
	for _, mhost := range conf.MetricsForHost {
		errs = append(errs, mhost.validate()...)
	}
	if len(errs) != 0 {
		log.Println("some errors found")
		for _, err := range errs {
			log.Println(err.Error())
		}
		log.Panicln()
	}
}

func Config() RootConfig {
	//if b, err := json.MarshalIndent(rootConfig, "", " "); err != nil {
	//	log.Println("Error marshalling:", err)
	//} else {
	//	os.Stdout.Write(b)
	//	log.Println("\n\n\n\n\n")
	//}
	// TODO(denisacostaq@gmail.com): Make it a singleton
	return rootConfig
}

func NewConfigFromRawString(strConf string) (error) {
	const generalScopeErr = "error creating a config instance"
	viper.SetConfigType("toml")
	buff := bytes.NewBuffer([]byte(strConf))
	if err := viper.ReadConfig(buff); err != nil {
		errCause := fmt.Sprintln("can not read the buffer", err.Error())
		return common.ErrorFromThisScope(errCause, generalScopeErr)
	}
	rootConfig = RootConfig{}
	if err := viper.Unmarshal(&rootConfig); err != nil {
		errCause := fmt.Sprintln("can not decode the config data", err.Error())
		return common.ErrorFromThisScope(errCause, generalScopeErr)
	}
	rootConfig.validate()
	return nil
}

// TODO(denisacostaq@gmail.com): Fill some data structures for efficient lookup from ref to host for example
func NewConfigFromFilePath(path string) error {
	const generalScopeErr = "error creating a config instance"
	viper.SetConfigFile(path)
	if err := viper.ReadInConfig(); err != nil {
		errCause := fmt.Sprintln("error reading config file:", path, err.Error())
		return common.ErrorFromThisScope(errCause, generalScopeErr)
	}
	if err:= viper.Unmarshal(&rootConfig); err != nil {
		errCause := fmt.Sprintln("can not decode the config data", err.Error())
		return common.ErrorFromThisScope(errCause, generalScopeErr)
	}
	rootConfig.validate()
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

func (conf RootConfig) FindMetricByRef(ref string) (metric Metric, err error) {
	found := false
	for _, metric = range conf.Metrics {
		found = strings.Compare(metric.Name, ref) == 0
		if found {
			return
		}
	}
	if !found {
		errCause := fmt.Sprintln("can not find a host for Ref:", ref)
		err = errors.New(errCause)
	}
	return Metric{}, err
}

func (host Host) UriToGetMetric(metricInHost Link) string {
	return host.Location + ":" + strconv.Itoa(host.Port) + metricInHost.URL
}

func (host Host) UriToGetToken() string {
	return host.Location + ":" + strconv.Itoa(host.Port) + host.TokenKeyFromEndpoint
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