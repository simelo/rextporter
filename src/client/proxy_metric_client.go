package client

import (
	"compress/gzip"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/simelo/rextporter/src/common"
	"github.com/simelo/rextporter/src/config"
)

// ProxyMetricClient implements the getRemoteInfo method from `client.Client` interface by using some `.toml` config parameters
// like for example: where is the host. It get the exposed metrics from a service as is.
type ProxyMetricClient struct {
	BaseClient
	Name string
}

// NewProxyMetricClient will put all the required info to be able to do http requests to get the remote data.
func NewProxyMetricClient(service config.Service) (client *ProxyMetricClient, err error) {
	const generalScopeErr = "error creating a proxy client to get the metrics from remote endpoint"
	if strings.Compare(service.Mode, config.ServiceTypeProxy) != 0 {
		return client, errors.New("can not create a proxy metric client from a service of type " + service.Mode)
	}
	client = new(ProxyMetricClient)
	client.BaseClient.service = service
	client.Name = service.Name
	client.BaseClient.req, err = http.NewRequest("GET", service.URIToGetExposedMetric(), nil)
	if err != nil {
		errCause := fmt.Sprintln("can not create the request: ", err.Error())
		return nil, common.ErrorFromThisScope(errCause, generalScopeErr)
	}
	return client, nil
}

func (client *ProxyMetricClient) getRemoteInfo() (data []byte, err error) {
	const generalScopeErr = "error making a server request to get the metrics from remote endpoint"
	httpClient := &http.Client{}
	var resp *http.Response
	if resp, err = httpClient.Do(client.req); err != nil {
		errCause := fmt.Sprintln("can not do the request: ", err.Error())
		return nil, common.ErrorFromThisScope(errCause, generalScopeErr)
	}
	if resp.StatusCode != http.StatusOK {
		errCause := fmt.Sprintf("no success response, status %s", resp.Status)
		return nil, common.ErrorFromThisScope(errCause, generalScopeErr)
	}
	defer resp.Body.Close()
	var reader io.ReadCloser
	// BUG(denisacostaq@gmail.com): close this reader.
	// defer reader.Close()
	switch resp.Header.Get("Content-Encoding") {
	case "gzip":
		reader, err = gzip.NewReader(resp.Body)
		if err != nil {
			errCause := fmt.Sprintln("can not create gzip reader.", err.Error())
			return nil, common.ErrorFromThisScope(errCause, generalScopeErr)
		}
	default:
		reader = resp.Body
	}
	// FIXME(denisacostaq@gmail.com): write an integration test for plain text and compressed content
	if data, err = ioutil.ReadAll(reader); err != nil {
		errCause := fmt.Sprintln("can not read the body: ", err.Error())
		return nil, common.ErrorFromThisScope(errCause, generalScopeErr)
	}
	return data, nil
}

// GetExposedMetrics returns the metrics in the defined service.
func (client *ProxyMetricClient) GetExposedMetrics() (data []byte, err error) {
	const generalScopeErr = "error getting metrics data"
	if data, err = client.getRemoteInfo(); err != nil {
		return data, common.ErrorFromThisScope(err.Error(), generalScopeErr)
	}
	return data, nil
}
