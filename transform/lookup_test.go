package transform

import (
	"github.com/andreluzz/cas-xog/constant"
	"github.com/andreluzz/cas-xog/model"
	"github.com/beevik/etree"
	"testing"
)

func TestExecuteToReturnStaticLookupTransformed(t *testing.T) {
	file := model.DriverFile{
		Code: "LOOKUP_CAS_XOG",
		Type: constant.LOOKUP,
	}

	xog := etree.NewDocument()
	xog.ReadFromFile(packageMockFolder + "lookup_static_full_xog.xml")
	err := Execute(xog, nil, &file)

	if err != nil {
		t.Fatalf("Error transforming static lookup XOG file. Debug: %s", err.Error())
	}

	if readMockResultAndCompare(xog, "lookup_static_result.xml") == false {
		t.Errorf("Error transforming static lookup XOG file. Invalid result XML.")
	}
}

func TestExecuteToReturnStaticLookupTargetPartition(t *testing.T) {
	file := model.DriverFile{
		Code:            "LOOKUP_CAS_XOG",
		Type:            constant.LOOKUP,
		TargetPartition: "NIKU.ROOT",
		Path:            "testTarget.xml",
	}

	xog := etree.NewDocument()
	xog.ReadFromFile(packageMockFolder + "lookup_static_full_xog.xml")
	err := Execute(xog, nil, &file)

	if err != nil {
		t.Fatalf("Error transforming static lookup XOG file. Debug: %s", err.Error())
	}

	if readMockResultAndCompare(xog, "lookup_static_partitions_result.xml") == false {
		t.Errorf("Error transforming static lookup XOG file. Invalid result XML.")
	}
}

func TestExecuteToReturnStaticLookupSourceAndTargetPartition(t *testing.T) {
	file := model.DriverFile{
		Code:            "LOOKUP_CAS_XOG",
		Type:            constant.LOOKUP,
		SourcePartition: "NIKU.ROOT",
		TargetPartition: "partition10",
		Path:            "testSourceAndTarget.xml",
	}

	xog := etree.NewDocument()
	xog.ReadFromFile(packageMockFolder + "lookup_static_full_xog.xml")
	err := Execute(xog, nil, &file)

	if err != nil {
		t.Fatalf("Error transforming static lookup XOG file. Debug: %s", err.Error())
	}

	if readMockResultAndCompare(xog, "lookup_static_source_target_partition_result.xml") == false {
		t.Errorf("Error transforming static lookup XOG file. Invalid result XML.")
	}
}

func TestExecuteToReturnDynamicLookupPartitionsTransformed(t *testing.T) {
	file := model.DriverFile{
		Code: "LOOKUP_CAS_XOG",
		Type: constant.LOOKUP,
	}

	xog := etree.NewDocument()
	xog.ReadFromFile(packageMockFolder + "lookup_dynamic_full_xog.xml")
	err := Execute(xog, nil, &file)

	if err != nil {
		t.Fatalf("Error transforming static lookup XOG file. Debug: %s", err.Error())
	}

	if readMockResultAndCompare(xog, "lookup_dynamic_result.xml") == false {
		t.Errorf("Error transforming static lookup XOG file. Invalid result XML.")
	}
}

func TestExecuteToReturnDynamicLookupReplacedNSQL(t *testing.T) {
	file := model.DriverFile{
		Code: "LOOKUP_CAS_XOG",
		Type: constant.LOOKUP,
		NSQL: "select * from inv_investments",
	}

	xog := etree.NewDocument()
	xog.ReadFromFile(packageMockFolder + "lookup_dynamic_full_xog.xml")
	err := Execute(xog, nil, &file)

	if err != nil {
		t.Fatalf("Error transforming static lookup XOG file. Debug: %s", err.Error())
	}

	if readMockResultAndCompare(xog, "lookup_dynamic_replace_nsql_result.xml") == false {
		t.Errorf("Error transforming static lookup XOG file. Invalid result XML.")
	}
}

func TestExecuteToReturnDynamicLookupOnlyStructure(t *testing.T) {
	file := model.DriverFile{
		Code: "LOOKUP_CAS_XOG",
		Type: constant.LOOKUP,
		OnlyStructure: true,
	}

	model.LoadXMLReadList("../xogRead.xml")

	xog := etree.NewDocument()
	xog.ReadFromFile(packageMockFolder + "lookup_dynamic_full_xog.xml")
	err := Execute(xog, nil, &file)

	if err != nil {
		t.Fatalf("Error transforming dynamic lookup XOG file. Debug: %s", err.Error())
	}

	if readMockResultAndCompare(xog, "lookup_dynamic_only_structure_result.xml") == false {
		t.Errorf("Error transforming static lookup XOG file. Invalid result XML.")
	}
}
