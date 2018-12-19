package memconfig

import (
	"github.com/simelo/rextporter/src/core"
)

// HTTPAuth implements the core.RextAuth interface
type HTTPAuth struct {
	authType string
	endpoint string
	options  core.RextKeyValueStore
}

// SetAuthType return the auth type
func (auth *HTTPAuth) SetAuthType(authType string) {
	auth.authType = authType
}

// GetAuthType return the auth type
func (auth HTTPAuth) GetAuthType() string {
	return auth.authType
}

// GetOptions return key/value pairs for extra options
func (auth *HTTPAuth) GetOptions() core.RextKeyValueStore {
	if auth.options == nil {
		auth.options = NewOptionsMap()
	}
	return auth.options
}

// NewHTTPAuth create a auth
func NewHTTPAuth(aType, url string, options core.RextKeyValueStore) core.RextAuthDef {
	return &HTTPAuth{
		authType: aType,
		endpoint: url,
		options:  options,
	}
}
