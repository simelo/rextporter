package memconfig

import (
	"github.com/simelo/rextporter/src/core"
	log "github.com/sirupsen/logrus"
)

// HTTPAuth implements the core.RextAuth interface
type HTTPAuth struct {
	authType string
	endpoint string
	options  core.RextKeyValueStore
}

// Clone make a deep copy of NodeSolver or return an error if any
func (auth HTTPAuth) Clone() (cAuth core.RextAuthDef, err error) {
	var cOpts core.RextKeyValueStore
	if cOpts, err = auth.GetOptions().Clone(); err != nil {
		log.WithError(err).Errorln("Can not clone options in HTTPAuth")
		return cAuth, err
	}
	cAuth = NewHTTPAuth(auth.authType, auth.endpoint, cOpts)
	return cAuth, err
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

// Validate the auth, return true if any error is found
func (auth HTTPAuth) Validate() (haveError bool) {
	return core.ValidateAuth(&auth)
}

// NewHTTPAuth create a auth
func NewHTTPAuth(aType, url string, options core.RextKeyValueStore) core.RextAuthDef {
	return &HTTPAuth{
		authType: aType,
		endpoint: url,
		options:  options,
	}
}
