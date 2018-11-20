package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
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
		mode="{{.Mode}}"
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

type Service struct {
	Name     string
	Port     uint16
	Mode     string
	BasePath string
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
		Services: []Service{Service{Name: "myMonitoredServer", Port: fakeNodePort, Mode: "rest_api", BasePath: ""}},
	}
}

// func (suite *HealthSuit) TestDefaultGeneratedConfigWorks() {
// 	// NOTE(denisacostaq@gmail.com): Giving
// 	suite.require = require.New(suite.T())
// 	port := testrand.RandomPort()
// 	srv := exporter.ExportMetrics("", "/metrics1", port)
// 	suite.require.NotNil(srv)
// 	// NOTE(denisacostaq@gmail.com): Wait for server starts
// 	time.Sleep(time.Second * 2)
// 	conf := config.Config()

// 	// NOTE(denisacostaq@gmail.com): When
// 	resp, err := http.Get(fmt.Sprintf("http://127.0.0.1:%d/metrics1", port))

// 	// NOTE(denisacostaq@gmail.com): Assert
// 	suite.Nil(err)
// 	defer func() { suite.Nil(resp.Body.Close()) }()
// 	suite.Equal(http.StatusOK, resp.StatusCode)
// 	var data []byte
// 	data, err = ioutil.ReadAll(resp.Body)
// 	suite.Nil(err)
// 	suite.require.Len(conf.Services, 1)
// 	suite.require.Len(conf.Services[0].Metrics, 1)
// 	metricName := "skycoin_" + conf.Services[0].Name + "_" + conf.Services[0].Metrics[0].Name
// 	suite.Equal(metricName, "skycoin_skycoin_seq")
// 	suite.require.True(metricHealthIsOk(metricName, string(data)))
// 	var usingAVariableToMakeLinterHappy = context.Context(nil)
// 	suite.Nil(srv.Shutdown(usingAVariableToMakeLinterHappy))
// }

func (suite *HealthSuit) TestMetricMonitorHealth() {
	// NOTE(denisacostaq@gmail.com): Giving
	suite.require = require.New(suite.T())
	mainConfigDir := testrand.RFolderPath()
	servicesDir := testrand.RFolderPath()
	myMonitoredServerMetricsDir := testrand.RFolderPath()
	metricsForServicesDir := testrand.RFolderPath()
	port := testrand.RandomPort()
	suite.createDirectoriesWithFullDepth([]string{mainConfigDir, servicesDir, myMonitoredServerMetricsDir, metricsForServicesDir})
	// suite.createDirectoriesWithFullDepth([]string{servicesDir})
	suite.mainConfFilePath = filepath.Join(mainConfigDir, testrand.RName())
	suite.servicesConfFilePath = filepath.Join(servicesDir, testrand.RName())
	suite.metricsConfFilePath = filepath.Join(myMonitoredServerMetricsDir, testrand.RName())
	suite.metricsForServicesConfFilePath = filepath.Join(metricsForServicesDir, testrand.RName())
	suite.mainConfTmplContent = mainConfigFileContenTemplate
	suite.metricsConfTmplContent = metricsConfigFileContenTemplate
	suite.metricsForServiceConfTmplContent = metricsForServicesConfFileContenTemplate
	suite.callSetUpTest()
	suite.createMainConfig()
	srv := exporter.ExportMetrics(suite.mainConfFilePath, "/metrics2", port)
	suite.require.NotNil(srv)
	// NOTE(denisacostaq@gmail.com): Wait for server starts
	time.Sleep(time.Second * 2)
	conf := config.Config()

	// NOTE(denisacostaq@gmail.com): When
	resp, err := http.Get(fmt.Sprintf("http://127.0.0.1:%d/metrics2", port))

	// NOTE(denisacostaq@gmail.com): Assert
	defer func() { suite.Nil(resp.Body.Close()) }()
	suite.Nil(err)
	suite.Equal(http.StatusOK, resp.StatusCode)
	var data []byte
	data, err = ioutil.ReadAll(resp.Body)
	suite.Nil(err)
	suite.Len(conf.Services, 1)
	suite.Len(conf.Services[0].Metrics, 1)
	metricName := "skycoin_" + conf.Services[0].Name + "_" + conf.Services[0].Metrics[0].Name
	suite.Equal(metricName, "skycoin_myMonitoredServer_open_connections_is_a_fake_name_for_test_purpose")
	suite.True(metricHealthIsOk(metricName, string(data)))
	var usingAVariableToMakeLinterHappy = context.Context(nil)
	suite.Nil(srv.Shutdown(usingAVariableToMakeLinterHappy))
}

func (suite *HealthSuit) TestMetricMonitorHealthCanSetUpFlag() {
	// NOTE(denisacostaq@gmail.com): Giving
	suite.require = require.New(suite.T())
	mainConfigDir := testrand.RFolderPath()
	servicesDir := testrand.RFolderPath()
	port := testrand.RandomPort()
	myMonitoredServerMetricsDir := testrand.RFolderPath()
	metricsForServicesDir := testrand.RFolderPath()
	suite.createDirectoriesWithFullDepth([]string{mainConfigDir, servicesDir, myMonitoredServerMetricsDir, metricsForServicesDir})
	suite.mainConfFilePath = filepath.Join(mainConfigDir, testrand.RName())
	suite.servicesConfFilePath = filepath.Join(servicesDir, testrand.RName())
	suite.metricsConfFilePath = filepath.Join(myMonitoredServerMetricsDir, testrand.RName())
	suite.metricsForServicesConfFilePath = filepath.Join(metricsForServicesDir, testrand.RName())
	suite.mainConfTmplContent = mainConfigFileContenTemplate
	suite.metricsConfTmplContent = metricsConfigFileContenTemplateNotAccesiblevalue
	suite.metricsForServiceConfTmplContent = metricsForServicesConfFileContenTemplate
	// TODO(denisacostaq@gmail.com): Do not call this directly, use a setup scheme
	suite.callSetUpTest()
	suite.createMainConfig()
	srv := exporter.ExportMetrics(suite.mainConfFilePath, "/metrics3", port)
	suite.require.NotNil(srv)
	// NOTE(denisacostaq@gmail.com): Wait for server starts
	time.Sleep(time.Second * 2)
	conf := config.Config()

	// NOTE(denisacostaq@gmail.com): When
	resp, err := http.Get(fmt.Sprintf("http://127.0.0.1:%d/metrics3", port))

	// NOTE(denisacostaq@gmail.com): Assert
	defer func() { suite.Nil(resp.Body.Close()) }()
	suite.Nil(err)
	suite.Equal(http.StatusOK, resp.StatusCode)
	var data []byte
	data, err = ioutil.ReadAll(resp.Body)
	suite.Nil(err)
	suite.require.Len(conf.Services, 1)
	suite.require.Len(conf.Services[0].Metrics, 1)
	metricName := "skycoin_" + conf.Services[0].Name + "_" + conf.Services[0].Metrics[0].Name
	suite.Equal(metricName, "skycoin_myMonitoredServer_can_not_be_updated")
	suite.False(metricHealthIsOk(metricName, string(data)))
	suite.Nil(srv.Shutdown(context.Context(nil)))
}

func (suite *HealthSuit) TestMetricMonitorAsProxy() {
	// NOTE(denisacostaq@gmail.com): Giving
	suite.require = require.New(suite.T())
	port := testrand.RandomPort()
	mainConfigDir := testrand.RFolderPath()
	servicesDir := testrand.RFolderPath()
	myMonitoredServerMetricsDir := testrand.RFolderPath()
	metricsForServicesDir := testrand.RFolderPath()
	suite.createDirectoriesWithFullDepth([]string{mainConfigDir, servicesDir, myMonitoredServerMetricsDir, metricsForServicesDir})
	suite.mainConfFilePath = filepath.Join(mainConfigDir, testrand.RName())
	suite.servicesConfFilePath = filepath.Join(servicesDir, testrand.RName())
	suite.metricsConfFilePath = filepath.Join(myMonitoredServerMetricsDir, testrand.RName())
	suite.metricsForServicesConfFilePath = filepath.Join(metricsForServicesDir, testrand.RName())
	suite.mainConfTmplContent = mainConfigFileContenTemplate
	suite.metricsConfTmplContent = metricsConfigFileContenTemplate
	suite.metricsForServiceConfTmplContent = metricsForServicesConfFileContenTemplate
	suite.callSetUpTest()
	suite.metricsForServicesConfData =
		map[string]string{
			"myMonitoredServer":        suite.metricsConfFilePath,
			"myMonitoredAsProxyServer": suite.metricsConfFilePath}
	suite.servicesConfData = ServicesConfData{
		Services: []Service{Service{Name: "myMonitoredAsProxyServer", Port: fakeNodePort, Mode: "forward_metrics", BasePath: "/metrics"}},
	}
	suite.createMainConfig()
	srv := exporter.ExportMetrics(suite.mainConfFilePath, "/metrics4", port)
	suite.require.NotNil(srv)
	conf := config.Config()
	// NOTE(denisacostaq@gmail.com): Wait for server starts
	time.Sleep(time.Second * 2)

	// NOTE(denisacostaq@gmail.com): When
	var resp *http.Response
	resp, err := http.Get(fmt.Sprintf("http://127.0.0.1:%d/metrics4", port))
	suite.require.NotNil(resp)

	// NOTE(denisacostaq@gmail.com): Assert
	suite.Nil(err)
	defer func() { suite.Nil(resp.Body.Close()) }()
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

func (suite *HealthSuit) TestMetricMonitorAsProxyWithNonMetricsEndpoint() {
	// NOTE(denisacostaq@gmail.com): Giving
	suite.require = require.New(suite.T())
	port := testrand.RandomPort()
	mainConfigDir := testrand.RFolderPath()
	servicesDir := testrand.RFolderPath()
	myMonitoredServerMetricsDir := testrand.RFolderPath()
	metricsForServicesDir := testrand.RFolderPath()
	suite.createDirectoriesWithFullDepth([]string{mainConfigDir, servicesDir, myMonitoredServerMetricsDir, metricsForServicesDir})
	suite.mainConfFilePath = filepath.Join(mainConfigDir, testrand.RName())
	suite.servicesConfFilePath = filepath.Join(servicesDir, testrand.RName()+".toml")
	suite.metricsConfFilePath = filepath.Join(myMonitoredServerMetricsDir, testrand.RName()+".toml")
	suite.metricsForServicesConfFilePath = filepath.Join(metricsForServicesDir, testrand.RName()+".toml")
	suite.mainConfTmplContent = mainConfigFileContenTemplate
	suite.callSetUpTest()
	suite.metricsForServicesConfData =
		map[string]string{
			"myMonitoredAsProxyServer": suite.metricsConfFilePath}
	suite.servicesConfData = ServicesConfData{
		Services: []Service{Service{Name: "myMonitoredAsProxyServer", Port: fakeNodePort, Mode: "forward_metrics", BasePath: "/api/v1/health"}},
	}
	suite.metricsConfTmplContent = metricsConfigFileContenTemplate
	suite.metricsForServiceConfTmplContent = metricsForServicesConfFileContenTemplate
	suite.createMainConfig()
	srv := exporter.ExportMetrics(suite.mainConfFilePath, "/metrics5", port)
	suite.require.NotNil(srv)
	conf := config.Config()
	// NOTE(denisacostaq@gmail.com): Wait for server starts
	time.Sleep(time.Second * 2)

	// NOTE(denisacostaq@gmail.com): When
	resp, err := http.Get(fmt.Sprintf("http://127.0.0.1:%d/metrics5", port))

	// NOTE(denisacostaq@gmail.com): Assert
	suite.Nil(err)
	defer func() { suite.Nil(resp.Body.Close()) }()
	suite.Equal(http.StatusOK, resp.StatusCode)
	var data []byte
	data, err = ioutil.ReadAll(resp.Body)
	suite.Nil(err)
	suite.require.Len(conf.Services, 1)
	suite.require.Len(conf.Services[0].Metrics, 0)
	metricName := conf.Services[0].Name + "_skycoin_wallet2_seq2"
	suite.Equal(metricName, "myMonitoredAsProxyServer_skycoin_wallet2_seq2")
	suite.require.False(metricHealthIsOk(metricName, string(data)))
	var usingAVariableToMakeLinterHappy = context.Context(nil)
	suite.require.Nil(srv.Shutdown(usingAVariableToMakeLinterHappy))
}

func (suite *HealthSuit) TestMetricMonitorAsProxyWithMetricsNamesOverlap() {
	// NOTE(denisacostaq@gmail.com): Giving
	suite.require = require.New(suite.T())
	port := testrand.RandomPort()
	mainConfigDir := testrand.RFolderPath()
	servicesDir := testrand.RFolderPath()
	myMonitoredServerMetricsDir := testrand.RFolderPath()
	metricsForServicesDir := testrand.RFolderPath()
	suite.createDirectoriesWithFullDepth([]string{mainConfigDir, servicesDir, myMonitoredServerMetricsDir, metricsForServicesDir})
	suite.mainConfFilePath = filepath.Join(mainConfigDir, testrand.RName())
	suite.servicesConfFilePath = filepath.Join(servicesDir, testrand.RName()+".toml")
	suite.metricsConfFilePath = filepath.Join(myMonitoredServerMetricsDir, testrand.RName()+".toml")
	suite.metricsForServicesConfFilePath = filepath.Join(metricsForServicesDir, testrand.RName()+".toml")
	suite.mainConfTmplContent = mainConfigFileContenTemplate
	suite.callSetUpTest()
	suite.metricsForServicesConfData =
		map[string]string{
			"myMonitoredAsProxyServer2": suite.metricsConfFilePath}
	suite.servicesConfData = ServicesConfData{
		Services: []Service{Service{Name: "myMonitoredAsProxyServer2", Port: fakeNodePort, Mode: "forward_metrics", BasePath: "/a_few_metrics"}},
	}
	suite.metricsConfTmplContent = metricsConfigFileContenTemplate
	suite.metricsForServiceConfTmplContent = metricsForServicesConfFileContenTemplate
	suite.createMainConfig()
	srv := exporter.ExportMetrics(suite.mainConfFilePath, "/metrics6", port)
	suite.require.NotNil(srv)
	conf := config.Config()
	// NOTE(denisacostaq@gmail.com): Wait for server starts
	time.Sleep(time.Second * 2)

	// NOTE(denisacostaq@gmail.com): When
	resp, err := http.Get(fmt.Sprintf("http://127.0.0.1:%d/metrics6", port))

	// NOTE(denisacostaq@gmail.com): Assert
	suite.Nil(err)
	defer func() { suite.Nil(resp.Body.Close()) }()
	suite.Equal(http.StatusOK, resp.StatusCode)
	var data []byte
	data, err = ioutil.ReadAll(resp.Body)
	suite.Nil(err)
	suite.require.Len(conf.Services, 1)
	suite.require.Len(conf.Services[0].Metrics, 0)
	metricName := conf.Services[0].Name + "_seq"
	suite.require.Equal(metricName, "myMonitoredAsProxyServer2_seq")
	suite.require.True(metricHealthIsOk(metricName, string(data)))
	metricName = conf.Services[0].Name + "_main_seq"
	suite.require.Equal(metricName, "myMonitoredAsProxyServer2_main_seq")
	suite.require.True(metricHealthIsOk(metricName, string(data)))
	suite.Nil(srv.Shutdown(context.Context(nil)))
}
