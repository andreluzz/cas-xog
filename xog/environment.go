package xog

import (
	"errors"
	"github.com/beevik/etree"
)

var envXml *etree.Document
var SourceEnv, TargetEnv *EnvType

func LoadEnvironmentsList() {
	envXml = etree.NewDocument()
	envXml.ReadFromFile("xogEnv.xml")
}

type EnvType struct {
	Name 		string
	URL 		string
	Username 	string
	Password 	string
	Session 	string
	Copy		bool
}

func (e *EnvType) init(envIndex string) error {
	envElement := envXml.FindElement("//env["+envIndex+"]").Copy()
	if envElement == nil {
		return errors.New("trying to initiate an invalid environment")
	}

	e.Name = envElement.SelectAttrValue("name", "Environment name not defined")
	e.Username = envElement.FindElement("//username").Text()
	e.Password = envElement.FindElement("//password").Text()
	e.URL = envElement.FindElement("//endpoint").Text()

	session, err := Login(e)
	if err != nil {
		return err
	}
	e.Session = session
	e.Copy = false

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
	err := Logout(e)
	e.clear()
	return err
}

func (e *EnvType) clear() error {
	e.Username = ""
	e.Password = ""
	e.URL = ""
	e.Session = ""
	e.Copy = false
	return nil
}

func (e *EnvType) copyEnv() *EnvType {
	ne := &EnvType{
		Name:  		e.Name,
		Username:   e.Username,
		Password: 	e.Password,
		URL: 		e.URL,
		Session: 	e.Session,
		Copy:		true,
	}
	return ne
}