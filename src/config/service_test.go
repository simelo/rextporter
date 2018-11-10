package config

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

type serviceConfSuite struct {
	suite.Suite
	ServiceConf Service
}

func (suite *serviceConfSuite) SetupTest() {
	suite.ServiceConf = Service{
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
	}
}

func TestServiceConfSuite(t *testing.T) {
	suite.Run(t, new(serviceConfSuite))
}

func (suite *serviceConfSuite) TestEnsureDefaultSuiteServiceConfIsValid() {
	// NOTE(denisacostaq@gmail.com): Giving
	// default
	serviceConf := suite.ServiceConf

	// NOTE(denisacostaq@gmail.com): When
	// test start

	// NOTE(denisacostaq@gmail.com): Assert
	suite.Len(serviceConf.validate(), 0)
}

func (suite *serviceConfSuite) TestNotEmptyName() {
	// NOTE(denisacostaq@gmail.com): Giving
	var serviceConf = suite.ServiceConf
	serviceConf.Name = string("")

	// NOTE(denisacostaq@gmail.com): When

	// NOTE(denisacostaq@gmail.com): Assert
	suite.Len(serviceConf.validate(), 1)
}

func (suite *serviceConfSuite) TestValidateLocation() {
	// NOTE(denisacostaq@gmail.com): Giving
	var serviceConf = suite.ServiceConf
	serviceConf.Location.Location = string("")

	// NOTE(denisacostaq@gmail.com): When

	// NOTE(denisacostaq@gmail.com): Assert
	suite.NotEmpty(serviceConf.validate()) // empty, invalid and url + port invalid
}

// TODO(denisacostaq@gmail.com): validate port if change type from uint16 to ...

func (suite *serviceConfSuite) TestCsrfAuthButEmptyTokenKeyFromEndpoint() {
	// NOTE(denisacostaq@gmail.com): Giving
	var serviceConf = suite.ServiceConf
	serviceConf.AuthType = "CSRF"
	serviceConf.TokenKeyFromEndpoint = ""

	// NOTE(denisacostaq@gmail.com): When

	// NOTE(denisacostaq@gmail.com): Assert
	suite.Len(serviceConf.validate(), 1)
}

func (suite *serviceConfSuite) TestCsrfAuthButEmptyTokenHeaderKey() {
	// NOTE(denisacostaq@gmail.com): Giving
	var serviceConf = suite.ServiceConf
	serviceConf.AuthType = "CSRF"
	serviceConf.TokenHeaderKey = ""

	// NOTE(denisacostaq@gmail.com): When

	// NOTE(denisacostaq@gmail.com): Assert
	suite.Len(serviceConf.validate(), 1)
}

func (suite *serviceConfSuite) TestCsrfAuthButEmptyGenTokenEndpoint() {
	// NOTE(denisacostaq@gmail.com): Giving
	var serviceConf = suite.ServiceConf
	serviceConf.AuthType = "CSRF"
	serviceConf.GenTokenEndpoint = ""

	// NOTE(denisacostaq@gmail.com): When

	// NOTE(denisacostaq@gmail.com): Assert
	suite.Len(serviceConf.validate(), 1)
}
