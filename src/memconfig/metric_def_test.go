package memconfig

import (
	"testing"

	"github.com/simelo/rextporter/src/core"
	"github.com/stretchr/testify/suite"
)

type label struct {
	name  string
	vPath string
}

func (l label) equals(l2 label) bool {
	return l.name == l2.name && l.vPath == l2.vPath
}

func (l label) kvs() core.RextKeyValueStore {
	kvs := NewOptionsMap()
	kvs.SetString(l.name, l.vPath)
	return kvs
}

func ls2kvs(ls []label) []core.RextKeyValueStore {
	kvss := make([]core.RextKeyValueStore, len(ls))
	for idxL := range ls {
		kvss[idxL] = ls[idxL].kvs()
	}
	return kvss
}

func newMetricDef(suite *metricDefConfSuit) core.RextMetricDef {
	return NewMetricDef(
		suite.metricName,
		suite.metricType,
		suite.metricDescription,
		suite.metricOptions,
		ls2kvs(suite.metricLabels),
	)
}

type metricDefConfSuit struct {
	suite.Suite
	metricDef                                 core.RextMetricDef
	metricName, metricType, metricDescription string
	metricLabels                              []label
	metricOptions                             core.RextKeyValueStore
}

func (suite *metricDefConfSuit) SetupTest() {
	suite.metricName = "MySupperMetric"
	suite.metricType = core.KeyTypeCounter
	suite.metricDescription = "This is all about ..."
	suite.metricLabels = []label{label{name: "ip", vPath: "/client_ip"}}
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
	suite.True(eqKvs(suite.Assert(), suite.metricOptions, metricDef.GetOptions()))
	suite.False(eqKvs(nil, opts, metricDef.GetOptions()))
	suite.True(eqKvss(suite.Assert(), ls2kvs(suite.metricLabels), metricDef.GetMetricLabels()))
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
func (suite *metricDefConfSuit) TestAbleToSetLabels() {
	// NOTE(denisacostaq@gmail.com): Giving
	orgLabels := suite.metricDef.GetMetricLabels()
	labels := []label{label{name: "ip2", vPath: "/client_ip2"}}
	suite.metricDef.SetMetricLabels(ls2kvs(labels))

	// NOTE(denisacostaq@gmail.com): When
	labels2 := suite.metricDef.GetMetricLabels()

	// NOTE(denisacostaq@gmail.com): Assert
	suite.True(eqKvss(suite.Assert(), ls2kvs(labels), labels2))
	suite.False(eqKvss(nil, orgLabels, labels2))
}
