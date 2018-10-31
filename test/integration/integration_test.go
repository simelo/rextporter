package client

import (
	"context"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/simelo/rextporter/src/exporter"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

var jsonResponse = `
{
    "blockchain": {
        "head": {
            "seq": 58894,
            "block_hash": "3961bea8c4ab45d658ae42effd4caf36b81709dc52a5708fdd4c8eb1b199a1f6",
            "previous_block_hash": "8eca94e7597b87c8587286b66a6b409f6b4bf288a381a56d7fde3594e319c38a",
            "timestamp": 1537581604,
            "fee": 485194,
            "version": 0,
            "tx_body_hash": "c03c0dd28841d5aa87ce4e692ec8adde923799146ec5504e17ac0c95036362dd",
            "ux_hash": "f7d30ecb49f132283862ad58f691e8747894c9fc241cb3a864fc15bd3e2c83d3"
        },
        "unspents": 38171,
        "unconfirmed": 1,
        "time_since_last_block": "4m46s"
    },
    "version": {
        "version": "0.24.1",
        "commit": "8798b5ee43c7ce43b9b75d57a1a6cd2c1295cd1e",
        "branch": "develop"
    },
    "open_connections": 8,
    "uptime": "6m30.629057248s",
    "csrf_enabled": true,
    "csp_enabled": true,
    "wallet_api_enabled": true,
    "gui_enabled": true,
    "unversioned_api_enabled": false,
    "json_rpc_enabled": false
}
`

func httpHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	if _, err := w.Write([]byte(jsonResponse)); err != nil {
		log.Panicln(err)
	}
}

type HealthSuit struct {
	suite.Suite
	testServer *httptest.Server
	srv        *http.Server
}

func (suite *HealthSuit) SetupSuite() {
	suite.testServer = stubSkycoin()
	suite.testServer.Start()

	require := require.New(suite.T())
	gopath := os.Getenv("GOPATH")
	configFilePath := gopath + "/src/github.com/simelo/rextporter/test/integration/simple.toml"
	suite.srv = exporter.ExportMetrics(configFilePath, 8081)
	for i := 0; i < 3; i++ {
		t := time.NewTimer(time.Second)
		<-t.C
	}
	require.NotNil(suite.srv)
}

func (suite *HealthSuit) TearDownSuite() {
	suite.testServer.Close()
	require := require.New(suite.T())
	var usingAVariableToMakeLinterHappy = context.Context(nil)
	require.Nil(suite.srv.Shutdown(usingAVariableToMakeLinterHappy))
}

func TestSkycoinHealthSuit(t *testing.T) {
	suite.Run(t, new(HealthSuit))
}

func stubSkycoin() *httptest.Server {
	l, err := net.Listen("tcp", "127.0.0.1:8080")
	if err != nil {
		log.Fatal(err)
	}
	testServer := httptest.NewUnstartedServer(http.HandlerFunc(httpHandler))
	testServer.Listener.Close()
	testServer.Listener = l
	return testServer
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
