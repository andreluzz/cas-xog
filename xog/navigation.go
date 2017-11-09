package xog

import (
	"os"
	"fmt"
	"strings"
	"github.com/andreluzz/cas-xog/common"
)

var startInstallingPackage int

func RenderHome() {

	common.InitLog()

	common.Info("\n")
	common.Info("--------------------------------------------\n")
	common.Info("##### CAS XOG Automation - Version %.1f #####\n", common.VERSION)
	common.Info("--------------------------------------------\n")

	startInstallingPackage = 0

	InitRead()

	LoadEnvironmentsList()

	LoadPackages()

	driverXOG = &common.Driver{}
	RenderDrivers()
}

func RenderInterface() bool {
	var inputAction string

	if startInstallingPackage == 1 {
		startInstallingPackage = -1
		inputAction = "p"
	} else {
		common.Info("\nChoose action")
		common.Info("\n(l = Load XOG Driver, r = Read XOGs, w = Write XOGs, m = Create Migration, p = Install Package or x = eXit): ")
		fmt.Scanln(&inputAction)
	}

	switch strings.ToLower(inputAction) {
	case "w", "r", "m":
		if driverXOG == nil || len(driverXOG.Files) <= 0 {
			common.Info("\n[CAS-XOG][red[ERROR]] - XOG driver not loaded. Try action 'l' to load a valid driver.\n")
			return false
		}
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
			common.Info("\n[CAS-XOG][red[PACKAGE]]: %s\n", err.Error())
			return false
		}
	case "l":
		RenderDrivers()
	case "x":
		common.Info("\n[CAS-XOG][blue[Action exit selected]] - Press enter key to exit...\n")
		scanExit := ""
		fmt.Scanln(&scanExit)
		os.Exit(0)
		return true
	default:
		common.Info("\n[CAS-XOG][red[ACTION]] - Action not implemented!\n")
	}

	return false
}