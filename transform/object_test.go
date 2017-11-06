package transform

import (
	"testing"
	"github.com/beevik/etree"
	"github.com/andreluzz/cas-xog/common"
)

func TestExecuteToReturnObjectNoPartitionNoElement(t *testing.T) {
	file := common.DriverFile{
		Code: "obj_sistema",
		Type: common.OBJECT,
	}

	xog := etree.NewDocument()
	xog.ReadFromFile(packageMockFolder + "object_full_xog.xml")
	err := Execute(xog, nil, file)

	if err != nil {
		t.Fatalf("Error transforming object XOG file. Debug: %s", err.Error())
	}

	if readMockResultAndCompare(xog, "object_result.xml") == false {
		t.Errorf("Error transforming object XOG file. Invalid result XML.")
	}
}

func TestExecuteToReturnObjectWithOneAttribute(t *testing.T) {
	file := common.DriverFile{
		Code: "obj_sistema",
		Type: common.OBJECT,
		Elements: []common.Element{
			{
				Type: "attribute",
				Code: "status",
			},
		},
	}

	xog := etree.NewDocument()
	xog.ReadFromFile(packageMockFolder + "object_full_xog.xml")
	err := Execute(xog, nil, file)

	if err != nil {
		t.Fatalf("Error transforming object XOG file. Debug: %s", err.Error())
	}

	if xog.FindElement("//customAttribute[@code='" + file.Elements[0].Code + "']") == nil {
		t.Errorf("Error transforming object XOG file. Attribute: %s not found.", file.Elements[0].Code)
	}

	count := len(xog.FindElements("//customAttribute"))
	if count > 1{
		t.Errorf("Error transforming object XOG file. Expected 1 got %d attributes", count)
	}

	if readMockResultAndCompare(xog, "object_result_only_one_attribute.xml") == false {
		t.Errorf("Error transforming object XOG file. Invalid result XML.")
	}
}

func TestExecuteToReturnObjectWithOneAction(t *testing.T) {
	file := common.DriverFile{
		Code: "obj_sistema",
		Type: common.OBJECT,
		Elements: []common.Element{
			{
				Type: "action",
				Code: "action_cas_xog",
			},
		},
	}

	xog := etree.NewDocument()
	xog.ReadFromFile(packageMockFolder + "object_full_xog.xml")
	err := Execute(xog, nil, file)

	if err != nil {
		t.Fatalf("Error transforming object XOG file. Debug: %s", err.Error())
	}

	if xog.FindElement("//action[@code='" + file.Elements[0].Code + "']") == nil {
		t.Errorf("Error transforming object XOG file. Action: %s not found.", file.Elements[0].Code)
	}

	count := len(xog.FindElements("//customAttribute"))
	if count > 0 {
		t.Errorf("Error transforming object XOG file. Expected 0 got %d customAttributes", count)
	}

	count = len(xog.FindElements("//link"))
	if count > 0 {
		t.Errorf("Error transforming object XOG file. Expected 0 got %d links", count)
	}

	if readMockResultAndCompare(xog, "object_result_only_one_action.xml") == false {
		t.Errorf("Error transforming object XOG file. Invalid result XML.")
	}
}

func TestExecuteToReturnObjectWithOneLink(t *testing.T) {
	file := common.DriverFile{
		Code: "obj_sistema",
		Type: common.OBJECT,
		Elements: []common.Element{
			{
				Type: "link",
				Code: "obj_sistema.lk_teste",
			},
		},
	}

	xog := etree.NewDocument()
	xog.ReadFromFile(packageMockFolder + "object_full_xog.xml")
	err := Execute(xog, nil, file)

	if err != nil {
		t.Fatalf("Error transforming object XOG file. Debug: %s", err.Error())
	}

	if xog.FindElement("//link[@code='" + file.Elements[0].Code + "']") == nil {
		t.Errorf("Error transforming object XOG file. Link: %s not found.", file.Elements[0].Code)
	}

	count := len(xog.FindElements("//customAttribute"))
	if count > 0 {
		t.Errorf("Error transforming object XOG file. Expected 0 got %d customAttributes", count)
	}

	count = len(xog.FindElements("//action"))
	if count > 0 {
		t.Errorf("Error transforming object XOG file. Expected 0 got %d actions", count)
	}

	if readMockResultAndCompare(xog, "object_result_only_one_link.xml") == false {
		t.Errorf("Error transforming object XOG file. Invalid result XML.")
	}
}

func TestExecuteToReturnObjectSourcePartition(t *testing.T) {
	file := common.DriverFile{
		Code: "obj_sistema",
		Type: common.OBJECT,
		SourcePartition: "NIKU.ROOT",
	}

	xog := etree.NewDocument()
	xog.ReadFromFile(packageMockFolder + "object_full_xog.xml")
	err := Execute(xog, nil, file)

	if err != nil {
		t.Fatalf("Error transforming object XOG file. Debug: %s", err.Error())
	}

	count := 0
	for _,e := range xog.FindElements("//*[@partitionCode]") {
		if e.SelectAttrValue("partitionCode", common.UNDEFINED) != file.SourcePartition {
			count ++
		}
	}

	if count > 0 {
		t.Errorf("Error transforming object XOG file. Expected 0 got %d elements from other partitions", count)
	}

	if readMockResultAndCompare(xog, "object_result_source_partition.xml") == false {
		t.Errorf("Error transforming object XOG file. Invalid result XML.")
	}
}

func TestExecuteToReturnObjectTargetPartition(t *testing.T) {
	file := common.DriverFile{
		Code: "obj_sistema",
		Type: common.OBJECT,
		TargetPartition: "NIKU.ROOT",
	}

	xog := etree.NewDocument()
	xog.ReadFromFile(packageMockFolder + "object_full_xog.xml")
	err := Execute(xog, nil, file)

	if err != nil {
		t.Fatalf("Error transforming object XOG file. Debug: %s", err.Error())
	}

	count := 0
	for _,e := range xog.FindElements("//*[@partitionCode]") {
		if e.SelectAttrValue("partitionCode", common.UNDEFINED) != file.TargetPartition {
			count ++
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
	file := common.DriverFile{
		Code: "obj_sistema",
		Type: common.OBJECT,
		SourcePartition: "partition10",
		TargetPartition: "NIKU.ROOT",
	}

	xog := etree.NewDocument()
	xog.ReadFromFile(packageMockFolder + "object_full_xog.xml")
	err := Execute(xog, nil, file)

	if err != nil {
		t.Fatalf("Error transforming object XOG file. Debug: %s", err.Error())
	}

	count := len(xog.FindElements("//object[@code='" + file.Code + "']/*[@partitionCode='" + file.SourcePartition + "']"))
	if count > 0 {
		t.Errorf("Error transforming object XOG file. Expected 0 got %d source partition (%s) elements", count, file.SourcePartition)
	}

	count = len(xog.FindElements("//object[@code='" + file.Code + "']/*[@partitionCode='" + file.TargetPartition + "']"))
	if count != 2 {
		t.Errorf("Error transforming object XOG file. Expected 2 got %d target partition (%s) elements", count, file.TargetPartition)
	}

	if readMockResultAndCompare(xog, "object_result_from_source_to_target_partition.xml") == false {
		t.Errorf("Error transforming object XOG file. Invalid result XML.")
	}
}

func TestExecuteToReturnObjectChangePartitionModel(t *testing.T) {
	file := common.DriverFile{
		Code: "obj_sistema",
		Type: common.OBJECT,
		PartitionModel: "NEW_PARTITION_MODEL",
	}

	xog := etree.NewDocument()
	xog.ReadFromFile(packageMockFolder + "object_full_xog.xml")
	err := Execute(xog, nil, file)

	if err != nil {
		t.Fatalf("Error transforming object XOG file. Debug: %s", err.Error())
	}

	for _,e := range xog.FindElements("*[@partitionModelCode]"){
		value := e.SelectAttrValue("partitionModelCode", common.UNDEFINED)
		if value != file.PartitionModel {
			t.Fatalf("Error transforming object XOG file. Expected %s got %s partition model", file.PartitionModel, value)
		}
	}

	if readMockResultAndCompare(xog, "object_result_change_partition_model.xml") == false {
		t.Errorf("Error transforming object XOG file. Invalid result XML.")
	}
}