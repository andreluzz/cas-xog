package transform

import (
	"testing"

	"github.com/andreluzz/cas-xog/constant"
	"github.com/andreluzz/cas-xog/model"
	"github.com/beevik/etree"
)

func TestExecuteToReturnObjectFull(t *testing.T) {
	file := model.DriverFile{
		Code: "obj_sistema",
		Type: constant.TypeObject,
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
		Type: constant.TypeObject,
		Elements: []model.Element{
			{
				Code: "test_action",
				Type: constant.ElementTypeAction,
			},
			{
				Code: "aprovador",
				Type: constant.ElementTypeAttribute,
			},
			{
				Code: "status",
				Type: constant.ElementTypeAttribute,
			},
		},
	}

	xog := etree.NewDocument()
	xog.ReadFromFile(packageMockFolder + "object_full_xog_elements.xml")
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
		Type:            constant.TypeObject,
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
		if e.SelectAttrValue("partitionCode", constant.Undefined) != file.TargetPartition {
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

func TestExecuteToReturnObjectSourcePartitionOnly(t *testing.T) {
	file := model.DriverFile{
		Code:            "obj_sistema",
		Type:            constant.TypeObject,
		SourcePartition: "partition20",
	}

	xog := etree.NewDocument()
	xog.ReadFromFile(packageMockFolder + "object_full_xog.xml")
	err := Execute(xog, nil, &file)

	if err != nil {
		t.Fatalf("Error transforming object XOG file. Debug: %s", err.Error())
	}

	if readMockResultAndCompare(xog, "object_result_only_source_partition.xml") == false {
		t.Errorf("Error transforming object XOG file. Invalid result XML.")
	}
}

func TestExecuteToReturnObjectChangeSourcePartitionToTarget(t *testing.T) {
	file := model.DriverFile{
		Code:            "obj_sistema",
		Type:            constant.TypeObject,
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
	if count != 2 {
		t.Errorf("Error transforming object XOG file. Expected 3 got %d target partition (%s) elements", count, file.TargetPartition)
	}

	if readMockResultAndCompare(xog, "object_result_from_source_to_target_partition.xml") == false {
		t.Errorf("Error transforming object XOG file. Invalid result XML.")
	}
}

func TestExecuteToReturnObjectChangePartitionModel(t *testing.T) {
	file := model.DriverFile{
		Code:           "obj_sistema",
		Type:           constant.TypeObject,
		PartitionModel: "NEW_PARTITION_MODEL",
	}

	xog := etree.NewDocument()
	xog.ReadFromFile(packageMockFolder + "object_full_xog.xml")
	err := Execute(xog, nil, &file)

	if err != nil {
		t.Fatalf("Error transforming object XOG file. Debug: %s", err.Error())
	}

	for _, e := range xog.FindElements("*[@partitionModelCode]") {
		value := e.SelectAttrValue("partitionModelCode", constant.Undefined)
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
		Type: constant.TypeObject,
		Elements: []model.Element{
			{
				XPath:  "//customAttribute[@code='analista']",
				Action: constant.ActionRemove,
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

func TestGetWildcardAttributesElements(t *testing.T) {
	file := model.DriverFile{
		Code: "obj_sistema",
		Type: constant.TypeObject,
		Elements: []model.Element{
			{
				Code: "acme_*",
				Type: constant.ElementTypeAttribute,
			},
			{
				Code: "status",
				Type: constant.ElementTypeAttribute,
			},
		},
	}

	xog := etree.NewDocument()
	xog.ReadFromFile(packageMockFolder + "object_full_xog_elements_wc.xml")
	getWildcardAttributesElements(xog, &file)

	totalElements := len(file.Elements)
	if totalElements != 4 {
		t.Errorf("Error getting wildcard attribute element. Should have 4 return %d", totalElements)
	}
}

func TestExecuteToReturnCompleteAsFalse(t *testing.T) {
	file := model.DriverFile{
		Code:         "obj_sistema",
		OnlyElements: true,
		Type:         constant.TypeObject,
		Elements: []model.Element{
			{
				Code: "acme_*",
				Type: constant.ElementTypeAttribute,
			},
			{
				Code: "status",
				Type: constant.ElementTypeAttribute,
			},
		},
	}

	xog := etree.NewDocument()
	xog.ReadFromFile(packageMockFolder + "object_full_xog_elements_wc.xml")
	aux := etree.NewDocument()
	aux.ReadFromFile(packageMockFolder + "object_full_aux.xml")
	err := Execute(xog, aux, &file)

	if err != nil {
		t.Fatalf("Error transforming object XOG file. Debug: %s", err.Error())
	}

	if readMockResultAndCompare(xog, "object_elements_result_wc.xml") == false {
		t.Errorf("Error transforming object XOG file. Invalid result XML.")
	}
}
