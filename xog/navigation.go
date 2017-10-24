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
	common.Info("##### CAS XOG Automation - Version 2.0 #####\n")
	common.Info("--------------------------------------------\n")

	startInstallingPackage = 0

	InitRead()

	LoadEnvironmentsList()

	LoadPackages()

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