package main

import (
	"context"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/alecthomas/template"
	"github.com/simelo/rextporter/src/common"
	"github.com/simelo/rextporter/src/exporter"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type HealthSuit struct {
	suite.Suite
}

func createConfigFile(tmplContent, path string, data interface{}) (err error) {
	generalScopeErr := "error creating config file for integration test"
	if strings.Compare(tmplContent, "") == 0 || strings.Compare(path, "") == 0 {
		return err
	}
	tmpl := template.New("fileConfig")
	var templateEngine *template.Template
	if templateEngine, err = tmpl.Parse(tmplContent); err != nil {
		errCause := "error parsing config: " + err.Error()
		return common.ErrorFromThisScope(errCause, generalScopeErr)
	}
	var configFile *os.File
	if configFile, err = os.Create(path); err != nil {
		errCause := "error creating config file: " + err.Error()
		return common.ErrorFromThisScope(errCause, generalScopeErr)
	}
	if err = templateEngine.Execute(configFile, data); err != nil {
		errCause := "error writing config file: " + err.Error()
		return common.ErrorFromThisScope(errCause, generalScopeErr)
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
		return common.ErrorFromThisScope(errCause, generalScopeErr)
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
		return common.ErrorFromThisScope(errCause, generalScopeErr)
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
		ServiceConfigPath      string
		MetricsForServicesPath string
	}
	confData := mainConfigData{
		ServiceConfigPath:      servicesConfigPath,
		MetricsForServicesPath: metricsForServicesConfPath,
	}
	if err = createConfigFile(tmplContent, mainConfPath, confData); err != nil {
		errCause := "error writing main config file: " + err.Error()
		return common.ErrorFromThisScope(errCause, generalScopeErr)
	}
	if err = createServicesConfPath(servicesConfigPath); err != nil {
		errCause := "error writing services config file: " + err.Error()
		return common.ErrorFromThisScope(errCause, generalScopeErr)
	}
	if err = createMetricsConfigPaths(myMonitoredServerMetricsPath); err != nil {
		errCause := "error writing my monitored server metrics config file: " + err.Error()
		return common.ErrorFromThisScope(errCause, generalScopeErr)
	}
	if err = createMetricsForServicesConfPath(metricsForServicesConfPath, map[string]string{"myMonitoredServer": myMonitoredServerMetricsPath}); err != nil {
		errCause := "error writing metrics for service config file: " + err.Error()
		return common.ErrorFromThisScope(errCause, generalScopeErr)
	}
	return err
}

func createMainConfigTestPaths() (mainConfFilePath string, err error) {
	const mainConfigFileContenTemplate = `
serviceConfigTransport = "file"
# render a template with a portable path
serviceConfigPath = "{{.ServiceConfigPath}}"
metricsForServicesPath = "{{.MetricsForServicesPath}}"
`
	mainConfigDir := filepath.Join(os.TempDir(), "sdsds", "675656", "aa")
	if err = os.MkdirAll(mainConfigDir, 0750); err != nil {
		return mainConfFilePath, err
	}
	mainConfFilePath = filepath.Join(mainConfigDir, "rrrr")
	servicesDir := filepath.Join(os.TempDir(), "test", "integration")
	if err = os.MkdirAll(servicesDir, 0750); err != nil {
		return mainConfFilePath, err
	}
	servicesConfPath := filepath.Join(servicesDir, "servicezz.toml")
	metricsForServicesDir := filepath.Join(os.TempDir(), "test", "trtr")
	if err = os.MkdirAll(metricsForServicesDir, 0750); err != nil {
		return mainConfFilePath, err
	}
	metricsForServicesConfPath := filepath.Join(metricsForServicesDir, "met4services.toml")
	myMonitoredServerMetricsDir := filepath.Join(os.TempDir(), "test", "integration")
	if err = os.MkdirAll(myMonitoredServerMetricsDir, 0750); err != nil {
		return mainConfFilePath, err
	}
	myMonitoredServerMetricsPath := filepath.Join(myMonitoredServerMetricsDir, "mymonitoredservermetri_s.toml")
	return mainConfFilePath,
		createMainConfig(
			mainConfigFileContenTemplate,
			mainConfFilePath,
			servicesConfPath,
			metricsForServicesConfPath,
			myMonitoredServerMetricsPath)
}

func TestSkycoinHealthSuit(t *testing.T) {
	suite.Run(t, new(HealthSuit))
}

func (suite *HealthSuit) TestMetricMonitorHealth() {
	// NOTE(denisacostaq@gmail.com): Giving
	require := require.New(suite.T())
	mainConfFilePath, err := createMainConfigTestPaths()
	require.Nil(err)
	srv := exporter.ExportMetrics(mainConfFilePath, "/metrics", 8081)
	require.NotNil(srv)
	// NOTE(denisacostaq@gmail.com): Wait for server starts
	t := time.NewTimer(time.Second * 2)
	<-t.C

	// NOTE(denisacostaq@gmail.com): When
	var resp *http.Response
	resp, err = http.Get("http://127.0.0.1:8081/metrics")

	// NOTE(denisacostaq@gmail.com): Assert
	suite.Nil(err)
	var data []byte
	data, err = ioutil.ReadAll(resp.Body)
	suite.Nil(err)
	suite.Equal(http.StatusOK, resp.StatusCode)
	suite.Contains(string(data), "open_connections_is_a_fake_name_for_test_purpose")
	var usingAVariableToMakeLinterHappy = context.Context(nil)
	require.Nil(srv.Shutdown(usingAVariableToMakeLinterHappy))
}

func (suite *HealthSuit) TestConfigWorks() {
	// NOTE(denisacostaq@gmail.com): Giving
	require := require.New(suite.T())
	srv := exporter.ExportMetrics("", "/metrics2", 8082)
	require.NotNil(srv)
	// NOTE(denisacostaq@gmail.com): Wait for server starts
	t := time.NewTimer(time.Second * 2)
	<-t.C

	// // NOTE(denisacostaq@gmail.com): When
	resp, err := http.Get("http://127.0.0.1:8082/metrics2")
	log.Println("resp, err", resp, err)
	// // NOTE(denisacostaq@gmail.com): Assert
	suite.Nil(err)
	suite.Equal(http.StatusOK, resp.StatusCode)
	var usingAVariableToMakeLinterHappy = context.Context(nil)
	require.Nil(srv.Shutdown(usingAVariableToMakeLinterHappy))
}
