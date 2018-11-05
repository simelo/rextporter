package config

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

type serverConfSuit struct {
	suite.Suite
	ServerConf Server
}

func (suite *serverConfSuit) SetupTest() {
	suite.ServerConf = Server{Location: "http://localhost:8080"}
}

func TestServerConfSuit(t *testing.T) {
	suite.Run(t, new(metricConfSuit))
}

func (suite *serverConfSuit) TestEnsureDefaultSuitServerConfIsValid() {
	// NOTE(denisacostaq@gmail.com): Giving
	// default
	serverConf := suite.ServerConf

	// NOTE(denisacostaq@gmail.com): When
	// test start

	// NOTE(denisacostaq@gmail.com): Assert
	suite.Len(serverConf.validate(), 0)
}

func (suite *serverConfSuit) TestNotEmptyName() {
	// NOTE(denisacostaq@gmail.com): Giving
	var serverConf = suite.ServerConf
	serverConf.Location = string("")

	// NOTE(denisacostaq@gmail.com): When

	// NOTE(denisacostaq@gmail.com): Assert
	suite.Len(serverConf.validate(), 1)
}
