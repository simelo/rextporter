package memconfig

import (
	"testing"

	"github.com/simelo/rextporter/src/core"
	"github.com/stretchr/testify/suite"
)

func newDecoder(suite *decoderSuit) core.RextDecoderDef {
	return NewDecoder(
		suite.decoderType,
		suite.options,
	)
}

type decoderSuit struct {
	suite.Suite
	decoder     core.RextDecoderDef
	decoderType string
	options     core.RextKeyValueStore
}

func (suite *decoderSuit) SetupTest() {
	suite.decoderType = "dfdf"
	suite.options = NewOptionsMap()
	suite.options.SetString("k1", "v1")
	suite.options.SetString("k2", "v2")
	suite.decoder = newDecoder(suite)
}

func TestDecoderSuit(t *testing.T) {
	suite.Run(t, new(decoderSuit))
}

func (suite *decoderSuit) TestNewDecoderDef() {
	// NOTE(denisacostaq@gmail.com): Giving

	// NOTE(denisacostaq@gmail.com): When
	decoderDef := newDecoder(suite)
	opts, err := suite.options.Clone()
	suite.Nil(err)
	suite.options.SetString("k1", "v2")

	// NOTE(denisacostaq@gmail.com): Assert
	suite.Equal(suite.decoderType, decoderDef.GetType())
	suite.True(eqKvs(suite.Assert(), suite.options, decoderDef.GetOptions()))
	suite.False(eqKvs(nil, opts, decoderDef.GetOptions()))
}

func (suite *decoderSuit) TestInitializeEmptyOptionsInFly() {
	// NOTE(denisacostaq@gmail.com): Giving

	// NOTE(denisacostaq@gmail.com): When
	decoderDef := Decoder{}

	// NOTE(denisacostaq@gmail.com): Assert
	suite.NotNil(decoderDef.GetOptions())
}

func (suite *decoderSuit) TestValidationClonedShouldBeValid() {
	// NOTE(denisacostaq@gmail.com): Giving

	// NOTE(denisacostaq@gmail.com): When
	cDecoderConf, err := suite.decoder.Clone()
	suite.Nil(err)
	suite.Equal(suite.decoder, cDecoderConf)
	hasError := cDecoderConf.Validate()

	// NOTE(denisacostaq@gmail.com): Assert
	suite.False(hasError)
}

func (suite *decoderSuit) TestValidationTypeShouldNotBeEmpty() {
	// NOTE(denisacostaq@gmail.com): Giving
	decoderDef := NewDecoder("", nil)

	// NOTE(denisacostaq@gmail.com): When
	hasError := decoderDef.Validate()

	// NOTE(denisacostaq@gmail.com): Assert
	suite.True(hasError)
}
