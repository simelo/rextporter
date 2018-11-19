package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/alecthomas/template"
	"github.com/simelo/rextporter/src/exporter"
	"github.com/simelo/rextporter/src/util"
	"github.com/simelo/rextporter/test/integration/testrand"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type HealthSuit struct {
	suite.Suite
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

func createServicesConfPath(servicesConfPath string) (err error) {
	const servicesConfigFileContenTemplate = `
	# Service configuration.
	[[services]]
		name = "myMonitoredServer"
		scheme = "http"
		port = 8080
		basePath = ""
		authType = "CSRF"
		tokenHeaderKey = "X-CSRF-Token"
		genTokenEndpoint = "/api/v1/csrf"
		tokenKeyFromEndpoint = "csrf_token"
		
		[services.location]
			location = "localhost"
`
	generalScopeErr := "error creating service config file for integration test"
	if err = createConfigFile(servicesConfigFileContenTemplate, servicesConfPath, nil); err != nil {
		errCause := "error writing service config file: " + err.Error()
		return util.ErrorFromThisScope(errCause, generalScopeErr)
	}
	return err
}

func createMetricsForServicesConfPath(metricsForServicesConfPath string, metricsForServices map[string]string) (err error) {
	const metricsForServicesConfFileContenTemplate = `
	serviceNameToMetricsConfPath = [{{range $key, $value := .}}
	{ {{$key}} = "{{$value}}" },{{end}}
]
`
	return createConfigFile(
		metricsForServicesConfFileContenTemplate,
		metricsForServicesConfPath,
		metricsForServices)
}

func createMetricsConfig(tmplContent, path string) (err error) {
	generalScopeErr := "error creating metrics config file for integration test"
	if err = createConfigFile(tmplContent, path, nil); err != nil {
		errCause := "error writing metrics config file: " + err.Error()
		return util.ErrorFromThisScope(errCause, generalScopeErr)
	}
	return err
}

func createMetricsConfigPaths(metricsConfigPath string) (err error) {
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
	return createMetricsConfig(metricsConfigFileContenTemplate, metricsConfigPath)
}

func createMainConfig(tmplContent, mainConfPath, servicesConfigPath, metricsForServicesConfPath, myMonitoredServerMetricsPath string) (err error) {
	generalScopeErr := "error creating main config file for integration test"
	type mainConfigData struct {
		ServicesConfigPath     string
		MetricsForServicesPath string
	}
	confData := mainConfigData{
		ServicesConfigPath:     servicesConfigPath,
		MetricsForServicesPath: metricsForServicesConfPath,
	}
	if err = createConfigFile(tmplContent, mainConfPath, confData); err != nil {
		errCause := "error writing main config file: " + err.Error()
		return util.ErrorFromThisScope(errCause, generalScopeErr)
	}
	if err = createServicesConfPath(servicesConfigPath); err != nil {
		errCause := "error writing services config file: " + err.Error()
		return util.ErrorFromThisScope(errCause, generalScopeErr)
	}
	if err = createMetricsConfigPaths(myMonitoredServerMetricsPath); err != nil {
		errCause := "error writing my monitored server metrics config file: " + err.Error()
		return util.ErrorFromThisScope(errCause, generalScopeErr)
	}
	if err = createMetricsForServicesConfPath(metricsForServicesConfPath, map[string]string{"myMonitoredServer": myMonitoredServerMetricsPath}); err != nil {
		errCause := "error writing metrics for service config file: " + err.Error()
		return util.ErrorFromThisScope(errCause, generalScopeErr)
	}
	return err
}

func createMainConfigTestPaths() (mainConfFilePath string, err error) {
	const mainConfigFileContenTemplate = `
serviceConfigTransport = "file"
# render a template with a portable path
serviceConfigPath = "{{.ServicesConfigPath}}"
metricsForServicesPath = "{{.MetricsForServicesPath}}"
`
	mainConfigDir := testrand.RFolderPath()
	if err = os.MkdirAll(mainConfigDir, 0750); err != nil {
		return mainConfFilePath, err
	}
	mainConfFilePath = filepath.Join(mainConfigDir, testrand.RName())
	servicesDir := testrand.RFolderPath()
	if err = os.MkdirAll(servicesDir, 0750); err != nil {
		return mainConfFilePath, err
	}
	servicesConfPath := filepath.Join(servicesDir, testrand.RName())
	metricsForServicesDir := testrand.RFolderPath()
	if err = os.MkdirAll(metricsForServicesDir, 0750); err != nil {
		return mainConfFilePath, err
	}
	metricsForServicesConfPath := filepath.Join(metricsForServicesDir, testrand.RName()+".toml")
	metricsDir := testrand.RFolderPath()
	if err = os.MkdirAll(metricsDir, 0750); err != nil {
		return mainConfFilePath, err
	}
	myMonitoredServerMetricsDir := testrand.RFolderPath()
	if err = os.MkdirAll(myMonitoredServerMetricsDir, 0750); err != nil {
		return mainConfFilePath, err
	}
	myMonitoredServerMetricsPath := filepath.Join(myMonitoredServerMetricsDir, testrand.RName()+".toml")
	return mainConfFilePath,
		createMainConfig(
			mainConfigFileContenTemplate,
			mainConfFilePath,
			servicesConfPath,
			metricsForServicesConfPath,
			myMonitoredServerMetricsPath)
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

func (suite *HealthSuit) TestMetricMonitorHealth() {
	// NOTE(denisacostaq@gmail.com): Giving
	require := require.New(suite.T())
	mainConfFilePath, err := createMainConfigTestPaths()
	require.Nil(err)
	port := testrand.RandomPort()
	srv := exporter.ExportMetrics(mainConfFilePath, "/metrics", port)
	require.NotNil(srv)
	// NOTE(denisacostaq@gmail.com): Wait for server starts
	t := time.NewTimer(time.Second * 2)
	<-t.C

	// NOTE(denisacostaq@gmail.com): When
	var resp *http.Response
	resp, err = http.Get(fmt.Sprintf("http://127.0.0.1:%d/metrics", port))

	// NOTE(denisacostaq@gmail.com): Assert
	suite.Nil(err)
	var data []byte
	data, err = ioutil.ReadAll(resp.Body)
	suite.Nil(err)
	suite.Equal(http.StatusOK, resp.StatusCode)
	suite.Contains(string(data), "open_connections_is_a_fake_name_for_test_purpose")
	require.Nil(srv.Shutdown(context.Context(nil)))
}

func (suite *HealthSuit) TestConfigWorks() {
	// NOTE(denisacostaq@gmail.com): Giving
	require := require.New(suite.T())
	port := testrand.RandomPort()
	srv := exporter.ExportMetrics("", "/metrics2", port)
	require.NotNil(srv)
	// NOTE(denisacostaq@gmail.com): Wait for server starts
	t := time.NewTimer(time.Second * 2)
	<-t.C

	// // NOTE(denisacostaq@gmail.com): When
	var resp *http.Response
	resp, err := http.Get(fmt.Sprintf("http://127.0.0.1:%d/metrics2", port))
	log.Println("resp, err", resp, err)
	// // NOTE(denisacostaq@gmail.com): Assert
	suite.Nil(err)
	suite.Equal(http.StatusOK, resp.StatusCode)
	require.Nil(srv.Shutdown(context.Context(nil)))
}
