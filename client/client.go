package client

import (
	"github.com/denisacostaq/rextporter/config"
	"net/http"
	"log"
	"encoding/json"
	"io/ioutil"
	"github.com/oliveagle/jsonpath"
	"strings"
)

type Error struct {
	Msg string
	Code int
}

type token struct {
	CsrfToken string `json:"csrf_token"`
}

func (err Error) Error() string {
	return err.Msg
}

func GetToken(conf config.Host) (string, error) {
	genTokenUri := config.UriToGetToken(conf)
	var resp *http.Response
	var err error
	if resp, err = http.Get(genTokenUri); err != nil {
		log.Println("Error getting the token:", err)
		return "", Error{Msg: "Unable to get the token"}
	}
	var tk token
	if err = json.NewDecoder(resp.Body).Decode(&tk); err != nil {
		log.Println("Error decoding the token:", err)
		return "", Error{Msg: "Unable to get the token"}
	}
	return tk.CsrfToken, nil
}

func GetMetric(host config.Host, link config.Link, tk string) interface{} {
	var req *http.Request
	var err error
	req, err = http.NewRequest(link.HttpMethod, config.UriToGetMetric(host, link), nil)
	if err != nil {
		log.Println("Can not create the request:", err)
	}
	req.Header.Add(host.TokenKey, tk)
	client := &http.Client{}
	var resp *http.Response
	resp, err = client.Do(req)
	if err != nil {
		log.Println("Can not do the request:", err)
	}
	var data []byte
	data, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("Can not read the body:", err)
	}
	var json_data interface{}
	err = json.Unmarshal(data, &json_data)
	if err != nil {
		log.Println("Can not decode the body:", err)
	}
	var val interface{}
	jpath := "$" + strings.Replace(link.Path, "/", ".", -1)
	val, err = jsonpath.JsonPathLookup(json_data, jpath)
	if err != nil {
		log.Println("Can not locate the path", err)
	}
	return val
}
