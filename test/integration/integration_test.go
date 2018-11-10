package main

import (
	"context"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/alecthomas/template"
	"github.com/simelo/rextporter/src/common"
	"github.com/simelo/rextporter/src/config"
	"github.com/simelo/rextporter/src/exporter"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

const serviceApiRestConfigFileContenTemplate = `
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

const serviceProxyConfigFileContenTemplate = `
	# Service configuration.
	[[services]]
		name = "myMonitoredServer"
		mode = "proxy"
		scheme = "http"
		port = 8080
		basePath = "/metrics"
		
		[services.location]
			location = "localhost"
`

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

func createServiceConfigPaths(serviceConfigPath, serviceContent string) (err error) {
	return createServiceConfig(serviceContent, serviceConfigPath)
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

func createMainConfig(tmplContent, path, metricsConfigPath, serviceConfigPath, serviceContent string) (err error) {
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
	if err = createServiceConfigPaths(serviceConfigPath, serviceContent); err != nil {
		errCause := "error writing service config file: " + err.Error()
		return common.ErrorFromThisScope(errCause, generalScopeErr)
	}
	if err = createMetricsConfigPaths(metricsConfigPath); err != nil {
		errCause := "error writing metrics config file: " + err.Error()
		return common.ErrorFromThisScope(errCause, generalScopeErr)
	}
	return err
}

func createMainConfigTestPaths(serviceContent string) (mainConfFilePath string, err error) {
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
		createMainConfig(mainConfigFileContenTemplate, mainConfFilePath, metricsConfigPath, serviceConfigPath, serviceContent)
}

func createMainConfigCustomPaths(serviceContent string) (mainConfFilePath string, err error) {
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
		createMainConfig(mainConfigFileContenTemplate, mainConfFilePath, "", serviceConfigPath, serviceContent)
}

type HealthSuit struct {
	suite.Suite
}

func TestSkycoinHealthSuit(t *testing.T) {
	suite.Run(t, new(HealthSuit))
}

func metricHealthIsOk(metricName, metricData string) bool {
	if !strings.Contains(metricData, metricName) {
		log.WithField("metricName", metricName).Errorln("metric name not found")
		return false
	}
	metricHealth := metricName + "_up"
	if !strings.Contains(metricData, metricHealth) {
		log.WithField("metricHealth", metricHealth).Errorln("metric health name not found")
		return false
	}
	lines := strings.Split(metricData, "\n")
	var linesWhoMentionMetric []string
	for _, line := range lines {
		if strings.Contains(line, metricHealth) {
			linesWhoMentionMetric = append(linesWhoMentionMetric, line)
		}
	}
	var targetLine string
	for _, line := range linesWhoMentionMetric {
		if strings.Contains(line, "# TYPE ") || strings.Contains(line, "# HELP ") {
			continue
		} else {
			targetLine = line
			break
		}
	}
	if strings.Compare(targetLine, "") == 0 {
		log.Errorln("can not find target line")
		return false
	}
	targetFields := strings.Split(targetLine, " ")
	if val, err := strconv.Atoi(targetFields[1]); err != nil || val != 0 {
		if err != nil {
			log.WithError(err).Errorln("unable to convert the value")
		}
		if val != 0 {
			log.WithField("val", val).Errorln("flag is set")
		}
		return false
	}
	return true
}

func (suite *HealthSuit) TestMetricMonitorHealth() {
	// NOTE(denisacostaq@gmail.com): Giving
	require := require.New(suite.T())
	mainConfFilePath, err := createMainConfigTestPaths(serviceApiRestConfigFileContenTemplate)
	require.Nil(err)
	srv := exporter.ExportMetrics(mainConfFilePath, "/metrics", 8081)
	require.NotNil(srv)
	conf := config.Config()
	// NOTE(denisacostaq@gmail.com): Wait for server starts
	t := time.NewTimer(time.Second * 2)
	<-t.C

	// NOTE(denisacostaq@gmail.com): When
	var resp *http.Response
	resp, err = http.Get("http://127.0.0.1:8081/metrics")

	// NOTE(denisacostaq@gmail.com): Assert
	suite.Nil(err)
	suite.Equal(http.StatusOK, resp.StatusCode)
	var data []byte
	data, err = ioutil.ReadAll(resp.Body)
	suite.Nil(err)
	require.Len(conf.Services, 1)
	require.Len(conf.Metrics, 1)
	metricName := "skycoin_" + conf.Services[0].Name + "_" + conf.Metrics[0].Name
	require.Equal(metricName, "skycoin_myMonitoredServer_open_connections_is_a_fake_name_for_test_purpose")
	require.True(metricHealthIsOk(metricName, string(data)))
	var usingAVariableToMakeLinterHappy = context.Context(nil)
	require.Nil(srv.Shutdown(usingAVariableToMakeLinterHappy))
}

func (suite *HealthSuit) TestMetricMonitorAsProxy() {
	// NOTE(denisacostaq@gmail.com): Giving
	require := require.New(suite.T())
	mainConfFilePath, err := createMainConfigCustomPaths(serviceProxyConfigFileContenTemplate)
	require.Nil(err)
	log.Errorln("mainConfFilePath", mainConfFilePath)
	srv := exporter.ExportMetrics(mainConfFilePath, "/metrics3", 8082)
	require.NotNil(srv)
	conf := config.Config()
	log.Errorln("services2", len(conf.Services))
	// NOTE(denisacostaq@gmail.com): Wait for server starts
	time.Sleep(time.Second * 2)

	// NOTE(denisacostaq@gmail.com): When
	var resp *http.Response
	resp, err = http.Get("http://127.0.0.1:8082/metrics3")

	// NOTE(denisacostaq@gmail.com): Assert
	suite.Nil(err)
	suite.Equal(http.StatusOK, resp.StatusCode)
	var data []byte
	data, err = ioutil.ReadAll(resp.Body)
	suite.Nil(err)
	require.Len(conf.Services, 1)
	require.Len(conf.Metrics, 1)
	metricName := conf.Services[0].Name + "_skycoin_wallet2_seq2"
	require.Equal(metricName, "myMonitoredServer_skycoin_wallet2_seq2")
	require.True(metricHealthIsOk(metricName, string(data)))
	var usingAVariableToMakeLinterHappy = context.Context(nil)
	require.Nil(srv.Shutdown(usingAVariableToMakeLinterHappy))
}

func (suite *HealthSuit) TestConfigWorks() {
	// NOTE(denisacostaq@gmail.com): Giving
	require := require.New(suite.T())
	mainConfFilePath, err := createMainConfigCustomPaths(serviceApiRestConfigFileContenTemplate)
	require.Nil(err)
	log.Errorln("mainConfFilePath", mainConfFilePath)
	srv := exporter.ExportMetrics(mainConfFilePath, "/metrics2", 8082)
	require.NotNil(srv)
	conf := config.Config()
	log.Errorln("services2", len(conf.Services))
	// NOTE(denisacostaq@gmail.com): Wait for server starts
	time.Sleep(time.Second * 2)

	// NOTE(denisacostaq@gmail.com): When
	var resp *http.Response
	resp, err = http.Get("http://127.0.0.1:8082/metrics2")

	// NOTE(denisacostaq@gmail.com): Assert
	suite.Nil(err)
	suite.Equal(http.StatusOK, resp.StatusCode)
	var data []byte
	data, err = ioutil.ReadAll(resp.Body)
	suite.Nil(err)
	require.Len(conf.Services, 1)
	require.Len(conf.Metrics, 1)
	metricName := "skycoin_" + conf.Services[0].Name + "_" + conf.Metrics[0].Name
	require.Equal(metricName, "skycoin_myMonitoredServer_seq")
	require.True(metricHealthIsOk(metricName, string(data)))
	var usingAVariableToMakeLinterHappy = context.Context(nil)
	require.Nil(srv.Shutdown(usingAVariableToMakeLinterHappy))
}
