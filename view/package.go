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

func InstallPackage(environments *model.Environments, selectedPackage *model.Package, selectedVersion *model.Version) error {
	start := time.Now()

	outputResults := map[string]int{constant.OUTPUT_SUCCESS: 0, constant.OUTPUT_WARNING: 0, constant.OUTPUT_ERROR: 0, constant.OUTPUT_IGNORED: 0}

	driverPath := constant.FOLDER_PACKAGE + selectedPackage.Folder + selectedPackage.DriverFileName
	if selectedVersion.DriverFileName != "" {
		driverPath = constant.FOLDER_PACKAGE + selectedPackage.Folder + selectedVersion.Folder + selectedVersion.DriverFileName
	}

	total, err := xog.LoadDriver(driverPath)
	if err != nil {
		return err
	}

	os.RemoveAll(constant.FOLDER_DEBUG)
	os.MkdirAll(constant.FOLDER_DEBUG, os.ModePerm)
	os.RemoveAll(constant.FOLDER_WRITE)
	os.MkdirAll(constant.FOLDER_WRITE, os.ModePerm)

	log.Info("\n------------------------------------------------------------------")
	log.Info("\n[blue[Initiated at]]: %s", start.Format("Mon _2 Jan 2006 - 15:04:05"))
	log.Info("\nProcessing Package: [blue[%s]] (%s)", selectedPackage.Name, selectedVersion.Name)
	log.Info("\n------------------------------------------------------------------\n")

	driver := xog.GetLoadedDriver()
	typePadLength := driver.MaxTypeNameLen()

	for i, f := range driver.Files {
		formattedType := util.RightPad(f.GetXMLType(), " ", typePadLength)
		if f.IgnoreReading {
			log.Info("\n[CAS-XOG][yellow[Processed ignored]] %03d/%03d | [blue[%s]] | file: %s", i+1, total, formattedType, f.Path)
			outputResults[constant.OUTPUT_IGNORED] += 1
			continue
		}
		log.Info("\n[CAS-XOG][blue[Processing]] %03d/%03d | [blue[%s]] | file: %s", i+1, total, formattedType, f.Path)
		packageFolder := constant.FOLDER_PACKAGE + selectedPackage.Folder + selectedVersion.Folder + f.Type + "/"
		writeFolder := constant.FOLDER_WRITE + f.Type
		output := xog.ProcessPackageFile(f, selectedVersion, packageFolder, writeFolder)
		status, color := util.GetStatusColorFromOutput(output.Code)
		log.Info("\r[CAS-XOG][%s[Processed %s]] %03d/%03d | [blue[%s]] | file: %s %s", color, status, i+1, total, formattedType, f.Path, util.GetOutputDebug(output.Code, output.Debug))
		outputResults[output.Code] += 1
	}

	elapsed := time.Since(start)

	log.Info("\n\n-----------------------------------------------------------------------------")
	log.Info("\nStats: total = %d | failure = %d | success = %d | warning = %d | ignored = %d", total, outputResults[constant.OUTPUT_ERROR], outputResults[constant.OUTPUT_SUCCESS], outputResults[constant.OUTPUT_WARNING], outputResults[constant.OUTPUT_IGNORED])
	log.Info("\n[blue[Concluded in]]: %.3f seconds", elapsed.Seconds())
	log.Info("\n-----------------------------------------------------------------------------\n")

	outputResults = map[string]int{constant.OUTPUT_SUCCESS: 0, constant.OUTPUT_WARNING: 0, constant.OUTPUT_ERROR: 0, constant.OUTPUT_IGNORED: 0}
	start = time.Now()

	log.Info("\n------------------------------------------------------------------")
	log.Info("\n[blue[Initiated at]]: %s", start.Format("Mon _2 Jan 2006 - 15:04:05"))
	log.Info("\nInstalling Package: [blue[%s]] (%s)", selectedPackage.Name, selectedVersion.Name)
	log.Info("\nTarget environment: [blue[%s]]", environments.Target.Name)
	if len(selectedVersion.Definitions) > 0 {
		log.Info("\nDefinitions:")
		for _, d := range selectedVersion.Definitions {
			log.Info("\n   %s: %s", d.Description, d.Value)
		}
	}
	log.Info("\n------------------------------------------------------------------\n")
	log.Info("\n[CAS-XOG]Start package install? (y = Yes, n = No) [n]: ")
	input := "n"
	fmt.Scanln(&input)
	if input != "y" {
		return nil
	}

	start = time.Now()

	for i, f := range driver.Files {
		formattedType := util.RightPad(f.GetXMLType(), " ", typePadLength)
		log.Info("\n[CAS-XOG][blue[Installing     ]] %03d/%03d | [blue[%s]] | file: %s", i+1, total, formattedType, f.Path)
		output := xog.InstallPackageFile(&f, environments, util.SoapCall)
		status, color := util.GetStatusColorFromOutput(output.Code)
		log.Info("\r[CAS-XOG][%s[Install %s]] %03d/%03d | [blue[%s]] | file: %s %s", color, status, i+1, total, formattedType, f.Path, util.GetOutputDebug(output.Code, output.Debug))
		outputResults[output.Code] += 1
	}

	environments.Logout()
	elapsed = time.Since(start)

	log.Info("\n\n------------------------------------------------------------------")
	log.Info("\nStats: total = %d | failure = %d | success = %d | warning = %d", total, outputResults[constant.OUTPUT_ERROR], outputResults[constant.OUTPUT_SUCCESS], outputResults[constant.OUTPUT_WARNING])
	log.Info("\n[blue[Concluded in]]: %.3f seconds", elapsed.Seconds())
	log.Info("\n------------------------------------------------------------------\n")

	return nil
}

func Packages() (bool, *model.Package, *model.Version) {

	availablePackages := xog.GetAvailablePackages()

	if len(availablePackages) <= 0 {
		log.Info("\n[CAS-XOG][yellow[WARNING]] - No package available, check your packages folder!\n")
		return false, nil, nil
	}

	log.Info("\n")
	log.Info("Available packages:\n")
	for i, p := range availablePackages {
		log.Info("%d - %s\n", i+1, p.Name)
	}
	log.Info("Choose package to install [1] or b = Back to options menu: ")
	input := "1"
	fmt.Scanln(&input)

	if input == "b" {
		return false, nil, nil
	}

	packageIndex, err := strconv.Atoi(input)

	if err != nil || packageIndex-1 < 0 || packageIndex > len(availablePackages) {
		log.Info("\n[CAS-XOG][red[ERROR]] - Invalid package!\n")
		return false, nil, nil
	}

	selectedPackage := availablePackages[packageIndex-1]
	log.Info("\n[CAS-XOG] [blue[Package %s selected]]\n", selectedPackage.Name)

	versionIndex := 1
	if len(selectedPackage.Versions) > 1 {
		log.Info("Available package versions:\n")
		for i, v := range selectedPackage.Versions {
			log.Info("%d - %s\n", i+1, v.Name)
		}
		log.Info("Choose version to install [1]: ")
		input := "1"
		fmt.Scanln(&input)
		versionIndex, err = strconv.Atoi(input)

		if err != nil {
			log.Info("\n[CAS-XOG][red[ERROR]] - Package definition error: %s\n", err.Error())
			return false, nil, nil
		}
	}

	selectedVersion := selectedPackage.Versions[versionIndex-1]
	if len(selectedVersion.Definitions) > 0 {
		log.Info("\n[CAS-XOG] [blue[Package required definitions:]]\n")
		for i, d := range selectedVersion.Definitions {
			log.Info("%s [%s]:", d.Description, d.Default)
			input := d.Default
			fmt.Scanln(&input)
			if input == "" {
				log.Info("\n[CAS-XOG][red[ERROR]] - Invalid definition!\n")
				return false, nil, nil
			}
			selectedVersion.Definitions[i].Value = input
		}
	}

	return true, &selectedPackage, &selectedVersion
}
