package common

import (
	"bytes"
	"github.com/beevik/etree"
	"io/ioutil"
	"net/http"
)

func SoapCall(request, endpoint string) (*etree.Document, error) {
	httpClient := new(http.Client)
	resp, err := httpClient.Post(endpoint+"/niku/xog", "text/xml; charset=utf-8", bytes.NewBufferString(request))
	if err != nil {
		return nil, err
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	doc := etree.NewDocument()
	doc.ReadFromBytes(body)

	return doc, nil
}
