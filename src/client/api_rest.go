package client

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/oliveagle/jsonpath"
	"github.com/simelo/rextporter/src/config"
	"github.com/simelo/rextporter/src/util"
)

// APIRest have common data to be shared through embedded struct in those type who implement the client.Client interface
type APIRest struct {
	req                  *http.Request
	tokenClient          TokenClient
	tokenHeaderKey       string
	tokenKeyFromEndpoint string
	token                string
}

// CreateAPIRest create client to make request to a rest API
func CreateAPIRest(metric config.Metric, service config.Service) (client APIRest, err error) {
	const generalScopeErr = "error creating api rest client"
	if client.req, err = http.NewRequest(metric.HTTPMethod, service.URIToGetMetric(metric), nil); err != nil {
		errCause := fmt.Sprintln("can not create the request client: ", err.Error())
		return client, util.ErrorFromThisScope(errCause, generalScopeErr)
	}
	if client.tokenClient, err = newTokenClient(service.URIToGetToken()); err != nil {
		errCause := fmt.Sprintln("create token client: ", err.Error())
		return client, util.ErrorFromThisScope(errCause, generalScopeErr)
	}
	client.tokenHeaderKey = service.TokenHeaderKey
	client.tokenKeyFromEndpoint = service.TokenKeyFromEndpoint
	return client, nil
}

// GetData can retrieve data from a rest API with a retry pollicy for token expiration.
func (cl APIRest) GetData() (data []byte, err error) {
	const generalScopeErr = "error making a server request to get metric from remote endpoint"
	cl.req.Header.Set(cl.tokenHeaderKey, cl.token)
	getData := func() (data []byte, err error) {
		httpClient := &http.Client{}
		var resp *http.Response
		if resp, err = httpClient.Do(cl.req); err != nil {
			errCause := fmt.Sprintln("can not do the request: ", err.Error())
			return nil, util.ErrorFromThisScope(errCause, generalScopeErr)
		}
		defer resp.Body.Close()
		if resp.StatusCode != http.StatusOK {
			errCause := fmt.Sprintf("no success response, status %s", resp.Status)
			return nil, util.ErrorFromThisScope(errCause, generalScopeErr)
		}
		if data, err = ioutil.ReadAll(resp.Body); err != nil {
			errCause := fmt.Sprintln("can not read the body: ", err.Error())
			return nil, util.ErrorFromThisScope(errCause, generalScopeErr)
		}
		return data, nil
	}
	if data, err = getData(); err != nil {
		// log.Println("can not do the request:", err.Error(), "trying with a new token...")
		if err = cl.resetToken(); err != nil {
			errCause := fmt.Sprintln("can not reset the token: ", err.Error())
			return nil, util.ErrorFromThisScope(errCause, generalScopeErr)
		}
		if data, err = getData(); err != nil {
			errCause := fmt.Sprintln("can not do the request after a token reset neither: ", err.Error())
			return nil, util.ErrorFromThisScope(errCause, generalScopeErr)
		}
	}
	return data, nil
}

func (cl APIRest) resetToken() (err error) {
	const generalScopeErr = "error making resetting the token"
	cl.token = ""
	var data []byte
	if data, err = cl.tokenClient.GetData(); err != nil {
		errCause := fmt.Sprintln("can make the request to get a token: ", err.Error())
		return util.ErrorFromThisScope(errCause, generalScopeErr)
	}
	var jsonData interface{}
	if err = json.Unmarshal(data, &jsonData); err != nil {
		errCause := fmt.Sprintln("can not decode the body: ", string(data), " ", err.Error())
		return util.ErrorFromThisScope(errCause, generalScopeErr)
	}
	var val interface{}
	jPath := "$" + strings.Replace(cl.tokenKeyFromEndpoint, "/", ".", -1)
	if val, err = jsonpath.JsonPathLookup(jsonData, jPath); err != nil {
		errCause := fmt.Sprintln("can not locate the path: ", err.Error())
		return util.ErrorFromThisScope(errCause, generalScopeErr)
	}
	tk, ok := val.(string)
	if !ok {
		errCause := fmt.Sprintln("unable the get the token as a string value")
		return util.ErrorFromThisScope(errCause, generalScopeErr)
	}
	cl.token = tk
	if len(cl.token) == 0 {
		errCause := fmt.Sprintln("unable the get a not null(empty) token")
		return util.ErrorFromThisScope(errCause, generalScopeErr)
	}
	return nil
}
