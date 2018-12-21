package memconfig

import (
	"github.com/simelo/rextporter/src/core"
	log "github.com/sirupsen/logrus"
)

// NodeSolver implements the interface core.RextNodeSolver
type NodeSolver struct {
	// FIXME(denisacostaq@gmail.com): lowercase
	MType    string
	nodePath string
	options  core.RextKeyValueStore
}

// Clone make a deep copy of NodeSolver or return an error if any
func (ns NodeSolver) Clone() (cNs core.RextNodeSolver, err error) {
	var cOpts core.RextKeyValueStore
	if cOpts, err = ns.GetOptions().Clone(); err != nil {
		log.WithError(err).Errorln("can not clone options in node solver")
		return cNs, err
	}
	cNs = NewNodeSolver(ns.MType, ns.nodePath, cOpts)
	return cNs, err
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

func (ns NodeSolver) Validate() bool {
	return core.ValidateNodeSolver(&ns)
}

// NewNodeSolver create a new node solver
func NewNodeSolver(mType, nodePath string, options core.RextKeyValueStore) core.RextNodeSolver {
	return &NodeSolver{
		MType:    mType,
		nodePath: nodePath,
		options:  options,
	}
}
