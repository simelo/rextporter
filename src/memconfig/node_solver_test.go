package memconfig

import (
	"testing"

	"github.com/simelo/rextporter/src/core"
	"github.com/stretchr/testify/suite"
)

func newNodeSolver(suite *nodeSolverSuit) core.RextNodeSolver {
	return NewNodeSolver(
		suite.nodeSolverType,
		suite.nodePath,
		suite.options,
	)
}

type nodeSolverSuit struct {
	suite.Suite
	nodeSolver               core.RextNodeSolver
	nodeSolverType, nodePath string
	options                  core.RextKeyValueStore
}

func (suite *nodeSolverSuit) SetupTest() {
	suite.nodeSolverType = core.RextNodeSolverTypeJSONPath
	suite.nodePath = "/tmp/a"
	suite.options = NewOptionsMap()
	suite.options.SetString("k1", "v1")
	suite.options.SetString("k2", "v2")
	suite.nodeSolver = newNodeSolver(suite)
}

func TestNodeSolverSuit(t *testing.T) {
	suite.Run(t, new(nodeSolverSuit))
}

func (suite *nodeSolverSuit) TestNewNodeSolver() {
	// NOTE(denisacostaq@gmail.com): Giving

	// NOTE(denisacostaq@gmail.com): When
	nodeSolver := newNodeSolver(suite)
	opts, err := suite.options.Clone()
	suite.Nil(err)
	suite.options.SetString("k1", "v2")

	// NOTE(denisacostaq@gmail.com): Assert
	suite.Equal(suite.nodeSolverType, nodeSolver.GetType())
	suite.Equal(suite.nodePath, nodeSolver.GetNodePath())
	suite.True(eqKvs(suite.Assert(), suite.options, nodeSolver.GetOptions()))
	suite.False(eqKvs(nil, opts, nodeSolver.GetOptions()))
}

func (suite *nodeSolverSuit) TestAbleToSetNodePath() {
	// NOTE(denisacostaq@gmail.com): Giving
	orgNodePath := suite.nodeSolver.GetNodePath()
	nodePath := "fgfg78"
	suite.nodeSolver.SetNodePath(nodePath)

	// NOTE(denisacostaq@gmail.com): When
	nodePath2 := suite.nodeSolver.GetNodePath()

	// NOTE(denisacostaq@gmail.com): Assert
	suite.Equal(nodePath, nodePath2)
	suite.NotEqual(orgNodePath, nodePath2)
}
