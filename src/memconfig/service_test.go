package memconfig

import (
	"testing"

	"github.com/simelo/rextporter/src/core"
	"github.com/stretchr/testify/suite"
)

type serviceConfSuit struct {
	suite.Suite
	srvConf    core.RextDataSource
	baseURL    string
	location   string
	httpMethod string
}

func newService(suite *serviceConfSuit) core.RextDataSource {
	return NewServiceConf(suite.baseURL, suite.location, suite.httpMethod)
}

func (suite *serviceConfSuit) SetupTest() {
	suite.baseURL = "/hosted_in"
	suite.location = "http://localhost:9000"
	suite.httpMethod = "GET"
	suite.srvConf = newService(suite)
}

func TestServiceConfSuit(t *testing.T) {
	suite.Run(t, new(serviceConfSuit))
}

func (suite *serviceConfSuit) TestNewMetricDef() {
	// NOTE(denisacostaq@gmail.com): Giving

	// NOTE(denisacostaq@gmail.com): When
	service := newService(suite)

	// NOTE(denisacostaq@gmail.com): Assert
	suite.Equal(suite.location, service.GetResourceLocation())
	suite.Equal(suite.httpMethod, service.GetMethod())
}

func (suite *serviceConfSuit) TestAbleToSetResourceLocation() {
	// NOTE(denisacostaq@gmail.com): Giving
	loc := "fgfg78"
	suite.srvConf.SetResourceLocation(loc)

	// NOTE(denisacostaq@gmail.com): When
	loc2 := suite.srvConf.GetResourceLocation()

	// NOTE(denisacostaq@gmail.com): Assert
	suite.Equal(loc, loc2)
}

func (suite *serviceConfSuit) TestAbleToSetMethod() {
	// NOTE(denisacostaq@gmail.com): Giving
	method := "fgfg78"
	suite.srvConf.SetMethod(method)

	// NOTE(denisacostaq@gmail.com): When
	method2 := suite.srvConf.GetMethod()

	// NOTE(denisacostaq@gmail.com): Assert
	suite.Equal(method, method2)
}
