package memconfig

import (
	"github.com/simelo/rextporter/src/core"
	log "github.com/sirupsen/logrus"
)

// HTTPAuth implements the core.RextAuth interface
type HTTPAuth struct {
	authType string
	endpoint string
}

// GetAuthType return the auth type
func (auth HTTPAuth) GetAuthType() string {
	return auth.authType
}

// GetOptions return key/value pairs for extra options
func (auth HTTPAuth) GetOptions() core.RextKeyValueStore {
	return &HTTPAuthKeyValueConf{httpAuth: auth}
}

// HTTPAuthKeyValueConf allow help to implement the core.RextKeyValueStore for http auth
type HTTPAuthKeyValueConf struct {
	httpAuth HTTPAuth
}

const (
	// OptionsKeyAuthEndpoint key used to save auth endpoint info
	OptionsKeyAuthEndpoint = "e84d8928-c074-4f19-897e-84fdaf61520a"
)

// GetString return a string value for a giving key
func (kv HTTPAuthKeyValueConf) GetString(key string) (val string, err error) {
	switch key {
	case OptionsKeyAuthEndpoint:
		val = kv.httpAuth.endpoint
	default:
		err = core.ErrKeyNotFound
		log.WithField("key", key).Errorln(err.Error())
	}
	return val, err
}

// SetString set a string value for a giving key
func (kv *HTTPAuthKeyValueConf) SetString(key, value string) (alreadyExist bool, err error) {
	alreadyExist = false
	if val, err := kv.GetString(key); err != nil && len(val) != 0 {
		alreadyExist = true
	}
	switch key {
	case OptionsKeyAuthEndpoint:
		kv.httpAuth.endpoint = value
	default:
		err = core.ErrKeyNotFound
		log.WithField("key", key).Errorln(err.Error())
	}
	return alreadyExist, err
}

// GetObject return an interface{} value for a giving key
func (kv HTTPAuthKeyValueConf) GetObject(key string) (val interface{}, err error) {
	switch key {
	case OptionsKeyAuthEndpoint:
		val = kv.httpAuth.endpoint
	default:
		err = core.ErrKeyNotFound
		log.WithField("key", key).Errorln(err.Error())
	}
	return nil, err
}

// SetObject set an interface{} value for a giving key
func (kv *HTTPAuthKeyValueConf) SetObject(key string, value interface{}) (alreadyExist bool, err error) {
	alreadyExist = false
	if val, err := kv.GetObject(key); err != nil && val != nil {
		alreadyExist = true
	}
	switch key {
	case OptionsKeyAuthEndpoint:
		strVal, okStrVal := value.(string)
		if okStrVal {
			kv.httpAuth.endpoint = strVal
		} else {
			log.WithField("val", value).Errorln("value not convertibe to string")
			err = core.ErrInvalidType
		}
	default:
		err = core.ErrKeyNotFound
		log.WithField("key", key).Errorln(err.Error())
	}
	return alreadyExist, err
}

// GetKeys return all the valid key for this storage
// FIXME(denisacostaq@gmail.com): This is a bug, no all keys have to be present all the time, the
// olemis implementation is better about this
func (kv HTTPAuthKeyValueConf) GetKeys() []string {
	return []string{OptionsKeyAuthEndpoint}
}

// Clone make a copy of HTTPAuthKeyValueConf, including the internal service
func (kv *HTTPAuthKeyValueConf) Clone() (core.RextKeyValueStore, error) {
	retKv := *kv
	return &retKv, nil
}

// NewHTTPAuth create a auth
func NewHTTPAuth(aType, url string, options core.RextKeyValueStore) *HTTPAuth {
	return &HTTPAuth{
		authType: aType,
		endpoint: url,
	}
}
