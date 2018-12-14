package memconfig

import (
	"testing"

	"github.com/simelo/rextporter/src/core"
	"github.com/stretchr/testify/suite"
)

func newScrapperConf(suite *scrapperConfSuit) core.RextServiceScraper {
	return NewScrapperConf(
		suite.metricOptions,
	)
}

type scrapperConfSuit struct {
	suite.Suite
	scrapperConf  core.RextServiceScraper
	metricOptions core.RextKeyValueStore
}

func (suite *scrapperConfSuit) SetupTest() {
	suite.metricOptions = NewOptionsMap()
	suite.metricOptions.SetString("k1", "v1")
	suite.metricOptions.SetString("k2", "v2")
	suite.scrapperConf = newScrapperConf(suite)
}

func TestScrapperConfSuit(t *testing.T) {
	suite.Run(t, new(scrapperConfSuit))
}

func (suite *scrapperConfSuit) TestNewScrapperConf() {
	// NOTE(denisacostaq@gmail.com): Giving

	// NOTE(denisacostaq@gmail.com): When
	metricDef := newScrapperConf(suite)
	opts, err := suite.metricOptions.Clone()
	suite.Nil(err)
	suite.metricOptions.SetString("k1", "v2")

	// NOTE(denisacostaq@gmail.com): Assert
	suite.True(eqKvs(suite.Assert(), suite.metricOptions, metricDef.GetOptions()))
	suite.False(eqKvs(nil, opts, metricDef.GetOptions()))
}
