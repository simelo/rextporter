package client

import (
	"compress/gzip"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/simelo/rextporter/src/config"
	"github.com/simelo/rextporter/src/util"
)

// ProxyMetricClient implements the getRemoteInfo method from `client.Client` interface by using some `.toml` config parameters
// like for example: where is the host. It get the exposed metrics from a service as is.
type ProxyMetricClient struct {
	req         *http.Request
	ServiceName string
}

// NewProxyMetricClient will put all the required info to be able to do http requests to get the remote data.
func NewProxyMetricClient(service config.Service) (client ProxyMetricClient, err error) {
	const generalScopeErr = "error creating a forward_metrics client to get the metrics from remote endpoint"
	if !util.StrSliceContains(service.Modes, config.ServiceTypeProxy) {
		return client, errors.New("can not create a forward_metrics metric client from a service whitout type " + config.ServiceTypeProxy)
	}
	client.req, err = http.NewRequest("GET", service.URIToGetExposedMetric(), nil)
	if err != nil {
		errCause := fmt.Sprintln("can not create the request: ", err.Error())
		return client, util.ErrorFromThisScope(errCause, generalScopeErr)
	}
	client.ServiceName = service.Name
	return client, nil
}

func (client ProxyMetricClient) GetData() (data []byte, err error) {
	const generalScopeErr = "error making a server request to get the metrics from remote endpoint"
	httpClient := &http.Client{}
	var resp *http.Response
	if resp, err = httpClient.Do(client.req); err != nil {
		errCause := fmt.Sprintln("can not do the request: ", err.Error())
		return nil, util.ErrorFromThisScope(errCause, generalScopeErr)
	}
	if resp.StatusCode != http.StatusOK {
		errCause := fmt.Sprintf("no success response, status %s", resp.Status)
		return nil, util.ErrorFromThisScope(errCause, generalScopeErr)
	}
	defer resp.Body.Close()
	var reader io.ReadCloser
	var isGzipContent = false
	defer func() {
		if isGzipContent {
			reader.Close()
		}
	}()
	switch resp.Header.Get("Content-Encoding") {
	case "gzip":
		isGzipContent = true
		reader, err = gzip.NewReader(resp.Body)
		if err != nil {
			errCause := fmt.Sprintln("can not create gzip reader.", err.Error())
			return nil, util.ErrorFromThisScope(errCause, generalScopeErr)
		}
		defer reader.Close()
	default:
		reader = resp.Body
	}
	// FIXME(denisacostaq@gmail.com): write an integration test for plain text and compressed content
	if data, err = ioutil.ReadAll(reader); err != nil {
		errCause := fmt.Sprintln("can not read the body: ", err.Error())
		return nil, util.ErrorFromThisScope(errCause, generalScopeErr)
	}
	return data, nil
}
