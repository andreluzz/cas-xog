package xog

import (
	"encoding/xml"
	"errors"
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

var packagesDriversFileInfo []model.Driver
var availablePackages []model.Package

func LoadPackages() error {
	err := unzipPackages()
	if err != nil {
		return err
	}

	_, dirErr := os.Stat(constant.FOLDER_PACKAGE)
	if os.IsNotExist(dirErr) {
		return nil
	}

	err = loadAvailablePackages()
	if err != nil {
		return err
	}
	return nil
}

func unzipPackages() error {
	userPackagesFolder := "packages/"
	packagesFiles, _ := ioutil.ReadDir(userPackagesFolder)
	if len(packagesFiles) == 0 {
		return nil
	}

	os.RemoveAll(constant.FOLDER_PACKAGE)
	os.MkdirAll(constant.FOLDER_PACKAGE, os.ModePerm)

	for _, f := range packagesFiles {
		_, err := util.Unzip(userPackagesFolder+f.Name(), constant.FOLDER_PACKAGE)
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
			driverXOG = new(model.Driver)
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

	return err
}

func GetPackagesDriversFileInfoList() []model.Driver {
	return packagesDriversFileInfo
}

func GetAvailablePackages() []model.Package {
	return availablePackages
}

func ProcessPackageFile(file model.DriverFile, selectedPackage *model.Package, selectedVersion *model.Version) model.Output {
	output := model.Output{Code: constant.OUTPUT_SUCCESS, Debug: ""}
	packageFolder := constant.FOLDER_PACKAGE + selectedPackage.Folder + selectedVersion.Folder + "/" + file.Type + "/"
	writeFolder := constant.FOLDER_WRITE + file.Type
	err := transform.ProcessPackageFile(file, packageFolder, writeFolder, selectedVersion.Definitions)
	if err != nil {
		output.Code = constant.OUTPUT_ERROR
		output.Debug = err.Error()
	}
	return output
}

func InstallPackageFile(file *model.DriverFile, environments *model.Environments) model.Output {
	output := model.Output{Code: constant.OUTPUT_SUCCESS, Debug: ""}

	util.ValidateFolder(constant.FOLDER_DEBUG + file.Type)

	err := file.InitXML(constant.WRITE, constant.FOLDER_WRITE)
	if err != nil {
		output.Code = constant.OUTPUT_ERROR
		output.Debug = err.Error()
		return output
	}

	iniTagRegexpStr, endTagRegexpStr := file.TagCDATA()
	if iniTagRegexpStr != "" && endTagRegexpStr != "" {
		responseString := transform.IncludeCDATA(file.GetXML(), iniTagRegexpStr, endTagRegexpStr)
		file.SetXML(responseString)
	}

	err = file.RunXML(constant.WRITE, constant.FOLDER_WRITE, environments)
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
