package transform

import (
	"github.com/andreluzz/cas-xog/constant"
	"github.com/andreluzz/cas-xog/model"
	"github.com/beevik/etree"
	"strings"
	"testing"
)

func TestProcessPackageToReplaceTargetPartitionModel(t *testing.T) {
	file := model.DriverFile{
		Type: constant.OBJECT,
		Path: "package_change_partition.xml",
	}

	def := []model.Definition{
		{
			Action:  constant.PACKAGE_ACTION_CHANGE_PARTITION_MODEL,
			Default: "PartitionModel1",
			Value:   "NEW_PARTITION_MODEL",
		},
	}

	folder := "../" + constant.FOLDER_WRITE + file.Type
	err := ProcessPackageFile(file, packageMockFolder, folder, def)

	if err != nil {
		t.Fatalf("Error processing package file. Debug: %s", err.Error())
	}

	result := etree.NewDocument()
	err = result.ReadFromFile(folder + "/" + file.Path)
	if err != nil {
		t.Fatalf("Error processing package file. Debug: %s", err.Error())
	}

	element := result.FindElement("//object[@partitionModelCode]")
	partitionModelCodeValue := element.SelectAttrValue("partitionModelCode", constant.UNDEFINED)

	if partitionModelCodeValue != def[0].Value {
		t.Errorf("Error processing package file. Expected %s got %s partitionModelCode.", def[0].Value, partitionModelCodeValue)
	}
}

func TestProcessPackageToDiscardObjectWithoutPartitionModel(t *testing.T) {
	file := model.DriverFile{
		Type: constant.OBJECT,
		Path: "package_object_with_no_partition_model.xml",
	}

	def := []model.Definition{
		{
			Action:  constant.PACKAGE_ACTION_CHANGE_PARTITION_MODEL,
			Default: "PartitionModel1",
			Value:   "NEW_PARTITION_MODEL",
		},
	}

	folder := "../" + constant.FOLDER_WRITE + file.Type
	err := ProcessPackageFile(file, packageMockFolder, folder, def)

	if err != nil {
		t.Fatalf("Error processing package file. Debug: %s", err.Error())
	}

	result := etree.NewDocument()
	err = result.ReadFromFile(folder + "/" + file.Path)
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
		Type: constant.OBJECT,
		Path: "package_change_partition.xml",
	}

	def := []model.Definition{
		{
			Action:  constant.PACKAGE_ACTION_CHANGE_PARTITION,
			Default: "partition20",
			Value:   "partition20",
		},
		{
			Action:  constant.PACKAGE_ACTION_CHANGE_PARTITION,
			Default: "partition10",
			Value:   "NIKU.ROOT",
		},
	}

	folder := "../" + constant.FOLDER_WRITE + file.Type
	err := ProcessPackageFile(file, packageMockFolder, folder, def)

	if err != nil {
		t.Fatalf("Error processing package file. Debug: %s", err.Error())
	}

	result := etree.NewDocument()
	err = result.ReadFromFile(folder + "/" + file.Path)
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
		Type: constant.PROCESS,
		Path: "package_replace_string.xml",
	}

	def := []model.Definition{
		{
			Action:  constant.PACKAGE_ACTION_REPLACE_STRING,
			From:    "Test cas-xog 002",
			To:      "Test cas-xog ##DEFINITION_VALUE##",
			Value:   "V0005",
			Default: "V0005",
		},
		{
			Action:  constant.PACKAGE_ACTION_REPLACE_STRING,
			From:    "Test cas-xog 002",
			To:      "Test cas-xog ##DEFINITION_VALUE##",
			Value:   "002",
			Default: "V0005",
		},
		{
			Action:  constant.PACKAGE_ACTION_REPLACE_STRING,
			From:    "Test cas-xog 002",
			To:      "Test cas-xog ##DEFINITION_VALUE##",
			Value:   "",
			Default: "002",
		},
		{
			Action: constant.PACKAGE_ACTION_REPLACE_STRING,
			From:   "Test cas-xog 002",
			To:     "Test cas-xog ##DEFINITION_VALUE##",
			Value:  "V0005",
		},
	}

	folder := "../" + constant.FOLDER_WRITE + file.Type
	err := ProcessPackageFile(file, packageMockFolder, folder, def)

	if err != nil {
		t.Fatalf("Error processing package file. Debug: %s", err.Error())
	}

	result := etree.NewDocument()
	err = result.ReadFromFile(folder + "/" + file.Path)
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

	err := ProcessPackageFile(model.DriverFile{}, "", "", nil)

	if err == nil {
		t.Errorf("Error processing package file. Debug: not validating if driver file is null")
	}
}
