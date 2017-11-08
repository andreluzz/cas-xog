package xog

import (
	"os"
	"fmt"
	"errors"
	"strconv"
	"github.com/beevik/etree"
	"github.com/andreluzz/cas-xog/common"
)

var envXml *etree.Document
var SourceEnv *EnvType
var TargetEnv *EnvType
var targetEnvInput string

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

func RenderEnvironments(action string) bool {
	if action == "m" {
		return true
	}

	sourceInput := "-1"
	targetInput := "-1"

	common.Info("\n")
	common.Info("Available environments:\n")

	availableEnvironments := envXml.FindElements("./xogenvs/env")

	if availableEnvironments == nil || len(availableEnvironments) == 0 {
		common.Info("\n[CAS-XOG][red[ERROR]] - None available environments found!\n")
		common.Info("\n[CAS-XOG][red[FATAL]] - Check your xogEnv.xml file. Press enter key to exit...")
		scanExit := ""
		fmt.Scanln(&scanExit)
		os.Exit(0)
	}

	for i, e := range availableEnvironments {
		common.Info("%d - %s\n", i+1, e.SelectAttrValue("name", "Unknown environment name"))
	}

	if action == "r" {
		common.Info("Choose reading environment [1]: ")
		sourceInput = "1"
		fmt.Scanln(&sourceInput)

		var err error
		envIndex, err := strconv.Atoi(sourceInput)

		if err != nil || envIndex <= 0 || envIndex > len(availableEnvironments) {
			common.Info("\n[CAS-XOG][red[ERROR]] - Invalid reading environment index!\n")
			return false
		}

		common.Info("[CAS-XOG]Processing environment login")
		SourceEnv = new(EnvType)
		err = SourceEnv.init(sourceInput)
		if err != nil {
			common.Info("\n[CAS-XOG][red[ERROR]] - %s", err.Error())
			common.Info("\n[CAS-XOG][red[FATAL]] - Check your xogEnv.xml file. Press enter key to exit...")
			scanExit := ""
			fmt.Scanln(&scanExit)
			os.Exit(0)
		}

		common.Info("\r[CAS-XOG][green[Login successfully]] - Environment: %s \n", SourceEnv.Name)
	}

	common.Info("Choose writing environment [1]: ")
	targetInput = "1"
	fmt.Scanln(&targetInput)

	envIndex, err := strconv.Atoi(targetInput)

	if action == "r" {
		targetEnvInput = targetInput
	} else if action == "w" && targetEnvInput != targetInput {
		common.Info("\n[CAS-XOG][yellow[Warning]]: Trying to write files read from a different target environment!")
		common.Info("\n[CAS-XOG]Do you want to continue anyway? (y = Yes, n = No) [n]: ")
		input := "n"
		fmt.Scanln(&input)
		if input == "n" || input != "y" {
			return false
		}
	}

	if err != nil || envIndex <= 0 || envIndex > len(availableEnvironments) {
		common.Info("\n[CAS-XOG][red[ERROR]] - Invalid writing environment index!\n")
		return false
	}

	TargetEnv = new(EnvType)
	if sourceInput == targetInput {
		TargetEnv = SourceEnv.copyEnv()
	} else {
		common.Info("[CAS-XOG]Processing environment login")
		err = TargetEnv.init(targetInput)
		if err != nil {
			common.Info("\n[CAS-XOG][red[ERROR]] - %s", err.Error())
			common.Info("\n[CAS-XOG][red[FATAL]] - Check your xogEnv.xml file. Press enter key to exit...")
			scanExit := ""
			fmt.Scanln(&scanExit)
			os.Exit(0)
		}
		common.Info("\r[CAS-XOG]Environment: %s - [green[Login successfully]]\n", TargetEnv.Name)
	}

	return true
}