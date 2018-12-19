package memconfig

import (
	"testing"

	"github.com/simelo/rextporter/src/core"
	"github.com/stretchr/testify/suite"
)

func newRootConfig(suite *rootConfigSuite) core.RextRoot {
	return NewRootConfig(suite.services)
}

type rootConfigSuite struct {
	suite.Suite
	services   []core.RextServiceDef
	rootConfig core.RextRoot
}

func (suite *rootConfigSuite) SetupTest() {
	suite.services = []core.RextServiceDef{}
	suite.rootConfig = newRootConfig(suite)
}

func TestRootConfig(t *testing.T) {
	suite.Run(t, new(rootConfigSuite))
}

func (suite *rootConfigSuite) TestNewRootConf() {
	// NOTE(denisacostaq@gmail.com): Giving

	// NOTE(denisacostaq@gmail.com): When
	rootConfig := newRootConfig(suite)
	services := rootConfig.GetServices()

	// NOTE(denisacostaq@gmail.com): Assert
	suite.Equal(len(suite.services), len(services))
}

func (suite *rootConfigSuite) TestAbleToAddService() {
	// NOTE(denisacostaq@gmail.com): Giving
	orgServices := suite.rootConfig.GetServices()
	service := &Service{}
	suite.rootConfig.AddService(service)

	// NOTE(denisacostaq@gmail.com): When
	services2 := suite.rootConfig.GetServices()

	// NOTE(denisacostaq@gmail.com): Assert
	suite.Equal(len(orgServices)+1, len(services2))
}
