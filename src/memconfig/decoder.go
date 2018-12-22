package memconfig

import (
	"github.com/simelo/rextporter/src/core"
	log "github.com/sirupsen/logrus"
)

// MetricDef implements the interface core.RextMetricDef
type Decoder struct {
	mType   string
	options core.RextKeyValueStore
}

// Clone make a deep copy of Decoder or return an error if any
func (d Decoder) Clone() (cD core.RextDecoderDef, err error) {
	var cOpts core.RextKeyValueStore
	if cOpts, err = d.GetOptions().Clone(); err != nil {
		log.WithError(err).Errorln("can not clone options in decoder")
		return cD, err
	}
	cD = NewDecoder(d.mType, cOpts)
	return cD, err
}

// GetType return the decoder type
func (d Decoder) GetType() string {
	return d.mType
}

// GetOptions return key/value pairs for extra options
func (d *Decoder) GetOptions() core.RextKeyValueStore {
	if d.options == nil {
		d.options = NewOptionsMap()
	}
	return d.options
}

// Validate the decoder, return true if any error is found
func (d Decoder) Validate() bool {
	return core.ValidateDecoder(&d)
}

// NewDecoder create a new decoder
func NewDecoder(mType string, options core.RextKeyValueStore) *Decoder {
	return &Decoder{
		mType:   mType,
		options: options,
	}
}
