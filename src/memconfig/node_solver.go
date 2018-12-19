package memconfig

import (
	"github.com/simelo/rextporter/src/core"
)

// NodeSolver implements the interface core.RextNodeSolver
type NodeSolver struct {
	// FIXME(denisacostaq@gmail.com): lowercase
	MType    string
	nodePath string
	options  core.RextKeyValueStore
}

// GetMetricName return solver type
func (ns NodeSolver) GetType() string {
	return ns.MType
}

// SetNodePath set the node path
func (ns *NodeSolver) SetNodePath(nodePath string) {
	ns.nodePath = nodePath
}

// GetMetricType return the node path
func (ns NodeSolver) GetNodePath() string {
	return ns.nodePath
}

// GetOptions return key/value pairs for extra options
func (ns NodeSolver) GetOptions() core.RextKeyValueStore {
	return ns.options
}

// NewNodeSolver create a new node solver
func NewNodeSolver(mType, nodePath string, options core.RextKeyValueStore) core.RextNodeSolver {
	return &NodeSolver{
		MType:    mType,
		nodePath: nodePath,
		options:  options,
	}
}
