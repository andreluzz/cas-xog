package model

import (
	"encoding/xml"
	"errors"
	"github.com/andreluzz/cas-xog/util"
	"github.com/beevik/etree"
	"io/ioutil"
)

var environments *Environments

func LoadEnvironmentsList(path string) (*Environments, error) {
	environments = &Environments{}

	xmlFile, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, errors.New("Error loading driver file - " + err.Error())
	}

	xml.Unmarshal(xmlFile, environments)
	environments.Source = &EnvType{}
	environments.Target = &EnvType{}
	return environments, err
}

type EnvType struct {
	Name         string `xml:"name,attr"`
	URL          string `xml:"endpoint"`
	Username     string `xml:"username"`
	Password     string `xml:"password"`
	Session      string
	Copy         bool
	RequestLogin bool
}

func (e *EnvType) Init(envIndex int) {
	available := environments.Available[envIndex].copyEnv()

	e.Name = available.Name
	e.Username = available.Username
	e.Password = available.Password
	e.URL = available.URL
	e.RequestLogin = false

	if e.Username == "" || e.Password == "" {
		e.RequestLogin = true
	}
}

func (e *EnvType) Login(envIndex int) error {
	var err error

	environments.Available[envIndex].Username = e.Username
	environments.Available[envIndex].Password = e.Password

	e.Session, err = login(e)
	if err != nil {
		return err
	}
	return nil
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
	Available []*EnvType `xml:"env"`
	Target    *EnvType
	Source    *EnvType
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
