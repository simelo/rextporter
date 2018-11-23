package scrapper

import (
	"encoding/json"
	"fmt"

	"github.com/simelo/rextporter/src/util"
)

type JsonParser struct {
}

func (parser JsonParser) DecodeBody(body []byte) (val interface{}, err error) {
	generalScopeErr := "error decoding body as json"
	if err = json.Unmarshal(body, &val); err != nil {
		errCause := fmt.Sprintf("can not decode the body: %s. Err: %s", string(body), err.Error())
		return nil, util.ErrorFromThisScope(errCause, generalScopeErr)
	}
	return val, err
}
