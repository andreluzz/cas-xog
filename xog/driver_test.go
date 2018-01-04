package xog

import (
	"errors"
	"github.com/andreluzz/cas-xog/constant"
	"github.com/andreluzz/cas-xog/model"
	"github.com/andreluzz/cas-xog/util"
	"io/ioutil"
	"os"
	"testing"
)

func deleteTestFolders() {
	os.RemoveAll(constant.FOLDER_DEBUG)
	os.RemoveAll(constant.FOLDER_READ)
	os.RemoveAll(constant.FOLDER_WRITE)
}

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

	sourceFolder, outputFolder := CreateFileFolder(constant.READ, fileType, "filename.xml")

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

	sourceFolder, outputFolder = CreateFileFolder(constant.WRITE, fileType, "filename.xml")

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

	sourceFolder, outputFolder = CreateFileFolder(constant.MIGRATE, fileType, "filename.xml")

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

func TestProcessDriverFileWrite(t *testing.T) {
	model.LoadXMLReadList("../xogRead.xml")

	LoadDriver("../mock/xog/xog.driver")
	file := GetLoadedDriver().Files[17]

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

	soapMock := func(request, endpoint string) (string, error) {
		file, _ := ioutil.ReadFile("../mock/xog/soap/soap_success_write_response.xml")
		return util.BytesToString(file), nil
	}

	sourceFolder := "../mock/xog/soap/"
	util.ValidateFolder(sourceFolder + file.Type)
	outputFolder := constant.FOLDER_DEBUG
	util.ValidateFolder(outputFolder + file.Type)

	output := ProcessDriverFile(&file, constant.WRITE, sourceFolder, outputFolder, mockEnvironments, soapMock)
	if output.Code != constant.OUTPUT_SUCCESS {
		t.Errorf("Error processing driver file. Debug: %s", output.Debug)
	}

}

func TestProcessDriverFileActionExportToExcel(t *testing.T) {
	model.LoadXMLReadList("../xogRead.xml")

	file := model.DriverFile{
		Type:          constant.RESOURCE_INSTANCE,
		Code:          "*",
		Path:          "instances.xml",
		InstanceTag:   "Resource",
		ExcelFile:     "test.xlsx",
		ExportToExcel: true,
		MatchExcel: []model.MatchExcel{
			{
				AttributeName: "resourceId",
			},
			{
				AttributeName: "displayName",
				XPath:         "//PersonalInformation",
			},
			{
				AttributeName: "emailAddress",
				XPath:         "//PersonalInformation",
			},
			{
				AttributeName: "firstName",
				XPath:         "//PersonalInformation",
			},
			{
				AttributeName: "lastName",
				XPath:         "//PersonalInformation",
			},
			{
				AttributeName: "unitPath",
				XPath:         "//OBSAssoc[@id='corpLocationOBS']",
			},
			{
				AttributeName: "unitPath",
				XPath:         "//OBSAssoc[@id='resourcePool']",
			},
			{
				XPath: "//ColumnValue[@name='partition_code']",
			},
		},
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

	soapMock := func(request, endpoint string) (string, error) {
		file, _ := ioutil.ReadFile("../mock/xog/soap/soap_read_resources_instance_response.xml")
		return util.BytesToString(file), nil
	}

	sourceFolder := "../" + constant.FOLDER_READ
	util.ValidateFolder(sourceFolder + file.Type)
	outputFolder := "../" + constant.FOLDER_DEBUG
	util.ValidateFolder(outputFolder + file.Type)

	output := ProcessDriverFile(&file, constant.READ, sourceFolder, outputFolder, mockEnvironments, soapMock)
	if output.Code != constant.OUTPUT_SUCCESS {
		t.Fatalf("Error processing driver file. Action migrate with errors. Debug: %s", output.Debug)
	}

	os.RemoveAll(constant.FOLDER_MIGRATION)
}

func TestProcessDriverFileActionMigrate(t *testing.T) {
	packageMockFolder := "../" + constant.FOLDER_MOCK + "migration/"
	file := model.DriverFile{
		Type:          constant.MIGRATION,
		Template:      packageMockFolder + "template.xml",
		ExcelFile:     packageMockFolder + "data.xlsx",
		InstanceTag:   "instance",
		ExcelStartRow: "1",
		MatchExcel: []model.MatchExcel{
			{
				Col:           1,
				AttributeName: "instanceCode",
			},
			{
				Col:   1,
				XPath: "//ColumnValue[@name='code']",
			},
			{
				Col:   2,
				XPath: "//ColumnValue[@name='name']",
			},
			{
				Col:   3,
				XPath: "//ColumnValue[@name='status_novo']",
			},
			{
				Col:         4,
				XPath:       "//ColumnValue[@name='multivalue_status']",
				MultiValued: true,
				Separator:   ";",
			},
			{
				Col:   5,
				XPath: "//ColumnValue[@name='analista']",
			},
		},
	}
	output := ProcessDriverFile(&file, constant.MIGRATE, "", "", nil, nil)
	if output.Code != constant.OUTPUT_SUCCESS {
		t.Errorf("Error processing driver file. Action migrate with errors. Debug: %s", output.Debug)
	}

	file = model.DriverFile{
		Type:          constant.MIGRATION,
		Template:      packageMockFolder + "template.xml",
		ExcelFile:     packageMockFolder + "data.xlsx",
		InstanceTag:   "instance",
		ExcelStartRow: "1",
		MatchExcel: []model.MatchExcel{
			{
				Col:           1,
				XPath:         "invalid_xpath",
				AttributeName: "name",
			},
		},
	}
	output = ProcessDriverFile(&file, constant.MIGRATE, "", "", nil, nil)
	if output.Code != constant.OUTPUT_ERROR {
		t.Errorf("Error processing driver file. Action migrate with errors not being validated.")
	}
}

func TestProcessDriverFileReturnInitXMLError(t *testing.T) {
	model.LoadXMLReadList("../xogRead.xml")

	file := model.DriverFile{
		Type: constant.UNDEFINED,
	}
	output := ProcessDriverFile(&file, constant.READ, "", "", nil, nil)
	if output.Code != constant.OUTPUT_ERROR {
		t.Errorf("Error processing driver file. Not treating invalid InitXML. Debug: %s", output.Debug)
	}
}

func TestProcessDriverFileReturnRunXMLError(t *testing.T) {
	model.LoadXMLReadList("../xogRead.xml")

	file := model.DriverFile{
		Type: constant.PROCESS,
		Code: "code",
		Path: "test.xml",
	}

	soapMock := func(request, endpoint string) (string, error) {
		return "", errors.New("soap mock error")
	}

	environments := model.Environments{
		Source: &model.EnvType{},
		Target: &model.EnvType{},
	}

	output := ProcessDriverFile(&file, constant.READ, constant.FOLDER_DEBUG, "", &environments, soapMock)
	if output.Code != constant.OUTPUT_ERROR {
		t.Errorf("Error processing driver file. Not treating invalid RunXML. Debug: %s", output.Debug)
	}
}

func TestProcessDriverFileReturnValidateXMLError(t *testing.T) {
	model.LoadXMLReadList("../xogRead.xml")

	file := model.DriverFile{
		Type: constant.PROCESS,
		Code: "code",
		Path: "test.xml",
	}

	soapMock := func(request, endpoint string) (string, error) {
		return "", nil
	}

	environments := model.Environments{
		Source: &model.EnvType{},
		Target: &model.EnvType{},
	}

	output := ProcessDriverFile(&file, constant.READ, constant.FOLDER_DEBUG, "", &environments, soapMock)
	if output.Code != constant.OUTPUT_ERROR {
		t.Errorf("Error processing driver file. Not treating invalid validate check. Debug: %s", output.Debug)
	}
}

func TestProcessDriverFileActionMigrateInvalidFileType(t *testing.T) {
	file := model.DriverFile{
		Type: constant.LOOKUP,
	}
	output := ProcessDriverFile(&file, constant.MIGRATE, "", "", nil, nil)
	if output.Code != constant.OUTPUT_WARNING {
		t.Errorf("Error processing driver file. Not treating invalid action and file type. Debug: %s", output.Debug)
	}
}

func TestProcessDriverFileActionReadInvalidFileTypeMigration(t *testing.T) {
	file := model.DriverFile{
		Type: constant.MIGRATION,
	}
	output := ProcessDriverFile(&file, constant.READ, "", "", nil, nil)
	if output.Code != constant.OUTPUT_WARNING {
		t.Errorf("Error processing driver file. Not treating invalid action and file type. Debug: %s", output.Debug)
	}
}

func TestProcessDriverFileActionRead(t *testing.T) {
	model.LoadXMLReadList("../xogRead.xml")

	LoadDriver("../mock/xog/xog.driver")
	file := GetLoadedDriver().Files[17]

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

	sourceFolder := "../mock/xog/soap/"
	util.ValidateFolder(sourceFolder + file.Type)
	outputFolder := constant.FOLDER_DEBUG
	util.ValidateFolder(outputFolder + file.Type)

	soapMock := func(request, endpoint string) (string, error) {
		file, _ := ioutil.ReadFile("../mock/xog/soap/soap_success_read_response.xml")
		return util.BytesToString(file), nil
	}

	output := ProcessDriverFile(&file, constant.READ, sourceFolder, outputFolder, mockEnvironments, soapMock)
	if output.Code != constant.OUTPUT_SUCCESS {
		t.Errorf("Error processing driver file. Debug: %s", output.Debug)
	}
}

func TestProcessDriverFileActionReadNeedAux(t *testing.T) {
	model.LoadXMLReadList("../xogRead.xml")

	file := model.DriverFile{
		Type:            constant.PROCESS,
		CopyPermissions: "code",
		Code:            "code",
		Path:            "test.xml",
	}

	mockEnvironments := &model.Environments{
		Source: &model.EnvType{
			Name:    "Mock Source Env",
			URL:     "Mock URL",
			Session: "Mock session",
		},
		Target: &model.EnvType{
			Name:    "Mock Source Env",
			URL:     "Aux Mock URL",
			Session: "Mock session",
		},
	}

	sourceFolder := constant.FOLDER_READ
	util.ValidateFolder(sourceFolder + file.Type)
	outputFolder := constant.FOLDER_DEBUG
	util.ValidateFolder(outputFolder + file.Type)

	soapMock := func(request, endpoint string) (string, error) {
		file, _ := ioutil.ReadFile("../mock/xog/soap/soap_success_read_process_response.xml")
		return util.BytesToString(file), nil
	}

	output := ProcessDriverFile(&file, constant.READ, sourceFolder, outputFolder, mockEnvironments, soapMock)
	if output.Code != constant.OUTPUT_SUCCESS {
		t.Errorf("Error processing driver file. Debug: %s", output.Debug)
	}
}

func TestProcessDriverFileActionReadNeedAuxErrorInvalidCheck(t *testing.T) {
	model.LoadXMLReadList("../xogRead.xml")

	file := model.DriverFile{
		Type:    constant.VIEW,
		Code:    "code",
		ObjCode: "project",
		Path:    "test.xml",
	}

	mockEnvironments := &model.Environments{
		Source: &model.EnvType{
			Name:    "Mock Source Env",
			URL:     "Mock URL",
			Session: "Mock session",
		},
		Target: &model.EnvType{
			Name:    "Mock Source Env",
			URL:     "Aux_Mock_URL",
			Session: "Mock session",
		},
	}

	sourceFolder := constant.FOLDER_READ
	util.ValidateFolder(sourceFolder + file.Type)
	outputFolder := constant.FOLDER_DEBUG
	util.ValidateFolder(outputFolder + file.Type)

	soapMock := func(request, endpoint string) (string, error) {
		if endpoint == "Aux_Mock_URL" {
			return "", nil
		}
		file, _ := ioutil.ReadFile("../mock/xog/soap/soap_success_read_process_response.xml")
		return util.BytesToString(file), nil
	}

	output := ProcessDriverFile(&file, constant.READ, sourceFolder, outputFolder, mockEnvironments, soapMock)
	if output.Code != constant.OUTPUT_ERROR {
		t.Errorf("Error processing driver file. Not validating aux response. Debug: %s", output.Debug)
	}
}

func TestProcessDriverFileActionReadAuxValidateError(t *testing.T) {
	model.LoadXMLReadList("../xogRead.xml")

	file := model.DriverFile{
		Type:            constant.PROCESS,
		CopyPermissions: "code",
		Code:            "code",
		Path:            "test.xml",
	}

	mockEnvironments := &model.Environments{
		Source: &model.EnvType{
			Name:    "Mock Source Env",
			URL:     "Mock URL",
			Session: "Mock session",
		},
		Target: &model.EnvType{
			Name:    "Mock Source Env",
			URL:     "Aux Mock URL",
			Session: "Mock session",
		},
	}

	sourceFolder := constant.FOLDER_READ
	util.ValidateFolder(sourceFolder + file.Type)
	outputFolder := constant.FOLDER_DEBUG
	util.ValidateFolder(outputFolder + file.Type)

	soapMock := func(request, endpoint string) (string, error) {
		file, _ := ioutil.ReadFile("../mock/xog/soap/soap_read_process_no_output_response.xml")
		return util.BytesToString(file), nil
	}

	output := ProcessDriverFile(&file, constant.READ, sourceFolder, outputFolder, mockEnvironments, soapMock)
	if output.Code != constant.OUTPUT_ERROR {
		t.Errorf("Error processing driver file. Not treating aux output validatin error. Debug: %s", output.Debug)
	}
}

func TestProcessDriverFileActionReadTransformError(t *testing.T) {
	model.LoadXMLReadList("../xogRead.xml")

	LoadDriver("../mock/xog/xog.driver")
	file := GetLoadedDriver().Files[17]

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

	sourceFolder := constant.FOLDER_READ
	util.ValidateFolder(sourceFolder + file.Type)
	outputFolder := constant.FOLDER_DEBUG
	util.ValidateFolder(outputFolder + file.Type)

	soapMock := func(request, endpoint string) (string, error) {
		return `<XOGOutput>
        	<Object type="contentPack"/>
        	<Status elapsedTime="0.789 seconds" state="SUCCESS"/>
        	<Statistics failureRecords="0" insertedRecords="0" totalNumberOfRecords="1" updatedRecords="1"/>
        	<Records/>
    	</XOGOutput>`, nil
	}

	output := ProcessDriverFile(&file, constant.READ, sourceFolder, outputFolder, mockEnvironments, soapMock)
	if output.Code != constant.OUTPUT_ERROR {
		t.Errorf("Error processing driver file. Debug: %s", output.Debug)
	}

	deleteTestFolders()
}
