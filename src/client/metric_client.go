package client

import (
	"github.com/denisacostaq/rextporter/src/config"
	"github.com/denisacostaq/rextporter/src/common"
	"net/http"
	"encoding/json"
	"io/ioutil"
	"github.com/oliveagle/jsonpath"
	"strings"
	"fmt"
	"log"
)

type token struct {
	CsrfToken string `json:"csrf_token"`
}

type BaseClient struct {
	req *http.Request
	host config.Host
}

type MetricClient struct {
	BaseClient
	link config.Link
	token string
}

func NewMetricClient(link config.Link) (client *MetricClient, err error) {
	const generalScopeErr = "error creating a client to get a metric from remote endpoint"
	client = new(MetricClient)
	client.link = link
	client.host, err = config.Config().FindHostByRef(link.HostRef)
	if err != nil {
		errCause := fmt.Sprintln("can not find a host", err.Error())
		return nil, common.ErrorFromThisScope(errCause, generalScopeErr)
	}
	client.req, err = http.NewRequest(link.HttpMethod, client.host.UriToGetMetric(link), nil)
	if err != nil {
		errCause := fmt.Sprintln("can not create the request:", err.Error())
		return nil, common.ErrorFromThisScope(errCause, generalScopeErr)
	}
	return client, nil
}

func (client *MetricClient) ResetToken() (err error) {
	const generalScopeErr = "error making resetting the token"
	client.token = ""
	var clientToken *TokenClient
	clientToken, err = NewTokenClient(client.host)
	var data []byte
	if data, err = clientToken.GetRemoteInfo(); err != nil {
		errCause := fmt.Sprintln("can make the request to get a token:", err.Error())
		return common.ErrorFromThisScope(errCause, generalScopeErr)
	}
	var tk token
	if err = json.Unmarshal(data, &tk); err != nil {
		errCause := fmt.Sprintln("error decoding the server response:", err.Error())
		return common.ErrorFromThisScope(errCause, generalScopeErr)
	}
	client.token = tk.CsrfToken
	return nil
}

func (client *MetricClient) GetRemoteInfo() (data []byte, err error) {
	const generalScopeErr = "error making a server request to get metric from remote endpoint"
	client.req.Header.Set(client.host.TokenKey, client.token)
	doRequest := func() (*http.Response, error) {
		httpClient := &http.Client{}
		var resp *http.Response
		if resp, err = httpClient.Do(client.req); err != nil {
			errCause := fmt.Sprintln("can not do the request:", err.Error())
			return nil, common.ErrorFromThisScope(errCause, generalScopeErr)
		}
		return resp, nil
	}
	var resp *http.Response
	if resp, err = doRequest(); err != nil {
		log.Println("can not do the request:", err.Error(), "trying with a new token...")
		if err = client.ResetToken(); err != nil {
			errCause := fmt.Sprintln("can not reset the token:", err.Error())
			return nil, common.ErrorFromThisScope(errCause, generalScopeErr)
		}
		if resp, err = doRequest(); err != nil {
			errCause := fmt.Sprintln("can not do the request after a token reset neither:", err.Error())
			return nil, common.ErrorFromThisScope(errCause, generalScopeErr)
		}
	}
	if data, err = ioutil.ReadAll(resp.Body); err != nil {
		errCause := fmt.Sprintln("can not read the body:", err.Error())
		return nil, common.ErrorFromThisScope(errCause, generalScopeErr)
	}
	return data, nil
}

func (client *MetricClient) GetMetric() (val interface{}, err error) {
	const generalScopeErr = "error getting metric data"
	var data []byte
	if data, err = client.GetRemoteInfo(); err != nil {
		return nil, common.ErrorFromThisScope(err.Error(), generalScopeErr)
	}
	var jsonData interface{}
	if err = json.Unmarshal(data, &jsonData); err != nil {
		errCause := fmt.Sprintln("can not decode the body:", string(data), err.Error())
		return nil, common.ErrorFromThisScope(errCause, generalScopeErr)
	}
	jpath := "$" + strings.Replace(client.link.Path, "/", ".", -1)
	if val, err = jsonpath.JsonPathLookup(jsonData, jpath); err != nil {
		errCause := fmt.Sprintln("can not locate the path:", err.Error())
		return nil, common.ErrorFromThisScope(errCause, generalScopeErr)
	}
	return val, nil
}
