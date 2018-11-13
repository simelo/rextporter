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

type templateData struct {
	ServiceConfigPath      string
	MetricsForServicesPath string
}

type metricsForServiceTemplateData struct {
	ServiceNameToMetricsConfPath map[string]string
}

type metricsForServiceConfigTemplateData struct {
	TmplData metricsForServiceTemplateData
}

type mainConfigData struct {
	mainConfigPath                  string
	tmplData                        templateData
	metricsForServiceConfigTmplData metricsForServiceConfigTemplateData
}

func (confData mainConfigData) ServiceConfigPath() string {
	return confData.tmplData.ServiceConfigPath
}

func (confData mainConfigData) metricsForServicesPath() string {
	return confData.tmplData.MetricsForServicesPath
}

func (confData mainConfigData) MetricsConfigPath(serviceName string) string {
	return confData.metricsForServiceConfigTmplData.TmplData.ServiceNameToMetricsConfPath[serviceName]
}

func (confData mainConfigData) MainConfigPath() string {
	return confData.mainConfigPath
}

const mainConfigFileContentTemplate = `
servicesConfigTransport = "file" # "file" | "consulCatalog"
servicesConfigPath = "{{.ServiceConfigPath}}"
metricsForServicesPath = "{{.MetricsForServicesPath}}"
`

const serviceConfigFileContentTemplate = `
# Services configuration.
[[services]]
  name = "skycoin"
  mode = "apiRest"
  scheme = "http"
  port = 8000
  basePath = ""
  authType = "CSRF"
  tokenHeaderKey = "X-CSRF-Token"
  genTokenEndpoint = "/api/v1/csrf"
  tokenKeyFromEndpoint = "csrf_token"

  [services.location]
    location = "localhost"
`
const skycoinMetricsConfigFileContentTemplate = `
# All metrics to be measured.
[[metrics]]
  name = "seq"
  url = "/api/v1/health"
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

const metricsForServiceMappingConfFileContentTemplate = `
serviceNameToMetricsConfPath = [{{range $key, $value := .}}
	{ {{$key}} = "{{$value}}" },{{end}}
]
`

var (
	systemVendorName                 = "simelo"
	systemProgramName                = "rextporter"
	mainConfigFileName               = "main.toml"
	serviceConfigFileName            = "services.toml"
	metricsForServicesConfigFileName = "metricsForServices.toml"
	skycoinMetricsConfigFileName     = "skycoinMetrics.toml"
	walletMetricsConfigFileName      = "walletMetrics.toml"
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

// createMetricsConfigFile creates the metrics file or return an error if any,
// if the file already exist does no thin.
func createMetricsConfigFile(metricConfPath string) (err error) {
	generalScopeErr := "error creating metrics config file"
	if existFile(metricConfPath) {
		return nil
	}
	tmpl := template.New("metricsConfig")
	var templateEngine *template.Template
	if templateEngine, err = tmpl.Parse(skycoinMetricsConfigFileContentTemplate); err != nil {
		errCause := "error parsing metrics config: " + err.Error()
		return common.ErrorFromThisScope(errCause, generalScopeErr)
	}
	if err = creteFullPathForFile(metricConfPath); err != nil {
		errCause := "error creating directory for metrics file: " + err.Error()
		return common.ErrorFromThisScope(errCause, generalScopeErr)
	}
	var metricsConfigFile *os.File
	if metricsConfigFile, err = os.Create(metricConfPath); err != nil {
		errCause := "error creating metrics config file: " + err.Error()
		return common.ErrorFromThisScope(errCause, generalScopeErr)
	}
	if err = templateEngine.Execute(metricsConfigFile, nil); err != nil {
		errCause := "error writing metrics config file: " + err.Error()
		return common.ErrorFromThisScope(errCause, generalScopeErr)
	}
	return err
}

func (confData mainConfigData) existMetricsForServicesConfigFile() bool {
	return existFile(confData.tmplData.MetricsForServicesPath)
}

// createMetricsForServicesConfFile creates the metrics for services file or return an error if any,
// if the file already exist does no thin.
func (confData mainConfigData) createMetricsForServicesConfFile() (err error) {
	generalScopeErr := "error creating metrics for services config file"
	if confData.existMetricsForServicesConfigFile() {
		return nil
	}
	tmpl := template.New("metricsForServiceConfig")
	var templateEngine *template.Template
	if templateEngine, err = tmpl.Parse(metricsForServiceMappingConfFileContentTemplate); err != nil {
		errCause := "error parsing metrics for services config: " + err.Error()
		return common.ErrorFromThisScope(errCause, generalScopeErr)
	}
	if err = creteFullPathForFile(confData.metricsForServicesPath()); err != nil {
		errCause := "error creating directory for metrics for services file: " + err.Error()
		return common.ErrorFromThisScope(errCause, generalScopeErr)
	}
	var metricsForServiceConfigFile *os.File
	if metricsForServiceConfigFile, err = os.Create(confData.metricsForServicesPath()); err != nil {
		errCause := "error creating metrics for services config file: " + err.Error()
		return common.ErrorFromThisScope(errCause, generalScopeErr)
	}
	if err = templateEngine.Execute(metricsForServiceConfigFile, confData.metricsForServiceConfigTmplData.TmplData.ServiceNameToMetricsConfPath); err != nil {
		errCause := "error writing metrics for services config file: " + err.Error()
		return common.ErrorFromThisScope(errCause, generalScopeErr)
	}
	for key, val := range confData.metricsForServiceConfigTmplData.TmplData.ServiceNameToMetricsConfPath {
		if err = createMetricsConfigFile(val); err != nil {
			errCause := fmt.Sprintf("error creating metrics config file for service %s: %s", key, err.Error())
			return common.ErrorFromThisScope(errCause, generalScopeErr)
		}
	}
	return err
}

func (confData mainConfigData) existMainConfigFile() bool {
	return existFile(confData.MainConfigPath())
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

func serviceDefaultConfigPath(conf *configdir.Config) (path string) {
	return fileDefaultConfigPath(serviceConfigFileName, conf)
}

func mainDefaultConfigPath(conf *configdir.Config) (path string) {
	return fileDefaultConfigPath(mainConfigFileName, conf)
}

func metricsForServicesDefaultConfigPath(conf *configdir.Config) (path string) {
	return fileDefaultConfigPath(metricsForServicesConfigFileName, conf)
}

func skycoinMetricsConfigPath(conf *configdir.Config) (path string) {
	return fileDefaultConfigPath(skycoinMetricsConfigFileName, conf)
}

func walletMetricsConfigPath(conf *configdir.Config) (path string) {
	return fileDefaultConfigPath(walletMetricsConfigFileName, conf)
}

func defaultTmplData(conf *configdir.Config) (tmplData templateData) {
	tmplData = templateData{
		ServiceConfigPath:      serviceDefaultConfigPath(conf),
		MetricsForServicesPath: metricsForServicesDefaultConfigPath(conf),
	}
	return tmplData
}

func defaultMetricsForServiceTmplData(conf *configdir.Config) (tmplData metricsForServiceConfigTemplateData) {
	tmplData = metricsForServiceConfigTemplateData{
		TmplData: metricsForServiceTemplateData{
			ServiceNameToMetricsConfPath: map[string]string{
				"skycoin": skycoinMetricsConfigPath(conf),
				"wallet":  walletMetricsConfigPath(conf),
			},
		},
	}
	return tmplData
}

func tmplDataFromMainFile(mainConfigFilePath string) (tmpl templateData, err error) {
	generalScopeErr := "error filling template data"
	viper.SetConfigFile(mainConfigFilePath)
	viper.SetConfigType("toml")
	if err := viper.ReadInConfig(); err != nil {
		errCause := fmt.Sprintln("error reading config file: ", mainConfigFilePath, err.Error())
		return tmpl, common.ErrorFromThisScope(errCause, generalScopeErr)
	}
	var mainConf templateData
	if err := viper.Unmarshal(&mainConf); err != nil {
		errCause := fmt.Sprintln("can not decode the config data: ", err.Error())
		return tmpl, common.ErrorFromThisScope(errCause, generalScopeErr)
	}
	tmpl.ServiceConfigPath, tmpl.MetricsForServicesPath = mainConf.ServiceConfigPath, mainConf.MetricsForServicesPath
	return tmpl, err
}

func (tmpl templateData) metricsForServicesTmplDataFromFile() (metricsForServicesTmpl metricsForServiceConfigTemplateData, err error) {
	generalScopeErr := "error filling template data"
	viper.SetConfigFile(tmpl.MetricsForServicesPath)
	viper.SetConfigType("toml")
	if err := viper.ReadInConfig(); err != nil {
		errCause := fmt.Sprintln("error reading config file: ", tmpl.MetricsForServicesPath, err.Error())
		return metricsForServicesTmpl, common.ErrorFromThisScope(errCause, generalScopeErr)
	}
	if err := viper.Unmarshal(&(metricsForServicesTmpl.TmplData)); err != nil {
		errCause := fmt.Sprintln("can not decode the config data: ", err.Error())
		return metricsForServicesTmpl, common.ErrorFromThisScope(errCause, generalScopeErr)
	}
	return metricsForServicesTmpl, err
}

func metricsForServicesTmplData(conf *configdir.Config) metricsForServiceConfigTemplateData {
	return defaultMetricsForServiceTmplData(conf)
}

func newMainConfigData(path string) (mainConf mainConfigData, err error) {
	generalScopeErr := "can not create main config instance"
	if isADirectoryPath(path) {
		path = filepath.Join(path, mainConfigFileName)
	}
	var tmplData templateData
	var metricsForServiceTmplData metricsForServiceConfigTemplateData
	if strings.Compare(path, "") == 0 || !existFile(path) {
		// TODO(denisacostaq@gmail.com): move homeConf to fn defaultTmplData
		var homeConf *configdir.Config
		if homeConf, err = homeConfigFolder(); err != nil {
			errCause := "error looking for config folder under home: " + err.Error()
			return mainConf, common.ErrorFromThisScope(errCause, generalScopeErr)
		}
		path = mainDefaultConfigPath(homeConf)
		tmplData = defaultTmplData(homeConf)
		metricsForServiceTmplData = metricsForServicesTmplData(homeConf)
	} else {
		if tmplData, err = tmplDataFromMainFile(path); err != nil {
			errCause := "error reading template data from file: " + err.Error()
			return mainConf, common.ErrorFromThisScope(errCause, generalScopeErr)
		}
		if metricsForServiceTmplData, err = tmplData.metricsForServicesTmplDataFromFile(); err != nil {
			errCause := "error reading template data from file: " + err.Error()
			return mainConf, common.ErrorFromThisScope(errCause, generalScopeErr)
		}
	}
	if strings.Compare(tmplData.ServiceConfigPath, "") == 0 || strings.Compare(tmplData.MetricsForServicesPath, "") == 0 {
		var homeConf *configdir.Config
		if homeConf, err = homeConfigFolder(); err != nil {
			errCause := "error looking for config folder under home: " + err.Error()
			return mainConf, common.ErrorFromThisScope(errCause, generalScopeErr)
		}
		tmpTmplData := defaultTmplData(homeConf)
		if strings.Compare(tmplData.ServiceConfigPath, "") == 0 {
			tmplData.ServiceConfigPath = tmpTmplData.ServiceConfigPath
		}
		if strings.Compare(tmplData.MetricsForServicesPath, "") == 0 {
			tmplData.MetricsForServicesPath = tmpTmplData.MetricsForServicesPath
		}
		metricsForServiceTmplData = metricsForServicesTmplData(homeConf)
	}
	mainConf = mainConfigData{
		mainConfigPath:                  path,
		tmplData:                        tmplData,
		metricsForServiceConfigTmplData: metricsForServiceTmplData,
	}
	if err = mainConf.createMainConfigFile(); err != nil {
		errCause := "error creating main config file: " + err.Error()
		return mainConf, common.ErrorFromThisScope(errCause, generalScopeErr)
	}
	if err = mainConf.createServiceConfigFile(); err != nil {
		errCause := "error creating service config file: " + err.Error()
		return mainConf, common.ErrorFromThisScope(errCause, generalScopeErr)
	}
	if err = mainConf.createMetricsForServicesConfFile(); err != nil {
		errCause := "error creating metrics for services config file: " + err.Error()
		return mainConf, common.ErrorFromThisScope(errCause, generalScopeErr)
	}

	return mainConf, err
}
