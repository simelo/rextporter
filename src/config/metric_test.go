package config

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

type metricConfSuit struct {
	suite.Suite
	MetricConf *Metric
}

func (suite *metricConfSuit) SetupTest() {
	rootConfig.Services = []Service{
		Service{
			Name:                 "MySupperServer",
			Mode:                 "apiRest",
			Scheme:               "http",
			Location:             Server{Location: "http://localhost:8080"},
			Port:                 8080,
			BasePath:             "/skycoin/node",
			AuthType:             "CSRF",
			TokenHeaderKey:       "X-CSRF-Token",
			GenTokenEndpoint:     "/api/v1/csrf",
			TokenKeyFromEndpoint: "csrf_token",
		},
	}
	rootConfig.Metrics = []Metric{
		Metric{
			Name:             "MySupperMetric",
			URL:              "/api/v1/health",
			HTTPMethod:       "GET",
			Path:             "/blockchain/head/seq",
			Options:          MetricOptions{Type: KeyTypeCounter, Description: "It is all about ..."},
			HistogramOptions: HistogramOptions{},
		}}
	suite.MetricConf = &(rootConfig.Metrics[0])
}

func TestMetricConfSuit(t *testing.T) {
	suite.Run(t, new(metricConfSuit))
}

func (suite *metricConfSuit) TestEnsureDefaultSuitMetricConfIsValid() {
	// NOTE(denisacostaq@gmail.com): Giving
	// default
	metricConf := suite.MetricConf

	// NOTE(denisacostaq@gmail.com): When
	// test start

	// NOTE(denisacostaq@gmail.com): Assert
	suite.Len(metricConf.validate(), 0)
}

func (suite *metricConfSuit) TestNotEmptyName() {
	// NOTE(denisacostaq@gmail.com): Giving
	var metricConf = suite.MetricConf
	metricConf.Name = string("")

	// NOTE(denisacostaq@gmail.com): When

	// NOTE(denisacostaq@gmail.com): Assert
	suite.Len(metricConf.validate(), 1)
}

func (suite *metricConfSuit) TestNotEmptyURL() {
	// NOTE(denisacostaq@gmail.com): Giving
	var metricConf = suite.MetricConf
	metricConf.URL = string("")

	// NOTE(denisacostaq@gmail.com): When

	// NOTE(denisacostaq@gmail.com): Assert
	suite.Len(metricConf.validate(), 1)
}

func (suite *metricConfSuit) TestNotEmptyHTTPMethod() {
	// NOTE(denisacostaq@gmail.com): Giving
	var metricConf = suite.MetricConf
	metricConf.HTTPMethod = string("")

	// NOTE(denisacostaq@gmail.com): When

	// NOTE(denisacostaq@gmail.com): Assert
	suite.Len(metricConf.validate(), 1)
}

func (suite *metricConfSuit) TestNotEmptyPath() {
	// NOTE(denisacostaq@gmail.com): Giving
	var metricConf = suite.MetricConf
	metricConf.Path = string("")

	// NOTE(denisacostaq@gmail.com): When

	// NOTE(denisacostaq@gmail.com): Assert
	suite.Len(metricConf.validate(), 1)
}

// TODO(denisacostaq@gmail.com): test define buckets but declare type counter for example
