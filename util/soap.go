package util

import (
	"bytes"
	"crypto/tls"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
)

// SoapOptions define soap connections options
type SoapOptions struct {
	InsecureSkipVerify bool
}

//Soap defines a soap interface to simplify the unit tests
type Soap func(request, endpoint, proxy string, opts ...SoapOptions) (string, error)

//SoapCall executes a soap call to the defined environment executing the xog xml
func SoapCall(request, endpoint, proxy string, opts ...SoapOptions) (string, error) {
	client := &http.Client{
		Timeout: time.Second * 600,
	}
	if opts != nil && len(opts) > 0 && opts[0].InsecureSkipVerify {
		customTransport := &(*http.DefaultTransport.(*http.Transport))
		customTransport.TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
		client.Transport = customTransport
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
