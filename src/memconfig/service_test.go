package memconfig

import (
	"testing"

	"github.com/simelo/rextporter/src/core"
	"github.com/simelo/rextporter/src/core/mocks"
	"github.com/stretchr/testify/suite"
)

type serviceConfSuit struct {
	suite.Suite
	srvConf   core.RextServiceDef
	basePath  string
	protocol  string
	auth      core.RextAuthDef
	resources []core.RextResourceDef
	options   core.RextKeyValueStore
}

func newService(suite *serviceConfSuit) core.RextServiceDef {
	return NewServiceConf(suite.basePath, suite.protocol, suite.auth, suite.resources, suite.options)
}

func (suite *serviceConfSuit) SetupTest() {
	suite.basePath = "/hosted_in/root"
	suite.protocol = "file"
	suite.auth = &HTTPAuth{}
	suite.auth.GetOptions()
	suite.resources = []core.RextResourceDef{&ResourceDef{}}
	suite.resources[0].GetOptions()
	suite.options = NewOptionsMap()
	suite.options.SetString(core.OptKeyRextServiceDefJobName, "v1")
	suite.options.SetString(core.OptKeyRextServiceDefInstanceName, "v2")
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
	orgResource := suite.srvConf.GetSources()
	resource := &ResourceDef{}
	suite.srvConf.AddResource(resource)

	// NOTE(denisacostaq@gmail.com): When
	resource2 := suite.srvConf.GetSources()

	// NOTE(denisacostaq@gmail.com): Assert
	suite.Equal(len(orgResource)+1, len(resource2))
}

func (suite *serviceConfSuit) TestInitializeEmptyOptionsInFly() {
	// NOTE(denisacostaq@gmail.com): Giving

	// NOTE(denisacostaq@gmail.com): When
	srvConf := Service{}

	// NOTE(denisacostaq@gmail.com): Assert
	suite.NotNil(srvConf.GetOptions())
}

func (suite *serviceConfSuit) TestValidationClonedShouldBeValid() {
	// NOTE(denisacostaq@gmail.com): Giving
	cSrvConf, err := suite.srvConf.Clone()
	suite.Nil(err)
	suite.Equal(suite.srvConf, cSrvConf)
	setUpFakeValidationOn3rdPartyOverService(cSrvConf)

	// NOTE(denisacostaq@gmail.com): When
	hasError := cSrvConf.Validate()

	// NOTE(denisacostaq@gmail.com): Assert
	suite.False(hasError)
}

func (suite *serviceConfSuit) TestValidationJobNameShouldNotBeEmpty() {
	// NOTE(denisacostaq@gmail.com): Giving
	srvConf, err := suite.srvConf.Clone()

	// NOTE(denisacostaq@gmail.com): When
	opts := srvConf.GetOptions()
	pe, err := opts.SetString(core.OptKeyRextServiceDefJobName, "")
	suite.True(pe)
	suite.Nil(err)
	hasError := srvConf.Validate()

	// NOTE(denisacostaq@gmail.com): Assert
	suite.True(hasError)
}

func (suite *serviceConfSuit) TestValidationInstanceNameShouldNotBeEmpty() {
	// NOTE(denisacostaq@gmail.com): Giving
	srvConf, err := suite.srvConf.Clone()

	// NOTE(denisacostaq@gmail.com): When
	opts := srvConf.GetOptions()
	pe, err := opts.SetString(core.OptKeyRextServiceDefInstanceName, "")
	suite.True(pe)
	suite.Nil(err)
	hasError := srvConf.Validate()

	// NOTE(denisacostaq@gmail.com): Assert
	suite.True(hasError)
}

func (suite *serviceConfSuit) TestValidationEmptyProtocol() {
	// NOTE(denisacostaq@gmail.com): Giving

	// NOTE(denisacostaq@gmail.com): When
	cSrvConf, err := suite.srvConf.Clone()
	suite.Nil(err)
	cSrvConf.SetProtocol("")
	hasError := cSrvConf.Validate()

	// NOTE(denisacostaq@gmail.com): Assert
	suite.True(hasError)
}

func (suite *serviceConfSuit) TestValidationShouldGoDownTroughFields() {
	// NOTE(denisacostaq@gmail.com): Giving
	cSrvConf, err := suite.srvConf.Clone()
	suite.Nil(err)
	mockAuth := new(mocks.RextAuthDef)
	mockAuth.On("Validate").Return(false)
	mockResource1 := new(mocks.RextResourceDef)
	mockResource1.On("Validate").Return(false)
	mockResource2 := new(mocks.RextResourceDef)
	mockResource2.On("Validate").Return(false)
	cSrvConf.SetAuthForBaseURL(mockAuth)
	cSrvConf.AddResources(mockResource1, mockResource2)

	// NOTE(denisacostaq@gmail.com): When
	cSrvConf.Validate()

	// NOTE(denisacostaq@gmail.com): Assert
	mockAuth.AssertCalled(suite.T(), "Validate")
	suite.Len(cSrvConf.GetSources(), 3)
	mockResource1.AssertCalled(suite.T(), "Validate")
	mockResource2.AssertCalled(suite.T(), "Validate")
}

func setUpFakeValidationOn3rdPartyOverService(srv core.RextServiceDef) {
	authStub := new(mocks.RextAuthDef)
	authStub.On("Validate").Return(false)
	sourceStub := new(mocks.RextResourceDef)
	sourceStub.On("Validate").Return(false)
	srv.SetAuthForBaseURL(authStub)
	srv.AddResource(sourceStub)
}
