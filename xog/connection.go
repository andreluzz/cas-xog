package xog

import (
	"errors"
	"github.com/beevik/etree"
	"github.com/andreluzz/cas-xog/common"
)

func Login(env *EnvType) (string, error) {
	request := etree.NewDocument()
	request.SetRoot(GetLoginXML())

	request.FindElement("//obj:Username").SetText(env.Username)
	request.FindElement("//obj:Password").SetText(env.Password)

	body, err := request.WriteToString()

	if err != nil {
		return "", errors.New("Problems getting login xml: " + err.Error())
	}

	resp, err := common.SoapCall(body, env.URL)

	if err != nil {
		return "", errors.New("Problems trying to log into environment: " + env.Name + " | Debug: " + err.Error())
	}

	sessionElement := resp.FindElement("//SessionID")

	if sessionElement == nil {
		responseError, _ := resp.WriteToString()
		return "", errors.New("Problems trying to log into environment: " + env.Name + "\n" + responseError)
	}

	return sessionElement.Text(), nil
}

func Logout(env *EnvType) error {
	request := etree.NewDocument()
	request.SetRoot(GetLogoutXML())
	request.FindElement("//obj:SessionID").SetText(env.Session)

	body, err := request.WriteToString()

	if err != nil {
		return errors.New("Problems getting logout xml: " + err.Error())
	}

	_, err = common.SoapCall(body, env.URL)

	if err != nil {
		return errors.New("Problems trying to logout from environment: " + env.Name + " | Debug: " + err.Error())
	}

	return nil
}
