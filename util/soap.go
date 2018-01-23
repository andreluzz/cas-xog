package util

import (
	"bytes"
	"io/ioutil"
	"net/http"
)

//Soap defines a soap interface to simplify the unit tests
type Soap func(request, endpoint string) (string, error)

//SoapCall executes a soap call to the defined environment executing the xog xml
func SoapCall(request, endpoint string) (string, error) {
	httpClient := new(http.Client)
	resp, err := httpClient.Post(endpoint+"/niku/xog", "text/xml; charset=utf-8", bytes.NewBufferString(request))
	if err != nil {
		return "", err
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return BytesToString(body), nil
}
