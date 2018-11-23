package scrapper

// Scrapper receive some data as input and should return the metric val
type Scrapper interface {
	GetMetric() (val interface{}, err error)
}

func getData(cl client.Client, p client.Parser) (data interface{}, err error) {
	const generalScopeErr = "error getting data"
	var body []byte
	if body, err = cl.GetData() {
		errCause := "client can not get data"
		return val, util.ErrorFromThisScope(errCause, generalScopeErr)
	}
	if data, err = p.DecodeBody(body); err != nil {
		errCause := "client can not decode the body"
		return val, util.ErrorFromThisScope(errCause, generalScopeErr)
	}
	return data, err
}