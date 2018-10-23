package client

import (
	"github.com/denisacostaq/rextporter/config"
	"github.com/denisacostaq/rextporter/common"
	"net/http"
	"encoding/json"
	"io/ioutil"
	"github.com/oliveagle/jsonpath"
	"strings"
	"fmt"
	"errors"
)

type token struct {
	CsrfToken string `json:"csrf_token"`
}

func GetToken(conf config.Host) (string, error) {
	const generalScopeErr = "error getting the token"
	genTokenUri := config.UriToGetToken(conf)
	var resp *http.Response
	var err error
	if resp, err = http.Get(genTokenUri); err != nil {
		errCause := fmt.Sprintln("server response error:", err.Error())
		return "", common.ErrorFromThisScope(errCause, generalScopeErr)
	}
	var tk token
	if err = json.NewDecoder(resp.Body).Decode(&tk); err != nil {
		errCause := fmt.Sprintln("error decoding the server response:", err.Error())
		return "", common.ErrorFromThisScope(errCause, generalScopeErr)
	}
	return tk.CsrfToken, nil
}


func getRemoteInfo(link config.Link, tk string) (data []byte, err error) {
	const generalScopeErr = "error getting remote data"
	var host config.Host
	host, err = findHostByRef(link.HostRef)
	if err != nil {
		errCause := fmt.Sprintln("can not find a host", err.Error())
		return nil, common.ErrorFromThisScope(errCause, generalScopeErr)
	}
	var req *http.Request
	req, err = http.NewRequest(link.HttpMethod, config.UriToGetMetric(host, link), nil)
	if err != nil {
		errCause := fmt.Sprintln("can not create the request:", err.Error())
		return nil, common.ErrorFromThisScope(errCause, generalScopeErr)
	}
	req.Header.Add(host.TokenKey, tk)
	client := &http.Client{}
	var resp *http.Response
	resp, err = client.Do(req)
	if err != nil {
		errCause := fmt.Sprintln("can not do the request:", err.Error())
		return nil, common.ErrorFromThisScope(errCause, generalScopeErr)
	}
	data, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		errCause := fmt.Sprintln("can not read the body:", err.Error())
		return nil, common.ErrorFromThisScope(errCause, generalScopeErr)
	}
	return data, nil
}

func GetMetric(link config.Link, tk string) (interface{}, error) {
	const generalScopeErr = "error getting metric data"
	data, err := getRemoteInfo(link, tk)
	if err != nil {
		return nil, common.ErrorFromThisScope(err.Error(), generalScopeErr)
	}
	var jsonData interface{}
	err = json.Unmarshal(data, &jsonData)
	if err != nil {
		errCause := fmt.Sprintln("can not decode the body:", err.Error())
		return nil, common.ErrorFromThisScope(errCause, generalScopeErr)
	}
	var val interface{}
	jpath := "$" + strings.Replace(link.Path, "/", ".", -1)
	val, err = jsonpath.JsonPathLookup(jsonData, jpath)
	if err != nil {
		errCause := fmt.Sprintln("can not locate the path:", err.Error())
		return nil, common.ErrorFromThisScope(errCause, generalScopeErr)
	}
	return val, nil
}

func findHostByRef(ref string) (host config.Host, err error) {
	found := false
	for _, host = range config.Config().Hosts {
		found = strings.Compare(host.Ref, ref) == 0
		if found {
			return
		}
	}
	if !found {
		errCause := fmt.Sprintln("can not find a host for Ref:", ref)
		err = errors.New(errCause)
	}
	return config.Host{}, err
}
