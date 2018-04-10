package view

import (
	"fmt"
	"github.com/andreluzz/cas-xog/constant"
	"github.com/andreluzz/cas-xog/log"
	"github.com/andreluzz/cas-xog/model"
	"github.com/andreluzz/cas-xog/util"
	"github.com/andreluzz/cas-xog/xog"
	"os"
	"strings"
)

var startInstallingPackage int
var environments *model.Environments

//Home display the system header and initializes variables
func Home(version string) {
	var err error

	log.InitLog()

	log.Info("\n")
	log.Info("------------------------------------------------\n")
	log.Info("##### CAS XOG Automation - Version %s #####\n", version)
	log.Info("------------------------------------------------\n")

	startInstallingPackage = 0

	model.LoadXMLReadList("xogRead.xml")

	environments, err = model.LoadEnvironmentsList("xogEnv.xml")
	if err != nil {
		log.Info("\n[CAS-XOG][red[Error]]: %s\n", err.Error())
	}

	renderDrivers()
}

//Interface display the main menu for user to interact and choose available actions
func Interface() bool {
	var inputAction string

	if startInstallingPackage == 1 {
		startInstallingPackage = -1
		inputAction = "p"
	} else {
		log.Info("\nChoose action")
		log.Info("\n(l = Load XOG Driver, r = Read XOGs, w = Write XOGs, m = Create Migration, p = Install Package or x = eXit): ")
		fmt.Scanln(&inputAction)
	}

	action := strings.ToLower(inputAction)
	switch action {
	case constant.Write, constant.Read, constant.Migrate:
		if xog.ValidateLoadedDriver() == false {
			log.Info("\n[CAS-XOG][red[ERROR]] - XOG driver not loaded. Try action 'l' to load a valid driver.\n")
			return false
		}
		if !Environments(action, environments) {
			return false
		}
		driver := xog.GetLoadedDriver()

		if action == constant.Read && driver.AutomaticWrite {
			log.Info("\n[CAS-XOG][yellow[Warning]]: This driver is configured to write automatically!")
			log.Info("\n[CAS-XOG]Do you want to proceed? (y = Yes, n = No) [n]: ")
			input := "n"
			fmt.Scanln(&input)
			if input != "y" {
				ProcessDriverFiles(driver, action, environments)
				environments.Logout(util.SoapCall)
				return false
			}
		}

		ProcessDriverFiles(driver, action, environments)

		if action == constant.Read && driver.AutomaticWrite {
			ProcessDriverFiles(driver, constant.Write, environments)
		}

		environments.Logout(util.SoapCall)
	case constant.Package:
		xog.LoadPackages(constant.FolderPackage, "packages/")
		output, selectedPackage, selectedVersion := renderPackages()
		if !output {
			return false
		}
		if !Environments(action, environments) {
			return false
		}
		err := InstallPackage(environments, selectedPackage, selectedVersion)
		if err != nil {
			log.Info("\n[CAS-XOG][red[PACKAGE]]: %s\n", err.Error())
			return false
		}
	case constant.Load:
		renderDrivers()
	case constant.Exit:
		log.Info("\n[CAS-XOG][blue[Action exit selected]] - Press enter key to exit...\n")
		scanExit := ""
		fmt.Scanln(&scanExit)
		os.Exit(0)
		return true
	default:
		log.Info("\n[CAS-XOG][red[ACTION]] - Action not implemented!\n")
	}

	return false
}
