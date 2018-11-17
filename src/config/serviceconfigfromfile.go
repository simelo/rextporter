package config

import (
	"fmt"
	"strings"

	"github.com/simelo/rextporter/src/util"
	"github.com/spf13/viper"
)

// ServiceConfigFromFile get a service config from a file toml
type ServiceConfigFromFile struct {
	filePath string
}

// NewServiceConfigFromFile create a config reader configure to read config from the file in path parameter
func NewServiceConfigFromFile(path string) (conf *ServiceConfigFromFile) {
	conf = &ServiceConfigFromFile{}
	conf.filePath = path
	return conf
}

// GetConfig read the file 'filePath' and returns the service config or an error if any
func (conf ServiceConfigFromFile) GetConfig() (service Service, err error) {
	generalScopeErr := "error reading config from file"
	if strings.Compare(conf.filePath, "") == 0 {
		errCause := fmt.Sprintln("file path should not be empty, are you using the 'NewServiceConfigFromFile' function to get an instance?")
		return service, util.ErrorFromThisScope(errCause, generalScopeErr)
	}
	viper.SetConfigFile(conf.filePath)
	if err := viper.ReadInConfig(); err != nil {
		errCause := fmt.Sprintln("error reading config file: ", conf.filePath, err.Error())
		return service, util.ErrorFromThisScope(errCause, generalScopeErr)
	}
	if err := viper.Unmarshal(&service); err != nil {
		errCause := fmt.Sprintln("can not decode the config data: ", err.Error())
		return service, util.ErrorFromThisScope(errCause, generalScopeErr)
	}
	return service, err
}
