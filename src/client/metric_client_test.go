package client

import (
	"testing"
	"net"
	"net/http/httptest"
	"github.com/denisacostaq/rextporter/config"
	"net/http"
	"log"
	"encoding/json"
	//"github.com/stretchr/testify/assert"
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

var systemConfig = `
# All hots to be monitored.
[[hosts]]
ref = "hostname1"
location = "http://127.0.0.1"
port = 8080
authType = "CSRF"
tokenKey = "X-CSRF-Token"
genTokenEndpoint = "/api/v1/csrf"
genTokenKey = "csrf_token"

# All metrics to be measured.
[[metrics]]
name = "uptime"

  [metrics.options]
  type = "Counter"
  description = "I am running since"

# Now you should define what metrics to take care of in what host
[[metricsForHost]]
hostRef = "hostname1"
metric = "uptime"
url = "/api/v1/health"
httpMethod = "GET"
path = "/blockchain/head/seq"
`

func httpHandler(w http.ResponseWriter, r *http.Request) {
		//switch r.RequestURI {
		//case "/latest/meta-data/instance-id":
		//	resp = "i-12345"
		//case "/latest/meta-data/placement/availability-zone":
		//	resp = "us-west-2a"
		//default:
		//	http.Error(w, "not found", http.StatusNotFound)
		//	return
		//}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(jsonResponse))
}

type SkycoinStatsSuit struct {
	suite.Suite
	testServer *httptest.Server
}

func (suite *SkycoinStatsSuit) SetupTest() {
	suite.testServer = stubSkycoin()
	suite.testServer.Start()
	config.NewConfig(systemConfig)
}

func (suite *SkycoinStatsSuit) TearDownTest() {
	suite.testServer.Close()
}

func TestSkycoinStatsSuit(t *testing.T) {
	suite.Run(t, new(SkycoinStatsSuit))
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


func (suit *SkycoinStatsSuit) TestSomething() {
	//assert :=
	//	assert.New(t)

	log.SetFlags(log.LstdFlags | log.Lshortfile)
	conf := config.Config()
	if /*b*/_, err := json.MarshalIndent(conf, "", " "); err != nil {
		log.Println("Error marshalling:", err)
	} else {
		//os.Stdout.Write(b)
	}
	for _, host := range conf.Hosts {
		// cl, err := client.NewTokenClient(host)
		// log.Println("tk:", tk)
		links := conf.FilterLinksByHost(host)
		for _, link := range links {
			if cl, err := NewMetricClient(link); err != nil {
				log.Println(err.Error())
			} else {
				a, _ := cl.GetMetric()

				//assert.NotNil(e)
				//assert.Nil(e)
				suit.Equal( float64(58894), a, "Error en esto")
				log.Println(a)
			}
		}
	}

}
