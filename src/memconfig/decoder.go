package memconfig

import (
	"github.com/simelo/rextporter/src/core"
)

// MetricDef implements the interface core.RextMetricDef
type Decoder struct {
	mType   string
	options core.RextKeyValueStore
}

// GetType return the decoder type
func (d Decoder) GetType() string {
	return d.mType
}

// GetOptions return key/value pairs for extra options
func (d *Decoder) GetOptions() core.RextKeyValueStore {
	return d.options
}

// NewDecoder create a new decoder
func NewDecoder(mType string, options core.RextKeyValueStore) *Decoder {
	return &Decoder{
		mType:   mType,
		options: options,
	}
}
