package memconfig

import (
	"testing"

	"github.com/davecgh/go-spew/spew"
	"github.com/simelo/rextporter/src/core"
	"github.com/stretchr/testify/suite"
)

func newMetricsExtractor(suite *metricsExtractorSuit) core.RextMetricsExtractor {
	return NewMetricsExtractor(
		suite.metrics,
		suite.metricOptions,
	)
}

type metricsExtractorSuit struct {
	suite.Suite
	metricsExtractor core.RextMetricsExtractor
	metrics          []*MetricDef
	metricOptions    core.RextKeyValueStore
}

func (suite *metricsExtractorSuit) SetupTest() {
	suite.metricOptions = NewOptionsMap()
	suite.metricOptions.SetString("k1", "v1")
	suite.metricOptions.SetString("k2", "v2")
	metricName := "MySupperMetric"
	metricType := core.KeyTypeCounter
	metricDescription := "This is all about ..."
	metricLabels := []label{label{name: "ip", vPath: "/client_ip"}}
	metricOptions := NewOptionsMap()
	metricOptions.SetString("k1", "v1")
	metricOptions.SetString("k2", "v2")
	suite.metrics = []*MetricDef{
		NewMetricDef(
			metricName,
			metricType,
			metricDescription,
			suite.metricOptions,
			ls2kvs(metricLabels),
		),
	}
	suite.metricsExtractor = newMetricsExtractor(suite)
}

func TestMetricsExtractorSuit(t *testing.T) {
	suite.Run(t, new(metricsExtractorSuit))
}

func (suite *metricsExtractorSuit) TestNewMetricsExtractor() {
	// NOTE(denisacostaq@gmail.com): Giving

	// NOTE(denisacostaq@gmail.com): When
	metricsExtractorDef := newMetricsExtractor(suite)
	opts, err := suite.metricOptions.Clone()
	suite.Nil(err)
	suite.metricOptions.SetString("k1", "v2")

	// NOTE(denisacostaq@gmail.com): Assert
	spew.Dump(metricsExtractorDef.GetOptions())
	suite.True(eqKvs(suite.Assert(), suite.metricOptions, metricsExtractorDef.GetOptions()))
	suite.False(eqKvs(nil, opts, metricsExtractorDef.GetOptions()))
}
