package config

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

type metricConfSuit struct {
	suite.Suite
	MetricConf Metric
}

func (suite *metricConfSuit) SetupTest() {
	suite.MetricConf = Metric{
		Name:    "MySupperMetric",
		Options: MetricOptions{Type: KeyTypeCounter, Description: "It is all about ..."},
	}
}

func TestMetricConfSuit(t *testing.T) {
	suite.Run(t, new(metricConfSuit))
}

func (suite *metricConfSuit) TestEnsureDefaultSuitHostConfIsValid() {
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
	var hostConf = suite.MetricConf
	hostConf.Name = string("")

	// NOTE(denisacostaq@gmail.com): When

	// NOTE(denisacostaq@gmail.com): Assert
	suite.Len(hostConf.validate(), 1)
}
