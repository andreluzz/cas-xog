package xog

import (
	"os"
	"errors"
	"strings"
	"io/ioutil"
	"encoding/xml"
	"path/filepath"
	"github.com/andreluzz/cas-xog/common"
	"fmt"
	"strconv"
	"time"
)

type Definition struct {
	Type string `xml:"type,attr"`
	Description string `xml:"description,attr"`
	Value string `xml:"default,attr"`
}

type Version struct {
	Name string `xml:"name,attr"`
	Path string `xml:"path,attr"`
	Definitions []Definition `xml:"definition"`
}

type Package struct {
	Name string `xml:"name,attr"`
	DriverPath string `xml:"driverPath,attr"`
	Versions []Version `xml:"version"`
}

var availablePackages []Package
var selectedPackage Package
var selectedVersion Version

func LoadAvailablePackages() error {
	err := filepath.Walk("./_packages", func(path string, info os.FileInfo, err error) error {
		if strings.Contains(path, ".package") {
			xmlPackageFile, err := ioutil.ReadFile(path)
			if err != nil {
				return errors.New("Error loading package file: " + err.Error())
			}
			driverXOG = new(common.Driver)
			pkg := new(Package)
			xml.Unmarshal(xmlPackageFile, pkg)
			availablePackages = append(availablePackages, *pkg)
		}
		return err
	})

	if err != nil {
		return err
	}

	return nil
}

func InstallPackage() error {
	start := time.Now()
	common.Debug("\n------------------------------------------------------------------")
	common.Debug("\n[blue[Initiated at]]: %s", start.Format("Mon _2 Jan 2006 - 15:04:05"))
	common.Debug("\nInstalling Package: [blue[%s]] (%s)", selectedPackage.Name, selectedVersion.Name)
	common.Debug("\nDefinitions:")
	for _, d := range selectedVersion.Definitions {
		common.Debug("\n   %s: %s", d.Type, d.Value)
	}
	common.Debug("\n------------------------------------------------------------------\n")

	err := LoadDriver(selectedPackage.DriverPath)
	if err != nil {
		return err
	}

	for i, f := range driverXOG.Files {
		common.Debug("\n[CAS-XOG][blue[Installing]] %03d/%03d | file: %s", i+1, len(driverXOG.Files), f.Path)
		
		common.Debug("\r[CAS-XOG][green[Installed]] %03d/%03d | file: %s", i+1, len(driverXOG.Files), f.Path)
	}

	TargetEnv.logout()
	elapsed := time.Since(start)

	common.Debug("\n\n------------------------------------------------------------------")
	common.Debug("\n[blue[Package installed in]]: %.3f seconds", elapsed.Seconds())
	common.Debug("\n------------------------------------------------------------------\n")

	return nil
}

func RenderPackages() bool {
	common.Debug("\n")
	common.Debug("Available packages:\n")
	for i, p := range availablePackages {
		common.Debug("%d - %s\n", i+1, p.Name)
	}
	common.Debug("Choose package to install [1]: ")
	var input string = "1"
	fmt.Scanln(&input)

	packageIndex, err := strconv.Atoi(input)

	if err != nil || packageIndex-1 < 0 || packageIndex > len(availablePackages) {
		common.Debug("\n[CAS-XOG][red[ERROR]] - Invalid package!\n")
		return false
	}

	selectedPackage = availablePackages[packageIndex-1]
	common.Debug("\n[CAS-XOG] [blue[Package %s selected]]\n", selectedPackage.Name)

	versionIndex := 1
	if len(selectedPackage.Versions) > 1 {
		common.Debug("Available package versions:\n")
		for i, v := range selectedPackage.Versions {
			common.Debug("%d - %s\n", i+1, v.Name)
		}
		common.Debug("Choose version to install [1]: ")
		var input string = "1"
		fmt.Scanln(&input)
		versionIndex, err = strconv.Atoi(input)

		if err != nil {
			common.Debug("\n[CAS-XOG][red[ERROR]] - Package definition error: %s\n", err.Error())
			return false
		}
	}

	common.Debug("\n[CAS-XOG] [blue[Package required definitions:]]\n")
	selectedVersion = selectedPackage.Versions[versionIndex-1]
	for i, d := range selectedVersion.Definitions {
		common.Debug("%s [%s]:", d.Description, d.Value)
		input = d.Value
		fmt.Scanln(&input)
		if input == "" {
			common.Debug("\n[CAS-XOG][red[ERROR]] - Invalid definition!\n")
			return false
		}
		selectedVersion.Definitions[i].Value = input
	}

	return true
}