package memconfig

import (
	"testing"

	"github.com/simelo/rextporter/src/core"
	"github.com/stretchr/testify/suite"
)

func newResourceDef(suite *resourceDefSuit) core.RextResourceDef {
	return NewResourceDef(
		suite.mType,
		suite.resourceURI,
		suite.auth,
		suite.metrics,
		suite.decoder,
		suite.options,
	)
}

type resourceDefSuit struct {
	suite.Suite
	mType       string
	resourceDef core.RextResourceDef
	resourceURI string
	auth        core.RextAuthDef
	decoder     core.RextDecoderDef
	metrics     []core.RextMetricDef
	options     core.RextKeyValueStore
}

func (suite *resourceDefSuit) SetupTest() {
	suite.auth = &HTTPAuth{}
	suite.mType = "tt"
	suite.resourceURI = "ddrer"
	suite.decoder = &Decoder{}
	suite.metrics = []core.RextMetricDef{}
	suite.options = NewOptionsMap()
	suite.options.SetString("k1", "v1")
	suite.options.SetString("k2", "v2")
	suite.resourceDef = newResourceDef(suite)
}

func TestResourceDefSuitSuit(t *testing.T) {
	suite.Run(t, new(resourceDefSuit))
}

func (suite *resourceDefSuit) TestNewResourceDefSuit() {
	// NOTE(denisacostaq@gmail.com): Giving

	// NOTE(denisacostaq@gmail.com): When
	resourceDef := newResourceDef(suite)
	opts, err := suite.options.Clone()
	suite.Nil(err)
	suite.options.SetString("k1", "v2")
	basePath := "dssds"

	// NOTE(denisacostaq@gmail.com): Assert
	suite.Equal(suite.mType, resourceDef.GetType())
	suite.Equal(suite.auth, resourceDef.GetAuth(nil))
	suite.Equal(basePath+suite.resourceURI, resourceDef.GetResourcePATH(basePath))
	suite.Equal(suite.decoder, resourceDef.GetDecoder())
	suite.True(eqKvs(suite.Assert(), suite.options, resourceDef.GetOptions()))
	suite.False(eqKvs(nil, opts, resourceDef.GetOptions()))
}

func (suite *resourceDefSuit) TestAbleToSetDecoder() {
	// NOTE(denisacostaq@gmail.com): Giving
	orgDecoder := suite.resourceDef.GetDecoder()
	decoder := &Decoder{mType: "t1"}
	suite.resourceDef.SetDecoder(decoder)

	// NOTE(denisacostaq@gmail.com): When
	decoder2 := suite.resourceDef.GetDecoder()

	// NOTE(denisacostaq@gmail.com): Assert
	suite.Equal(decoder, decoder2)
	suite.NotEqual(orgDecoder, decoder2)
}

func (suite *resourceDefSuit) TestAbleToSetURI() {
	// NOTE(denisacostaq@gmail.com): Giving
	basePath := "fffd"
	orgPath := suite.resourceDef.GetResourcePATH(basePath)
	resourceURI := "uri1"
	suite.resourceDef.SetResourceURI(resourceURI)

	// NOTE(denisacostaq@gmail.com): When
	resourceURL := suite.resourceDef.GetResourcePATH(basePath)

	// NOTE(denisacostaq@gmail.com): Assert
	suite.Equal(basePath+resourceURI, resourceURL)
	suite.NotEqual(orgPath, resourceURL)
}

func (suite *resourceDefSuit) TestAbleToAddMetricDef() {
	// NOTE(denisacostaq@gmail.com): Giving
	orgMetrics := suite.resourceDef.GetMetricDefs()
	metric := &MetricDef{}
	suite.resourceDef.AddMetricDef(metric)

	// NOTE(denisacostaq@gmail.com): When
	metric2 := suite.resourceDef.GetMetricDefs()

	// NOTE(denisacostaq@gmail.com): Assert
	suite.Equal(len(orgMetrics)+1, len(metric2))
}

func (suite *resourceDefSuit) TestAbleToSetType() {
	// NOTE(denisacostaq@gmail.com): Giving
	orgType := suite.resourceDef.GetType()
	mType := "dgfg"
	suite.resourceDef.SetType(mType)

	// NOTE(denisacostaq@gmail.com): When
	mType2 := suite.resourceDef.GetType()

	// NOTE(denisacostaq@gmail.com): Assert
	suite.Equal(mType, mType2)
	suite.NotEqual(orgType, mType2)
}

func (suite *resourceDefSuit) TestInitializeEmptyOptionsInFly() {
	// NOTE(denisacostaq@gmail.com): Giving

	// NOTE(denisacostaq@gmail.com): When
	resDef := ResourceDef{}

	// NOTE(denisacostaq@gmail.com): Assert
	suite.NotNil(resDef.GetOptions())
}
