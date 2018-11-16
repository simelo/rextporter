package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/shibukawa/configdir"
	"github.com/simelo/rextporter/src/common"
	"github.com/simelo/rextporter/src/util/file"
	"github.com/spf13/viper"
)

type templateData struct {
	ServiceConfigPath string
	MetricsConfigPath string
}

type mainConfigData struct {
	mainConfigPath string
	tmplData       templateData
}

func (confData mainConfigData) ServiceConfigPath() string {
	return confData.tmplData.ServiceConfigPath
}

func (confData mainConfigData) MetricsConfigPath() string {
	return confData.tmplData.MetricsConfigPath
}

func (confData mainConfigData) MainConfigPath() string {
	return confData.mainConfigPath
}

const mainConfigFileContentTemplate = `
serviceConfigTransport = "file" # "file" | "consulCatalog"
# render a template with a portable path
serviceConfigPath = "{{.ServiceConfigPath}}"
metricsConfigPath = "{{.MetricsConfigPath}}"
`

const serviceConfigFileContentTemplate = `
# Service configuration.
name = "wallet"
scheme = "http"
port = 8000
basePath = ""
authType = "CSRF"
tokenHeaderKey = "X-CSRF-Token"
genTokenEndpoint = "/api/v1/csrf.json"
tokenKeyFromEndpoint = "csrf_token"

[location]
  location = "localhost"
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
# from https://github.com/simelo/rextporter/pull/17


@denisacostaq services should be match against whole templates , rather than individual metrics. 
The match is not for hosts directly . The match is for service types . Works like this

metrics <- m:1 -> templates <- m:n -> services <- 1:n -> (physical | virtual) hosts
`

var (
	systemVendorName      = "simelo"
	systemProgramName     = "rextporter"
	mainConfigFileName    = "main.toml"
	serviceConfigFileName = "service.toml"
	metricsConfigFileName = "metrics.toml"
)

func (confData mainConfigData) existServiceConfigFile() bool {
	return file.ExistFile(confData.ServiceConfigPath())
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
	if err = file.CreteFullPathForFile(confData.ServiceConfigPath()); err != nil {
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

func (confData mainConfigData) existMetricsConfigFile() bool {
	return file.ExistFile(confData.tmplData.MetricsConfigPath)
}

// createMetricsConfigFile creates the metrics file or return an error if any,
// if the file already exist does no thin.
func (confData mainConfigData) createMetricsConfigFile() (err error) {
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
	if err = file.CreteFullPathForFile(confData.MetricsConfigPath()); err != nil {
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
	return file.ExistFile(confData.MetricsConfigPath())
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
	if err = file.CreteFullPathForFile(confData.MainConfigPath()); err != nil {
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

func metricsDefaultConfigPath(conf *configdir.Config) (path string) {
	return file.DefaultConfigPath(metricsConfigFileName, conf)
}

func serviceDefaultConfigPath(conf *configdir.Config) (path string) {
	return file.DefaultConfigPath(serviceConfigFileName, conf)
}

func mainDefaultConfigPath(conf *configdir.Config) (path string) {
	return file.DefaultConfigPath(mainConfigFileName, conf)
}

func defaultTmplData(conf *configdir.Config) (tmplData templateData) {
	tmplData = templateData{
		ServiceConfigPath: serviceDefaultConfigPath(conf),
		MetricsConfigPath: metricsDefaultConfigPath(conf),
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
	tmpl.ServiceConfigPath, tmpl.MetricsConfigPath = mainConf.ServiceConfigPath, mainConf.MetricsConfigPath
	return tmpl, err
}

func newMainConfigData(path string) (mainConf mainConfigData, err error) {
	generalScopeErr := "can not create main config instance"
	if file.IsADirectoryPath(path) {
		path = filepath.Join(path, mainConfigFileName)
	}
	var tmplData templateData
	if strings.Compare(path, "") == 0 || !file.ExistFile(path) {
		// TODO(denisacostaq@gmail.com): move homeConf to fn defaultTmplData
		var homeConf *configdir.Config
		if homeConf, err = file.HomeConfigFolder(systemVendorName, systemProgramName); err != nil {
			errCause := "error looking for config folder under home: " + err.Error()
			return mainConf, common.ErrorFromThisScope(errCause, generalScopeErr)
		}
		path = mainDefaultConfigPath(homeConf)
		tmplData = defaultTmplData(homeConf)
	} else {
		if tmplData, err = tmplDataFromMainFile(path); err != nil {
			errCause := "error reading template data from file: " + err.Error()
			return mainConf, common.ErrorFromThisScope(errCause, generalScopeErr)
		}
	}
	if strings.Compare(tmplData.ServiceConfigPath, "") == 0 || strings.Compare(tmplData.MetricsConfigPath, "") == 0 {
		var homeConf *configdir.Config
		if homeConf, err = file.HomeConfigFolder(systemVendorName, systemProgramName); err != nil {
			errCause := "error looking for config folder under home: " + err.Error()
			return mainConf, common.ErrorFromThisScope(errCause, generalScopeErr)
		}
		tmpTmplData := defaultTmplData(homeConf)
		if strings.Compare(tmplData.ServiceConfigPath, "") == 0 {
			tmplData.ServiceConfigPath = tmpTmplData.ServiceConfigPath
		}
		if strings.Compare(tmplData.MetricsConfigPath, "") == 0 {
			tmplData.MetricsConfigPath = tmpTmplData.MetricsConfigPath
		}
	}
	mainConf = mainConfigData{
		mainConfigPath: path,
		tmplData:       tmplData,
	}
	if err = mainConf.createMainConfigFile(); err != nil {
		errCause := "error creating main config file: " + err.Error()
		return mainConf, common.ErrorFromThisScope(errCause, generalScopeErr)
	}
	if err = mainConf.createMetricsConfigFile(); err != nil {
		errCause := "error creating metrics config file: " + err.Error()
		return mainConf, common.ErrorFromThisScope(errCause, generalScopeErr)
	}
	if err = mainConf.createServiceConfigFile(); err != nil {
		errCause := "error creating service config file: " + err.Error()
		return mainConf, common.ErrorFromThisScope(errCause, generalScopeErr)
	}
	return mainConf, err
}
