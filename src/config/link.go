package config

import (
	"errors"
	"fmt"

	"github.com/simelo/rextporter/src/common"
)

// Link is a concept who map properties of a Metric in a Host, for example, you can define
// some hosts some metrics and in Link your specific the properties of a giving metric in
// a giving host, for example, the Url and the json path(Path) from where you can scrap the information.
type Link struct {
	HostRef    string `json:"host_ref"`
	MetricRef  string `json:"metric_ref"`
	URL        string `json:"url"`
	HTTPMethod string `json:"http_method"`
	Path       string `json:"path,omitempty"`
}

// MetricName will return a name for a metric in a host
func (link Link) MetricName() string {
	return link.HostRef + "_" + link.MetricRef
}

// MetricDescription will look for the metric associated trough ref and return the description
func (link Link) MetricDescription() (string, error) {
	const generalScopeErr = "error getting metric description"
	var metric Metric
	var err error
	if metric, err = Config().findMetricByRef(link.MetricRef); err != nil {
		errCause := fmt.Sprintln("can not find the metric: ", err.Error())
		return "", common.ErrorFromThisScope(errCause, generalScopeErr)
	}
	return metric.Options.Description, err
}

// FindMetricType will return the metric type through the metric related with ref
func (link Link) FindMetricType() (metricType string, err error) {
	const generalScopeErr = "error looking for metric type"
	var metric Metric
	if metric, err = link.FindMetric(); err != nil {
		errCause := fmt.Sprintln("can not find metric by ref: ", link.MetricRef, err.Error())
		return metricType, common.ErrorFromThisScope(errCause, generalScopeErr)
	}
	metricType = metric.Options.Type
	return metricType, err
}

// FindMetric will return the metric through the metric related with ref
func (link Link) FindMetric() (metric Metric, err error) {
	const generalScopeErr = "error looking for metric"
	if metric, err = Config().findMetricByRef(link.MetricRef); err != nil {
		errCause := fmt.Sprintln("can not find metric by ref: ", link.MetricRef, err.Error())
		return metric, common.ErrorFromThisScope(errCause, generalScopeErr)
	}
	return metric, err
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
	if len(link.HTTPMethod) == 0 {
		errs = append(errs, errors.New("HttpMethod is required in Link(metric fo host)"))
	}
	if len(link.Path) == 0 {
		errs = append(errs, errors.New("path is required in Link(metric fo host)"))
	}
	host, hostNotFound := Config().FindHostByRef(link.HostRef)
	if hostNotFound != nil {
		errs = append(errs, hostNotFound)
	} else {
		if !isValidURL(host.URIToGetMetric(link)) {
			errs = append(errs, errors.New("can not create a valid uri under link"))
		}
		errs = append(errs, host.validate()...)
	}
	metric, metricNotFound := Config().findMetricByRef(link.MetricRef)
	if metricNotFound != nil {
		errs = append(errs, metricNotFound)
	} else {
		errs = append(errs, metric.validate()...)
	}
	return errs
}
