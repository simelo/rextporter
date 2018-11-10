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

func createServiceConfig(tmplContent, path string) (err error) {
	generalScopeErr := "error creating service config file for integration test"
	if err = createConfigFile(tmplContent, path, nil); err != nil {
		errCause := "error writing service config file: " + err.Error()
		return common.ErrorFromThisScope(errCause, generalScopeErr)
	}
	return err
}

func createServiceConfigPaths(serviceConfigPath string) (err error) {
	const serviceConfigFileContenTemplate = `
	# Service configuration.
	[[services]]
		name = "myMonitoredServer"
		mode = "apiRest"
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
	return createServiceConfig(serviceConfigFileContenTemplate, serviceConfigPath)
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

func createMainConfig(tmplContent, path, metricsConfigPath, serviceConfigPath string) (err error) {
	generalScopeErr := "error creating main config file for integration test"
	type mainConfigData struct {
		ServiceConfigPath string
		MetricsConfigPath string
	}
	confData := mainConfigData{
		MetricsConfigPath: metricsConfigPath,
		ServiceConfigPath: serviceConfigPath,
	}
	if err = createConfigFile(tmplContent, path, confData); err != nil {
		errCause := "error writing main config file: " + err.Error()
		return common.ErrorFromThisScope(errCause, generalScopeErr)
	}
	if err = createServiceConfigPaths(serviceConfigPath); err != nil {
		errCause := "error writing service config file: " + err.Error()
		return common.ErrorFromThisScope(errCause, generalScopeErr)
	}
	if err = createMetricsConfigPaths(metricsConfigPath); err != nil {
		errCause := "error writing metrics config file: " + err.Error()
		return common.ErrorFromThisScope(errCause, generalScopeErr)
	}
	return err
}

func createMainConfigTestPaths() (mainConfFilePath string, err error) {
	const mainConfigFileContenTemplate = `
serviceConfigTransport = "file"
# render a template with a portable path
serviceConfigPath = "{{.ServiceConfigPath}}"
metricsConfigPath = "{{.MetricsConfigPath}}"
`
	mainConfigDir := filepath.Join(os.TempDir(), "sdsds", "675656", "aa")
	if err = os.MkdirAll(mainConfigDir, 0750); err != nil {
		return mainConfFilePath, err
	}
	mainConfFilePath = filepath.Join(mainConfigDir, "rrrr")
	serviceDir := filepath.Join(os.TempDir(), "test", "integration")
	if err = os.MkdirAll(serviceDir, 0750); err != nil {
		return mainConfFilePath, err
	}
	serviceConfigPath := filepath.Join(serviceDir, "service.toml")
	metricsDir := filepath.Join(os.TempDir(), "integration", "test")
	if err = os.MkdirAll(metricsDir, 0750); err != nil {
		return mainConfFilePath, err
	}
	metricsConfigPath := filepath.Join(metricsDir, "metrics.toml")
	return mainConfFilePath,
		createMainConfig(mainConfigFileContenTemplate, mainConfFilePath, metricsConfigPath, serviceConfigPath)
}

func createMainConfigCustomPaths() (mainConfFilePath string, err error) {
	const mainConfigFileContenTemplate = `
serviceConfigTransport = "file"
# render a template with a portable path
serviceConfigPath = "{{.ServiceConfigPath}}"
metricsConfigPath = "{{.MetricsConfigPath}}"
`
	mainConfigDir := filepath.Join(os.TempDir(), "sdsds", "675656", "aa")
	if err = os.MkdirAll(mainConfigDir, 0750); err != nil {
		return mainConfFilePath, err
	}
	mainConfFilePath = filepath.Join(mainConfigDir, "rrrr")
	serviceConfigPath := filepath.Join(os.TempDir(), "sdsds", "epe.toml")
	return mainConfFilePath,
		createMainConfig(mainConfigFileContenTemplate, mainConfFilePath, "", serviceConfigPath)
}

type HealthSuit struct {
	suite.Suite
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
	mainConfFilePath, err := createMainConfigCustomPaths()
	require.Nil(err)
	srv := exporter.ExportMetrics(mainConfFilePath, "/metrics2", 8082)
	require.NotNil(srv)
	// NOTE(denisacostaq@gmail.com): Wait for server starts
	t := time.NewTimer(time.Second * 2)
	<-t.C

	// // NOTE(denisacostaq@gmail.com): When
	var resp *http.Response
	resp, err = http.Get("http://127.0.0.1:8082/metrics2")
	log.Println("resp, err", resp, err)
	// // NOTE(denisacostaq@gmail.com): Assert
	suite.Nil(err)
	suite.Equal(http.StatusOK, resp.StatusCode)
	var usingAVariableToMakeLinterHappy = context.Context(nil)
	require.Nil(srv.Shutdown(usingAVariableToMakeLinterHappy))
}
