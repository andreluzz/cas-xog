package xog

import (
	"github.com/andreluzz/cas-xog/constant"
	"github.com/andreluzz/cas-xog/model"
	"os"
	"path/filepath"
	"testing"
	"io/ioutil"
	"github.com/andreluzz/cas-xog/util"
)

func TestLoadPackages(t *testing.T) {
	folder := "../mock/xog/" + constant.FOLDER_PACKAGE
	LoadPackages(folder, "../mock/xog/mock_packages/")

	packages := GetAvailablePackages()
	if len(packages) == 0 {
		t.Fatalf("Error loading available packages, no packages loaded")
	}

	if packages[0].Name != "Mock Package" {
		t.Errorf("Error loading .package file, expected name 'Mock Package' received '%s'", packages[0].Name)
	}

	driversList := GetPackagesDriversFileInfoList()

	if len(driversList) != 2 {
		t.Errorf("Error loading available packages no driver loaded, expected 2 received %d", len(driversList))
	}
}

func TestLoadPackagesInvalidUserPackageFolder(t *testing.T) {
	folder := "../mock/xog/" + constant.FOLDER_PACKAGE
	LoadPackages(folder, "")

	packages := GetAvailablePackages()
	if len(packages) != 0 {
		t.Fatalf("Error loading available packages, invalid user package folder not cover expected 0 received %d", len(packages))
	}
}

func TestLoadAvailablePackages(t *testing.T) {
	folder := "../mock/xog/" + constant.FOLDER_PACKAGE
	unzipPackages(folder, "../mock/xog/mock_packages/")
	loadAvailablePackages(folder)

	packages := GetAvailablePackages()
	if len(packages) == 0 {
		t.Fatalf("Error loading available packages, no packages loaded")
	}

	if packages[0].Name != "Mock Package" {
		t.Errorf("Error loading .package file, expected name 'Mock Package' received '%s'", packages[0].Name)
	}

	driversList := GetPackagesDriversFileInfoList()

	if len(driversList) != 2 {
		t.Errorf("Error loading available packages no driver loaded, expected 2 received %d", len(driversList))
	}
}

func TestUnzipPackages(t *testing.T) {
	folder := "../mock/xog/" + constant.FOLDER_PACKAGE
	unzipPackages(folder, "../mock/xog/mock_packages/")

	total := 0
	filepath.Walk(folder+"mock-pkg/", func(path string, file os.FileInfo, err error) error {
		total += 1
		return err
	})

	if total != 20 {
		t.Errorf("Error unziping package, expected 20 files received %d", total)
	}
}

func TestProcessPackageFile(t *testing.T) {
	folder := "../mock/xog/" + constant.FOLDER_PACKAGE
	LoadPackages(folder, "../mock/xog/mock_packages/")

	selectedPackage := GetAvailablePackages()[0]
	driverPath := folder + selectedPackage.Folder + selectedPackage.DriverFileName

	LoadDriver(driverPath)

	file := GetLoadedDriver().Files[0]

	packageFolder := folder + selectedPackage.Folder + selectedPackage.Versions[0].Folder + file.Type + "/"
	writeFolder := constant.FOLDER_WRITE + file.Type

	output := ProcessPackageFile(file, &selectedPackage.Versions[0], packageFolder, writeFolder)
	if output.Code != constant.OUTPUT_SUCCESS {
		t.Errorf("Error processing package file. Debug: %s", output.Debug)
	}

	output = ProcessPackageFile(model.DriverFile{}, &selectedPackage.Versions[0], packageFolder, writeFolder)
	if output.Code != constant.OUTPUT_ERROR {
		t.Errorf("Error processing package file. Not validating invalid file")
	}
}

func TestInstallPackageFile(t *testing.T) {
	model.LoadXMLReadList("../xogRead.xml")

	folder := "../mock/xog/" + constant.FOLDER_PACKAGE
	LoadPackages(folder, "../mock/xog/mock_packages/")

	selectedPackage := GetAvailablePackages()[0]
	driverPath := folder + selectedPackage.Folder + selectedPackage.DriverFileName

	LoadDriver(driverPath)

	file := GetLoadedDriver().Files[4]

	packageFolder := folder + selectedPackage.Folder + selectedPackage.Versions[0].Folder + file.Type + "/"
	writeFolder := constant.FOLDER_WRITE + file.Type

	output := ProcessPackageFile(file, &selectedPackage.Versions[0], packageFolder, writeFolder)

	mockEnvironments := &model.Environments{
		Source: &model.EnvType{
			Name: "Mock Source Env",
			URL: "Mock URL",
			Session: "Mock session",
		},
		Target: &model.EnvType{
			Name: "Mock Source Env",
			URL: "Mock URL",
			Session: "Mock session",
		},
	}

	soapMock := func(request, endpoint string) (string, error) {
		file, _ := ioutil.ReadFile("../mock/xog/soap/soap_success_write_response.xml")
		return util.BytesToString(file), nil
	}

	output = InstallPackageFile(&file, mockEnvironments, soapMock)
	if output.Code != constant.OUTPUT_SUCCESS {
		t.Errorf("Error installing package file. Debug: %s", output.Debug)
	}

	soapMock = func(request, endpoint string) (string, error) {
		return "", nil
	}

	output = InstallPackageFile(&file, mockEnvironments, soapMock)
	if output.Code != constant.OUTPUT_ERROR {
		t.Errorf("Error installing package file. Not validating soap response")
	}
}