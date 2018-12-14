package memconfig

import (
	"github.com/simelo/rextporter/src/core"
	log "github.com/sirupsen/logrus"
)

// Service implements core.RextDataSource interface
type Service struct {
	baseURL    string
	location   string
	method     string
	extractors []core.RextMetricsExtractor
}

// SetBaseURL set the base url for the service
func (srv *Service) SetBaseURL(url string) {
	srv.baseURL = url
}

// GetMethod return the service http method
func (srv Service) GetMethod() string {
	return srv.method
}

// SetMethod set the http method for the scrapper
func (srv *Service) SetMethod(method string) {
	srv.method = method
}

// GetResourceLocation get the service resource location
func (srv Service) GetResourceLocation() string {
	return srv.location
}

// SetResourceLocation set a resource location for the service
func (srv *Service) SetResourceLocation(location string) error {
	srv.location = location
	return nil
}

// GetOptions return key/value pairs for extra options
func (srv Service) GetOptions() core.RextKeyValueStore {
	return &ServiceKeyValueConf{srv: srv}
}

// ActivateScraper activate the giving extractor in scrapper
func (srv *Service) ActivateScraper(me core.RextMetricsExtractor) (err error) {
	if extME, isME := me.(*MetricsExtractor); isME {
		srv.extractors = append(srv.extractors, extME)
	} else {
		err = core.ErrInvalidType
	}
	return err
}

// ServiceKeyValueConf allow help to implement the core.RextKeyValueStore for service srv
type ServiceKeyValueConf struct {
	srv Service
}

const (
	// OptionsKeyServiceBaseURL key used to save base url info
	OptionsKeyServiceBaseURL = "e71ea9fe-6a90-49c6-bfd0-07a34b0cc780"
	// OptionsKeyServiceLocation key used to save location info
	OptionsKeyServiceLocation = "4f31ac1f-94e5-4cf7-abc6-fe8f57946622"
	// OptionsKeyServiceMethod key used to save method info
	OptionsKeyServiceMethod = "506926e9-c4df-4784-8c2c-e0c178ef91b7"
)

// GetKeys return all the valid key for this storage
// FIXME(denisacostaq@gmail.com): This is a bug, no all keys have to be present all the time, the
// olemis implementation is better about this
func (kv ServiceKeyValueConf) GetKeys() []string {
	return []string{
		OptionsKeyServiceBaseURL,
		OptionsKeyServiceLocation,
		OptionsKeyServiceMethod,
	}
}

// GetString return a string value for a giving key
func (kv ServiceKeyValueConf) GetString(key string) (val string, err error) {
	switch key {
	case OptionsKeyServiceBaseURL:
		val = kv.srv.baseURL
	case OptionsKeyServiceLocation:
		val = kv.srv.location
	case OptionsKeyServiceMethod:
		val = kv.srv.method
	default:
		err = core.ErrKeyNotFound
		log.WithField("key", key).Errorln(err.Error())
	}
	return val, err
}

// SetString set a string value for a giving key
func (kv *ServiceKeyValueConf) SetString(key string, value string) (alreadyExist bool, err error) {
	alreadyExist = false
	if val, err := kv.GetString(key); err != nil && len(val) != 0 {
		alreadyExist = true
	}
	switch key {
	case OptionsKeyServiceBaseURL:
		kv.srv.baseURL = value
	case OptionsKeyServiceLocation:
		kv.srv.location = value
	case OptionsKeyServiceMethod:
		kv.srv.method = value
	default:
		err = core.ErrKeyNotFound
		log.WithField("key", key).Errorln(err.Error())
	}
	return alreadyExist, err
}

// GetObject return an interface{} value for a giving key
func (kv ServiceKeyValueConf) GetObject(key string) (val interface{}, err error) {
	switch key {
	case OptionsKeyServiceBaseURL:
		val = kv.srv.baseURL
	case OptionsKeyServiceLocation:
		val = kv.srv.location
	case OptionsKeyServiceMethod:
		val = kv.srv.method
	default:
		err = core.ErrKeyNotFound
		log.WithField("key", key).Errorln(err.Error())
	}
	return val, err
}

// SetObject set an interface{} value for a giving key
func (kv *ServiceKeyValueConf) SetObject(key string, value interface{}) (alreadyExist bool, err error) {
	alreadyExist = false
	if val, err := kv.GetObject(key); err != nil && val != nil {
		alreadyExist = true
	}
	switch key {
	case OptionsKeyServiceBaseURL:
		baseURL, okBaseURL := value.(string)
		if !okBaseURL {
			kv.srv.baseURL = baseURL
		} else {
			log.WithField("val", value).Errorln("value not convertible to string")
			err = core.ErrInvalidType
		}
	case OptionsKeyServiceLocation:
		location, okLocation := value.(string)
		if !okLocation {
			kv.srv.location = location
		} else {
			log.WithField("val", value).Errorln("value not convertible to string")
			err = core.ErrInvalidType
		}
	case OptionsKeyServiceMethod:
		method, okMethod := value.(string)
		if !okMethod {
			kv.srv.method = method
		} else {
			log.WithField("val", value).Errorln("value not convertible to string")
			err = core.ErrInvalidType
		}
	default:
		err = core.ErrKeyNotFound
		log.WithField("key", key).Errorln(err.Error())
	}
	return alreadyExist, err
}

// Clone make a copy of ServiceKeyValueConf, including the internal service
func (kv *ServiceKeyValueConf) Clone() (core.RextKeyValueStore, error) {
	retKv := *kv
	return &retKv, nil
}

// NewServiceConf create a new service
func NewServiceConf(baseURL, location, httpMethod string) core.RextDataSource {
	return &Service{
		baseURL:  baseURL,
		location: location,
		method:   httpMethod,
	}
}
