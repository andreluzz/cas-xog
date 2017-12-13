package transform

import (
	"github.com/andreluzz/cas-xog/constant"
	"github.com/andreluzz/cas-xog/model"
	"github.com/beevik/etree"
	"testing"
)

func TestExecuteToReturnObjectFull(t *testing.T) {
	file := model.DriverFile{
		Code: "obj_sistema",
		Type: constant.OBJECT,
	}

	xog := etree.NewDocument()
	xog.ReadFromFile(packageMockFolder + "object_full_xog.xml")
	err := Execute(xog, nil, &file)

	if err != nil {
		t.Fatalf("Error transforming object XOG file. Debug: %s", err.Error())
	}

	if readMockResultAndCompare(xog, "object_result.xml") == false {
		t.Errorf("Error transforming object XOG file. Invalid result XML.")
	}
}

func TestExecuteToReturnObjectElementAttribute(t *testing.T) {
	file := model.DriverFile{
		Code: "obj_sistema",
		Type: constant.OBJECT,
		Elements: []model.Element{
			{
				Code: "aprovador",
				Type: constant.ELEMENT_TYPE_ATTRIBUTE,
			},
			{
				Code: "status",
				Type: constant.ELEMENT_TYPE_ATTRIBUTE,
			},
		},
	}

	xog := etree.NewDocument()
	xog.ReadFromFile(packageMockFolder + "object_full_xog.xml")
	aux := etree.NewDocument()
	aux.ReadFromFile(packageMockFolder + "object_full_aux.xml")
	err := Execute(xog, aux, &file)

	if err != nil {
		t.Fatalf("Error transforming object XOG file. Debug: %s", err.Error())
	}

	if readMockResultAndCompare(xog, "object_elements_result.xml") == false {
		t.Errorf("Error transforming object XOG file. Invalid result XML.")
	}
}

func TestExecuteToReturnObjectTargetPartition(t *testing.T) {
	file := model.DriverFile{
		Code:            "obj_sistema",
		Type:            constant.OBJECT,
		TargetPartition: "NIKU.ROOT",
	}

	xog := etree.NewDocument()
	xog.ReadFromFile(packageMockFolder + "object_full_xog.xml")
	err := Execute(xog, nil, &file)

	if err != nil {
		t.Fatalf("Error transforming object XOG file. Debug: %s", err.Error())
	}

	count := 0
	for _, e := range xog.FindElements("//*[@partitionCode]") {
		if e.SelectAttrValue("partitionCode", constant.UNDEFINED) != file.TargetPartition {
			count++
		}
	}

	if count > 0 {
		t.Errorf("Error transforming object XOG file. Expected 0 got %d elements from other partitions", count)
	}

	if readMockResultAndCompare(xog, "object_result_target_partition.xml") == false {
		t.Errorf("Error transforming object XOG file. Invalid result XML.")
	}
}

func TestExecuteToReturnObjectChangeSourcePartitionToTarget(t *testing.T) {
	file := model.DriverFile{
		Code:            "obj_sistema",
		Type:            constant.OBJECT,
		SourcePartition: "partition10",
		TargetPartition: "NIKU.ROOT",
	}

	xog := etree.NewDocument()
	xog.ReadFromFile(packageMockFolder + "object_full_xog.xml")
	err := Execute(xog, nil, &file)

	if err != nil {
		t.Fatalf("Error transforming object XOG file. Debug: %s", err.Error())
	}

	count := len(xog.FindElements("//object[@code='" + file.Code + "']/*[@partitionCode='" + file.SourcePartition + "']"))
	if count > 0 {
		t.Errorf("Error transforming object XOG file. Expected 0 got %d source partition (%s) elements", count, file.SourcePartition)
	}

	count = len(xog.FindElements("//object[@code='" + file.Code + "']/*[@partitionCode='" + file.TargetPartition + "']"))
	if count != 3 {
		t.Errorf("Error transforming object XOG file. Expected 3 got %d target partition (%s) elements", count, file.TargetPartition)
	}

	if readMockResultAndCompare(xog, "object_result_from_source_to_target_partition.xml") == false {
		t.Errorf("Error transforming object XOG file. Invalid result XML.")
	}
}

func TestExecuteToReturnObjectChangePartitionModel(t *testing.T) {
	file := model.DriverFile{
		Code:           "obj_sistema",
		Type:           constant.OBJECT,
		PartitionModel: "NEW_PARTITION_MODEL",
	}

	xog := etree.NewDocument()
	xog.ReadFromFile(packageMockFolder + "object_full_xog.xml")
	err := Execute(xog, nil, &file)

	if err != nil {
		t.Fatalf("Error transforming object XOG file. Debug: %s", err.Error())
	}

	for _, e := range xog.FindElements("*[@partitionModelCode]") {
		value := e.SelectAttrValue("partitionModelCode", constant.UNDEFINED)
		if value != file.PartitionModel {
			t.Fatalf("Error transforming object XOG file. Expected %s got %s partition model", file.PartitionModel, value)
		}
	}

	if readMockResultAndCompare(xog, "object_result_change_partition_model.xml") == false {
		t.Errorf("Error transforming object XOG file. Invalid result XML.")
	}
}

func TestExecuteToReturnObjectRemoveAttribute(t *testing.T) {
	file := model.DriverFile{
		Code: "cas_environment",
		Type: constant.OBJECT,
		Elements: []model.Element{
			{
				XPath:  "//customAttribute[@code='analista']",
				Action: constant.ACTION_REMOVE,
			},
		},
	}

	xog := etree.NewDocument()
	xog.ReadFromFile(packageMockFolder + "object_full_xog.xml")
	err := Execute(xog, nil, &file)

	if err != nil {
		t.Fatalf("Error transforming object XOG file. Debug: %s", err.Error())
	}

	if readMockResultAndCompare(xog, "object_remove_attribute_result.xml") == false {
		t.Errorf("Error transforming object XOG file. Invalid result XML.")
	}
}
