package client

import (
	"github.com/simelo/rextporter/src/config"
	"fmt"
	"github.com/simelo/rextporter/src/common"
	"net/http"
	"io/ioutil"
)

type TokenClient struct {
	BaseClient
}


func NewTokenClient(host config.Host) (client *TokenClient, err error) {
	const generalScopeErr = "error creating a client to get a toke from remote endpoint for making future requests"
	client = new(TokenClient)
	client.host = host
	if client.req, err = http.NewRequest("GET", client.host.UriToGetToken(), nil); err != nil {
		errCause := fmt.Sprintln("can not create the request:", err.Error())
		return nil, common.ErrorFromThisScope(errCause, generalScopeErr)
	}
	return client, nil
}

func (client *TokenClient) GetRemoteInfo() (data []byte, err error) {
	const generalScopeErr = "error making a server request to get metric from remote endpoint"
	httpClient := &http.Client{}
	var resp *http.Response
	if resp, err = httpClient.Do(client.req); err != nil {
		errCause := fmt.Sprintln("can not do the request:", err.Error())
		return nil, common.ErrorFromThisScope(errCause, generalScopeErr)
	}
	if data, err = ioutil.ReadAll(resp.Body); err != nil {
		errCause := fmt.Sprintln("can not read the body:", err.Error())
		return nil, common.ErrorFromThisScope(errCause, generalScopeErr)
	}
	return data, nil
}
