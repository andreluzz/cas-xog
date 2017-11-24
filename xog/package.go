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

func GetAvailablePackages() []model.Package {
	return availablePackages
}

func GetPackagesDriversFileInfoList() []model.Driver {
	return packagesDriversFileInfo
}

func ProcessPackageFile(file model.DriverFile, selectedVersion *model.Version, packageFolder, writeFolder string) model.Output {
	output := model.Output{Code: constant.OUTPUT_SUCCESS, Debug: ""}
	err := transform.ProcessPackageFile(file, packageFolder, writeFolder, selectedVersion.Definitions)
	if err != nil {
		output.Code = constant.OUTPUT_ERROR
		output.Debug = err.Error()
	}
	return output
}

func InstallPackageFile(file *model.DriverFile, environments *model.Environments, soapFunc util.Soap) model.Output {
	output := model.Output{Code: constant.OUTPUT_SUCCESS, Debug: ""}

	util.ValidateFolder(constant.FOLDER_DEBUG + file.Type)

	file.InitXML(constant.WRITE, constant.FOLDER_WRITE)

	iniTagRegexpStr, endTagRegexpStr := file.TagCDATA()
	if iniTagRegexpStr != "" && endTagRegexpStr != "" {
		responseString := transform.IncludeCDATA(file.GetXML(), iniTagRegexpStr, endTagRegexpStr)
		file.SetXML(responseString)
	}

	err := file.RunXML(constant.WRITE, constant.FOLDER_WRITE, environments, soapFunc)
	xogResponse := etree.NewDocument()
	xogResponse.ReadFromString(file.GetXML())
	output, err = validate.Check(xogResponse)
	if err != nil {
		output.Code = constant.OUTPUT_ERROR
		output.Debug = err.Error()
		return output
	}
	file.Write(constant.FOLDER_DEBUG)
	return output
}
