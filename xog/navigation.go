package xog

import (
	"os"
	"fmt"
	"strconv"
	"strings"
	"io/ioutil"
	"github.com/andreluzz/cas-xog/common"
)

func RenderHome() {
	common.Debug("\n")
	common.Debug("--------------------------------------------\n")
	common.Debug("##### CAS XOG Automation - Version 2.0 #####\n")
	common.Debug("--------------------------------------------\n")

	InitRead()

	LoadEnvironmentsList()

	err := LoadAvailablePackages()
	if err != nil {
		common.Debug("\n[CAS-XOG][red[ERROR]] - Error loading packages: %s", err.Error())
		common.Debug("\n[CAS-XOG][red[FATAL]] - Check your _packages folder. Press enter key to exit...")
		scanExit := ""
		fmt.Scanln(&scanExit)
		os.Exit(0)
	}

	RenderDrivers()
}

func RenderInterface() bool {
	var inputAction string
	common.Debug("\nChoose action")
	common.Debug("\n(l = Load XOG Driver, r = Read XOGs, w = Write XOGs, m = Create Migration, p = Install Packages or x = eXit): ")
	fmt.Scanln(&inputAction)
	switch strings.ToLower(inputAction) {
	case "w", "r", "m":
		if !RenderEnvironments(strings.ToLower(inputAction)) {
			return false
		}
		ProcessDriverFiles(strings.ToLower(inputAction))
	case "p":
		if !RenderPackages() {
			return false
		}
		if !RenderEnvironments(strings.ToLower(inputAction)) {
			return false
		}
		err := InstallPackage()
		if err != nil {
			common.Debug("\n[CAS-XOG][red[ERROR]] %s\n", err.Error())
			return false
		}
	case "l":
		RenderDrivers()
	case "x":
		common.Debug("\n[CAS-XOG][blue[Action exit selected]]\nPress enter key to exit...")
		scanExit := ""
		fmt.Scanln(&scanExit)
		os.Exit(0)
		return true
	default:
		common.Debug("\n[CAS-XOG][red[ACTION]] - Action not implemented!\n")
	}

	return false
}

func RenderDrivers() {
	var driverIndex = 1
	driverPath := "drivers/"
	driverFileList, _ := ioutil.ReadDir(driverPath)

	if len(driverFileList) == 0 {
		common.Debug("\n[CAS-XOG][red[ERROR]] - XogDriver folders or file not found! Press enter key to exit...\n")
		scanexit := ""
		fmt.Scanln(&scanexit)
		os.Exit(0)
	}

	fmt.Println("")
	fmt.Println("Available drivers:")
	for k, f := range driverFileList {
		common.Debug("%d - %s\n", k+1, f.Name())
	}
	fmt.Print("Choose driver [1]: ")
	var input string = "1"
	fmt.Scanln(&input)

	var err error
	driverIndex, err = strconv.Atoi(input)

	if err != nil || driverIndex-1 < 0 || driverIndex > len(driverFileList) {
		common.Debug("\n[CAS-XOG][red[ERROR]] - Invalid XOG driver!\n")
		return
	}

	completePath := driverPath + driverFileList[driverIndex-1].Name()
	err = LoadDriver(completePath)
	if err != nil {
		common.Debug("\n[CAS-XOG][red[ERROR]] - %s", err.Error())
		common.Debug("\n[CAS-XOG][red[FATAL]] - Check your driver file. Press enter key to exit...")
		scanExit := ""
		fmt.Scanln(&scanExit)
		os.Exit(0)
	}

	common.Debug("\n[CAS-XOG][blue[Loaded XOG Driver file]]: %s | Total files: [green[%d]]\n", completePath, len(driverXOG.Files))
}

var targetEnvInput string

func RenderEnvironments(action string) bool {
	if action == "m" {
		return true
	}

	sourceInput := "-1"
	targetInput := "-1"

	common.Debug("\n")
	common.Debug("Available environments:\n")

	availableEnvironments := envXml.FindElements("./xogenvs/env")

	if availableEnvironments == nil || len(availableEnvironments) == 0 {
		common.Debug("\n[CAS-XOG][red[ERROR]] - None available environments found!\n")
		common.Debug("\n[CAS-XOG][red[FATAL]] - Check your xogEnv.xml file. Press enter key to exit...")
		scanExit := ""
		fmt.Scanln(&scanExit)
		os.Exit(0)
	}

	for i, e := range availableEnvironments {
		common.Debug("%d - %s\n", i+1, e.SelectAttrValue("name", "Unknown environment name"))
	}

	if action == "r" {
		common.Debug("Choose reading environment [1]: ")
		sourceInput = "1"
		fmt.Scanln(&sourceInput)

		var err error
		envIndex, err := strconv.Atoi(sourceInput)

		if err != nil || envIndex <= 0 || envIndex > len(availableEnvironments) {
			common.Debug("\n[CAS-XOG][red[ERROR]] - Invalid reading environment index!\n\n")
			return false
		}

		common.Debug("[CAS-XOG]Processing environment login")
		SourceEnv = new(EnvType)
		err = SourceEnv.init(sourceInput)
		if err != nil {
			common.Debug("\n[CAS-XOG][red[ERROR]] - %s", err.Error())
			common.Debug("\n[CAS-XOG][red[FATAL]] - Check your xogEnv.xml file. Press enter key to exit...")
			scanExit := ""
			fmt.Scanln(&scanExit)
			os.Exit(0)
		}

		common.Debug("\r[CAS-XOG]Environment: %s - [green[Login successfully]]\n", SourceEnv.Name)
	}

	common.Debug("Choose writing environment [1]: ")
	targetInput = "1"
	fmt.Scanln(&targetInput)

	envIndex, err := strconv.Atoi(targetInput)

	if action == "r" {
		targetEnvInput = targetInput
	} else if action == "w" && targetEnvInput != targetInput {
		common.Debug("\n[CAS-XOG][yellow[Warning]]: Trying to write files read from a different target environment!")
		common.Debug("\n[CAS-XOG]Do you want to continue anyway? (y = Yes, n = No) [n]: ")
		input := "n"
		fmt.Scanln(&input)
		if input == "n" || input != "y" {
			return false
		}
	}

	if err != nil || envIndex <= 0 || envIndex > len(availableEnvironments) {
		common.Debug("\n[CAS-XOG][red[ERROR]] - Invalid writing environment index!\n\n")
		return false
	}

	TargetEnv = new(EnvType)
	if sourceInput == targetInput {
		TargetEnv = SourceEnv.copyEnv()
	} else {
		common.Debug("[CAS-XOG]Processing environment login")
		err = TargetEnv.init(targetInput)
		if err != nil {
			common.Debug("\n[CAS-XOG][red[ERROR]] - %s", err.Error())
			common.Debug("\n[CAS-XOG][red[FATAL]] - Check your xogEnv.xml file. Press enter key to exit...")
			scanExit := ""
			fmt.Scanln(&scanExit)
			os.Exit(0)
		}
		common.Debug("\r[CAS-XOG]Environment: %s - [green[Login successfully]]\n", TargetEnv.Name)
	}

	return true
}
