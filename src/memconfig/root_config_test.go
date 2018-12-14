package memconfig

import (
	"testing"

	"github.com/simelo/rextporter/src/core"
	"github.com/stretchr/testify/suite"
)

func newRootConfig(suite *rootConfigSuite) core.RextEnv {
	return NewRootConfig(
		suite.options,
	)
}

type rootConfigSuite struct {
	suite.Suite
	options    core.RextKeyValueStore
	rootConfig core.RextEnv
}

func (suite *rootConfigSuite) SetupTest() {
	suite.options = NewOptionsMap()
	suite.options.SetString("k1", "v1")
	suite.options.SetString("k2", "v2")
	suite.rootConfig = newRootConfig(suite)
}

func TestRootConfig(t *testing.T) {
	suite.Run(t, new(rootConfigSuite))
}

func (suite *rootConfigSuite) TestNewMetricDef() {
	// NOTE(denisacostaq@gmail.com): Giving

	// NOTE(denisacostaq@gmail.com): When
	rootConfig := newRootConfig(suite)
	opts, err := suite.options.Clone()
	suite.Nil(err)
	suite.options.SetString("k1", "v2")

	// NOTE(denisacostaq@gmail.com): Assert
	suite.True(eqKvs(suite.Assert(), suite.options, rootConfig.GetOptions()))
	suite.False(eqKvs(nil, opts, rootConfig.GetOptions()))
}
