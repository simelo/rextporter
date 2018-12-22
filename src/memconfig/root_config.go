package memconfig

import (
	"github.com/simelo/rextporter/src/core"
	log "github.com/sirupsen/logrus"
)

// RootConfig implements core.RextRoot
type RootConfig struct {
	services []core.RextServiceDef
}

// Clone make a deep copy of RootConfig or return an error if any
func (rc RootConfig) Clone() (cRc core.RextRoot, err error) {
	var cSrvs []core.RextServiceDef
	for _, srv := range rc.GetServices() {
		var cSrv core.RextServiceDef
		if cSrv, err = srv.Clone(); err != nil {
			log.WithError(err).Errorln("can not clone services in root config")
			return cRc, err
		}
		cSrvs = append(cSrvs, cSrv)
	}
	cRc = NewRootConfig(cSrvs)
	return cRc, err
}

func (root RootConfig) GetServices() []core.RextServiceDef {
	services := make([]core.RextServiceDef, len(root.services))
	for idxSrv := range root.services {
		services[idxSrv] = root.services[idxSrv]
	}
	return services
}

func (root *RootConfig) AddService(srv core.RextServiceDef) {
	root.services = append(root.services, srv)
}

func (root RootConfig) Validate() bool {
	return core.ValidateRoot(&root)
}

func NewRootConfig(services []core.RextServiceDef) core.RextRoot {
	return &RootConfig{services: services}
}
