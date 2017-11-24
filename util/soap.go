package util

import (
	"bytes"
	"io/ioutil"
	"net/http"
)

type Soap func(request, endpoint string) (string, error)

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