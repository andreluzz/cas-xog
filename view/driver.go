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

//ProcessDriverFiles displays the feedback of drivers processing
func ProcessDriverFiles(driver *model.Driver, action string, environments *model.Environments) {
	start := time.Now()

	outputResults := map[string]int{constant.OutputSuccess: 0, constant.OutputWarning: 0, constant.OutputError: 0, constant.OutputIgnored: 0}

	log.Info("\n------------------------------------------------------------------")
	log.Info("\n[blue[Initiated at]]: %s", start.Format("Mon _2 Jan 2006 - 15:04:05"))
	log.Info("\nProcessing driver: %s", driver.FilePath)
	log.Info("\n------------------------------------------------------------------\n")

	processingString := "processing  "
	if action == "r" {
		os.RemoveAll(constant.FolderRead)
		os.MkdirAll(constant.FolderRead, os.ModePerm)
		os.RemoveAll(constant.FolderWrite)
		os.MkdirAll(constant.FolderWrite, os.ModePerm)
	} else if action == "w" {
		os.RemoveAll(constant.FolderDebug)
		os.MkdirAll(constant.FolderDebug, os.ModePerm)
		processingString = "processing   "
	} else if action == "m" {
		os.RemoveAll(constant.FolderMigration)
		os.MkdirAll(constant.FolderMigration, os.ModePerm)
		processingString = "processing    "
	}

	total := len(driver.Files)
	typePadLength := driver.MaxTypeNameLen()

	for i, f := range driver.Files {
		formattedType := util.RightPad(f.GetXMLType(), " ", typePadLength)
		if f.IgnoreReading && action == "r" {
			log.Info("\n[CAS-XOG][yellow[Read ignored]] %03d/%03d | [blue[%s]] | file: %s", i+1, total, formattedType, f.Path)
			outputResults[constant.OutputIgnored]++
			continue
		}
		sourceFolder, outputFolder := xog.CreateFileFolder(action, f.Type, f.Path)

		if f.Type == constant.TypeMigration {
			sourceFolder = constant.FolderMigration
		}

		splitFilename, _ := f.GetSplitWriteFilesPath(sourceFolder)
		if len(splitFilename) > 0 {
			totalSplit := len(splitFilename)
			for j, filename := range splitFilename {
				f.Path = filename
				log.Info("\n[CAS-XOG][blue[%s]] %03d/%03d | [blue[%s]] | Split: %03d/%03d | file: %s", processingString, i+1, total, formattedType, j+1, totalSplit, f.Path)
				output := xog.ProcessDriverFile(&f, action, sourceFolder, outputFolder, environments, util.SoapCall)
				status, color := util.GetStatusColorFromOutput(output.Code)
				log.Info("\r[CAS-XOG][%s[%s %s]] %03d/%03d | [blue[%s]] | Split: %03d/%03d | file: %s %s", color, util.GetActionLabel(action), status, i+1, total, formattedType, j+1, totalSplit, f.Path, util.GetOutputDebug(output.Code, output.Debug))
				outputResults[output.Code]++
			}
		} else {
			log.Info("\n[CAS-XOG][blue[%s]] %03d/%03d | [blue[%s]] | file: %s", processingString, i+1, total, formattedType, f.Path)

			output := xog.ProcessDriverFile(&f, action, sourceFolder, outputFolder, environments, util.SoapCall)
			status, color := util.GetStatusColorFromOutput(output.Code)
			log.Info("\r[CAS-XOG][%s[%s %s]] %03d/%03d | [blue[%s]] | file: %s %s", color, util.GetActionLabel(action), status, i+1, total, formattedType, f.Path, util.GetOutputDebug(output.Code, output.Debug))
			outputResults[output.Code]++
		}
	}

	elapsed := time.Since(start)

	environments.Logout(util.SoapCall)

	totalFilesProcessed := outputResults[constant.OutputError] + outputResults[constant.OutputSuccess] + outputResults[constant.OutputWarning] + outputResults[constant.OutputIgnored]

	log.Info("\n\n-----------------------------------------------------------------------------")
	log.Info("\nStats: total = %d | failure = %d | success = %d | warning = %d | ignored = %d", totalFilesProcessed, outputResults[constant.OutputError], outputResults[constant.OutputSuccess], outputResults[constant.OutputWarning], outputResults[constant.OutputIgnored])
	log.Info("\n[blue[Concluded in]]: %.3f seconds", elapsed.Seconds())
	log.Info("\n-----------------------------------------------------------------------------\n")
}

func renderDrivers() {
	folder := "drivers/"
	currentFolder := constant.Undefined
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
			if d.Folder != constant.Undefined && currentFolder != d.Folder {
				log.Info("[blue[Folder %s:]]\n", d.Folder)
			}
			if d.Folder != constant.Undefined {
				log.Info("    %d - %s\n", k+1, d.Info.Name())
			} else {
				log.Info("%d - %s\n", k+1, d.Info.Name())
			}
		}
		currentFolder = d.Folder
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
