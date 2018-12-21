package memconfig

import (
	"github.com/simelo/rextporter/src/core"
	log "github.com/sirupsen/logrus"
)

// Service implements core.RextDataSource interface
type Service struct {
	basePath string
	// FIXME(denisacostaq@gmail.com): how to use base path, what about protocol, port, url
	protocol string
	auth     core.RextAuthDef
	sources  []core.RextResourceDef
	options  core.RextKeyValueStore
}

// Clone make a deep copy of Service or return an error if any
func (srv Service) Clone() (cSrv core.RextServiceDef, err error) {
	var cAuth core.RextAuthDef
	if srv.auth != nil {
		if cAuth, err = srv.auth.Clone(); err != nil {
			log.WithError(err).Errorln("can not clone auth in service")
			return cSrv, err
		}
	}
	var cOpts core.RextKeyValueStore
	if cOpts, err = srv.GetOptions().Clone(); err != nil {
		log.WithError(err).Errorln("can not clone options in service")
		return cSrv, err
	}
	var cResources []core.RextResourceDef
	for _, resource := range srv.sources {
		var cResource core.RextResourceDef
		if cResource, err = resource.Clone(); err != nil {
			log.WithError(err).Errorln("can not clone resources in service")
			return cSrv, err
		}
		cResources = append(cResources, cResource)
	}
	cSrv = NewServiceConf(srv.basePath, srv.protocol, cAuth, cResources, cOpts)
	return cSrv, err
}

// SetBaseURL set the base path for the service
func (srv *Service) SetBasePath(path string) {
	srv.basePath = path
}

// GetProtocol return the service protocol
func (srv Service) GetProtocol() string {
	return srv.protocol
}

// GetBasePath return the base path
func (srv Service) GetBasePath() string {
	return srv.basePath
}

// SetMethod set the protocol for the scrapper
func (srv *Service) SetProtocol(protocol string) {
	srv.protocol = protocol
}

func (srv *Service) SetAuthForBaseURL(auth core.RextAuthDef) {
	srv.auth = auth
}

// GetAuthForBaseURL return the base auth
func (srv Service) GetAuthForBaseURL() core.RextAuthDef {
	return srv.auth
}

func (srv *Service) AddSource(source core.RextResourceDef) {
	srv.sources = append(srv.sources, source)
}

func (srv *Service) AddSources(sources ...core.RextResourceDef) {
	srv.sources = append(srv.sources, sources...)
}

func (srv Service) GetSources() []core.RextResourceDef {
	return srv.sources
}

// GetOptions return key/value pairs for extra options
func (srv *Service) GetOptions() core.RextKeyValueStore {
	if srv.options == nil {
		srv.options = NewOptionsMap()
	}
	return srv.options
}

// NewServiceConf create a new service
func NewServiceConf(basePath, protocol string, auth core.RextAuthDef, sources []core.RextResourceDef, options core.RextKeyValueStore) core.RextServiceDef {
	return &Service{
		basePath: basePath,
		protocol: protocol,
		auth:     auth,
		sources:  sources,
		options:  options,
	}
}
