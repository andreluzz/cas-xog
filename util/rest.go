package util

import (
	"bytes"
	"io"
	"io/ioutil"
	"net/http"
)

//Rest defines a rest interface to simplify the unit tests
type Rest func(jsonString []byte, endpoint, method, token string) ([]byte, error)

//RestCall executes a rest call to the defined environment executing a json
func RestCall(jsonString []byte, endpoint, method, token string) ([]byte, error) {
	var body io.Reader
	if jsonString != nil {
		body = bytes.NewBuffer(jsonString)
	}
	req, err := http.NewRequest(method, endpoint, body)
	if err != nil {
		return nil, err
	}
	req.Header.Add("authToken", token)
	req.Header.Add("content-type", "application/json")
	req.Header.Add("x-api-force-patch", "true")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	return ioutil.ReadAll(resp.Body)
}

//APIPostLogin send a post to get the auth token
func APIPostLogin(endpoint, username, password string) ([]byte, error) {
	req, err := http.NewRequest(http.MethodPost, endpoint, nil)
	if err != nil {
		return nil, err
	}
	req.SetBasicAuth(username, password)

	client := &http.Client{}
	resp, err := client.Do(req)

	return ioutil.ReadAll(resp.Body)
}
