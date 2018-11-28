package client

import (
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/simelo/rextporter/src/util"
)

// TokenClient implements the getRemoteInfo method from `client.Client` interface by using some .toml config parameters
// like for example: where is the host? it should be a GET, a POST or some other ... It works like an http wrapper to
// to get a token from the server.
// sa newTokenClient method.
type TokenClient struct {
	req *http.Request
}

func newTokenClient(uriToGenToken string) (client TokenClient, err error) {
	const generalScopeErr = "error creating a client to get a toke from remote endpoint for making future requests"
	// client = TokenClient{baseClient: baseClient{service: service}}
	// FIXME(denisacostaq@gmail.com): make the "GET" configurable.
	if client.req, err = http.NewRequest("GET", uriToGenToken, nil); err != nil {
		errCause := fmt.Sprintln("can not create the request: ", err.Error())
		return client, util.ErrorFromThisScope(errCause, generalScopeErr)
	}
	return client, nil
}

// GetData can get a token value from a remote server
func (client TokenClient) GetData() (data []byte, err error) {
	const generalScopeErr = "error making a server request to get token from remote endpoint"
	httpClient := &http.Client{}
	var resp *http.Response
	if resp, err = httpClient.Do(client.req); err != nil {
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
