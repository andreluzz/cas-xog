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
type Rest func(jsonString []byte, config APIConfig, params map[string]string) ([]byte, int, error)

// APIConfig definitions to realize rest api requests
type APIConfig struct {
	Endpoint string
	Method   string
	Token    string
	Proxy    string
	Cookie   string
	Client   string
}

//RestCall executes a rest call to the defined environment executing a json
func RestCall(jsonString []byte, config APIConfig, params map[string]string) ([]byte, int, error) {
	if config.Token == "" {
		return nil, -1, fmt.Errorf("invalid token")
	}
	var body io.Reader
	if jsonString != nil {
		body = bytes.NewBuffer(jsonString)
	}
	req, err := http.NewRequest(config.Method, config.Endpoint, body)
	if err != nil {
		return nil, -1, err
	}

	if config.Client != "" {
		req.Header.Add("x-api-ppm-client", config.Client)
		req.Header.Add("Authorization", config.Token)
	} else {
		req.Header.Add("authToken", config.Token)
	}

	req.Header.Add("content-type", "application/json")
	req.Header.Add("x-api-force-patch", "true")
	if config.Cookie != "" {
		req.Header.Add("Cookie", config.Cookie)
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

	if config.Proxy != "" {
		proxyURL, err := url.Parse(config.Proxy)
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
