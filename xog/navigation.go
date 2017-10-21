package xog

import (
	"os"
	"fmt"
	"strings"
	"github.com/andreluzz/cas-xog/common"
)

func RenderHome() {
	common.Debug("\n")
	common.Debug("--------------------------------------------\n")
	common.Debug("##### CAS XOG Automation - Version 2.0 #####\n")
	common.Debug("--------------------------------------------\n")

	InitRead()

	LoadEnvironmentsList()

	LoadPackages()

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
			common.Debug("\n[CAS-XOG][red[PACKAGE]]: %s\n", err.Error())
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