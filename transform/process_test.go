package transform

import (
	"github.com/andreluzz/cas-xog/constant"
	"github.com/andreluzz/cas-xog/model"
	"github.com/beevik/etree"
	"strings"
	"testing"
)

func TestExecuteToReturnProcess(t *testing.T) {
	file := model.DriverFile{
		Code: "PRC_0002",
		Type: constant.TypeProcess,
	}

	xog := etree.NewDocument()
	xog.ReadFromFile(packageMockFolder + "process_full_xog.xml")
	err := Execute(xog, nil, &file)

	if err != nil {
		t.Fatalf("Error transforming process XOG file. Debug: %s", err.Error())
	}

	if readMockResultAndCompare(xog, "process_result.xml") == false {
		t.Errorf("Error transforming process XOG file. Invalid result XML.")
	}
}

func TestExecuteToReturnProcessReplace(t *testing.T) {
	file := model.DriverFile{
		Code: "PRC_0002",
		Type: constant.TypeProcess,
		Replace: []model.FileReplace{
			{
				From: "Test cas-xog 002",
				To:   "Test CAS XOG after replace",
			},
		},
	}

	xog := etree.NewDocument()
	xog.ReadFromFile(packageMockFolder + "process_full_xog.xml")
	err := Execute(xog, nil, &file)

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
	file := model.DriverFile{
		Code:            "PRC_0002",
		Type:            constant.TypeProcess,
		CopyPermissions: "PRC_0001",
	}

	aux := etree.NewDocument()
	aux.ReadFromFile(packageMockFolder + "process_full_aux_with_security.xml")

	xog := etree.NewDocument()
	xog.ReadFromFile(packageMockFolder + "process_full_xog.xml")
	err := Execute(xog, aux, &file)

	if err != nil {
		t.Fatalf("Error transforming process XOG file. Debug: %s", err.Error())
	}

	if readMockResultAndCompare(xog, "process_result_security.xml") == false {
		t.Errorf("Error transforming process XOG file. Invalid result XML.")
	}
}

func TestExecuteToReturnErrorProcessCopyingPermissions(t *testing.T) {
	file := model.DriverFile{
		Code:            "PRC_0002",
		Type:            constant.TypeProcess,
		CopyPermissions: "PRC_0001",
	}

	aux := etree.NewDocument()
	aux.ReadFromFile(packageMockFolder + "process_full_aux_with_no_security.xml")

	xog := etree.NewDocument()
	xog.ReadFromFile(packageMockFolder + "process_full_xog.xml")
	err := Execute(xog, aux, &file)

	if err == nil {
		t.Errorf("Error transforming process XOG file. Debug: not validating if aux file has the security tag")
	}
}

func TestExecuteToReturnErrorProcessElementNotFound(t *testing.T) {
	file := model.DriverFile{
		Code:            "PRC_0002",
		Type:            constant.TypeProcess,
		CopyPermissions: "PRC_0001",
	}

	aux := etree.NewDocument()
	aux.ReadFromFile(packageMockFolder + "process_full_aux_with_security.xml")

	xog := etree.NewDocument()
	xog.ReadFromFile(packageMockFolder + "process_full_xog_empty.xml")
	err := Execute(xog, aux, &file)

	if err == nil {
		t.Errorf("Error transforming process XOG file. Debug: not validating if element process exist")
	}
}

func TestIncludeEscapeTextToReturnString(t *testing.T) {
	xog := etree.NewDocument()
	xog.ReadFromFile(packageMockFolder + "process_full_xog_escapetext.xml")

	includeEscapeText(xog)
	XOGString, _ := xog.WriteToString()
	result := etree.NewDocument()
	result.ReadFromString(XOGString)

	if readMockResultAndCompare(result, "process_result_escapetext.xml") == false {
		t.Errorf("Error including escapeText attribute to process XOG file. Invalid result XML.")
	}
}

func TestIncludeEscapeTextWithoutQueryToReturnXML(t *testing.T) {
	xog := etree.NewDocument()
	xog.ReadFromFile(packageMockFolder + "process_full_xog.xml")

	includeEscapeText(xog)
	XOGString, _ := xog.WriteToString()
	result := etree.NewDocument()
	result.ReadFromString(XOGString)

	if readMockResultAndCompare(result, "process_full_xog.xml") == false {
		t.Errorf("Error including escapeText attribute to process XOG file. Invalid result XML.")
	}
}
