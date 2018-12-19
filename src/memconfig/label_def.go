package memconfig

import (
	"github.com/simelo/rextporter/src/core"
)

// LabelDef implements the interface core.RextLabelDef
type LabelDef struct {
	name       string
	nodeSolver core.RextNodeSolver
}

// GetName return the label name
func (l LabelDef) GetName() string {
	return l.name
}

func (l *LabelDef) SetName(name string) {
	l.name = name
}

// GetNodeSolver return the solver for the label value
func (l LabelDef) GetNodeSolver() core.RextNodeSolver {
	return l.nodeSolver
}

// SetNodeSolver set the solver for the label value
func (l *LabelDef) SetNodeSolver(nodeSolver core.RextNodeSolver) {
	l.nodeSolver = nodeSolver
}

// NewLabelDef create a new label definition
func NewLabelDef(name string, nodeSolver core.RextNodeSolver) *LabelDef {
	return &LabelDef{
		name:       name,
		nodeSolver: nodeSolver,
	}
}
