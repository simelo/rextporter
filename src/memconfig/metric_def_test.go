package memconfig

import (
	"testing"

	"github.com/simelo/rextporter/src/core"
	"github.com/stretchr/testify/suite"
)

func newMetricDef(suite *metricDefConfSuit) core.RextMetricDef {
	return NewMetricDef(
		suite.metricName,
		suite.metricType,
		suite.metricDescription,
		suite.nodeSolver,
		suite.metricOptions,
		suite.metricLabels,
	)
}

type metricDefConfSuit struct {
	suite.Suite
	metricDef                                 core.RextMetricDef
	metricName, metricType, metricDescription string
	nodeSolver                                core.RextNodeSolver
	metricLabels                              []core.RextLabelDef
	metricOptions                             core.RextKeyValueStore
}

func (suite *metricDefConfSuit) SetupTest() {
	suite.metricName = "MySupperMetric"
	suite.metricType = core.KeyMetricTypeCounter
	suite.metricDescription = "This is all about ..."
	suite.nodeSolver = &NodeSolver{nodePath: "sds"}
	suite.metricLabels = make([]core.RextLabelDef, 1)
	suite.metricLabels[0] = &LabelDef{name: "ip"}
	suite.metricOptions = NewOptionsMap()
	suite.metricOptions.SetString("k1", "v1")
	suite.metricOptions.SetString("k2", "v2")
	suite.metricDef = newMetricDef(suite)
}

func TestMetricDefConfSuit(t *testing.T) {
	suite.Run(t, new(metricDefConfSuit))
}

func (suite *metricDefConfSuit) TestNewMetricDef() {
	// NOTE(denisacostaq@gmail.com): Giving

	// NOTE(denisacostaq@gmail.com): When
	metricDef := newMetricDef(suite)
	opts, err := suite.metricOptions.Clone()
	suite.Nil(err)
	suite.metricOptions.SetString("k1", "v2")

	// NOTE(denisacostaq@gmail.com): Assert
	suite.Equal(suite.metricName, metricDef.GetMetricName())
	suite.Equal(suite.metricType, metricDef.GetMetricType())
	suite.Equal(suite.metricDescription, metricDef.GetMetricDescription())
	suite.Equal(suite.nodeSolver, metricDef.GetNodeSolver())
	suite.True(eqKvs(suite.Assert(), suite.metricOptions, metricDef.GetOptions()))
	suite.False(eqKvs(nil, opts, metricDef.GetOptions()))
	suite.Equal(suite.metricLabels, metricDef.GetLabels())
}

func (suite *metricDefConfSuit) TestAbleToSetName() {
	// NOTE(denisacostaq@gmail.com): Giving
	orgName := suite.metricDef.GetMetricName()
	name := "fgfg78"
	suite.metricDef.SetMetricName(name)

	// NOTE(denisacostaq@gmail.com): When
	name2 := suite.metricDef.GetMetricName()

	// NOTE(denisacostaq@gmail.com): Assert
	suite.Equal(name, name2)
	suite.NotEqual(orgName, name2)
}

func (suite *metricDefConfSuit) TestAbleToSetType() {
	// NOTE(denisacostaq@gmail.com): Giving
	orgT := suite.metricDef.GetMetricType()
	tp := "fgfg78"
	suite.metricDef.SetMetricType(tp)

	// NOTE(denisacostaq@gmail.com): When
	tp2 := suite.metricDef.GetMetricType()

	// NOTE(denisacostaq@gmail.com): Assert
	suite.Equal(tp, tp2)
	suite.NotEqual(orgT, tp2)
}

func (suite *metricDefConfSuit) TestAbleToSetDescription() {
	// NOTE(denisacostaq@gmail.com): Giving
	orgDescription := suite.metricDef.GetMetricDescription()
	description := "fgfg78"
	suite.metricDef.SetMetricDescription(description)

	// NOTE(denisacostaq@gmail.com): When
	description2 := suite.metricDef.GetMetricDescription()

	// NOTE(denisacostaq@gmail.com): Assert
	suite.Equal(description, description2)
	suite.NotEqual(orgDescription, description2)
}

func (suite *metricDefConfSuit) TestAbleToSetNodeSolver() {
	// NOTE(denisacostaq@gmail.com): Giving
	orgNs := suite.metricDef.GetMetricName()
	ns := &NodeSolver{nodePath: "dfdfd"}
	suite.metricDef.SetNodeSolver(ns)

	// NOTE(denisacostaq@gmail.com): When
	ns2 := suite.metricDef.GetNodeSolver()

	// NOTE(denisacostaq@gmail.com): Assert
	suite.Equal(ns, ns2)
	suite.NotEqual(orgNs, ns2)
}
func (suite *metricDefConfSuit) TestAbleToAddLabel() {
	// NOTE(denisacostaq@gmail.com): Giving
	orgLabels := suite.metricDef.GetLabels()

	// NOTE(denisacostaq@gmail.com): When
	suite.metricDef.AddLabel(&LabelDef{})
	labels2 := suite.metricDef.GetLabels()

	// NOTE(denisacostaq@gmail.com): Assert
	suite.Equal(len(orgLabels)+1, len(labels2))
}

func (suite *metricDefConfSuit) TestInitializeEmptyOptionsInFly() {
	// NOTE(denisacostaq@gmail.com): Giving

	// NOTE(denisacostaq@gmail.com): When
	metricDef := MetricDef{}

	// NOTE(denisacostaq@gmail.com): Assert
	suite.NotNil(metricDef.GetOptions())
}
