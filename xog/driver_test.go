package xog

import (
	"github.com/andreluzz/cas-xog/constant"
	"os"
	"testing"
	"github.com/andreluzz/cas-xog/model"
	"io/ioutil"
	"github.com/andreluzz/cas-xog/util"
)

func TestValidateLoadedDriver(t *testing.T) {
	result := ValidateLoadedDriver()
	if result {
		t.Errorf("Error validating loaded driver when driverXOG is nil. Expected validation false received true")
	}
	LoadDriver("")
	result = ValidateLoadedDriver()
	if result {
		t.Errorf("Error validating loaded driver when driverXOG total files is 0. Expected validation false received true")
	}
	LoadDriver("../mock/xog/xog.driver")
	result = ValidateLoadedDriver()
	if !result {
		t.Errorf("Error validating loaded driver when driverXOG is valid. Expected validation true received false")
	}
}

func TestGetDriversList(t *testing.T) {
	driverList, err := GetDriversList("../mock/xog/")
	if err != nil {
		t.Errorf("Error getting drivers list from folder. Debug: %s", err.Error())
	}
	if driverList == nil || len(driverList) <= 0 {
		t.Errorf("Error getting drivers list from folder. Expected 3 received %d", len(driverList))
	}
}

func TestGetDriversListInvalidFolder(t *testing.T) {
	driverList, err := GetDriversList("")
	if err == nil {
		t.Errorf("Error getting drivers list from folder. Not validating invalid folder")
	}
	if driverList != nil || len(driverList) > 0 {
		t.Errorf("Error getting drivers list from folder. Expected 0 files")
	}
}

func TestCreateFileFolder(t *testing.T) {
	fileType := constant.PROCESS

	sourceFolder, outputFolder := CreateFileFolder(constant.READ, fileType)

	if sourceFolder != constant.FOLDER_READ {
		t.Errorf("Error creating file folder, expected source folder %s and received %s", constant.FOLDER_READ, sourceFolder)
	}
	if outputFolder != constant.FOLDER_WRITE {
		t.Errorf("Error creating file folder, expected output folder %s and received %s", constant.FOLDER_WRITE, outputFolder)
	}
	folder := constant.FOLDER_READ + fileType
	_, dirErr := os.Stat(folder)
	if os.IsNotExist(dirErr) {
		t.Errorf("Error creating file folder action %s, the read folder %s was not created", constant.READ, folder)
	}
	os.RemoveAll(folder)
	os.RemoveAll(sourceFolder)

	folder = outputFolder + fileType
	_, dirErr = os.Stat(folder)
	if os.IsNotExist(dirErr) {
		t.Errorf("Error creating file folder action %s, the output folder %s was not created", constant.READ, folder)
	}
	os.RemoveAll(folder)
	os.RemoveAll(outputFolder)

	sourceFolder, outputFolder = CreateFileFolder(constant.WRITE, fileType)

	if sourceFolder != constant.FOLDER_WRITE {
		t.Errorf("Error creating file folder, expected source folder %s and received %s", constant.FOLDER_READ, sourceFolder)
	}
	if outputFolder != constant.FOLDER_DEBUG {
		t.Errorf("Error creating file folder, expected output folder %s and received %s", constant.FOLDER_WRITE, outputFolder)
	}
	folder = outputFolder + fileType
	_, dirErr = os.Stat(folder)
	if os.IsNotExist(dirErr) {
		t.Errorf("Error creating file folder action %s, the folder %s was not created", constant.WRITE, folder)
	}
	os.RemoveAll(folder)
	os.RemoveAll(outputFolder)

	sourceFolder, outputFolder = CreateFileFolder(constant.MIGRATE, fileType)

	if sourceFolder != constant.FOLDER_WRITE {
		t.Errorf("Error creating file folder, expected source folder %s and received %s", constant.FOLDER_READ, sourceFolder)
	}
	if outputFolder != constant.FOLDER_MIGRATION {
		t.Errorf("Error creating file folder, expected output folder %s and received %s", constant.FOLDER_WRITE, outputFolder)
	}
	folder = outputFolder + fileType
	_, dirErr = os.Stat(folder)
	if os.IsNotExist(dirErr) {
		t.Errorf("Error creating file folder action %s, the folder %s was not created", constant.WRITE, folder)
	}
	os.RemoveAll(folder)
	os.RemoveAll(outputFolder)
}

func TestLoadDriver(t *testing.T) {
	total, err := LoadDriver("../mock/xog/xog.driver")

	if total != 29 {
		t.Errorf("Error loading driver expected %d and received %d", 29, total)
	}

	if err != nil {
		t.Errorf("Error loading driver. Debug: %s", err.Error())
	}

	driver := GetLoadedDriver()
	if driver.Files[3].Type != constant.VIEW {
		t.Errorf("Error loading driver. Incorrect execution order")
	}
}

func TestLoadDriverInvalidVersion(t *testing.T) {
	total, err := LoadDriver("../mock/xog/invalidVersion.driver")

	if total != 0 {
		t.Errorf("Error loading driver expected %d and received %d", 0, total)
	}

	if err == nil {
		t.Errorf("Error loading driver. Not catching error with invalid driver version")
	}
}

func TestLoadDriverInvalidTagFile(t *testing.T) {
	total, err := LoadDriver("../mock/xog/invalidTagFile.driver")

	if total != 0 {
		t.Errorf("Error loading driver expected %d and received %d", 0, total)
	}

	if err == nil {
		t.Errorf("Error loading driver. Not catching error with invalid tag '<file>'")
	}
}

func TestLoadDriverInvalidPath(t *testing.T) {
	total, err := LoadDriver("")

	if total != 0 {
		t.Errorf("Error loading driver expected %d and received %d", 0, total)
	}

	if err == nil {
		t.Errorf("Error loading driver. Not catching error with invalid file path")
	}
}

func TestProcessDriverFile(t *testing.T) {
	model.LoadXMLReadList("../xogRead.xml")

	LoadDriver("../mock/xog/xog.driver")
	file := GetLoadedDriver().Files[17]

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

	sourceFolder := "../mock/xog/soap/"
	util.ValidateFolder(sourceFolder+file.Type)
	outputFolder := constant.FOLDER_DEBUG
	util.ValidateFolder(outputFolder+file.Type)

	output := ProcessDriverFile(&file, constant.WRITE, sourceFolder, outputFolder, mockEnvironments, soapMock)
	if output.Code != constant.OUTPUT_SUCCESS {
		t.Errorf("Error installing package file. Debug: %s", output.Debug)
	}

}