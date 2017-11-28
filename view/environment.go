package view

import (
	"fmt"
	"github.com/andreluzz/cas-xog/log"
	"github.com/andreluzz/cas-xog/model"
	"github.com/howeyc/gopass"
	"os"
	"strconv"
	"github.com/andreluzz/cas-xog/util"
)

var targetEnvInput string

func Environments(action string, environments *model.Environments) bool {
	if action == "m" {
		return true
	}

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
		log.Info("Choose reading environment [1]: ")
		sourceInput = "1"
		fmt.Scanln(&sourceInput)

		envIndex, err := strconv.Atoi(sourceInput)
		envIndex--

		if err != nil || envIndex < 0 || envIndex > len(environments.Available) {
			log.Info("\n[CAS-XOG][red[ERROR]] - Invalid reading environment index!\n")
			return false
		}

		environments.Source.Init(envIndex)
		if environments.Source.RequestLogin {
			requestLogin(environments.Source)
		}
		log.Info("[CAS-XOG]Processing environment login")
		err = environments.Source.Login(envIndex, util.SoapCall)

		if err != nil {
			log.Info("\n[CAS-XOG][red[ERROR]] - %s", err.Error())
			log.Info("\n[CAS-XOG][red[FATAL]] - Check your xogEnv.xml file. Press enter key to exit...")
			scanExit := ""
			fmt.Scanln(&scanExit)
			os.Exit(0)
		}
		log.Info("\r[CAS-XOG][green[Login successfully]] - Environment: %s \n", environments.Source.Name)
	}

	log.Info("Choose writing environment [1]: ")
	targetInput = "1"
	fmt.Scanln(&targetInput)

	envIndex, err := strconv.Atoi(targetInput)
	envIndex--

	if err != nil || envIndex < 0 || envIndex > len(environments.Available) {
		log.Info("\n[CAS-XOG][red[ERROR]] - Invalid writing environment index!\n")
		return false
	}

	if action == "r" {
		targetEnvInput = targetInput
	} else if action == "w" && targetEnvInput != targetInput {
		log.Info("\n[CAS-XOG][yellow[Warning]]: Trying to write files read from a different target environment!")
		log.Info("\n[CAS-XOG]Do you want to continue anyway? (y = Yes, n = No) [n]: ")
		input := "n"
		fmt.Scanln(&input)
		if input != "y" {
			return false
		}
	}

	if sourceInput == targetInput {
		environments.CopyTargetFromSource()
	} else {
		environments.Target.Init(envIndex)
		if environments.Target.RequestLogin {
			requestLogin(environments.Target)
		}
		log.Info("[CAS-XOG]Processing environment login")
		err := environments.Target.Login(envIndex, util.SoapCall)
		if err != nil {
			log.Info("\n[CAS-XOG][red[ERROR]] - %s", err.Error())
			log.Info("\n[CAS-XOG][red[FATAL]] - Check your xogEnv.xml file. Press enter key to exit...")
			scanExit := ""
			fmt.Scanln(&scanExit)
			os.Exit(0)
		}
		log.Info("\r[CAS-XOG][green[Login successfully]] - Environment: %s \n", environments.Target.Name)
	}

	return true
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
