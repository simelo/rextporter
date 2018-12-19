package memconfig

import (
	"github.com/simelo/rextporter/src/core"
)

// ResourceDef implements the interface core.RextResourceDef
type ResourceDef struct {
	mType       string
	resourceURI string
	auth        core.RextAuthDef
	decoder     core.RextDecoderDef
	metrics     []core.RextMetricDef
	options     core.RextKeyValueStore
}

func (rd ResourceDef) GetResourcePATH(basePath string) string {
	return basePath + rd.resourceURI
}

// GetType return the path type
func (rd ResourceDef) GetType() string {
	return rd.mType
}

// SetType set the type
func (rd *ResourceDef) SetType(t string) {
	rd.mType = t
}

func (rd *ResourceDef) SetResourceURI(uri string) {
	rd.resourceURI = uri
}

func (rd ResourceDef) GetAuth(defAuth core.RextAuthDef) core.RextAuthDef {
	if rd.auth == nil {
		return defAuth
	}
	return rd.auth
}

func (rd *ResourceDef) SetAuth(auth core.RextAuthDef) {
	rd.auth = auth
}

func (rd *ResourceDef) SetDecoder(decoder core.RextDecoderDef) {
	rd.decoder = decoder
}

func (rd ResourceDef) GetDecoder() core.RextDecoderDef {
	return rd.decoder
}

func (rd *ResourceDef) AddMetricDef(mtrDef core.RextMetricDef) {
	rd.metrics = append(rd.metrics, mtrDef)
}

func (rd ResourceDef) GetMetricDefs() []core.RextMetricDef {
	return rd.metrics
}

// GetOptions return key/value pairs for extra options
func (m *ResourceDef) GetOptions() core.RextKeyValueStore {
	if m.options == nil {
		m.options = NewOptionsMap()
	}
	return m.options
}

// NewResourceDef create a new metric definition
func NewResourceDef(mType, resourceURI string, auth core.RextAuthDef, metrics []core.RextMetricDef, decoder core.RextDecoderDef, options core.RextKeyValueStore) core.RextResourceDef {
	return &ResourceDef{
		mType:       mType,
		resourceURI: resourceURI,
		auth:        auth,
		decoder:     decoder,
		metrics:     metrics,
		options:     options,
	}
}
