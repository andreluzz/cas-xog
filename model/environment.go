package model

import (
	"encoding/json"
	"encoding/xml"
	"errors"
	"io/ioutil"

	"github.com/andreluzz/cas-xog/util"
	"github.com/beevik/etree"
)

var environments *Environments

//LoadEnvironmentsList loads the list of user-defined environments to use when executing xog
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

//EnvType defines an environment attributes
type EnvType struct {
	Name         string `xml:"name,attr"`
	URL          string `xml:"endpoint"`
	Username     string `xml:"username"`
	Password     string `xml:"password"`
	Proxy        string `xml:"proxy"`
	Cookie       string `xml:"cookie"`
	Session      string
	AuthToken    string
	Copy         bool
	RequestLogin bool
}

//Init loads a specific environment from user environments list
func (e *EnvType) Init(envIndex int) {
	available := environments.Available[envIndex].copyEnv()

	e.Name = available.Name
	e.Username = available.Username
	e.Password = available.Password
	e.URL = available.URL
	e.Proxy = available.Proxy
	e.Cookie = available.Cookie
	e.RequestLogin = false

	if e.Username == "" || e.Password == "" {
		e.RequestLogin = true
	}
}

//Login executes an soap call to retrieve the session id from the environment
func (e *EnvType) Login(envIndex int, soapFunc util.Soap, restFunc util.Rest) error {
	var err error

	environments.Available[envIndex].Username = e.Username
	environments.Available[envIndex].Password = e.Password

	e.Session, err = login(e, soapFunc)
	if err != nil {
		return err
	}

	e.AuthToken, err = loginAPI(e)
	if err != nil {
		return err
	}
	return nil
}

func (e *EnvType) logout(soapFunc util.Soap) error {
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
	err := logout(e, soapFunc)
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
		Proxy:    e.Proxy,
		Cookie:   e.Cookie,
		Copy:     true,
	}
	return ne
}

func (e *EnvType) clear() error {
	e.Username = ""
	e.Password = ""
	e.URL = ""
	e.Session = ""
	e.Proxy = ""
	e.Cookie = ""
	e.Copy = false
	return nil
}

//Environments defines a list of available environments
type Environments struct {
	Available []*EnvType `xml:"env"`
	Target    *EnvType
	Source    *EnvType
}

//CopyTargetFromSource copy the data from source to target environment
func (e *Environments) CopyTargetFromSource() {
	e.Target = e.Source.copyEnv()
}

//Logout closes the session id in the environment
func (e *Environments) Logout(soapFunc util.Soap) error {
	err := e.Source.logout(soapFunc)
	if err != nil {
		return err
	}
	err = e.Target.logout(soapFunc)
	if err != nil {
		return err
	}
	return nil
}

type apiLogin struct {
	Token string `json:"authToken"`
}

func loginAPI(env *EnvType) (string, error) {
	response, err := util.APIPostLogin(env.URL+"/ppm/rest/v1/auth/login", env.Username, env.Password, env.Proxy, env.Cookie)
	if err != nil {
		return "", errors.New("Problems trying to get API Token from environment: " + env.Name + " | Debug: " + err.Error())
	}

	api := &apiLogin{}
	json.Unmarshal(response, api)

	return api.Token, nil
}

func login(env *EnvType, soapFunc util.Soap) (string, error) {
	loginEnvelopeElement := docXogReadXML.FindElement("//xogtype[@type='login']/soapenv:Envelope").Copy()
	request := etree.NewDocument()
	request.SetRoot(loginEnvelopeElement)

	request.FindElement("//obj:Username").SetText(env.Username)
	request.FindElement("//obj:Password").SetText(env.Password)

	body, err := request.WriteToString()

	if err != nil {
		return "", errors.New("Problems getting login xml: " + err.Error())
	}

	response, err := soapFunc(body, env.URL, env.Proxy)
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

func logout(env *EnvType, soapFunc util.Soap) error {
	logoutEnvelopeElement := docXogReadXML.FindElement("//xogtype[@type='logout']/soapenv:Envelope").Copy()
	request := etree.NewDocument()
	request.SetRoot(logoutEnvelopeElement)
	request.FindElement("//obj:SessionID").SetText(env.Session)

	body, err := request.WriteToString()

	if err != nil {
		return errors.New("Problems getting logout xml: " + err.Error())
	}

	_, err = soapFunc(body, env.URL, env.Proxy)

	if err != nil {
		return errors.New("Problems trying to logout from environment: " + env.Name + " | Debug: " + err.Error())
	}

	return nil
}
