package util

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
)

//Soap defines a soap interface to simplify the unit tests
type Soap func(request, endpoint, proxy string) (string, error)

//SoapCall executes a soap call to the defined environment executing the xog xml
func SoapCall(request, endpoint, proxy string) (string, error) {
	client := &http.Client{
		Timeout: time.Second * 600,
	}

	if proxy != "" {
		proxyURL, err := url.Parse(proxy)
		if err != nil {
			return "", err
		}
		client.Transport = &http.Transport{
			Proxy: http.ProxyURL(proxyURL),
		}
	}

	resp, err := client.Post(endpoint+"/niku/xog", "text/xml; charset=utf-8", bytes.NewBufferString(request))
	if err != nil {
		return "", err
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return BytesToString(body), nil
}
