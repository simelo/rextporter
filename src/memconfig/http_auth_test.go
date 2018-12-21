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
	suite.options.SetString(core.OptKeyRextAuthDefTokenHeaderKey, "v1")
	suite.options.SetString(core.OptKeyRextAuthDefTokenGenEndpoint, "v2")
	suite.options.SetString(core.OptKeyRextAuthDefTokenKeyFromEndpoint, "v3")
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

func (suite *authConfSuit) TestValidationClonedShouldBeValid() {
	// NOTE(denisacostaq@gmail.com): Giving

	// NOTE(denisacostaq@gmail.com): When
	cAuthConf, err := suite.authConf.Clone()
	suite.Nil(err)
	hasError := cAuthConf.Validate()

	// NOTE(denisacostaq@gmail.com): Assert
	suite.False(hasError)
}

func (suite *authConfSuit) TestValidationTypeShouldNotBeEmpty() {
	// NOTE(denisacostaq@gmail.com): Giving
	authDef, err := suite.authConf.Clone()
	suite.Nil(err)

	// NOTE(denisacostaq@gmail.com): When
	authDef.SetAuthType("")
	hasError := authDef.Validate()

	// NOTE(denisacostaq@gmail.com): Assert
	suite.True(hasError)
}

func (suite *authConfSuit) TestValidationTokenHeaderKeyShouldNotBeEmptyInCSRF() {
	// NOTE(denisacostaq@gmail.com): Giving
	authDef, err := suite.authConf.Clone()
	suite.Nil(err)

	// NOTE(denisacostaq@gmail.com): When
	opts := authDef.GetOptions()
	pe, err := opts.SetString(core.OptKeyRextAuthDefTokenHeaderKey, "")
	suite.True(pe)
	suite.Nil(err)
	hasError := authDef.Validate()

	// NOTE(denisacostaq@gmail.com): Assert
	suite.True(hasError)
}

func (suite *authConfSuit) TestValidationTokenGenEndpointShouldNotBeEmptyInCSRF() {
	// NOTE(denisacostaq@gmail.com): Giving
	authDef, err := suite.authConf.Clone()
	suite.Nil(err)

	// NOTE(denisacostaq@gmail.com): When
	opts := authDef.GetOptions()
	pe, err := opts.SetString(core.OptKeyRextAuthDefTokenGenEndpoint, "")
	suite.True(pe)
	suite.Nil(err)
	hasError := authDef.Validate()

	// NOTE(denisacostaq@gmail.com): Assert
	suite.True(hasError)
}

func (suite *authConfSuit) TestValidationTokenKeyFromEndpointShouldNotBeEmptyInCSRF() {
	// NOTE(denisacostaq@gmail.com): Giving
	authDef, err := suite.authConf.Clone()
	suite.Nil(err)

	// NOTE(denisacostaq@gmail.com): When
	opts := authDef.GetOptions()
	pe, err := opts.SetString(core.OptKeyRextAuthDefTokenKeyFromEndpoint, "")
	suite.True(pe)
	suite.Nil(err)
	hasError := authDef.Validate()

	// NOTE(denisacostaq@gmail.com): Assert
	suite.True(hasError)
}

func (suite *authConfSuit) TestValidationTokenValsCanBeEmptyInNotCSRF() {
	// NOTE(denisacostaq@gmail.com): Giving
	authDef, err := suite.authConf.Clone()
	suite.Nil(err)
	authDef.SetAuthType("tt3")

	// NOTE(denisacostaq@gmail.com): When
	opts := authDef.GetOptions()
	pe, err := opts.SetString(core.OptKeyRextAuthDefTokenKeyFromEndpoint, "")
	suite.True(pe)
	suite.Nil(err)
	pe, err = opts.SetString(core.OptKeyRextAuthDefTokenGenEndpoint, "")
	suite.True(pe)
	suite.Nil(err)
	pe, err = opts.SetString(core.OptKeyRextAuthDefTokenHeaderKey, "")
	suite.True(pe)
	suite.Nil(err)
	hasError := authDef.Validate()

	// NOTE(denisacostaq@gmail.com): Assert
	suite.False(hasError)
}
