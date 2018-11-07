package main

import (
	"context"
	"io/ioutil"
	"net/http"
	"os"
	"testing"

	"github.com/alecthomas/template"
	"github.com/denisacostaq/rextporter/src/common"
	"github.com/simelo/rextporter/src/exporter"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type HealthSuit struct {
	suite.Suite
	srv *http.Server
}

func createMainConfig(path string) (err error) {
	generalScopeErr := "error creating main config file for integration test"
	const mainConfigFileContenTemplate = `
serviceConfigTransport = "file" # "file" | "consulCatalog"
# render a template with a portable path
serviceConfigPath = "{{.ServiceConfigPath}}"
metricsConfigPath = "{{.MetricsConfigPath}}"
`
	tmpl := template.New("mainConfig")
	var templateEngine *template.Template
	if templateEngine, err = tmpl.Parse(mainConfigFileContenTemplate); err != nil {
		errCause := "error parsing main config: " + err.Error()
		return common.ErrorFromThisScope(errCause, generalScopeErr)
	}
	var mainConfigFile *os.File
	if mainConfigFile, err = os.Create(path); err != nil {
		errCause := "error creating main config file: " + err.Error()
		return common.ErrorFromThisScope(errCause, generalScopeErr)
	}
	type mainConfigData struct {
		ServiceConfigPath string
		MetricsConfigPath string
	}
	gopath := os.Getenv("GOPATH")
	confData := mainConfigData{
		MetricsConfigPath: gopath + "/src/github.com/simelo/rextporter/test/integration/metrics.toml",
		ServiceConfigPath: gopath + "/src/github.com/simelo/rextporter/test/integration/service.toml",
	}
	if err = templateEngine.Execute(mainConfigFile, confData); err != nil {
		errCause := "error writing main config file: " + err.Error()
		return common.ErrorFromThisScope(errCause, generalScopeErr)
	}
	return err
}

func (suite *HealthSuit) SetupSuite() {
	require := require.New(suite.T())
	gopath := os.Getenv("GOPATH")
	mainConfFilePath := gopath + "/src/github.com/simelo/rextporter/test/integration/main.toml"
	require.Nil(createMainConfig(mainConfFilePath))
	suite.srv = exporter.ExportMetrics(mainConfFilePath, 8081)
	require.NotNil(suite.srv)
}

func (suite *HealthSuit) TearDownSuite() {
	require := require.New(suite.T())
	var usingAVariableToMakeLinterHappy = context.Context(nil)
	require.Nil(suite.srv.Shutdown(usingAVariableToMakeLinterHappy))
}

func TestSkycoinHealthSuit(t *testing.T) {
	suite.Run(t, new(HealthSuit))
}

func (suite *HealthSuit) TestMetricMonitorHealth() {
	// NOTE(denisacostaq@gmail.com): Giving

	// NOTE(denisacostaq@gmail.com): When
	resp, err := http.Get("http://127.0.0.1:8081/metrics")
	suite.Nil(err)
	var data []byte
	data, err = ioutil.ReadAll(resp.Body)
	suite.Nil(err)

	// NOTE(denisacostaq@gmail.com): Assert
	suite.Equal(http.StatusOK, resp.StatusCode)
	suite.Contains(string(data), "open_connections_is_a_fake_name_for_test_purpose")
}
