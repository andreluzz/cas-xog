package xog

import (
	"encoding/xml"
	"github.com/andreluzz/cas-xog/constant"
	"github.com/andreluzz/cas-xog/model"
	"github.com/andreluzz/cas-xog/transform"
	"github.com/andreluzz/cas-xog/util"
	"github.com/andreluzz/cas-xog/validate"
	"github.com/beevik/etree"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

var availablePackages []model.Package
var packagesDriversFileInfo []model.Driver

//LoadPackages load the available packages in the user's folder
func LoadPackages(systemPackageFolder, userPackageFolder string) {
	availablePackages = nil
	packagesDriversFileInfo = nil
	unzipPackages(systemPackageFolder, userPackageFolder)

	_, dirErr := os.Stat(systemPackageFolder)
	if os.IsNotExist(dirErr) {
		return
	}

	loadAvailablePackages(systemPackageFolder)
}

func unzipPackages(systemPackageFolder, userPackageFolder string) {
	os.RemoveAll(systemPackageFolder)

	packagesFiles, _ := ioutil.ReadDir(userPackageFolder)
	if len(packagesFiles) == 0 {
		return
	}

	os.MkdirAll(systemPackageFolder, os.ModePerm)

	for _, f := range packagesFiles {
		util.Unzip(userPackageFolder+f.Name(), systemPackageFolder)
	}
}

func loadAvailablePackages(folder string) {
	availablePackages = nil
	packagesDriversFileInfo = nil
	filepath.Walk(folder, func(path string, info os.FileInfo, err error) error {
		if strings.Contains(path, ".package") {
			xmlPackageFile, _ := ioutil.ReadFile(path)
			pkg := new(model.Package)
			xml.Unmarshal(xmlPackageFile, pkg)
			availablePackages = append(availablePackages, *pkg)
		} else if strings.Contains(path, ".driver") {
			driver := new(model.Driver)
			driver.Info = info
			driver.PackageDriver = true
			driver.FilePath = path
			packagesDriversFileInfo = append(packagesDriversFileInfo, *driver)
		}
		return err
	})
}

//GetAvailablePackages returns a list of available packages
func GetAvailablePackages() []model.Package {
	return availablePackages
}

//ProcessPackageFile validates if the driver needs transformation and creates the write xog files according to the installation environment
func ProcessPackageFile(file *model.DriverFile, selectedVersion *model.Version, packageFolder, writeFolder string, environments *model.Environments, soapFunc util.Soap) model.Output {
	if file.PackageTransform && file.NeedPackageTransform() {
		file.InitXML(constant.Read, constant.Undefined)
		file.RunAuxXML(environments.Target, soapFunc)
	}

	return transform.ProcessPackageFile(file, packageFolder, writeFolder, selectedVersion.Definitions)
}

//InstallPackageFile execute the soap call to install the driver and returns the output
func InstallPackageFile(file *model.DriverFile, environments *model.Environments, soapFunc util.Soap) model.Output {
	output := model.Output{Code: constant.OutputSuccess, Debug: constant.Undefined}

	util.ValidateFolder(constant.FolderDebug + file.Type + util.GetPathFolder(file.Path))

	file.InitXML(constant.Write, constant.FolderWrite)

	iniTagRegexpStr, endTagRegexpStr := file.TagCDATA()
	if iniTagRegexpStr != constant.Undefined && endTagRegexpStr != constant.Undefined {
		responseString := transform.IncludeCDATA(file.GetXML(), iniTagRegexpStr, endTagRegexpStr)
		file.SetXML(responseString)
	}

	err := file.RunXML(constant.Write, constant.FolderWrite, environments, soapFunc)
	xogResponse := etree.NewDocument()
	xogResponse.ReadFromString(file.GetXML())
	output, err = validate.Check(xogResponse)
	if err != nil {
		return output
	}
	file.Write(constant.FolderDebug)
	return output
}
