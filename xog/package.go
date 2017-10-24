package xog

import (
	"os"
	"fmt"
	"time"
	"errors"
	"strconv"
	"strings"
	"io/ioutil"
	"encoding/xml"
	"path/filepath"
	"github.com/andreluzz/cas-xog/common"
	"github.com/andreluzz/cas-xog/transform"
)

var packagesDriversFileInfo []common.Driver
var availablePackages []common.Package
var selectedPackage common.Package
var selectedVersion common.Version

func LoadPackages() {
	err := unzipPackages()
	if err != nil {
		common.Info("\n[CAS-XOG][red[ERROR]] - Error unzipping packages: %s", err.Error())
		common.Info("\n[CAS-XOG][red[FATAL]] - Check your packages folder. Press enter key to exit...")
		scanExit := ""
		fmt.Scanln(&scanExit)
		os.Exit(0)
	}

	_, dirErr := os.Stat(common.FOLDER_PACKAGE)
	if os.IsNotExist(dirErr) {
		return
	}

	err = loadAvailablePackages()
	if err != nil {
		common.Info("\n[CAS-XOG][red[ERROR]] - Error loading packages: %s", err.Error())
		common.Info("\n[CAS-XOG][red[FATAL]] - Check your _packages folder. Press enter key to exit...")
		scanExit := ""
		fmt.Scanln(&scanExit)
		os.Exit(0)
	}
}

func unzipPackages() error {
	userPackagesFolder := "packages/"
	packagesFiles, _ := ioutil.ReadDir(userPackagesFolder)
	if len(packagesFiles) == 0 {
		return nil
	}

	os.RemoveAll(common.FOLDER_PACKAGE)
	os.MkdirAll(common.FOLDER_PACKAGE, os.ModePerm)

	for _, f := range packagesFiles {
		_, err := common.Unzip(userPackagesFolder + f.Name(), common.FOLDER_PACKAGE)
		if err != nil {
			return err
		}
	}

	return nil
}

func loadAvailablePackages() error {
	availablePackages = nil
	err := filepath.Walk("./_packages", func(path string, info os.FileInfo, err error) error {
		if strings.Contains(path, ".package") {
			xmlPackageFile, err := ioutil.ReadFile(path)
			if err != nil {
				return errors.New("Error loading package file: " + err.Error())
			}
			driverXOG = new(common.Driver)
			pkg := new(common.Package)
			xml.Unmarshal(xmlPackageFile, pkg)
			availablePackages = append(availablePackages, *pkg)
		} else if strings.Contains(path, ".driver") {
			driver := new(common.Driver)
			driver.Info = info
			driver.PackageDriver = true
			driver.FilePath = path
			packagesDriversFileInfo = append(packagesDriversFileInfo, *driver)
		}
		return err
	})

	return err
}

func InstallPackage() error {
	start := time.Now()

	output = map[string]int{common.OUTPUT_SUCCESS: 0, common.OUTPUT_WARNING: 0, common.OUTPUT_ERROR: 0}

	common.Info("\n------------------------------------------------------------------")
	common.Info("\n[blue[Initiated at]]: %s", start.Format("Mon _2 Jan 2006 - 15:04:05"))
	common.Info("\nInstalling Package: [blue[%s]] (%s)", selectedPackage.Name, selectedVersion.Name)
	common.Info("\nTarget environment: [blue[%s]]", TargetEnv.Name)
	if len(selectedVersion.Definitions) > 0 {
		common.Info("\nDefinitions:")
		for _, d := range selectedVersion.Definitions {
			common.Info("\n   %s: %s", d.Action, d.Value)
		}
	}
	common.Info("\n------------------------------------------------------------------\n")
	common.Info("\n[CAS-XOG]Start package install? (y = Yes, n = No) [n]: ")
	input := "n"
	fmt.Scanln(&input)
	if input == "n" || input != "y" {
		return nil
	}

	driverPath := common.FOLDER_PACKAGE + selectedPackage.Folder + selectedPackage.DriverFileName
	if selectedVersion.DriverFileName != "" {
		driverPath = common.FOLDER_PACKAGE + selectedPackage.Folder + selectedVersion.Folder + selectedVersion.DriverFileName
	}

	err := LoadDriver(driverPath)
	if err != nil {
		return err
	}

	os.RemoveAll(common.FOLDER_DEBUG)
	os.MkdirAll(common.FOLDER_DEBUG, os.ModePerm)
	os.RemoveAll(common.FOLDER_WRITE)
	os.MkdirAll(common.FOLDER_WRITE, os.ModePerm)

	for i, f := range driverXOG.Files {
		common.Info("\n[CAS-XOG][blue[Processing]] %03d/%03d | file: %s", i+1, len(driverXOG.Files), f.Path)
		f.PackageFolder = common.FOLDER_PACKAGE + selectedPackage.Folder + selectedVersion.Folder

		transform.ProcessPackage(f, selectedVersion.Definitions)

		common.Info("\r[CAS-XOG]Writing %03d/%03d | file: %s   ", i+1, len(driverXOG.Files), f.Path)

		action := "w"
		folder := common.FOLDER_DEBUG

		//check if target folder type dir exists
		_, dirErr := os.Stat(folder + f.Type)
		if os.IsNotExist(dirErr) {
			_ = os.Mkdir(folder+f.Type, os.ModePerm)
		}

		_, validateOutput, err := loadAndValidate(action, folder, &f, TargetEnv)

		if err != nil {
			debug(i+1, len(driverXOG.Files), action, validateOutput.Code, f.Path, err.Error())
			continue
		}

		debug(i+1, len(driverXOG.Files), action, validateOutput.Code, f.Path, validateOutput.Debug)
	}

	TargetEnv.logout()
	elapsed := time.Since(start)

	common.Info("\n\n------------------------------------------------------------------")
	common.Info("\nStats: total = %d | failure = %d | success = %d | warning = %d", len(driverXOG.Files), output[common.OUTPUT_ERROR], output[common.OUTPUT_SUCCESS], output[common.OUTPUT_WARNING])
	common.Info("\n[blue[Concluded in]]: %.3f seconds", elapsed.Seconds())
	common.Info("\n------------------------------------------------------------------\n")

	return nil
}

func RenderPackages() bool {

	if len(availablePackages) <= 0 {
		common.Info("\n[CAS-XOG][yellow[WARNING]] - No package available, check your packages folder!\n")
		return false
	}

	common.Info("\n")
	common.Info("Available packages:\n")
	for i, p := range availablePackages {
		common.Info("%d - %s\n", i+1, p.Name)
	}
	common.Info("Choose package to install [1]: ")
	input := "1"
	fmt.Scanln(&input)

	packageIndex, err := strconv.Atoi(input)

	if err != nil || packageIndex-1 < 0 || packageIndex > len(availablePackages) {
		common.Info("\n[CAS-XOG][red[ERROR]] - Invalid package!\n")
		return false
	}

	selectedPackage = availablePackages[packageIndex-1]
	common.Info("\n[CAS-XOG] [blue[Package %s selected]]\n", selectedPackage.Name)

	versionIndex := 1
	if len(selectedPackage.Versions) > 1 {
		common.Info("Available package versions:\n")
		for i, v := range selectedPackage.Versions {
			common.Info("%d - %s\n", i+1, v.Name)
		}
		common.Info("Choose version to install [1]: ")
		input := "1"
		fmt.Scanln(&input)
		versionIndex, err = strconv.Atoi(input)

		if err != nil {
			common.Info("\n[CAS-XOG][red[ERROR]] - Package definition error: %s\n", err.Error())
			return false
		}
	}

	common.Info("\n[CAS-XOG] [blue[Package required definitions:]]\n")
	selectedVersion = selectedPackage.Versions[versionIndex-1]
	for i, d := range selectedVersion.Definitions {
		common.Info("%s [%s]:", d.Description, d.Default)
		input := d.Default
		fmt.Scanln(&input)
		if input == "" {
			common.Info("\n[CAS-XOG][red[ERROR]] - Invalid definition!\n")
			return false
		}
		selectedVersion.Definitions[i].Value = input
	}

	return true
}