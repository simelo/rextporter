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

const mainConfigFileContenTemplate = `
serviceConfigTransport = "file"
# render a template with a portable path
servicesConfigPath = "{{.ServicesConfigPath}}"
metricsForServicesPath = "{{.MetricsForServicesPath}}"
`

const servicesAPIRestConfigFileContenTemplate = `
	# Service configuration.
	[[services]]
		name = "myMonitoredServer"
		mode="apiRest"
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

const servicesProxyConfigFileContenTemplate = `
# Service configuration.
[[services]]
	name = "myMonitoredAsProxyServer"
	mode="proxy"
	scheme = "http"
	port = 8080
	basePath = "/metrics"
	
	[services.location]
		location = "localhost"
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

const metricsConfigFileContenTemplateNotAccesiblevalue = `
# All metrics to be measured.
[[metrics]]
	name = "can_not_be_updated"
	url = "/api/v1/health"
	httpMethod = "GET"
	path = "fake_json_path_he_he"

	[metrics.options]
		type = "Gauge"
		description = "Track the open connections in the system"	
`

const metricsForServicesConfFileContenTemplate = `
	serviceNameToMetricsConfPath = [{{range $key, $value := .}}
	{ {{$key}} = "{{$value}}" },{{end}}
]
`

type HealthSuit struct {
	suite.Suite
	require                          *require.Assertions
	mainConfFilePath                 string
	mainConfTmplContent              string
	servicesConfTmplContent          string
	servicesConfFilePath             string
	metricsConfTmplContent           string
	metricsConfFilePath              string
	metricsForServiceConfTmplContent string
	metricsForServicesConfFilePath   string
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

func (suite *HealthSuit) createServicesConfPath() (err error) {
	generalScopeErr := "error creating service config file for integration test"
	if err = createConfigFile(suite.servicesConfTmplContent, suite.servicesConfFilePath, nil); err != nil {
		errCause := "error writing service config file: " + err.Error()
		return common.ErrorFromThisScope(errCause, generalScopeErr)
	}
	return err
}

func (suite *HealthSuit) createMainConfPath(tmplData interface{}) (err error) {
	generalScopeErr := "error creating main config file for integration test"
	if err = createConfigFile(suite.mainConfTmplContent, suite.mainConfFilePath, tmplData); err != nil {
		errCause := "error writing service config file: " + err.Error()
		return common.ErrorFromThisScope(errCause, generalScopeErr)
	}
	return err
}

func (suite *HealthSuit) createMetricsForServicesConfPath(metricsForServices interface{}) (err error) {
	return createConfigFile(
		suite.metricsForServiceConfTmplContent,
		suite.metricsForServicesConfFilePath,
		metricsForServices)
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
		suite.require.Nil(common.ErrorFromThisScope(errCause, generalScopeErr))
	}
	if err := suite.createServicesConfPath(); err != nil {
		errCause := "error writing services config file: " + err.Error()
		suite.require.Nil(common.ErrorFromThisScope(errCause, generalScopeErr))
	}
	if err := suite.createMetricsConfigPaths(); err != nil {
		errCause := "error writing my monitored server metrics config file: " + err.Error()
		suite.require.Nil(common.ErrorFromThisScope(errCause, generalScopeErr))
	}
	if err := suite.createMetricsForServicesConfPath(
		map[string]string{
			"myMonitoredServer":        suite.metricsConfFilePath,
			"myMonitoredAsProxyServer": suite.metricsConfFilePath}); err != nil {
		errCause := "error writing metrics for service config file: " + err.Error()
		suite.require.Nil(common.ErrorFromThisScope(errCause, generalScopeErr))
	}
}

func (suite *HealthSuit) createDirectoriesWithFullDepth(dirs []string) {
	for _, dir := range dirs {
		suite.require.Nil(os.MkdirAll(dir, 0750))
	}
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

func TestSkycoinHealthSuit(t *testing.T) {
	suite.Run(t, new(HealthSuit))
}

func (suite *HealthSuit) TestDefaultGeneratedConfigWorks() {
	// NOTE(denisacostaq@gmail.com): Giving
	suite.require = require.New(suite.T())
	srv := exporter.ExportMetrics("", "/metrics1", 8081)
	suite.require.NotNil(srv)
	// NOTE(denisacostaq@gmail.com): Wait for server starts
	time.Sleep(time.Second * 2)
	conf := config.Config()

	// NOTE(denisacostaq@gmail.com): When
	resp, err := http.Get("http://127.0.0.1:8081/metrics1")
	if err == nil {
		defer func() {
			suite.Nil(resp.Body.Close())
		}()
	}
	// NOTE(denisacostaq@gmail.com): Assert
	suite.Nil(err)
	suite.Equal(http.StatusOK, resp.StatusCode)
	var data []byte
	data, err = ioutil.ReadAll(resp.Body)
	suite.Nil(err)
	suite.require.Len(conf.Services, 1)
	suite.require.Len(conf.Services[0].Metrics, 1)
	metricName := "skycoin_" + conf.Services[0].Name + "_" + conf.Services[0].Metrics[0].Name
	suite.require.Equal(metricName, "skycoin_skycoin_seq")
	suite.require.True(metricHealthIsOk(metricName, string(data)))
	var usingAVariableToMakeLinterHappy = context.Context(nil)
	suite.require.Nil(srv.Shutdown(usingAVariableToMakeLinterHappy))
}

func (suite *HealthSuit) TestMetricMonitorHealth() {
	// NOTE(denisacostaq@gmail.com): Giving
	suite.require = require.New(suite.T())
	mainConfigDir := filepath.Join(os.TempDir(), "fdf", "vcvcv", "aa")
	servicesDir := filepath.Join(os.TempDir(), "dfdfd", "3333")
	myMonitoredServerMetricsDir := filepath.Join(os.TempDir(), "test", "aaaa")
	metricsForServicesDir := filepath.Join(os.TempDir(), "tdddddest", "trtr")
	suite.createDirectoriesWithFullDepth([]string{mainConfigDir, servicesDir, myMonitoredServerMetricsDir, metricsForServicesDir})
	suite.mainConfFilePath = filepath.Join(mainConfigDir, "qqqqqqq")
	suite.servicesConfFilePath = filepath.Join(servicesDir, "servicezz.toml")
	suite.metricsConfFilePath = filepath.Join(myMonitoredServerMetricsDir, "__dd.toml")
	suite.metricsForServicesConfFilePath = filepath.Join(metricsForServicesDir, "dfdf.toml")
	suite.mainConfTmplContent = mainConfigFileContenTemplate
	suite.servicesConfTmplContent = servicesAPIRestConfigFileContenTemplate
	suite.metricsConfTmplContent = metricsConfigFileContenTemplate
	suite.metricsForServiceConfTmplContent = metricsForServicesConfFileContenTemplate
	suite.createMainConfig()
	srv := exporter.ExportMetrics(suite.mainConfFilePath, "/metrics2", 8082)
	suite.require.NotNil(srv)
	// NOTE(denisacostaq@gmail.com): Wait for server starts
	time.Sleep(time.Second * 2)
	conf := config.Config()

	// NOTE(denisacostaq@gmail.com): When
	resp, err := http.Get("http://127.0.0.1:8082/metrics2")
	if err == nil {
		defer func() {
			suite.Nil(resp.Body.Close())
		}()
	}

	// NOTE(denisacostaq@gmail.com): Assert
	suite.Nil(err)
	suite.Equal(http.StatusOK, resp.StatusCode)
	var data []byte
	data, err = ioutil.ReadAll(resp.Body)
	suite.Nil(err)
	suite.require.Len(conf.Services, 1)
	suite.require.Len(conf.Services[0].Metrics, 1)
	metricName := "skycoin_" + conf.Services[0].Name + "_" + conf.Services[0].Metrics[0].Name
	suite.require.Equal(metricName, "skycoin_myMonitoredServer_open_connections_is_a_fake_name_for_test_purpose")
	suite.require.True(metricHealthIsOk(metricName, string(data)))
	var usingAVariableToMakeLinterHappy = context.Context(nil)
	suite.require.Nil(srv.Shutdown(usingAVariableToMakeLinterHappy))
}

func (suite *HealthSuit) TestMetricMonitorHealthCanSetUpFlag() {
	// NOTE(denisacostaq@gmail.com): Giving
	suite.require = require.New(suite.T())
	mainConfigDir := filepath.Join(os.TempDir(), "fdf", "vcvcv", "aa")
	servicesDir := filepath.Join(os.TempDir(), "dfdfd", "3333")
	myMonitoredServerMetricsDir := filepath.Join(os.TempDir(), "test", "aaaa")
	metricsForServicesDir := filepath.Join(os.TempDir(), "tdddddest", "trtr")
	suite.createDirectoriesWithFullDepth([]string{mainConfigDir, servicesDir, myMonitoredServerMetricsDir, metricsForServicesDir})
	suite.mainConfFilePath = filepath.Join(mainConfigDir, "qqqqqqq")
	suite.servicesConfFilePath = filepath.Join(servicesDir, "servicezz.toml")
	suite.metricsConfFilePath = filepath.Join(myMonitoredServerMetricsDir, "__dd.toml")
	suite.metricsForServicesConfFilePath = filepath.Join(metricsForServicesDir, "dfdf.toml")
	suite.mainConfTmplContent = mainConfigFileContenTemplate
	suite.servicesConfTmplContent = servicesAPIRestConfigFileContenTemplate
	suite.metricsConfTmplContent = metricsConfigFileContenTemplateNotAccesiblevalue
	suite.metricsForServiceConfTmplContent = metricsForServicesConfFileContenTemplate
	suite.createMainConfig()
	srv := exporter.ExportMetrics(suite.mainConfFilePath, "/metrics3", 8083)
	suite.require.NotNil(srv)
	// NOTE(denisacostaq@gmail.com): Wait for server starts
	time.Sleep(time.Second * 2)
	conf := config.Config()

	// NOTE(denisacostaq@gmail.com): When
	resp, err := http.Get("http://127.0.0.1:8083/metrics3")
	if err == nil {
		defer func() {
			suite.Nil(resp.Body.Close())
		}()
	}

	// NOTE(denisacostaq@gmail.com): Assert
	suite.Nil(err)
	suite.Equal(http.StatusOK, resp.StatusCode)
	var data []byte
	data, err = ioutil.ReadAll(resp.Body)
	suite.Nil(err)
	suite.require.Len(conf.Services, 1)
	suite.require.Len(conf.Services[0].Metrics, 1)
	metricName := "skycoin_" + conf.Services[0].Name + "_" + conf.Services[0].Metrics[0].Name
	suite.require.Equal(metricName, "skycoin_myMonitoredServer_can_not_be_updated")
	suite.require.False(metricHealthIsOk(metricName, string(data)))
	var usingAVariableToMakeLinterHappy = context.Context(nil)
	suite.require.Nil(srv.Shutdown(usingAVariableToMakeLinterHappy))
}

func (suite *HealthSuit) TestMetricMonitorAsProxy() {
	// NOTE(denisacostaq@gmail.com): Giving
	suite.require = require.New(suite.T())
	mainConfigDir := filepath.Join(os.TempDir(), "sdsds", "675656", "aa")
	servicesDir := filepath.Join(os.TempDir(), "test", "integration")
	myMonitoredServerMetricsDir := filepath.Join(os.TempDir(), "test", "integration")
	metricsForServicesDir := filepath.Join(os.TempDir(), "test", "trtr")
	suite.createDirectoriesWithFullDepth([]string{mainConfigDir, servicesDir, myMonitoredServerMetricsDir, metricsForServicesDir})
	suite.mainConfFilePath = filepath.Join(mainConfigDir, "rrrr")
	suite.servicesConfFilePath = filepath.Join(servicesDir, "servicezz.toml")
	suite.metricsConfFilePath = filepath.Join(myMonitoredServerMetricsDir, "mymonitoredservermetri_s.toml")
	suite.metricsForServicesConfFilePath = filepath.Join(metricsForServicesDir, "met4services.toml")
	suite.mainConfTmplContent = mainConfigFileContenTemplate
	suite.servicesConfTmplContent = servicesProxyConfigFileContenTemplate
	suite.metricsConfTmplContent = metricsConfigFileContenTemplate
	suite.metricsForServiceConfTmplContent = metricsForServicesConfFileContenTemplate
	suite.createMainConfig()
	srv := exporter.ExportMetrics(suite.mainConfFilePath, "/metrics4", 8084)
	suite.require.NotNil(srv)
	conf := config.Config()
	// NOTE(denisacostaq@gmail.com): Wait for server starts
	time.Sleep(time.Second * 2)

	// NOTE(denisacostaq@gmail.com): When
	resp, err := http.Get("http://127.0.0.1:8084/metrics4")
	if err == nil {
		defer func() {
			suite.Nil(resp.Body.Close())
		}()
	}

	// NOTE(denisacostaq@gmail.com): Assert
	suite.Nil(err)
	suite.Equal(http.StatusOK, resp.StatusCode)
	var data []byte
	data, err = ioutil.ReadAll(resp.Body)
	suite.Nil(err)
	suite.require.Len(conf.Services, 1)
	suite.require.Len(conf.Services[0].Metrics, 0)
	metricName := conf.Services[0].Name + "_skycoin_wallet2_seq2"
	suite.require.Equal(metricName, "myMonitoredAsProxyServer_skycoin_wallet2_seq2")
	suite.require.True(metricHealthIsOk(metricName, string(data)))
	var usingAVariableToMakeLinterHappy = context.Context(nil)
	suite.require.Nil(srv.Shutdown(usingAVariableToMakeLinterHappy))
}
