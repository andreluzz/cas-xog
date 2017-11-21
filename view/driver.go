package view

import (
	"fmt"
	"github.com/andreluzz/cas-xog/constant"
	"github.com/andreluzz/cas-xog/log"
	"github.com/andreluzz/cas-xog/model"
	"github.com/andreluzz/cas-xog/util"
	"github.com/andreluzz/cas-xog/xog"
	"os"
	"strconv"
	"time"
)

func ProcessDriverFiles(driver *model.Driver, action string, environments *model.Environments) {
	start := time.Now()

	outputResults := map[string]int{constant.OUTPUT_SUCCESS: 0, constant.OUTPUT_WARNING: 0, constant.OUTPUT_ERROR: 0}

	log.Info("\n------------------------------------------------------------------")
	log.Info("\n[blue[Initiated at]]: %s", start.Format("Mon _2 Jan 2006 - 15:04:05"))
	log.Info("\nProcessing driver: %s", driver.FilePath)
	log.Info("\n------------------------------------------------------------------\n")

	if action == "r" {
		os.RemoveAll(constant.FOLDER_READ)
		os.MkdirAll(constant.FOLDER_READ, os.ModePerm)
		os.RemoveAll(constant.FOLDER_WRITE)
		os.MkdirAll(constant.FOLDER_WRITE, os.ModePerm)
	} else if action == "w" {
		os.RemoveAll(constant.FOLDER_DEBUG)
		os.MkdirAll(constant.FOLDER_DEBUG, os.ModePerm)
	} else if action == "m" {
		os.RemoveAll(constant.FOLDER_MIGRATION)
		os.MkdirAll(constant.FOLDER_MIGRATION, os.ModePerm)
	}

	total := len(driver.Files)

	for i, f := range driver.Files {
		log.Info("\n[CAS-XOG][blue[Processing]] %03d/%03d | file: %s", i+1, total, f.Path)
		output := xog.ProcessDriverFile(&f, action, environments)
		status, color := util.GetStatusColorFromOutput(output.Code)
		log.Info("\r[CAS-XOG][%s[%s %s]] %03d/%03d | file: %s %s", color, util.GetActionLabel(action), status, i+1, total, f.Path, util.GetOutputDebug(output.Debug))
		outputResults[output.Code] += 1
	}

	elapsed := time.Since(start)

	environments.Logout()

	log.Info("\n\n------------------------------------------------------------------")
	log.Info("\nStats: total = %d | failure = %d | success = %d | warning = %d", len(driver.Files), outputResults[constant.OUTPUT_ERROR], outputResults[constant.OUTPUT_SUCCESS], outputResults[constant.OUTPUT_WARNING])
	log.Info("\n[blue[Concluded in]]: %.3f seconds", elapsed.Seconds())
	log.Info("\n------------------------------------------------------------------\n")
}

func Drivers() {
	folder := "drivers/"
	driversList, err := xog.GetDriversList(folder)
	if err != nil {
		log.Info("\n[CAS-XOG][red[ERROR]] - %s\n", err.Error())
		return
	}

	fmt.Println("")
	fmt.Println("Available drivers:")
	for k, d := range driversList {
		if d.PackageDriver {
			log.Info("%d - [blue[Package driver:]] %s\n", k+1, d.Info.Name())
		} else {
			log.Info("%d - %s\n", k+1, d.Info.Name())
		}
	}
	if startInstallingPackage == 0 {
		fmt.Print("Choose driver [1] or p = Install Package: ")
	} else {
		fmt.Print("Choose driver [1]: ")
	}

	input := "1"
	fmt.Scanln(&input)

	if input == "p" && startInstallingPackage == 0 {
		startInstallingPackage = 1
		return
	}
	startInstallingPackage = -1

	driverIndex, err := strconv.Atoi(input)

	if err != nil || driverIndex-1 < 0 || driverIndex > len(driversList) {
		log.Info("\n[CAS-XOG][red[ERROR]] - Invalid XOG driver!\n")
		return
	}

	total, err := xog.LoadDriver(driversList[driverIndex-1].FilePath)
	if err != nil {
		log.Info("\n[CAS-XOG][red[ERROR]] - %s\n", err.Error())
		return
	}

	log.Info("\n[CAS-XOG][blue[Loaded XOG Driver file]]: %s | Total files: [green[%d]]\n", driversList[driverIndex-1].FilePath, total)
}
