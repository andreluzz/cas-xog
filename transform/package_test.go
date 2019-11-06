package transform

import (
	"io/ioutil"
	"strings"
	"testing"

	"github.com/andreluzz/cas-xog/constant"
	"github.com/andreluzz/cas-xog/model"
	"github.com/andreluzz/cas-xog/util"
	"github.com/beevik/etree"
)

func TestProcessPackageToReplaceTargetPartitionModel(t *testing.T) {
	file := model.DriverFile{
		Type: constant.TypeObject,
		Path: "package_change_partition.xml",
	}

	def := []model.Definition{
		{
			Action:  constant.PackageActionChangePartitionModel,
			Default: "PartitionModel1",
			Value:   "NEW_PARTITION_MODEL",
		},
	}

	folder := "../" + constant.FolderWrite + file.Type
	output := ProcessPackageFile(&file, packageMockFolder, folder, def)

	if output.Code != constant.OutputSuccess {
		t.Fatalf("Error processing package file. Code: %s | Debug: %s", output.Code, output.Debug)
	}

	result := etree.NewDocument()
	err := result.ReadFromFile(folder + "/" + file.Path)
	if err != nil {
		t.Fatalf("Error processing package file. Debug: %s", err.Error())
	}

	element := result.FindElement("//object[@partitionModelCode]")
	partitionModelCodeValue := element.SelectAttrValue("partitionModelCode", constant.Undefined)

	if partitionModelCodeValue != def[0].Value {
		t.Errorf("Error processing package file. Expected %s got %s partitionModelCode.", def[0].Value, partitionModelCodeValue)
	}
}

func TestProcessPackageToDiscardObjectWithoutPartitionModel(t *testing.T) {
	file := model.DriverFile{
		Type: constant.TypeObject,
		Path: "package_object_with_no_partition_model.xml",
	}

	def := []model.Definition{
		{
			Action:  constant.PackageActionChangePartitionModel,
			Default: "PartitionModel1",
			Value:   "NEW_PARTITION_MODEL",
		},
	}

	folder := "../" + constant.FolderWrite + file.Type
	output := ProcessPackageFile(&file, packageMockFolder, folder, def)

	if output.Code != constant.OutputSuccess {
		t.Fatalf("Error processing package file. Code: %s | Debug: %s", output.Code, output.Debug)
	}

	result := etree.NewDocument()
	err := result.ReadFromFile(folder + "/" + file.Path)
	if err != nil {
		t.Fatalf("Error processing package file. Debug: %s", err.Error())
	}

	element := result.FindElement("//object[@partitionModelCode]")
	if element != nil {
		t.Errorf("Error processing package file. Expected no element with partitionModelCode.")
	}
}

func TestProcessPackageToReplaceTargetPartition(t *testing.T) {
	file := model.DriverFile{
		Type: constant.TypeObject,
		Path: "package_change_partition.xml",
	}

	def := []model.Definition{
		{
			Action:  constant.PackageActionChangePartition,
			Default: "partition20",
			Value:   "partition20",
		},
		{
			Action:  constant.PackageActionChangePartition,
			Default: "partition10",
			Value:   "NIKU.ROOT",
		},
	}

	folder := "../" + constant.FolderWrite + file.Type
	output := ProcessPackageFile(&file, packageMockFolder, folder, def)

	if output.Code != constant.OutputSuccess {
		t.Fatalf("Error processing package file. Code: %s | Debug: %s", output.Code, output.Debug)
	}

	result := etree.NewDocument()
	err := result.ReadFromFile(folder + "/" + file.Path)
	if err != nil {
		t.Fatalf("Error processing package file. Debug: %s", err.Error())
	}

	count := len(result.FindElements("//object[@partitionCode='partition10']"))
	if count > 0 {
		t.Errorf("Error processing package file. Expected 0 got %d elements with old partitionModel.", count)
	}
}

func TestProcessPackageToProcessDefinitionReplaceString(t *testing.T) {
	file := model.DriverFile{
		Type: constant.TypeProcess,
		Path: "package_replace_string.xml",
	}

	def := []model.Definition{
		{
			Action:  constant.PackageActionReplaceString,
			From:    "Test cas-xog 002",
			To:      "Test cas-xog ##DEFINITION_VALUE##",
			Value:   "V0005",
			Default: "V0005",
		},
		{
			Action:  constant.PackageActionReplaceString,
			From:    "Test cas-xog 002",
			To:      "Test cas-xog ##DEFINITION_VALUE##",
			Value:   "002",
			Default: "V0005",
		},
		{
			Action:  constant.PackageActionReplaceString,
			From:    "Test cas-xog 002",
			To:      "Test cas-xog ##DEFINITION_VALUE##",
			Value:   "",
			Default: "002",
		},
		{
			Action: constant.PackageActionReplaceString,
			From:   "Test cas-xog 002",
			To:     "Test cas-xog ##DEFINITION_VALUE##",
			Value:  "V0005",
		},
	}

	folder := "../" + constant.FolderWrite + file.Type
	output := ProcessPackageFile(&file, packageMockFolder, folder, def)

	if output.Code != constant.OutputSuccess {
		t.Fatalf("Error processing package file. Code: %s | Debug: %s", output.Code, output.Debug)
	}

	result := etree.NewDocument()
	err := result.ReadFromFile(folder + "/" + file.Path)
	if err != nil {
		t.Fatalf("Error processing package file. Debug: %s", err.Error())
	}

	resultString, _ := result.WriteToString()
	count := strings.Count(resultString, def[0].From)
	if count > 0 {
		t.Errorf("Error processing package file. Expected 0 got %d substring that should have being replaced.", count)
	}
}

func TestProcessPackageToReturnErrorFileIsNil(t *testing.T) {
	output := ProcessPackageFile(nil, "", "", nil)

	if output.Code != constant.OutputError {
		t.Fatalf("Error processing package file. Code: %s | Debug: not validating if driver file is nil", output.Code)
	}
}

func TestProcessPackageToReturnErrorTypeCannotTransform(t *testing.T) {
	file := model.DriverFile{
		Type:             constant.TypePortlet,
		Path:             "package_change_partition.xml",
		PackageTransform: true,
	}

	folder := "../" + constant.FolderWrite + file.Type
	output := ProcessPackageFile(&file, packageMockFolder, folder, nil)

	if output.Code != constant.OutputWarning {
		t.Fatalf("Error processing package file. Code: %s | Debug: error validating if trying to transform an invalid type", output.Code)
	}
}

func TestProcessPackageToTransform(t *testing.T) {
	file := model.DriverFile{
		Type:             constant.TypeView,
		Code:             "cas_environmentProperties",
		ObjCode:          "cas_environment",
		Path:             "package_transform_view_source.xml",
		PackageTransform: true,
		Sections: []model.Section{
			{
				Action:         constant.ActionInsert,
				SourcePosition: "2",
			},
		},
	}

	soapMock := func(request, endpoint, proxy string) (string, error) {
		file, _ := ioutil.ReadFile("../mock/transform/package_transform_view_target.xml")
		return util.BytesToString(file), nil
	}
	file.RunAuxXML(&model.EnvType{}, soapMock)

	folder := "../" + constant.FolderWrite + file.Type
	output := ProcessPackageFile(&file, packageMockFolder, folder, nil)

	if output.Code != constant.OutputSuccess {
		t.Fatalf("Error processing package file. Code: %s | Debug: %s", output.Code, output.Debug)
	}
}

func TestProcessPackageToReturnErrorTransformValidate(t *testing.T) {
	file := model.DriverFile{
		Type:             constant.TypeView,
		Path:             "package_transform_view_source.xml",
		PackageTransform: true,
	}
	soapMock := func(request, endpoint, proxy string) (string, error) {
		return "", nil
	}
	file.RunAuxXML(&model.EnvType{}, soapMock)

	folder := "../" + constant.FolderWrite + file.Type
	output := ProcessPackageFile(&file, packageMockFolder, folder, nil)

	if output.Code != constant.OutputError {
		t.Fatalf("Error processing package file. Code: %s | Debug: not validating if transform validate without errors ", output.Code)
	}
}

func TestProcessPackageToReturnErrorTransformExecute(t *testing.T) {
	file := model.DriverFile{
		Type:             constant.TypeView,
		Code:             "viewCode",
		ObjCode:          "cas_environment",
		Path:             "package_transform_view_source.xml",
		TargetPartition:  "Partition1",
		PackageTransform: true,
		Sections: []model.Section{
			{
				Action:         constant.ActionInsert,
				SourcePosition: "2",
			},
		},
	}
	soapMock := func(request, endpoint, proxy string) (string, error) {
		file, _ := ioutil.ReadFile("../mock/transform/package_transform_view_target.xml")
		return util.BytesToString(file), nil
	}
	file.RunAuxXML(&model.EnvType{}, soapMock)

	folder := "../" + constant.FolderWrite + file.Type
	output := ProcessPackageFile(&file, packageMockFolder, folder, nil)

	if output.Code != constant.OutputError {
		t.Fatalf("Error processing package file, not validating if transform executed without errors. Code: %s | Debug: %s", output.Code, output.Debug)
	}
}

func TestProcessPackageToReturnErrorFilePathIsUndefined(t *testing.T) {
	file := model.DriverFile{
		Type: constant.TypePortlet,
		Path: constant.Undefined,
	}

	output := ProcessPackageFile(&file, "", "", nil)

	if output.Code != constant.OutputError {
		t.Fatalf("Error processing package file. Code: %s | Debug: %s", output.Code, output.Debug)
	}
}
