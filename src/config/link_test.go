package config

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

type LinkConfigSuit struct {
	suite.Suite
	LinkConf Link
}

func (suite *LinkConfigSuit) SetupTest() {
	hostConf := Host{
		Ref:                  "MySupperServer",
		Location:             "https://my.supper.server",
		Port:                 8080,
		AuthType:             "CSRF",
		TokenHeaderKey:       "X-CSRF-Token",
		GenTokenEndpoint:     "/api/v1/csrf",
		TokenKeyFromEndpoint: "csrf_token",
	}
	var hostsConf []Host
	hostsConf = append(hostsConf, hostConf)
	metricOptionsConf := MetricOptions{
		Type:        KeyTypeGauge,
		Description: "This is a very cool metric",
	}
	metricConf := Metric{
		Name:    "my_cool_metric",
		Options: metricOptionsConf,
	}
	var metricsConf []Metric
	metricsConf = append(metricsConf, metricConf)
	suite.LinkConf = Link{
		HostRef:    hostConf.Ref,
		MetricRef:  metricConf.Name,
		URL:        "/api/v1.0/status",
		HTTPMethod: "GET",
		Path:       "/load/connected_users",
	}
	var linksConf []Link
	linksConf = append(linksConf, suite.LinkConf)
	rootConfig = RootConfig{
		Hosts:          hostsConf,
		Metrics:        metricsConf,
		MetricsForHost: linksConf,
	}
}

func TestLinkConfigSuit(t *testing.T) {
	suite.Run(t, new(LinkConfigSuit))
}

func (suite *LinkConfigSuit) TestEnsureDefaultSuitHostConfIsValid() {
	// NOTE(denisacostaq@gmail.com): Giving
	// default
	linkConf := suite.LinkConf

	// NOTE(denisacostaq@gmail.com): When
	// test start

	// NOTE(denisacostaq@gmail.com): Assert
	suite.Len(linkConf.validate(), 0)
}

func (suite *LinkConfigSuit) TestNotEmptyHostRef() {
	// NOTE(denisacostaq@gmail.com): Giving
	var linkConf = suite.LinkConf
	linkConf.HostRef = string("")

	// NOTE(denisacostaq@gmail.com): When

	// NOTE(denisacostaq@gmail.com): Assert
	suite.Len(linkConf.validate(), 2) // required, can not find a host
}

func (suite *LinkConfigSuit) TestNotEmptyMetricRef() {
	// NOTE(denisacostaq@gmail.com): Giving
	var linkConf = suite.LinkConf
	linkConf.MetricRef = string("")

	// NOTE(denisacostaq@gmail.com): When

	// NOTE(denisacostaq@gmail.com): Assert
	suite.Len(linkConf.validate(), 2) // required, can not find a metric
}

func (suite *LinkConfigSuit) TestNotEmptyURL() {
	// NOTE(denisacostaq@gmail.com): Giving
	var linkConf = suite.LinkConf
	linkConf.URL = string("")

	// NOTE(denisacostaq@gmail.com): When

	// NOTE(denisacostaq@gmail.com): Assert
	suite.Len(linkConf.validate(), 1)
}

func (suite *LinkConfigSuit) TestNotEmptyHttpMethod() {
	// NOTE(denisacostaq@gmail.com): Giving
	var linkConf = suite.LinkConf
	linkConf.URL = string("")

	// NOTE(denisacostaq@gmail.com): When

	// NOTE(denisacostaq@gmail.com): Assert
	suite.Len(linkConf.validate(), 1)
}

func (suite *LinkConfigSuit) TestNotEmptyPath() {
	// NOTE(denisacostaq@gmail.com): Giving
	var linkConf = suite.LinkConf
	linkConf.Path = string("")

	// NOTE(denisacostaq@gmail.com): When

	// NOTE(denisacostaq@gmail.com): Assert
	suite.Len(linkConf.validate(), 1)
}

func (suite *LinkConfigSuit) TestCanNotFindHostByRef() {
	// NOTE(denisacostaq@gmail.com): Giving
	var linkConf = suite.LinkConf
	linkConf.HostRef = linkConf.HostRef + linkConf.HostRef

	// NOTE(denisacostaq@gmail.com): When

	// NOTE(denisacostaq@gmail.com): Assert
	suite.Len(linkConf.validate(), 1)
}

func (suite *LinkConfigSuit) TestCanNotFindMetricByRef() {
	// NOTE(denisacostaq@gmail.com): Giving
	var linkConf = suite.LinkConf
	linkConf.MetricRef = linkConf.MetricRef + linkConf.MetricRef

	// NOTE(denisacostaq@gmail.com): When

	// NOTE(denisacostaq@gmail.com): Assert
	suite.Len(linkConf.validate(), 1)
}
