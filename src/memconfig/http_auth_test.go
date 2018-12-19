package memconfig

import (
	"testing"

	"github.com/simelo/rextporter/src/core"
	"github.com/stretchr/testify/suite"
)

func newAuth(suite *authConfSuit) core.RextAuthDef {
	return NewHTTPAuth(suite.authType, suite.authURL, suite.options)
}

type authConfSuit struct {
	suite.Suite
	authConf          core.RextAuthDef
	authType, authURL string
	options           OptionsMap
}

func (suite *authConfSuit) SetupTest() {
	suite.options = NewOptionsMap()
	suite.options.SetString("k1", "v1")
	suite.options.SetString("k2", "v2")
	suite.authType = "CSRF"
	suite.authURL = "http://localhost:9000/hosted_in/auth"
	suite.authConf = newAuth(suite)
}

func TestAuthConfSuit(t *testing.T) {
	suite.Run(t, new(authConfSuit))
}

func (suite *authConfSuit) TestNewHTTPAuth() {
	// NOTE(denisacostaq@gmail.com): Giving

	// NOTE(denisacostaq@gmail.com): When
	authConf := newAuth(suite)

	// NOTE(denisacostaq@gmail.com): Assert
	suite.Equal(suite.authType, authConf.GetAuthType())
	// TODO(denisacostaq@gmail.com):
	suite.True(eqKvs(suite.Assert(), suite.options, authConf.GetOptions()))
}

func (suite *authConfSuit) TestAbleToSetType() {
	// NOTE(denisacostaq@gmail.com): Giving
	orgAuthType := suite.authConf.GetAuthType()
	authType := "fgfg78"
	suite.authConf.SetAuthType(authType)

	// NOTE(denisacostaq@gmail.com): When
	authType2 := suite.authConf.GetAuthType()

	// NOTE(denisacostaq@gmail.com): Assert
	suite.Equal(authType, authType2)
	suite.NotEqual(orgAuthType, authType2)
}

func (suite *authConfSuit) TestInitializeEmptyOptionsInFly() {
	// NOTE(denisacostaq@gmail.com): Giving

	// NOTE(denisacostaq@gmail.com): When
	authDef := HTTPAuth{}

	// NOTE(denisacostaq@gmail.com): Assert
	suite.NotNil(authDef.GetOptions())
}
