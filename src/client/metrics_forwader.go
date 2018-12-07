package client

import (
	"compress/gzip"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/simelo/rextporter/src/config"
	"github.com/simelo/rextporter/src/util"
)

// ProxyMetricClientCreator create a metrics fordwader client
type ProxyMetricClientCreator struct {
	dataPath     string
	JobName      string
	InstanceName string
}

// CreateProxyMetricClientCreator create a ProxyMetricClientCreator with required info to create a metrics fordwader client
func CreateProxyMetricClientCreator(service config.Service) (cf ProxyMetricClientCreator, err error) {
	if !util.StrSliceContains(service.Modes, config.ServiceTypeProxy) {
		return ProxyMetricClientCreator{}, errors.New("can not create a forward_metrics metric client from a service whitout type " + config.ServiceTypeProxy)
	}
	cf = ProxyMetricClientCreator{
		dataPath:     service.URIToGetExposedMetric(),
		JobName:      service.JobName(),
		InstanceName: service.InstanceName(),
	}
	return cf, err
}

// CreateClient create a metrics forwader client
func (pmc ProxyMetricClientCreator) CreateClient() (cl Client, err error) {
	const generalScopeErr = "error creating a forward_metrics client to get the metrics from remote endpoint"
	var req *http.Request
	if req, err = http.NewRequest("GET", pmc.dataPath, nil); err != nil {
		errCause := fmt.Sprintln("can not create the request: ", err.Error())
		return cl, util.ErrorFromThisScope(errCause, generalScopeErr)
	}
	cl = ProxyMetricClient{
		req: req,
	}
	return cl, err
}

// ProxyMetricClient implements the getRemoteInfo method from `client.Client` interface by using some `.toml` config parameters
// like for example: where is the host. It get the exposed metrics from a service as is.
type ProxyMetricClient struct {
	req *http.Request
}

// GetData can get raw metrics from a endpoint
func (client ProxyMetricClient) GetData(metricsCollector chan<- prometheus.Metric) (data []byte, err error) {
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
