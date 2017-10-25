package migration

import (
	"testing"
	"github.com/tealeg/xlsx"
	"github.com/beevik/etree"
	"github.com/andreluzz/cas-xog/common"
)

var packageMockFolder string

func init() {
	packageMockFolder = "../" + common.FOLDER_MOCK + "migration/"
}

func TestReadDataFromExcelToReturnErrorExcelStartRow(t *testing.T) {
	file := common.DriverFile{
		ExcelStartRow: "A",
	}
	_, err := ReadDataFromExcel(file)
	if err == nil {
		t.Errorf("Error reading data from excel to XOG file. Debug: not validating if ExcelStartRow is a number")
	}
}

func TestReadDataFromExcelToReturnErrorTemplateExists(t *testing.T) {
	file := common.DriverFile{}
	_, err := ReadDataFromExcel(file)
	if err == nil {
		t.Errorf("Error reading data from excel to XOG file. Debug: not validating if Template exists")
	}
}

func TestReadDataFromExcelToReturnErrorInstanceElementExists(t *testing.T) {
	file := common.DriverFile{
		Template: packageMockFolder + "template.xml",
		InstanceTag: "WrongInstance",
	}
	_, err := ReadDataFromExcel(file)
	if err == nil {
		t.Errorf("Error reading data from excel to XOG file. Debug: not validating if Instance element exists")
	}
}

func TestReadDataFromExcelToReturnErrorExcelFileExists(t *testing.T) {
	file := common.DriverFile{
		Template: packageMockFolder + "template.xml",
	}
	_, err := ReadDataFromExcel(file)
	if err == nil {
		t.Errorf("Error reading data from excel to XOG file. Debug: not validating if ExcelFile exists")
	}
}

func TestReadDataFromExcelToReturnErrorMatchElementExists(t *testing.T) {
	file := common.DriverFile{
		Template: packageMockFolder + "template.xml",
		ExcelFile:  packageMockFolder + "data.xlsx",
		InstanceTag: "instance",
		ExcelStartRow: "1",
		MatchExcel: []common.MatchExcel{
			{
				Col: 1,
				Tag: "WrongTagName",
				AttributeName: "name",
				AttributeValue: "code",
			},
		},
	}
	_, err := ReadDataFromExcel(file)
	if err == nil {
		t.Errorf("Error reading data from excel to XOG file. Debug: not validating if Match Tag element exists")
	}
}

func TestReadDataFromExcelToReturnXMLResult(t *testing.T) {
	file := common.DriverFile{
		Template:  packageMockFolder + "template.xml",
		ExcelFile:  packageMockFolder + "data.xlsx",
		InstanceTag: "instance",
		ExcelStartRow: "1",
		MatchExcel: []common.MatchExcel{
			{
				Col: 1,
				Tag: "instance",
				AttributeName: "instanceCode",
				IsAttribute: true,
			},
			{
				Col: 1,
				Tag: "ColumnValue",
				AttributeName: "name",
				AttributeValue: "code",
			},
			{
				Col: 2,
				Tag: "ColumnValue",
				AttributeName: "name",
				AttributeValue: "name",
			},
			{
				Col: 3,
				Tag: "ColumnValue",
				AttributeName: "name",
				AttributeValue: "status_novo",
			},
			{
				Col: 4,
				Tag: "ColumnValue",
				AttributeName: "name",
				AttributeValue: "multivalue_status",
				MultiValued: true,
				Separator: ";",
			},
			{
				Col: 5,
				Tag: "ColumnValue",
				AttributeName: "name",
				AttributeValue: "analista",
			},
		},
	}

	result, err := ReadDataFromExcel(file)

	if err != nil {
		t.Fatalf("Error reading data from excel to XOG file. Debug: %s", err.Error())
	}

	expectedResult := etree.NewDocument()
	expectedResult.ReadFromFile( packageMockFolder + "result.xml")
	expectedResult.IndentTabs()
	expectedResultString, _ := expectedResult.WriteToString()

	result.IndentTabs()
	resultString, _ := result.WriteToString()
	if resultString != expectedResultString {
		t.Errorf("Error reading data from excel to XOG file. Debug: incorrect result")
	}
}

func TestExportInstancesToExcelToReturnErrorExcelPath(t *testing.T) {
	file := common.DriverFile{
		Type:          common.CUSTOM_OBJECT_INSTANCE,
		ExportToExcel: true,
		InstanceTag:   "instance",
		MatchExcel: []common.MatchExcel{
			{
				AttributeName: "instanceCode",
			},
			{
				XPath: "//ColumnValue[@name='name']",
			},
			{
				XPath: "//ColumnValue[@name='status_novo']",
			},
			{
				XPath:       "//ColumnValue[@name='multivalue_status']",
				MultiValued: true,
			},
			{
				XPath: "//ColumnValue[@name='analista']",
			},
		},
	}

	result := etree.NewDocument()
	result.ReadFromFile( packageMockFolder + "result.xml")
	folder := ""
	err := ExportInstancesToExcel(result, file, folder)
	if err == nil {
		t.Fatalf("Error exporting instances to excel file. Debug: not validating if ExcelFile and Folder exists")
	}
}

func TestExportInstancesToExcelToReturnErrorXPath(t *testing.T) {
	file := common.DriverFile{
		Type:          common.CUSTOM_OBJECT_INSTANCE,
		ExportToExcel: true,
		InstanceTag:   "instance",
		MatchExcel: []common.MatchExcel{
			{
				XPath: "//WrongXPath",
			},
		},
	}

	result := etree.NewDocument()
	result.ReadFromFile( packageMockFolder + "result.xml")
	folder := ""
	err := ExportInstancesToExcel(result, file, folder)
	if err == nil {
		t.Fatalf("Error exporting instances to excel file. Debug: not validating if XPath element")
	}
}

func TestExportInstancesToExcelToReturnExcelFile(t *testing.T) {
	file := common.DriverFile{
		Type: common.CUSTOM_OBJECT_INSTANCE,
		ExportToExcel: true,
		ExcelFile: "instances.xlsx",
		InstanceTag: "instance",
		MatchExcel: []common.MatchExcel{
			{
				AttributeName: "instanceCode",
			},
			{
				XPath: "//ColumnValue[@name='name']",
			},
			{
				XPath: "//ColumnValue[@name='status_novo']",
			},
			{
				XPath: "//ColumnValue[@name='multivalue_status']",
				MultiValued: true,
				Separator: ";",
			},
			{
				XPath: "//ColumnValue[@name='analista']",
			},
		},
	}

	result := etree.NewDocument()
	result.ReadFromFile( packageMockFolder + "result.xml")
	folder := "../" + common.FOLDER_READ + file.Type + "/"
	err := ExportInstancesToExcel(result, file, folder)
	if err != nil {
		t.Fatalf("Error exporting instances to excel file. Debug: %s", err.Error())
	}

	f1, _:= xlsx.OpenFile( packageMockFolder + "data.xlsx")
	f2, _:= xlsx.OpenFile(folder + file.ExcelFile)

	f1Array, _ := f1.ToSlice()
	f2Array, _ := f2.ToSlice()

	assert := true
	for i, a := range f1Array[0][0] {
		if a != f2Array[0][0][i] {
			assert = false
			break
		}
	}

	if !assert {
		t.Errorf("Error exporting instances to excel file. Debug: incorrect result")
	}
}