package util

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
)

//Rest defines a rest interface to simplify the unit tests
type Rest func(jsonString []byte, endpoint, method, token, proxy, cookie string, params map[string]string) ([]byte, int, error)

//RestCall executes a rest call to the defined environment executing a json
func RestCall(jsonString []byte, endpoint, method, token, proxy, cookie string, params map[string]string) ([]byte, int, error) {
	if token == "" {
		return nil, -1, fmt.Errorf("invalid token")
	}
	var body io.Reader
	if jsonString != nil {
		body = bytes.NewBuffer(jsonString)
	}
	req, err := http.NewRequest(method, endpoint, body)
	if err != nil {
		return nil, -1, err
	}
	req.Header.Add("authToken", token)
	req.Header.Add("content-type", "application/json")
	req.Header.Add("x-api-force-patch", "true")
	if cookie != "" {
		req.Header.Add("Cookie", cookie)
	}

	q := req.URL.Query()
	if params != nil {
		for key, value := range params {
			q.Add(key, value)
		}
	}

	req.URL.RawQuery = q.Encode()

	client := &http.Client{
		Timeout: time.Second * 60,
	}

	if proxy != "" {
		proxyURL, err := url.Parse(proxy)
		if err != nil {
			return nil, -1, err
		}
		client.Transport = &http.Transport{
			Proxy: http.ProxyURL(proxyURL),
		}
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, -1, err
	}

	data, err := ioutil.ReadAll(resp.Body)
	return data, resp.StatusCode, err
}

//APIPostLogin send a post to get the auth token
func APIPostLogin(endpoint, username, password, proxy, cookie string) ([]byte, error) {
	req, err := http.NewRequest(http.MethodPost, endpoint, nil)
	if err != nil {
		return nil, err
	}
	req.SetBasicAuth(username, password)
	if cookie != "" {
		req.Header.Add("Cookie", cookie)
	}

	client := &http.Client{
		Timeout: time.Second * 60,
	}

	if proxy != "" {
		proxyURL, err := url.Parse(proxy)
		if err != nil {
			return nil, err
		}
		client.Transport = &http.Transport{
			Proxy: http.ProxyURL(proxyURL),
		}
	}

	resp, err := client.Do(req)

	return ioutil.ReadAll(resp.Body)
}
