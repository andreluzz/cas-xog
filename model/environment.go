package model

import (
	"errors"
	"github.com/andreluzz/cas-xog/util"
	"github.com/beevik/etree"
)

var envXml *etree.Document

func LoadEnvironmentsList(path string, environments *Environments) error {
	envXml = etree.NewDocument()
	err := envXml.ReadFromFile(path)
	environments.Available = envXml.FindElements("./xogenvs/env")
	return err
}

type EnvType struct {
	Name     string
	URL      string
	Username string
	Password string
	Session  string
	Copy     bool
}

func (e *EnvType) init(envIndex string, username string, password string) (error, bool) {
	noLogin := false
	envElement := envXml.FindElement("//env[" + envIndex + "]").Copy()
	if envElement == nil {
		return errors.New("trying to initiate an invalid environment"), noLogin
	}

	e.Name = envElement.SelectAttrValue("name", "Environment name not defined")

	envElemUsername := envElement.FindElement("//username")
	if envElemUsername != nil {
		e.Username = envElemUsername.Text()
	}

	envElemPassword := envElement.FindElement("//password")
	if envElemPassword != nil {
		e.Password = envElemPassword.Text()
	}

	if envElemUsername == nil || envElemPassword == nil {
		if username == "" || password == "" {
			noLogin = true
			return nil, noLogin
		} else {
			requestLoginElement := docXogReadXML.FindElement("//xogtype[@type='requestLogin']").Copy()

			reqLogElemUsername := requestLoginElement.FindElement("//username")
			reqLogElemPassword := requestLoginElement.FindElement("//password")

			reqLogElemUsername.SetText(username)
			reqLogElemPassword.SetText(password)

			envXml.FindElement("//env[" + envIndex + "]").AddChild(reqLogElemUsername)
			envXml.FindElement("//env[" + envIndex + "]").AddChild(reqLogElemPassword)

			e.Username = username
			e.Password = password
		}
	}

	e.URL = envElement.FindElement("//endpoint").Text()

	session, err := login(e)
	if err != nil {
		return err, noLogin
	}
	e.Session = session
	e.Copy = false

	return nil, noLogin
}

func (e *EnvType) logout() error {
	if e == nil {
		return nil
	}
	if e.Copy {
		e.clear()
		return nil
	}
	if e.Session == "" {
		e.clear()
		return nil
	}
	err := logout(e)
	e.clear()
	return err
}

func (e *EnvType) copyEnv() *EnvType {
	ne := &EnvType{
		Name:     e.Name,
		Username: e.Username,
		Password: e.Password,
		URL:      e.URL,
		Session:  e.Session,
		Copy:     true,
	}
	return ne
}

func (e *EnvType) clear() error {
	e.Username = ""
	e.Password = ""
	e.URL = ""
	e.Session = ""
	e.Copy = false
	return nil
}

type Environments struct {
	Target    *EnvType
	Source    *EnvType
	Available []*etree.Element
}

func (e *Environments) InitSource(index string, username string, password string) (error, bool) {
	e.Source = new(EnvType)
	return e.Source.init(index, username, password)
}

func (e *Environments) InitTarget(index string, username string, password string) (error, bool) {
	e.Target = new(EnvType)
	return e.Target.init(index, username, password)
}

func (e *Environments) CopyTargetFromSource() {
	e.Target = e.Source.copyEnv()
}

func (e *Environments) Logout() error {
	err := e.Source.logout()
	if err != nil {
		return err
	}
	err = e.Target.logout()
	if err != nil {
		return err
	}
	return nil
}

func login(env *EnvType) (string, error) {
	loginEnvelopeElement := docXogReadXML.FindElement("//xogtype[@type='login']/soapenv:Envelope").Copy()
	request := etree.NewDocument()
	request.SetRoot(loginEnvelopeElement)

	request.FindElement("//obj:Username").SetText(env.Username)
	request.FindElement("//obj:Password").SetText(env.Password)

	body, err := request.WriteToString()

	if err != nil {
		return "", errors.New("Problems getting login xml: " + err.Error())
	}

	response, err := util.SoapCall(body, env.URL)
	resp := etree.NewDocument()
	resp.ReadFromString(response)

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

func logout(env *EnvType) error {
	logoutEnvelopeElement := docXogReadXML.FindElement("//xogtype[@type='logout']/soapenv:Envelope").Copy()
	request := etree.NewDocument()
	request.SetRoot(logoutEnvelopeElement)
	request.FindElement("//obj:SessionID").SetText(env.Session)

	body, err := request.WriteToString()

	if err != nil {
		return errors.New("Problems getting logout xml: " + err.Error())
	}

	_, err = util.SoapCall(body, env.URL)

	if err != nil {
		return errors.New("Problems trying to logout from environment: " + env.Name + " | Debug: " + err.Error())
	}

	return nil
}
