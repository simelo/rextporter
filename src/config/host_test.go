package config

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

type hostConfSuit struct {
	suite.Suite
	HostConf Host
}

func (suite *hostConfSuit) SetupTest() {
	suite.HostConf = Host{
		Ref:                  "MySupperServer",
		Location:             "https://my.supper.server",
		Port:                 8080,
		AuthType:             "CSRF",
		TokenHeaderKey:       "X-CSRF-Token",
		GenTokenEndpoint:     "/api/v1/csrf",
		TokenKeyFromEndpoint: "csrf_token",
	}
}

func TesthostConfSuit(t *testing.T) {
	suite.Run(t, new(hostConfSuit))
}

func (suite *hostConfSuit) TestEnsureDefaultSuitHostConfIsValid() {
	// NOTE(denisacostaq@gmail.com): Giving
	// default
	hostConf := suite.HostConf

	// NOTE(denisacostaq@gmail.com): When
	// test start

	// NOTE(denisacostaq@gmail.com): Assert
	suite.Len(hostConf.validate(), 0)
}

func (suite *hostConfSuit) TestNotEmptyRef() {
	// NOTE(denisacostaq@gmail.com): Giving
	var hostConf = suite.HostConf
	hostConf.Ref = string("")

	// NOTE(denisacostaq@gmail.com): When

	// NOTE(denisacostaq@gmail.com): Assert
	suite.Len(hostConf.validate(), 1)
}

func (suite *hostConfSuit) TestNotEmptyLocation() {
	// NOTE(denisacostaq@gmail.com): Giving
	var hostConf = suite.HostConf
	hostConf.Location = string("")

	// NOTE(denisacostaq@gmail.com): When

	// NOTE(denisacostaq@gmail.com): Assert
	suite.Len(hostConf.validate(), 3) // empty, invalid and url + port invalid
}

func (suite *hostConfSuit) TestValidUrlLocation() {
	// NOTE(denisacostaq@gmail.com): Giving
	var hostConf = suite.HostConf
	hostConf.Location = string("invalid.url")

	// NOTE(denisacostaq@gmail.com): When

	// NOTE(denisacostaq@gmail.com): Assert
	suite.Len(hostConf.validate(), 1)
}

func (suite *hostConfSuit) TestCsrfAuthButEmptyTokenKeyFromEndpoint() {
	// NOTE(denisacostaq@gmail.com): Giving
	var hostConf = suite.HostConf
	hostConf.AuthType = "CSRF"
	hostConf.TokenKeyFromEndpoint = ""

	// NOTE(denisacostaq@gmail.com): When

	// NOTE(denisacostaq@gmail.com): Assert
	suite.Len(hostConf.validate(), 1)
}

func (suite *hostConfSuit) TestCsrfAuthButEmptyTokenHeaderKey() {
	// NOTE(denisacostaq@gmail.com): Giving
	var hostConf = suite.HostConf
	hostConf.AuthType = "CSRF"
	hostConf.TokenHeaderKey = ""

	// NOTE(denisacostaq@gmail.com): When

	// NOTE(denisacostaq@gmail.com): Assert
	suite.Len(hostConf.validate(), 1)
}

func (suite *hostConfSuit) TestCsrfAuthButEmptyGenTokenEndpoint() {
	// NOTE(denisacostaq@gmail.com): Giving
	var hostConf = suite.HostConf
	hostConf.AuthType = "CSRF"
	hostConf.GenTokenEndpoint = ""

	// NOTE(denisacostaq@gmail.com): When

	// NOTE(denisacostaq@gmail.com): Assert
	suite.Len(hostConf.validate(), 1)
}
