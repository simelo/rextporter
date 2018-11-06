package main

import (
	"context"
	"io/ioutil"
	"net/http"
	"os"
	"testing"

	"github.com/simelo/rextporter/src/exporter"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type HealthSuit struct {
	suite.Suite
	srv *http.Server
}

func (suite *HealthSuit) SetupSuite() {
	require := require.New(suite.T())
	gopath := os.Getenv("GOPATH")
	metricsConfFilePath := gopath + "/src/github.com/simelo/rextporter/test/integration/metrics.toml"
	serviceConfFilePath := gopath + "/src/github.com/simelo/rextporter/test/integration/service.toml"
	suite.srv = exporter.ExportMetrics(metricsConfFilePath, serviceConfFilePath, 8081)
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
