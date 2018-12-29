package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"
	"time"

	"github.com/denisacostaq/rextporter/test/util"
	"github.com/simelo/rextporter/src/core"
	"github.com/simelo/rextporter/src/exporter"
	"github.com/simelo/rextporter/src/toml2config"
	"github.com/simelo/rextporter/src/tomlconfig"
	"github.com/simelo/rextporter/test/util/testrand"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type SkycoinSuit struct {
	suite.Suite
	require            *require.Assertions
	rextporterEndpoint string
	rextporterServer   *http.Server
}

func (suite *SkycoinSuit) SetupSuite() {
	suite.require = require.New(suite.T())
	mainConfFilePath := "/usr/share/gocode/src/github.com/simelo/rextporter/test/integration/skycoin/tomlconfig/main.toml"
	tomlConf, err := tomlconfig.ReadConfigFromFileSystem(mainConfFilePath)
	suite.Nil(err)
	var conf core.RextRoot
	conf, err = toml2config.Fill(tomlConf)
	suite.Nil(err)
	listenPort := testrand.RandomPort()
	suite.rextporterEndpoint = fmt.Sprintf("http://localhost:%d%s", listenPort, "/metrics")
	suite.rextporterServer = exporter.MustExportMetrics("", "/metrics", listenPort, conf)
	suite.require.NotNil(suite.rextporterServer)
	// NOTE(denisacostaq@gmail.com): Wait for server starts
	time.Sleep(time.Second * 2)
}

func (suite *SkycoinSuit) TearDownSuite() {
	log.Info("Shutting down server...")
	suite.Nil(suite.rextporterServer.Shutdown(context.Context(nil)))
}

func TestSkycoinSuitSuit(t *testing.T) {
	suite.Run(t, new(SkycoinSuit))
}

func (suite *SkycoinSuit) TestMetricsAreExported() {
	// NOTE(denisacostaq@gmail.com): Giving

	// NOTE(denisacostaq@gmail.com): When
	resp, err := http.Get(suite.rextporterEndpoint)

	// NOTE(denisacostaq@gmail.com): Assert
	suite.Nil(err)
	suite.Equal(http.StatusOK, resp.StatusCode)
	suite.NotNil(resp.Body)
}

func (suite *SkycoinSuit) TestDefaultMetricsForAPIScrapperArePresent() {
	// NOTE(denisacostaq@gmail.com): Giving

	// NOTE(denisacostaq@gmail.com): When
	resp, err := http.Get(suite.rextporterEndpoint)

	// NOTE(denisacostaq@gmail.com): Assert
	suite.Nil(err)
	suite.Equal(http.StatusOK, resp.StatusCode)
	suite.NotNil(resp.Body)
	var respBody []byte
	respBody, err = ioutil.ReadAll(resp.Body)
	suite.Nil(err)
	suite.NotNil(respBody)
	mtrs := []string{
		"scrape_duration_seconds",
		"scrape_samples_scraped",
		"data_source_response_duration_seconds",
		"data_source_scrape_duration_seconds",
		"data_source_scrape_samples_scraped",
	}
	for _, mtr := range mtrs {
		var found bool
		found, err = util.FoundMetric(respBody, mtr)
		suite.Nil(err)
		suite.True(found)
	}
	var found bool
	found, err = util.FoundMetric(respBody, "data_source_scrape_samples_scrapeds")
	suite.Nil(err)
	suite.False(found)
}

func (suite *SkycoinSuit) TestDefaultMetricsForMetricsFordwaderArePresent() {
	// NOTE(denisacostaq@gmail.com): Giving

	// NOTE(denisacostaq@gmail.com): When
	resp, err := http.Get(suite.rextporterEndpoint)

	// NOTE(denisacostaq@gmail.com): Assert
	suite.Nil(err)
	suite.Equal(http.StatusOK, resp.StatusCode)
	suite.NotNil(resp.Body)
	var respBody []byte
	respBody, err = ioutil.ReadAll(resp.Body)
	suite.Nil(err)
	suite.NotNil(respBody)
	mtrs := []string{
		"fordwader_response_duration_seconds",
		"fordwader_scrape_duration_seconds",
	}
	for _, mtr := range mtrs {
		var found bool
		found, err = util.FoundMetric(respBody, mtr)
		suite.Nil(err)
		suite.True(found)
	}
	var found bool
	found, err = util.FoundMetric(respBody, "fordwader_scrape_duration_secondss")
	suite.Nil(err)
	suite.False(found)
}

func (suite *SkycoinSuit) TestHealthSeq() {
	// NOTE(denisacostaq@gmail.com): Giving

	// NOTE(denisacostaq@gmail.com): When
	resp, err := http.Get(suite.rextporterEndpoint)

	// NOTE(denisacostaq@gmail.com): Assert
	suite.Nil(err)
	suite.Equal(http.StatusOK, resp.StatusCode)
	suite.NotNil(resp.Body)
	var respBody []byte
	respBody, err = ioutil.ReadAll(resp.Body)
	suite.Nil(err)
	suite.NotNil(respBody)
	var found bool
	found, err = util.FoundMetric(respBody, "seq")
	suite.Nil(err)
	suite.True(found)
	var val float64
	val, err = util.GetCounterValue(respBody, "seq")
	suite.Nil(err)
	suite.Equal(float64(180), val)
}

func (suite *SkycoinSuit) TestHealthFee() {
	// NOTE(denisacostaq@gmail.com): Giving

	// NOTE(denisacostaq@gmail.com): When
	resp, err := http.Get(suite.rextporterEndpoint)

	// NOTE(denisacostaq@gmail.com): Assert
	suite.Nil(err)
	suite.Equal(http.StatusOK, resp.StatusCode)
	suite.NotNil(resp.Body)
	var respBody []byte
	respBody, err = ioutil.ReadAll(resp.Body)
	suite.Nil(err)
	suite.NotNil(respBody)
	var found bool
	found, err = util.FoundMetric(respBody, "fee")
	suite.Nil(err)
	suite.True(found)
	var val float64
	val, err = util.GetGaugeValue(respBody, "fee")
	suite.Nil(err)
	suite.Equal(float64(2265261), val)
}

func (suite *SkycoinSuit) TestHealthUnspents() {
	// NOTE(denisacostaq@gmail.com): Giving

	// NOTE(denisacostaq@gmail.com): When
	resp, err := http.Get(suite.rextporterEndpoint)

	// NOTE(denisacostaq@gmail.com): Assert
	suite.Nil(err)
	suite.Equal(http.StatusOK, resp.StatusCode)
	suite.NotNil(resp.Body)
	var respBody []byte
	respBody, err = ioutil.ReadAll(resp.Body)
	suite.Nil(err)
	suite.NotNil(respBody)
	var found bool
	found, err = util.FoundMetric(respBody, "Unspents")
	suite.Nil(err)
	suite.True(found)
	var val float64
	val, err = util.GetGaugeValue(respBody, "Unspents")
	suite.Nil(err)
	suite.Equal(float64(218), val)
}

func (suite *SkycoinSuit) TestHealthUnconfirmed() {
	// NOTE(denisacostaq@gmail.com): Giving

	// NOTE(denisacostaq@gmail.com): When
	resp, err := http.Get(suite.rextporterEndpoint)

	// NOTE(denisacostaq@gmail.com): Assert
	suite.Nil(err)
	suite.Equal(http.StatusOK, resp.StatusCode)
	suite.NotNil(resp.Body)
	var respBody []byte
	respBody, err = ioutil.ReadAll(resp.Body)
	suite.Nil(err)
	suite.NotNil(respBody)
	var found bool
	found, err = util.FoundMetric(respBody, "Unconfirmed")
	suite.Nil(err)
	suite.True(found)
	var val float64
	val, err = util.GetGaugeValue(respBody, "Unconfirmed")
	suite.Nil(err)
	suite.Equal(float64(1), val)
}

func (suite *SkycoinSuit) TestHealthOpenConnections() {
	// NOTE(denisacostaq@gmail.com): Giving

	// NOTE(denisacostaq@gmail.com): When
	resp, err := http.Get(suite.rextporterEndpoint)

	// NOTE(denisacostaq@gmail.com): Assert
	suite.Nil(err)
	suite.Equal(http.StatusOK, resp.StatusCode)
	suite.NotNil(resp.Body)
	var respBody []byte
	respBody, err = ioutil.ReadAll(resp.Body)
	suite.Nil(err)
	suite.NotNil(respBody)
	var found bool
	found, err = util.FoundMetric(respBody, "OpenConnections")
	suite.Nil(err)
	suite.True(found)
	var val float64
	val, err = util.GetGaugeValue(respBody, "OpenConnections")
	suite.Nil(err)
	suite.Equal(float64(0), val)
}

func (suite *SkycoinSuit) TestHealthOutgoingConnections() {
	// NOTE(denisacostaq@gmail.com): Giving

	// NOTE(denisacostaq@gmail.com): When
	resp, err := http.Get(suite.rextporterEndpoint)

	// NOTE(denisacostaq@gmail.com): Assert
	suite.Nil(err)
	suite.Equal(http.StatusOK, resp.StatusCode)
	suite.NotNil(resp.Body)
	var respBody []byte
	respBody, err = ioutil.ReadAll(resp.Body)
	suite.Nil(err)
	suite.NotNil(respBody)
	var found bool
	found, err = util.FoundMetric(respBody, "OutgoingConnections")
	suite.Nil(err)
	suite.True(found)
	var val float64
	val, err = util.GetGaugeValue(respBody, "OutgoingConnections")
	suite.Nil(err)
	suite.Equal(float64(0), val)
}

func (suite *SkycoinSuit) TestHealthIncomingConnections() {
	// NOTE(denisacostaq@gmail.com): Giving

	// NOTE(denisacostaq@gmail.com): When
	resp, err := http.Get(suite.rextporterEndpoint)

	// NOTE(denisacostaq@gmail.com): Assert
	suite.Nil(err)
	suite.Equal(http.StatusOK, resp.StatusCode)
	suite.NotNil(resp.Body)
	var respBody []byte
	respBody, err = ioutil.ReadAll(resp.Body)
	suite.Nil(err)
	suite.NotNil(respBody)
	var found bool
	found, err = util.FoundMetric(respBody, "IncomingConnections")
	suite.Nil(err)
	suite.True(found)
	var val float64
	val, err = util.GetGaugeValue(respBody, "IncomingConnections")
	suite.Nil(err)
	suite.Equal(float64(0), val)
}

func (suite *SkycoinSuit) TestHealthBurnFactor() {
	// NOTE(denisacostaq@gmail.com): Giving

	// NOTE(denisacostaq@gmail.com): When
	resp, err := http.Get(suite.rextporterEndpoint)

	// NOTE(denisacostaq@gmail.com): Assert
	suite.Nil(err)
	suite.Equal(http.StatusOK, resp.StatusCode)
	suite.NotNil(resp.Body)
	var respBody []byte
	respBody, err = ioutil.ReadAll(resp.Body)
	suite.Nil(err)
	suite.NotNil(respBody)
	var found bool
	found, err = util.FoundMetric(respBody, "BurnFactor")
	suite.Nil(err)
	suite.True(found)
	var val float64
	val, err = util.GetGaugeValue(respBody, "BurnFactor")
	suite.Nil(err)
	suite.Equal(float64(2), val)
}

func (suite *SkycoinSuit) TestHealthMaxTransactionSize() {
	// NOTE(denisacostaq@gmail.com): Giving

	// NOTE(denisacostaq@gmail.com): When
	resp, err := http.Get(suite.rextporterEndpoint)

	// NOTE(denisacostaq@gmail.com): Assert
	suite.Nil(err)
	suite.Equal(http.StatusOK, resp.StatusCode)
	suite.NotNil(resp.Body)
	var respBody []byte
	respBody, err = ioutil.ReadAll(resp.Body)
	suite.Nil(err)
	suite.NotNil(respBody)
	var found bool
	found, err = util.FoundMetric(respBody, "MaxTransactionSize")
	suite.Nil(err)
	suite.True(found)
	var val float64
	val, err = util.GetGaugeValue(respBody, "MaxTransactionSize")
	suite.Nil(err)
	suite.Equal(float64(32768), val)
}
