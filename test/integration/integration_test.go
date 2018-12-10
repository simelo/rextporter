package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/alecthomas/template"
	"github.com/simelo/rextporter/src/config"
	"github.com/simelo/rextporter/src/exporter"
	"github.com/simelo/rextporter/src/util"
	"github.com/simelo/rextporter/test/integration/testrand"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

const mainConfigFileContenTemplate = `
serviceConfigTransport = "file"
# render a template with a portable path
servicesConfigPath = "{{.ServicesConfigPath}}"
metricsForServicesPath = "{{.MetricsForServicesPath}}"
`

const servicesConfigFileContenTemplate = `
	# Service configuration.{{range .Services}}
	[[services]]
		name = "{{.Name}}"
		metricsToForwardPath = "{{.ForwardPath}}"
		modes=[{{range .Modes}}"{{.}}" {{end}}]
		scheme = "http"
		port = {{.Port}}
		basePath = "{{.BasePath}}"
		authType = "CSRF"
		tokenHeaderKey = "X-CSRF-Token"
		genTokenEndpoint = "/api/v1/csrf"
		tokenKeyFromEndpoint = "csrf_token"
		
		[services.location]
			location = "localhost"

{{end}}
`

const metricsConfigFileContenTemplate = `
# All metrics to be measured.
[[metrics]]
	name = "open_connections_is_a_fake_name_for_test_purpose"
	url = "/api/v1/health"
	httpMethod = "GET"
	path = "open_connections"

	[metrics.options]
		type = "Gauge"
		description = "Track the open connections in the system"	
`

const metricsForServicesConfFileContenTemplate = `
	serviceNameToMetricsConfPath = [{{range $key, $value := .}}
	{ {{$key}} = "{{$value}}" },{{end}}
]
`

type Service struct {
	Name        string
	Port        uint16
	ForwardPath string
	Modes       []string
	BasePath    string
}

type ServicesConfData struct {
	Services []Service
}

type HealthSuit struct {
	suite.Suite
	require                          *require.Assertions
	mainConfFilePath                 string
	mainConfTmplContent              string
	servicesConfFilePath             string
	servicesConfData                 ServicesConfData
	metricsConfTmplContent           string
	metricsConfFilePath              string
	metricsForServiceConfTmplContent string
	metricsForServicesConfData       map[string]string
	metricsForServicesConfFilePath   string
}

var fakeNodePort uint16

func createConfigFile(tmplContent, path string, data interface{}) (err error) {
	generalScopeErr := "error creating config file for integration test"
	if len(tmplContent) == 0 || len(path) == 0 {
		return err
	}
	tmpl := template.New("fileConfig")
	var templateEngine *template.Template
	if templateEngine, err = tmpl.Parse(tmplContent); err != nil {
		errCause := "error parsing config: " + err.Error()
		return util.ErrorFromThisScope(errCause, generalScopeErr)
	}
	var configFile *os.File
	if configFile, err = os.Create(path); err != nil {
		errCause := "error creating config file: " + err.Error()
		return util.ErrorFromThisScope(errCause, generalScopeErr)
	}
	if err = templateEngine.Execute(configFile, data); err != nil {
		errCause := "error writing config file: " + err.Error()
		return util.ErrorFromThisScope(errCause, generalScopeErr)
	}
	return err
}

func (suite *HealthSuit) createServicesConfPath() (err error) {
	generalScopeErr := "error creating service config file for integration test"
	if err = createConfigFile(servicesConfigFileContenTemplate, suite.servicesConfFilePath, suite.servicesConfData); err != nil {
		errCause := "error writing service config file: " + err.Error()
		return util.ErrorFromThisScope(errCause, generalScopeErr)
	}
	return err
}

func (suite *HealthSuit) createMainConfPath(tmplData interface{}) (err error) {
	generalScopeErr := "error creating main config file for integration test"
	if err = createConfigFile(suite.mainConfTmplContent, suite.mainConfFilePath, tmplData); err != nil {
		errCause := "error writing service config file: " + err.Error()
		return util.ErrorFromThisScope(errCause, generalScopeErr)
	}
	return err
}

func (suite *HealthSuit) createMetricsForServicesConfPath() (err error) {
	return createConfigFile(
		suite.metricsForServiceConfTmplContent,
		suite.metricsForServicesConfFilePath,
		suite.metricsForServicesConfData)
}

func (suite *HealthSuit) createMetricsConfigPaths() (err error) {
	return createConfigFile(
		suite.metricsConfTmplContent,
		suite.metricsConfFilePath,
		nil)
}

func (suite *HealthSuit) createMainConfig() {
	generalScopeErr := "error creating main config file for integration test"
	type mainConfigData struct {
		ServicesConfigPath     string
		MetricsForServicesPath string
	}
	confData := mainConfigData{
		ServicesConfigPath:     suite.servicesConfFilePath,
		MetricsForServicesPath: suite.metricsForServicesConfFilePath,
	}
	if err := suite.createMainConfPath(confData); err != nil {
		errCause := "error writing main config file: " + err.Error()
		suite.Nil(util.ErrorFromThisScope(errCause, generalScopeErr))
	}
	if err := suite.createServicesConfPath(); err != nil {
		errCause := "error writing services config file: " + err.Error()
		suite.Nil(util.ErrorFromThisScope(errCause, generalScopeErr))
	}
	if err := suite.createMetricsConfigPaths(); err != nil {
		errCause := "error writing my monitored server metrics config file: " + err.Error()
		suite.Nil(util.ErrorFromThisScope(errCause, generalScopeErr))
	}
	if err := suite.createMetricsForServicesConfPath(); err != nil {
		errCause := "error writing metrics for service config file: " + err.Error()
		suite.Nil(util.ErrorFromThisScope(errCause, generalScopeErr))
	}
}

func (suite *HealthSuit) createDirectoriesWithFullDepth(dirs []string) {
	for _, dir := range dirs {
		suite.Nil(os.MkdirAll(dir, 0750))
	}
}

func readListenPortFromFile() (port uint16, err error) {
	var path string
	path, err = testrand.FilePathToSharePort()
	var file *os.File
	file, err = os.OpenFile(path, os.O_RDONLY, 0400)
	if err != nil {
		log.WithError(err).Errorln("error opening file")
		return 0, err
	}
	defer file.Close()
	_, err = fmt.Fscanf(file, "%d", &port)
	if err != nil {
		log.WithError(err).Errorln("error reading file")
		return port, err
	}
	return port, err
}

func TestSkycoinHealthSuit(t *testing.T) {
	suite.Run(t, new(HealthSuit))
}

func (suite *HealthSuit) SetupSuite() {
	require := require.New(suite.T())
	var port uint16
	var err error
	port, err = readListenPortFromFile()
	require.Nil(err)
	fakeNodePort = port
}

func (suite *HealthSuit) SetupTest() {
	suite.callSetUpTest()
}

func (suite *HealthSuit) callSetUpTest() {
	suite.metricsForServicesConfData =
		map[string]string{
			"myMonitoredServer":        suite.metricsConfFilePath,
			"myMonitoredAsProxyServer": suite.metricsConfFilePath}
	suite.servicesConfData = ServicesConfData{
		Services: []Service{Service{Name: "myMonitoredServer", Port: fakeNodePort, Modes: []string{"rest_api"}, BasePath: ""}},
	}
}

func (suite *HealthSuit) TestMetricMonitorHealth() {
	// NOTE(denisacostaq@gmail.com): Giving
	suite.require = require.New(suite.T())
	mainConfigDir := testrand.RFolderPath()
	servicesDir := testrand.RFolderPath()
	myMonitoredServerMetricsDir := testrand.RFolderPath()
	metricsForServicesDir := testrand.RFolderPath()
	port := testrand.RandomPort()
	suite.createDirectoriesWithFullDepth([]string{mainConfigDir, servicesDir, myMonitoredServerMetricsDir, metricsForServicesDir})
	suite.mainConfFilePath = filepath.Join(mainConfigDir, testrand.RName())
	suite.servicesConfFilePath = filepath.Join(servicesDir, testrand.RName())
	suite.metricsConfFilePath = filepath.Join(myMonitoredServerMetricsDir, testrand.RName())
	suite.metricsForServicesConfFilePath = filepath.Join(metricsForServicesDir, testrand.RName())
	suite.mainConfTmplContent = mainConfigFileContenTemplate
	suite.metricsConfTmplContent = metricsConfigFileContenTemplate
	suite.metricsForServiceConfTmplContent = metricsForServicesConfFileContenTemplate
	suite.callSetUpTest()
	suite.createMainConfig()
	conf := config.MustConfigFromFileSystem(suite.mainConfFilePath)
	srv := exporter.MustExportMetrics("/metrics2", port, conf)
	suite.require.NotNil(srv)
	// NOTE(denisacostaq@gmail.com): Wait for server starts
	time.Sleep(time.Second * 2)

	// NOTE(denisacostaq@gmail.com): When
	resp, err := http.Get(fmt.Sprintf("http://127.0.0.1:%d/metrics2", port))

	// NOTE(denisacostaq@gmail.com): Assert
	defer func() { suite.Nil(resp.Body.Close()) }()
	suite.Nil(err)
	suite.Equal(http.StatusOK, resp.StatusCode)
	suite.Len(conf.Services, 1)
	suite.Len(conf.Services[0].Metrics, 1)
	metricName := conf.Services[0].Metrics[0].Name
	suite.Equal(metricName, "open_connections_is_a_fake_name_for_test_purpose")
	var usingAVariableToMakeLinterHappy = context.Context(nil)
	suite.Nil(srv.Shutdown(usingAVariableToMakeLinterHappy))
}

// func (suite *HealthSuit) TestMetricMonitorAsProxy() {
// 	// NOTE(denisacostaq@gmail.com): Giving
// 	suite.require = require.New(suite.T())
// 	port := testrand.RandomPort()
// 	mainConfigDir := testrand.RFolderPath()
// 	servicesDir := testrand.RFolderPath()
// 	myMonitoredServerMetricsDir := testrand.RFolderPath()
// 	metricsForServicesDir := testrand.RFolderPath()
// 	suite.createDirectoriesWithFullDepth([]string{mainConfigDir, servicesDir, myMonitoredServerMetricsDir, metricsForServicesDir})
// 	suite.mainConfFilePath = filepath.Join(mainConfigDir, testrand.RName())
// 	suite.servicesConfFilePath = filepath.Join(servicesDir, testrand.RName())
// 	suite.metricsConfFilePath = filepath.Join(myMonitoredServerMetricsDir, testrand.RName())
// 	suite.metricsForServicesConfFilePath = filepath.Join(metricsForServicesDir, testrand.RName())
// 	suite.mainConfTmplContent = mainConfigFileContenTemplate
// 	suite.metricsConfTmplContent = metricsConfigFileContenTemplate
// 	suite.metricsForServiceConfTmplContent = metricsForServicesConfFileContenTemplate
// 	suite.callSetUpTest()
// 	suite.metricsForServicesConfData =
// 		map[string]string{
// 			"myMonitoredServer":        suite.metricsConfFilePath,
// 			"myMonitoredAsProxyServer": suite.metricsConfFilePath}
// 	suite.servicesConfData = ServicesConfData{
// 		Services: []Service{Service{
// 			Name:        "myMonitoredAsProxyServer",
// 			Port:        fakeNodePort,
// 			Modes:       []string{"forward_metrics"},
// 			ForwardPath: "/metrics"},
// 		},
// 	}
// 	suite.createMainConfig()
// 	conf := config.MustConfigFromFileSystem(suite.mainConfFilePath)
// 	srv := exporter.MustExportMetrics("/metrics4", port, conf)
// 	suite.require.NotNil(srv)

// 	// NOTE(denisacostaq@gmail.com): Wait for server starts
// 	time.Sleep(time.Second * 2)

// 	// NOTE(denisacostaq@gmail.com): When
// 	var resp *http.Response
// 	resp, err := http.Get(fmt.Sprintf("http://127.0.0.1:%d/metrics4", port))
// 	suite.require.NotNil(resp)

// 	// NOTE(denisacostaq@gmail.com): Assert
// 	suite.Nil(err)
// 	defer func() { suite.Nil(resp.Body.Close()) }()
// 	suite.Equal(http.StatusOK, resp.StatusCode)
// 	suite.require.Len(conf.Services, 1)
// 	suite.require.Len(conf.Services[0].Metrics, 0)
// 	metricName := conf.Services[0].Name + "_skycoin_wallet2_seq2"
// 	suite.require.Equal(metricName, "myMonitoredAsProxyServer_skycoin_wallet2_seq2")
// 	var usingAVariableToMakeLinterHappy = context.Context(nil)
// 	suite.require.Nil(srv.Shutdown(usingAVariableToMakeLinterHappy))
// }

// func (suite *HealthSuit) TestMetricMonitorAsProxyWithNonMetricsEndpoint() {
// 	// NOTE(denisacostaq@gmail.com): Giving
// 	suite.require = require.New(suite.T())
// 	port := testrand.RandomPort()
// 	mainConfigDir := testrand.RFolderPath()
// 	servicesDir := testrand.RFolderPath()
// 	myMonitoredServerMetricsDir := testrand.RFolderPath()
// 	metricsForServicesDir := testrand.RFolderPath()
// 	suite.createDirectoriesWithFullDepth([]string{mainConfigDir, servicesDir, myMonitoredServerMetricsDir, metricsForServicesDir})
// 	suite.mainConfFilePath = filepath.Join(mainConfigDir, testrand.RName())
// 	suite.servicesConfFilePath = filepath.Join(servicesDir, testrand.RName()+".toml")
// 	suite.metricsConfFilePath = filepath.Join(myMonitoredServerMetricsDir, testrand.RName()+".toml")
// 	suite.metricsForServicesConfFilePath = filepath.Join(metricsForServicesDir, testrand.RName()+".toml")
// 	suite.mainConfTmplContent = mainConfigFileContenTemplate
// 	suite.callSetUpTest()
// 	suite.metricsForServicesConfData =
// 		map[string]string{
// 			"myMonitoredAsProxyServer": suite.metricsConfFilePath}
// 	suite.servicesConfData = ServicesConfData{
// 		Services: []Service{Service{
// 			Name:        "myMonitoredAsProxyServer",
// 			Port:        fakeNodePort,
// 			Modes:       []string{"forward_metrics"},
// 			BasePath:    "/api/v1/health",
// 			ForwardPath: "/metrics"},
// 		},
// 	}
// 	suite.metricsConfTmplContent = metricsConfigFileContenTemplate
// 	suite.metricsForServiceConfTmplContent = metricsForServicesConfFileContenTemplate
// 	suite.createMainConfig()
// 	conf := config.MustConfigFromFileSystem(suite.mainConfFilePath)
// 	srv := exporter.MustExportMetrics("/metrics5", port, conf)
// 	suite.require.NotNil(srv)
// 	// NOTE(denisacostaq@gmail.com): Wait for server starts
// 	time.Sleep(time.Second * 2)

// 	// NOTE(denisacostaq@gmail.com): When
// 	resp, err := http.Get(fmt.Sprintf("http://127.0.0.1:%d/metrics5", port))

// 	// NOTE(denisacostaq@gmail.com): Assert
// 	suite.Nil(err)
// 	defer func() { suite.Nil(resp.Body.Close()) }()
// 	suite.Equal(http.StatusOK, resp.StatusCode)
// 	suite.require.Len(conf.Services, 1)
// 	suite.require.Len(conf.Services[0].Metrics, 0)
// 	metricName := "skycoin_wallet2_seq2"
// 	suite.Equal(metricName, "skycoin_wallet2_seq2")
// 	var usingAVariableToMakeLinterHappy = context.Context(nil)
// 	suite.require.Nil(srv.Shutdown(usingAVariableToMakeLinterHappy))
// }
