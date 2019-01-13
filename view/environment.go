package view

import (
	"fmt"
	"os"
	"strconv"

	"github.com/andreluzz/cas-xog/constant"
	"github.com/andreluzz/cas-xog/log"
	"github.com/andreluzz/cas-xog/model"
	"github.com/andreluzz/cas-xog/util"
	"github.com/howeyc/gopass"
)

var targetEnvInput string

//Environments displays the options for the user to choose the environment for reading and writing
func Environments(action string, environments *model.Environments) bool {
	if action == "m" {
		return true
	}

	result := false
	sourceInput := "-1"
	targetInput := "-1"

	log.Info("\n")
	log.Info("Available environments:\n")

	if environments.Available == nil || len(environments.Available) == 0 {
		log.Info("\n[CAS-XOG][red[ERROR]] - None available environments found!\n")
		log.Info("\n[CAS-XOG][red[FATAL]] - Check your xogEnv.xml file. Press enter key to exit...")
		scanExit := ""
		fmt.Scanln(&scanExit)
		os.Exit(0)
	}

	for i, e := range environments.Available {
		log.Info("%d - %s\n", i+1, e.Name)
	}

	if action == "r" {
		sourceInput, result = processingChooseEnvironment(environments, constant.Source, "reading", constant.Undefined)
		if result == false {
			return false
		}
	}

	targetInput, result = processingChooseEnvironment(environments, constant.Target, "writing", sourceInput)
	if result == false {
		return false
	}

	if action == "r" {
		targetEnvInput = targetInput
	}

	if action == "w" && targetEnvInput != targetInput {
		log.Info("\n[CAS-XOG][yellow[Warning]]: Trying to write files read from a different target environment!")
		log.Info("\n[CAS-XOG]Do you want to continue anyway? (y = Yes, n = No) [n]: ")
		input := "n"
		fmt.Scanln(&input)
		if input != "y" {
			return false
		}
	}

	return true
}

func processingChooseEnvironment(environments *model.Environments, envType, label, lastEnvInput string) (string, bool) {
	log.Info("Choose %s environment [1]: ", label)
	userScanInput := "1"
	fmt.Scanln(&userScanInput)

	envIndex, err := strconv.Atoi(userScanInput)
	envIndex--

	if err != nil || envIndex < 0 || envIndex >= len(environments.Available) {
		log.Info("\n[CAS-XOG][red[ERROR]] - Invalid reading environment index!\n")
		return userScanInput, false
	}
	var env *model.EnvType

	if envType == constant.Source {
		env = environments.Source
	} else {
		env = environments.Target
	}

	if userScanInput == lastEnvInput {
		environments.CopyTargetFromSource()
		return userScanInput, true
	}

	env.Init(envIndex)
	if env.RequestLogin {
		requestLogin(env)
	}
	log.Info("[CAS-XOG]Processing environment login")
	err = env.Login(envIndex, util.SoapCall, util.RestCall)

	if err != nil {
		log.Info("\n[CAS-XOG][red[ERROR]] - %s", err.Error())
		log.Info("\n[CAS-XOG][red[FATAL]] - Check your xogEnv.xml file. Press enter key to exit...")
		scanExit := ""
		fmt.Scanln(&scanExit)
		os.Exit(0)
	}
	log.Info("\r[CAS-XOG][green[Login successfully]] - Environment: %s \n", env.Name)

	return userScanInput, true
}

func requestLogin(envType *model.EnvType) {
	log.Info("\n[CAS-XOG][yellow[Login needed]] - Enter credentials for environment: %s \n", envType.Name)
	log.Info("Username: ")
	fmt.Scanln(&envType.Username)

	log.Info("Password: ")
	passwordTemp, _ := gopass.GetPasswdMasked()
	envType.Password = string(passwordTemp[:])
	log.Info("\n")
}
