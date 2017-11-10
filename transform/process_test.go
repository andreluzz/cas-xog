package transform

import (
	"strings"
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

func TestExecuteToReturnProcessReplace(t *testing.T) {
	file := common.DriverFile{
		Code: "PRC_0002",
		Type: common.PROCESS,
		Replace: []common.FileReplace {
			{
				From: "Test cas-xog 002",
				To: "Test CAS XOG after replace",
			},
		},
	}

	xog := etree.NewDocument()
	xog.ReadFromFile(packageMockFolder + "process_full_xog.xml")
	err := Execute(xog, nil, file)

	if err != nil {
		t.Fatalf("Error transforming process XOG file. Debug: %s", err.Error())
	}

	resultString, _ := xog.WriteToString()
	count := strings.Count(resultString, file.Replace[0].From)
	if count > 0 {
		t.Errorf("Error transforming process XOG file. Expected 0 got %d substring that should have being replaced.", count)
	}

	if readMockResultAndCompare(xog, "process_replace_result.xml") == false {
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

func TestIncludeCDATAToReturnXML(t *testing.T) {
	xog := etree.NewDocument()
	xog.ReadFromFile(packageMockFolder + "process_full_xog_cdata.xml")

	XOGBytes, err := IncludeCDATA(xog)

	if err != nil {
		t.Errorf("Error including CDATA tag to process XOG file. Debug: %s", err.Error())
	}

	result := etree.NewDocument()
	result.ReadFromBytes(XOGBytes)

	if readMockResultAndCompare(result, "process_result_cdata.xml") == false {
		t.Errorf("Error including CDATA tag to process XOG file. Invalid result XML.")
	}
}

func TestIncludeCDATAWithoutQueryToReturnXML(t *testing.T) {
	xog := etree.NewDocument()
	xog.ReadFromFile(packageMockFolder + "process_full_xog.xml")

	XOGBytes, err := IncludeCDATA(xog)

	if err != nil {
		t.Errorf("Error including CDATA tag to process XOG file. Debug: %s", err.Error())
	}

	result := etree.NewDocument()
	result.ReadFromBytes(XOGBytes)

	if readMockResultAndCompare(result, "process_full_xog.xml") == false {
		t.Errorf("Error including CDATA tag to process XOG file. Invalid result XML.")
	}
}