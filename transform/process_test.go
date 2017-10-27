package transform

import (
	"testing"
	"github.com/beevik/etree"
	"github.com/andreluzz/cas-xog/common"
)

func TestExecuteToReturnProcess(t *testing.T) {
	file := common.DriverFile{
		Code: "PRC_0002",
		Type: common.PROCESS,
	}

	xog := etree.NewDocument()
	xog.ReadFromFile(packageMockFolder + "process_full_xog.xml")
	err := Execute(xog, nil, file)

	if err != nil {
		t.Fatalf("Error transforming process XOG file. Debug: %s", err.Error())
	}

	if readMockResultAndCompare(xog, "process_result.xml") == false {
		t.Errorf("Error transforming process XOG file. Invalid result XML.")
	}
}

func TestExecuteToReturnProcessCopyingPermissions(t *testing.T) {
	file := common.DriverFile{
		Code: "PRC_0002",
		Type: common.PROCESS,
		CopyPermissions: "PRC_0001",
	}

	aux := etree.NewDocument()
	aux.ReadFromFile(packageMockFolder + "process_full_aux_with_security.xml")

	xog := etree.NewDocument()
	xog.ReadFromFile(packageMockFolder + "process_full_xog.xml")
	err := Execute(xog, aux, file)

	if err != nil {
		t.Fatalf("Error transforming process XOG file. Debug: %s", err.Error())
	}

	if readMockResultAndCompare(xog, "process_result_security.xml") == false {
		t.Errorf("Error transforming process XOG file. Invalid result XML.")
	}
}

func TestExecuteToReturnErrorProcessCopyingPermissions(t *testing.T) {
	file := common.DriverFile{
		Code: "PRC_0002",
		Type: common.PROCESS,
		CopyPermissions: "PRC_0001",
	}

	aux := etree.NewDocument()
	aux.ReadFromFile(packageMockFolder + "process_full_aux_with_no_security.xml")

	xog := etree.NewDocument()
	xog.ReadFromFile(packageMockFolder + "process_full_xog.xml")
	err := Execute(xog, aux, file)

	if err == nil {
		t.Errorf("Error transforming process XOG file. Debug: not validating if aux file has the security tag")
	}
}

func TestExecuteToReturnErrorProcessElementNotFound(t *testing.T) {
	file := common.DriverFile{
		Code: "PRC_0002",
		Type: common.PROCESS,
		CopyPermissions: "PRC_0001",
	}

	aux := etree.NewDocument()
	aux.ReadFromFile(packageMockFolder + "process_full_aux_with_security.xml")

	xog := etree.NewDocument()
	xog.ReadFromFile(packageMockFolder + "process_full_xog_empty.xml")
	err := Execute(xog, aux, file)

	if err == nil {
		t.Errorf("Error transforming process XOG file. Debug: not validating if element process exist")
	}
}