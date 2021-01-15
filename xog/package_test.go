package xog

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/andreluzz/cas-xog/constant"
	"github.com/andreluzz/cas-xog/model"
	"github.com/andreluzz/cas-xog/util"
)

func TestLoadPackages(t *testing.T) {
	folder := "../mock/xog/" + constant.FolderPackage
	LoadPackages(folder, "../mock/xog/mock_packages/")

	packages := GetAvailablePackages()
	if len(packages) == 0 {
		t.Fatalf("Error loading available packages, no packages loaded")
	}

	if packages[0].Name != "Mock Package" {
		t.Errorf("Error loading .package file, expected name 'Mock Package' received '%s'", packages[0].Name)
	}
}

func TestLoadPackagesInvalidUserPackageFolder(t *testing.T) {
	folder := "../mock/xog/" + constant.FolderPackage
	LoadPackages(folder, "")

	packages := GetAvailablePackages()
	if len(packages) != 0 {
		t.Fatalf("Error loading available packages, invalid user package folder not cover expected 0 received %d", len(packages))
	}
}

func TestLoadAvailablePackages(t *testing.T) {
	folder := "../mock/xog/" + constant.FolderPackage
	unzipPackages(folder, "../mock/xog/mock_packages/")
	loadAvailablePackages(folder)

	packages := GetAvailablePackages()
	if len(packages) == 0 {
		t.Fatalf("Error loading available packages, no packages loaded")
	}

	if packages[0].Name != "Mock Package" {
		t.Errorf("Error loading .package file, expected name 'Mock Package' received '%s'", packages[0].Name)
	}
}

func TestUnzipPackages(t *testing.T) {
	folder := "../mock/xog/" + constant.FolderPackage
	unzipPackages(folder, "../mock/xog/mock_packages/")

	total := 0
	filepath.Walk(folder+"mock-pkg/", func(path string, file os.FileInfo, err error) error {
		total++
		return err
	})

	if total != 24 {
		t.Errorf("Error unziping package, expected 24 files received %d", total)
	}
}

func TestProcessPackageFile(t *testing.T) {
	folder := "../mock/xog/" + constant.FolderPackage
	LoadPackages(folder, "../mock/xog/mock_packages/")

	selectedPackage := GetAvailablePackages()[0]
	driverPath := folder + selectedPackage.Folder + selectedPackage.DriverFileName

	LoadDriver(driverPath)

	file := GetLoadedDriver().Files[0]

	packageFolder := folder + selectedPackage.Folder + selectedPackage.Versions[0].Folder + file.Type + "/"
	writeFolder := constant.FolderWrite + file.Type

	output := ProcessPackageFile(&file, &selectedPackage.Versions[0], packageFolder, writeFolder, nil, nil)
	if output.Code != constant.OutputSuccess {
		t.Errorf("Error processing package file. Debug: %s", output.Debug)
	}

	output = ProcessPackageFile(&model.DriverFile{}, &selectedPackage.Versions[0], packageFolder, writeFolder, nil, nil)
	if output.Code != constant.OutputError {
		t.Errorf("Error processing package file. Not validating invalid file")
	}
}

func TestProcessAndTransformPackageFile(t *testing.T) {
	folder := "../mock/xog/" + constant.FolderPackage
	LoadPackages(folder, "../mock/xog/mock_packages/")

	selectedPackage := GetAvailablePackages()[0]
	driverPath := folder + selectedPackage.Folder + selectedPackage.DriverFileName

	LoadDriver(driverPath)

	file := GetLoadedDriver().Files[5]

	packageFolder := folder + selectedPackage.Folder + selectedPackage.Versions[0].Folder + file.Type + "/"
	writeFolder := constant.FolderWrite + file.Type

	soapMock := func(request, endpoint, proxy string, opts ...interface{}) (string, error) {
		file, _ := ioutil.ReadFile("../mock/xog/package_transform_view_target.xml")
		return util.BytesToString(file), nil
	}
	mockEnvironments := &model.Environments{
		Source: &model.EnvType{
			Name:    "Mock Source Env",
			URL:     "Mock URL",
			Session: "Mock session",
		},
		Target: &model.EnvType{
			Name:    "Mock Source Env",
			URL:     "Mock URL",
			Session: "Mock session",
		},
	}
	output := ProcessPackageFile(&file, &selectedPackage.Versions[0], packageFolder, writeFolder, mockEnvironments, soapMock)
	if output.Code != constant.OutputSuccess {
		t.Errorf("Error processing package file. Debug: %s", output.Debug)
	}

	output = ProcessPackageFile(&model.DriverFile{}, &selectedPackage.Versions[0], packageFolder, writeFolder, nil, nil)
	if output.Code != constant.OutputError {
		t.Errorf("Error processing package file. Not validating invalid file")
	}
}

func TestInstallPackageFile(t *testing.T) {
	model.LoadXMLReadList("../xogRead.xml")

	folder := "../mock/xog/" + constant.FolderPackage
	LoadPackages(folder, "../mock/xog/mock_packages/")

	selectedPackage := GetAvailablePackages()[0]
	driverPath := folder + selectedPackage.Folder + selectedPackage.DriverFileName

	LoadDriver(driverPath)

	file := GetLoadedDriver().Files[4]

	packageFolder := folder + selectedPackage.Folder + selectedPackage.Versions[0].Folder + file.Type + "/"
	writeFolder := constant.FolderWrite + file.Type

	output := ProcessPackageFile(&file, &selectedPackage.Versions[0], packageFolder, writeFolder, nil, nil)

	mockEnvironments := &model.Environments{
		Source: &model.EnvType{
			Name:    "Mock Source Env",
			URL:     "Mock URL",
			Session: "Mock session",
		},
		Target: &model.EnvType{
			Name:    "Mock Source Env",
			URL:     "Mock URL",
			Session: "Mock session",
		},
	}

	soapMock := func(request, endpoint, proxy string, opts ...interface{}) (string, error) {
		file, _ := ioutil.ReadFile("../mock/xog/soap/soap_success_write_response.xml")
		return util.BytesToString(file), nil
	}

	output = InstallPackageFile(&file, mockEnvironments, soapMock)
	if output.Code != constant.OutputSuccess {
		t.Errorf("Error installing package file. Debug: %s", output.Debug)
	}

	soapMock = func(request, endpoint, proxy string, opts ...interface{}) (string, error) {
		return "", nil
	}

	output = InstallPackageFile(&file, mockEnvironments, soapMock)
	if output.Code != constant.OutputError {
		t.Errorf("Error installing package file. Not validating soap response")
	}

	deleteTestFolders()
}
