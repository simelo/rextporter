package client

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/oliveagle/jsonpath"
	"github.com/simelo/rextporter/src/common"
	"github.com/simelo/rextporter/src/config"
)

type token struct {
	// FIXME(denisacostaq@gmail.com): you are just decoding with this key and ignoring the configured 'TokenKeyFromEndpoint'
	CsrfToken string `json:"csrf_token"`
}

// BaseClient have common data to be shared through embedded struct in those type who implement the
// client.Client interface
type BaseClient struct {
	req     *http.Request
	service config.Service
}

// MetricClient implements the getRemoteInfo method from `client.Client` interface by using some `.toml` config parameters
// like for example: where is the host? it should be a GET, a POST or some other? ...
// sa NewMetricClient method.
type MetricClient struct {
	BaseClient
	token       string
	metricJPath string
}

// NewMetricClient will put all the required info to be able to do http requests to get the remote data.
func NewMetricClient(metric config.Metric, conf config.RootConfig) (client *MetricClient, err error) {
	const generalScopeErr = "error creating a client to get a metric from remote endpoint"
	client = new(MetricClient)
	client.BaseClient.service = conf.Service
	client.metricJPath = metric.Path
	client.BaseClient.req, err = http.NewRequest(metric.HTTPMethod, client.service.URIToGetMetric(metric), nil)
	if err != nil {
		errCause := fmt.Sprintln("can not create the request: ", err.Error())
		return nil, common.ErrorFromThisScope(errCause, generalScopeErr)
	}
	return client, nil
}

func (client *MetricClient) resetToken() (err error) {
	const generalScopeErr = "error making resetting the token"
	client.token = ""
	var clientToken *TokenClient
	if clientToken, err = newTokenClient(client.service); err != nil {
		errCause := fmt.Sprintln("can not find a host: ", err.Error())
		return common.ErrorFromThisScope(errCause, generalScopeErr)
	}
	var data []byte
	if data, err = clientToken.getRemoteInfo(); err != nil {
		errCause := fmt.Sprintln("can make the request to get a token: ", err.Error())
		return common.ErrorFromThisScope(errCause, generalScopeErr)
	}
	var tk token
	if err = json.Unmarshal(data, &tk); err != nil {
		errCause := fmt.Sprintln("error decoding the server response: ", err.Error())
		return common.ErrorFromThisScope(errCause, generalScopeErr)
	}
	client.token = tk.CsrfToken
	if strings.Compare(client.token, "") == 0 {
		errCause := fmt.Sprintln("unable the get a not null(empty) token")
		return common.ErrorFromThisScope(errCause, generalScopeErr)
	}
	return nil
}

func (client *MetricClient) getRemoteInfo() (data []byte, err error) {
	const generalScopeErr = "error making a server request to get metric from remote endpoint"
	client.req.Header.Set(client.service.TokenHeaderKey, client.token)
	doRequest := func() (*http.Response, error) {
		httpClient := &http.Client{}
		var resp *http.Response
		if resp, err = httpClient.Do(client.req); err != nil {
			errCause := fmt.Sprintln("can not do the request: ", err.Error())
			return nil, common.ErrorFromThisScope(errCause, generalScopeErr)
		}
		return resp, nil
	}
	var resp *http.Response
	if resp, err = doRequest(); err != nil {
		// log.Println("can not do the request:", err.Error(), "trying with a new token...")
		if err = client.resetToken(); err != nil {
			errCause := fmt.Sprintln("can not reset the token: ", err.Error())
			return nil, common.ErrorFromThisScope(errCause, generalScopeErr)
		}
		if resp, err = doRequest(); err != nil {
			errCause := fmt.Sprintln("can not do the request after a token reset neither: ", err.Error())
			return nil, common.ErrorFromThisScope(errCause, generalScopeErr)
		}
	}
	if data, err = ioutil.ReadAll(resp.Body); err != nil {
		errCause := fmt.Sprintln("can not read the body: ", err.Error())
		return nil, common.ErrorFromThisScope(errCause, generalScopeErr)
	}
	return data, nil
}

// GetMetric returns the metric previously bound through config parameters like:
// url(endpoint), json path, type and so on.
func (client *MetricClient) GetMetric() (val interface{}, err error) {
	const generalScopeErr = "error getting metric data"
	var data []byte
	if data, err = client.getRemoteInfo(); err != nil {
		return nil, common.ErrorFromThisScope(err.Error(), generalScopeErr)
	}
	var jsonData interface{}
	if err = json.Unmarshal(data, &jsonData); err != nil {
		errCause := fmt.Sprintln("can not decode the body: ", string(data), " ", err.Error())
		return nil, common.ErrorFromThisScope(errCause, generalScopeErr)
	}
	jPath := "$" + strings.Replace(client.metricJPath, "/", ".", -1)
	if val, err = jsonpath.JsonPathLookup(jsonData, jPath); err != nil {
		errCause := fmt.Sprintln("can not locate the path: ", err.Error())
		return nil, common.ErrorFromThisScope(errCause, generalScopeErr)
	}
	return val, nil
}
