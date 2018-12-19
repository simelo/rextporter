package memconfig

import "github.com/simelo/rextporter/src/core"

// RootConfig implements core.RextRoot
type RootConfig struct {
	services []core.RextServiceDef
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

func NewRootConfig(services []core.RextServiceDef) core.RextRoot {
	return &RootConfig{services: services}
}
