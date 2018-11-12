package config

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/shibukawa/configdir"
	"github.com/simelo/rextporter/src/common"
	"github.com/spf13/viper"
)

type serviceConfigTmplData struct {
	metricsConfigPath string
}

type serviceConfigData struct {
	serviceConfPath string
	tmplData        serviceConfigTmplData
}

type mainConfigTmplData struct {
	serviceConfData serviceConfigData
}

type mainConfigData struct {
	mainConfPath string
	tmplData     mainConfigTmplData
}

func (serviceData serviceConfigData) MetricsConfigPath() string {
	return serviceData.serviceConfPath
}

func (mainConfData mainConfigData) MetricsConfigPath() string {
	return mainConfData.tmplData.serviceConfData.MetricsConfigPath()
}

func (confData mainConfigData) ServiceConfigPath() string {
	return confData.tmplData.serviceConfData.serviceConfPath
}

func (confData mainConfigData) MainConfigPath() string {
	return confData.mainConfPath
}

const mainConfigFileContentTemplate = `
serviceConfigTransport = "file" # "file" | "consulCatalog"
# render a template with a portable path
serviceConfigPath = "{{.ServiceConfigPath}}"
`

const serviceConfigFileContentTemplate = `
# Service configuration.
[[services]]
  name = "wallet1"
  scheme = "http"
  port = 8000
  basePath = ""
  authType = "CSRF"
  tokenHeaderKey = "X-CSRF-Token"
  genTokenEndpoint = "/api/v1/csrf.json"
  tokenKeyFromEndpoint = "csrf_token"

  [services.location]
		location = "localhost"
	[services.location]
		location = "{{.MetricsConfigPath}}"
`
const metricsConfigFileContentTemplate = `
# All metrics to be measured.
[[metrics]]
  name = "seq"
  url = "/api/v1/health.json"
  httpMethod = "GET"
  path = "/blockchain/head/seq"

  [metrics.options]
    type = "Counter"
    description = "I am running since"

# [[metrics]]
#   name = "openConnections"
#   url = "/api/v1/network/connections"
#   httpMethod = "GET"
#   path = "/"

#   [metrics.options]
#     type = "Histogram"
#     description = "Connections ammount"

#   [metrics.histogramOptions]
#     buckets = [
#       1,
#       2, 
#       3
#     ]




# TODO(denisacostaq@gmail.com):
# if you refer(under "metrics_for_host") to a not previously defined host or metric it will be raise an error and the process will not start
# if in all your definition you not use some host or metric the process will raise a warning and the process will start normally.
`

var (
	systemVendorName      = "simelo"
	systemProgramName     = "rextporter"
	mainConfigFileName    = "main.toml"
	serviceConfigFileName = "service.toml"
	metricsConfigFileName = "metrics.toml"
)

func createFullPath(path string) error {
	return os.MkdirAll(path, 0750)
}

func creteFullPathForFile(filePath string) (err error) {
	dir, _ := filepath.Split(filePath)
	return createFullPath(dir)
}

// isADirectoryPath get info about the path string not about a phisical resource in the filesystem.
// return true if the path is a directory path
func isADirectoryPath(path string) bool {
	dir, file := filepath.Split(path)
	return (strings.Compare(dir, "") != 0 && strings.Compare(file, "") == 0)
}

func existFile(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}

func (confData mainConfigData) existServiceConfigFile() bool {
	return existFile(confData.ServiceConfigPath())
}

// createServiceConfigFile creates the service file or return an error if any,
// if the file already exist does no thin.
func (confData mainConfigData) createServiceConfigFile() (err error) {
	generalScopeErr := "error creating service config file"
	if confData.existServiceConfigFile() {
		return nil
	}
	tmpl := template.New("serviceConfig")
	var templateEngine *template.Template
	if templateEngine, err = tmpl.Parse(serviceConfigFileContentTemplate); err != nil {
		errCause := "error parsing service config: " + err.Error()
		return common.ErrorFromThisScope(errCause, generalScopeErr)
	}
	if err = creteFullPathForFile(confData.ServiceConfigPath()); err != nil {
		errCause := "error creating directory for service file: " + err.Error()
		return common.ErrorFromThisScope(errCause, generalScopeErr)
	}
	var serviceConfigFile *os.File
	if serviceConfigFile, err = os.Create(confData.ServiceConfigPath()); err != nil {
		errCause := "error creating service config file: " + err.Error()
		return common.ErrorFromThisScope(errCause, generalScopeErr)
	}
	if err = templateEngine.Execute(serviceConfigFile, nil); err != nil {
		errCause := "error writing main config file: " + err.Error()
		return common.ErrorFromThisScope(errCause, generalScopeErr)
	}
	return err
}

func (confData serviceConfigData) existMetricsConfigFile() bool {
	return existFile(confData.tmplData.metricsConfigPath)
}

// createMetricsConfigFile creates the metrics file or return an error if any,
// if the file already exist does no thin.
func (confData serviceConfigData) createMetricsConfigFile() (err error) {
	generalScopeErr := "error creating metrics config file"
	if confData.existMetricsConfigFile() {
		return nil
	}
	tmpl := template.New("metricsConfig")
	var templateEngine *template.Template
	if templateEngine, err = tmpl.Parse(metricsConfigFileContentTemplate); err != nil {
		errCause := "error parsing metrics config: " + err.Error()
		return common.ErrorFromThisScope(errCause, generalScopeErr)
	}
	if err = creteFullPathForFile(confData.MetricsConfigPath()); err != nil {
		errCause := "error creating directory for metrics file: " + err.Error()
		return common.ErrorFromThisScope(errCause, generalScopeErr)
	}
	var metricsConfigFile *os.File
	if metricsConfigFile, err = os.Create(confData.MetricsConfigPath()); err != nil {
		errCause := "error creating metrics config file: " + err.Error()
		return common.ErrorFromThisScope(errCause, generalScopeErr)
	}
	if err = templateEngine.Execute(metricsConfigFile, nil); err != nil {
		errCause := "error writing metrics config file: " + err.Error()
		return common.ErrorFromThisScope(errCause, generalScopeErr)
	}
	return err
}

func (confData mainConfigData) existMainConfigFile() bool {
	return existFile(confData.mainConfPath)
}

// createMainConfigFile creates the main file or return an error if any,
// if the file already exist does no thin.
func (confData mainConfigData) createMainConfigFile() (err error) {
	generalScopeErr := "error creating main config file"
	if confData.existMainConfigFile() {
		return nil
	}
	tmpl := template.New("mainConfig")
	var templateEngine *template.Template
	if templateEngine, err = tmpl.Parse(mainConfigFileContentTemplate); err != nil {
		errCause := "error parsing main config: " + err.Error()
		return common.ErrorFromThisScope(errCause, generalScopeErr)
	}
	if err = creteFullPathForFile(confData.MainConfigPath()); err != nil {
		errCause := "error creating directory for main file: " + err.Error()
		return common.ErrorFromThisScope(errCause, generalScopeErr)
	}
	var mainConfigFile *os.File
	if mainConfigFile, err = os.Create(confData.MainConfigPath()); err != nil {
		errCause := "error creating main config file: " + err.Error()
		return common.ErrorFromThisScope(errCause, generalScopeErr)
	}
	if err = templateEngine.Execute(mainConfigFile, confData.tmplData); err != nil {
		errCause := "error writing main config file: " + err.Error()
		return common.ErrorFromThisScope(errCause, generalScopeErr)
	}
	return err
}

func homeConfigFolder() (*configdir.Config, error) {
	configDirs := configdir.New(systemVendorName, systemProgramName)
	folders := configDirs.QueryFolders(configdir.Global)
	if len(folders) <= 0 {
		return nil, errors.New("some strange error was happen, can not determine the home config folder")
	}
	return folders[0], nil
}

func fileDefaultConfigPath(fileName string, homeConf *configdir.Config) (path string) {
	return filepath.Join(homeConf.Path, fileName)
}

func metricsDefaultConfigPath(conf *configdir.Config) (path string) {
	return fileDefaultConfigPath(metricsConfigFileName, conf)
}

func serviceDefaultConfigPath(conf *configdir.Config) (path string) {
	return fileDefaultConfigPath(serviceConfigFileName, conf)
}

func mainDefaultConfigPath(conf *configdir.Config) (path string) {
	return fileDefaultConfigPath(mainConfigFileName, conf)
}

func defaultMainConfigTmplData(conf *configdir.Config) (tmplData mainConfigTmplData) {
	tmplData = mainConfigTmplData{
		serviceConfData: serviceConfigData{
			serviceConfPath: serviceDefaultConfigPath(conf),
			tmplData:        defaultServiceTmplData(conf),
		},
	}
	return tmplData
}

func defaultServiceTmplData(conf *configdir.Config) (tmplData serviceConfigTmplData) {
	tmplData = serviceConfigTmplData{
		metricsConfigPath: metricsDefaultConfigPath(conf),
	}
	return tmplData
}

func tmplDataFromMainFile(mainConfigFilePath string) (tmpl serviceConfigTmplData, err error) {
	generalScopeErr := "error filling template data"
	viper.SetConfigFile(mainConfigFilePath)
	viper.SetConfigType("toml")
	if err := viper.ReadInConfig(); err != nil {
		errCause := fmt.Sprintln("error reading config file: ", mainConfigFilePath, err.Error())
		return tmpl, common.ErrorFromThisScope(errCause, generalScopeErr)
	}
	var mainConf mainConfigTmplData
	if err := viper.Unmarshal(&mainConf); err != nil {
		errCause := fmt.Sprintln("can not decode the config data: ", err.Error())
		return tmpl, common.ErrorFromThisScope(errCause, generalScopeErr)
	}
	tmpl.metricsConfigPath = mainConf.serviceConfData.tmplData.metricsConfigPath
	return tmpl, err
}

func newMainConfigData(path string) (mainConf mainConfigData, err error) {
	generalScopeErr := "can not create main config instance"
	if isADirectoryPath(path) {
		path = filepath.Join(path, mainConfigFileName)
	}
	var tmplData mainConfigTmplData
	if strings.Compare(path, "") == 0 || !existFile(path) {
		// TODO(denisacostaq@gmail.com): move homeConf to fn defaultTmplData
		var homeConf *configdir.Config
		if homeConf, err = homeConfigFolder(); err != nil {
			errCause := "error looking for config folder under home: " + err.Error()
			return mainConf, common.ErrorFromThisScope(errCause, generalScopeErr)
		}
		path = mainDefaultConfigPath(homeConf)
		tmplData = defaultMainConfigTmplData(homeConf)
	} else {
		if tmplData.serviceConfData.tmplData, err = tmplDataFromMainFile(path); err != nil {
			errCause := "error reading template data from file: " + err.Error()
			return mainConf, common.ErrorFromThisScope(errCause, generalScopeErr)
		}
	}
	if strings.Compare(tmplData.serviceConfData.serviceConfPath, "") == 0 || strings.Compare(tmplData.serviceConfData.tmplData.metricsConfigPath, "") == 0 {
		var homeConf *configdir.Config
		if homeConf, err = homeConfigFolder(); err != nil {
			errCause := "error looking for config folder under home: " + err.Error()
			return mainConf, common.ErrorFromThisScope(errCause, generalScopeErr)
		}
		tmpTmplData := defaultMainConfigTmplData(homeConf)
		if strings.Compare(tmplData.serviceConfData.serviceConfPath, "") == 0 {
			tmplData.serviceConfData.serviceConfPath = tmpTmplData.serviceConfData.serviceConfPath
		}
		if strings.Compare(tmplData.serviceConfData.tmplData.metricsConfigPath, "") == 0 {
			tmplData.serviceConfData.tmplData.metricsConfigPath = tmpTmplData.serviceConfData.tmplData.metricsConfigPath
		}
	}
	mainConf = mainConfigData{
		mainConfPath: path,
		tmplData:     tmplData,
	}
	if err = mainConf.createMainConfigFile(); err != nil {
		errCause := "error creating main config file: " + err.Error()
		return mainConf, common.ErrorFromThisScope(errCause, generalScopeErr)
	}
	if err = mainConf.tmplData.serviceConfData.createMetricsConfigFile(); err != nil {
		errCause := "error creating metrics config file: " + err.Error()
		return mainConf, common.ErrorFromThisScope(errCause, generalScopeErr)
	}
	if err = mainConf.createServiceConfigFile(); err != nil {
		errCause := "error creating service config file: " + err.Error()
		return mainConf, common.ErrorFromThisScope(errCause, generalScopeErr)
	}
	return mainConf, err
}
