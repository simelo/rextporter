package memconfig

import (
	"testing"

	"github.com/simelo/rextporter/src/core"
	"github.com/stretchr/testify/suite"
)

func newAuth(suite *authConfSuit) core.RextAuth {
	return NewHTTPAuth(suite.authType, suite.authURL, suite.authOptions)
}

type authConfSuit struct {
	suite.Suite
	authConf          core.RextAuth
	authType, authURL string
	authOptions       OptionsMap
}

func (suite *authConfSuit) SetupTest() {
	suite.authOptions = NewOptionsMap()
	suite.authOptions.SetString("k1", "v1")
	suite.authOptions.SetString("k2", "v2")
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
	// suite.True(eqKvs(suite.Assert(), suite.authOptions, authConf.GetOptions()))
}
