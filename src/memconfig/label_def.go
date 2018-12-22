package memconfig

import (
	"github.com/simelo/rextporter/src/core"
	log "github.com/sirupsen/logrus"
)

// LabelDef implements the interface core.RextLabelDef
type LabelDef struct {
	name       string
	nodeSolver core.RextNodeSolver
}

// Clone make a deep copy of LabelDef or return an error if any
func (l LabelDef) Clone() (cL core.RextLabelDef, err error) {
	var cNs core.RextNodeSolver
	if l.GetNodeSolver() != nil {
		if cNs, err = l.GetNodeSolver().Clone(); err != nil {
			log.WithError(err).Errorln("can not clone node solver in label")
			return cL, err
		}
	}
	cL = NewLabelDef(l.name, cNs)
	return cL, err
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

// Validate the label, return true if any error is found
func (l LabelDef) Validate() bool {
	return core.ValidateLabel(&l)
}

// NewLabelDef create a new label definition
func NewLabelDef(name string, nodeSolver core.RextNodeSolver) *LabelDef {
	return &LabelDef{
		name:       name,
		nodeSolver: nodeSolver,
	}
}
