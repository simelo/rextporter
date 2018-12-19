package memconfig

import (
	"testing"

	"github.com/simelo/rextporter/src/core"
	"github.com/stretchr/testify/suite"
)

type serviceConfSuit struct {
	suite.Suite
	srvConf  core.RextServiceDef
	basePath string
	protocol string
	auth     core.RextAuthDef
	sources  []core.RextResourceDef
	options  core.RextKeyValueStore
}

func newService(suite *serviceConfSuit) core.RextServiceDef {
	return NewServiceConf(suite.basePath, suite.protocol, suite.auth, suite.sources, suite.options)
}

func (suite *serviceConfSuit) SetupTest() {
	suite.basePath = "/hosted_in/root"
	suite.protocol = "file"
	suite.auth = &HTTPAuth{}
	suite.sources = []core.RextResourceDef{&ResourceDef{}}
	suite.options = NewOptionsMap()
	suite.options.SetString("k1", "v1")
	suite.options.SetString("k2", "v2")
	suite.srvConf = newService(suite)
}

func TestServiceConfSuit(t *testing.T) {
	suite.Run(t, new(serviceConfSuit))
}

func (suite *serviceConfSuit) TestNewServiceDef() {
	// NOTE(denisacostaq@gmail.com): Giving

	// NOTE(denisacostaq@gmail.com): When
	serviceDef := newService(suite)
	opts, err := suite.options.Clone()
	suite.Nil(err)
	suite.options.SetString("k1", "v2")

	// NOTE(denisacostaq@gmail.com): Assert
	suite.Equal(suite.protocol, serviceDef.GetProtocol())
	suite.True(eqKvs(suite.Assert(), suite.options, serviceDef.GetOptions()))
	suite.False(eqKvs(nil, opts, serviceDef.GetOptions()))
}

func (suite *serviceConfSuit) TestAbleToSetProtocol() {
	// NOTE(denisacostaq@gmail.com): Giving
	protocol := "http"
	suite.srvConf.SetProtocol(protocol)

	// NOTE(denisacostaq@gmail.com): When
	protocol2 := suite.srvConf.GetProtocol()

	// NOTE(denisacostaq@gmail.com): Assert
	suite.Equal(protocol, protocol2)
}

func (suite *serviceConfSuit) TestAbleToSetBasePath() {
	// NOTE(denisacostaq@gmail.com): Giving
	basePath := "dfdf"
	suite.srvConf.SetBasePath(basePath)

	// NOTE(denisacostaq@gmail.com): When
	basePath2 := suite.srvConf.GetBasePath()

	// NOTE(denisacostaq@gmail.com): Assert
	suite.Equal(basePath, basePath2)
}

func (suite *serviceConfSuit) TestAbleToSetBaseAuth() {
	// NOTE(denisacostaq@gmail.com): Giving
	orgAuth := suite.srvConf.GetAuthForBaseURL()
	auth := &HTTPAuth{authType: "ds"}
	suite.srvConf.SetAuthForBaseURL(auth)

	// NOTE(denisacostaq@gmail.com): When
	auth2 := suite.srvConf.GetAuthForBaseURL()

	// NOTE(denisacostaq@gmail.com): Assert
	suite.Equal(auth, auth2)
	suite.NotEqual(orgAuth, auth2)
}

func (suite *serviceConfSuit) TestAbleToAddSource() {
	// NOTE(denisacostaq@gmail.com): Giving
	orgSources := suite.srvConf.GetSources()
	source := &ResourceDef{}
	suite.srvConf.AddSource(source)

	// NOTE(denisacostaq@gmail.com): When
	sources2 := suite.srvConf.GetSources()

	// NOTE(denisacostaq@gmail.com): Assert
	suite.Equal(len(orgSources)+1, len(sources2))
}

func (suite *serviceConfSuit) TestInitializeEmptyOptionsInFly() {
	// NOTE(denisacostaq@gmail.com): Giving

	// NOTE(denisacostaq@gmail.com): When
	resDef := Service{}

	// NOTE(denisacostaq@gmail.com): Assert
	suite.NotNil(resDef.GetOptions())
}
