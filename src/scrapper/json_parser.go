package scrapper

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/oliveagle/jsonpath"
	"github.com/simelo/rextporter/src/util"
)

type JsonParser struct {
}

func (p JsonParser) decodeBody(body []byte) (val interface{}, err error) {
	generalScopeErr := "error decoding body as json"
	if err = json.Unmarshal(body, &val); err != nil {
		errCause := fmt.Sprintf("can not decode the body: %s. Err: %s", string(body), err.Error())
		return nil, util.ErrorFromThisScope(errCause, generalScopeErr)
	}
	return val, err
}

func (p JsonParser) pathLookup(path string, val interface{}) (node interface{}, err error) {
	jPath := "$" + strings.Replace(path, "/", ".", -1)
	if node, err = jsonpath.JsonPathLookup(iBody, jPath); err != nil {
		errCause := fmt.Sprintln("can not locate the path: ", err.Error())
		return nil, util.ErrorFromThisScope(errCause, generalScopeErr)
	}
	return node, err
}
