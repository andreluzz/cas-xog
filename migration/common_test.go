package migration

import (
	"github.com/andreluzz/cas-xog/constant"
	"github.com/andreluzz/cas-xog/model"
	"github.com/beevik/etree"
	"github.com/tealeg/xlsx"
	"testing"
)

var packageMockFolder string

func init() {
	packageMockFolder = "../" + constant.FolderMock + "migration/"
}

func TestReadDataFromExcelToReturnErrorExcelStartRow(t *testing.T) {
	file := model.DriverFile{
		ExcelStartRow: "A",
	}
	_, err := ReadDataFromExcel(&file)
	if err == nil {
		t.Errorf("Error reading data from excel to XOG file. Debug: not validating if ExcelStartRow is a number")
	}
}

func TestReadDataFromExcelToReturnErrorTemplateExists(t *testing.T) {
	file := model.DriverFile{}
	_, err := ReadDataFromExcel(&file)
	if err == nil {
		t.Errorf("Error reading data from excel to XOG file. Debug: not validating if Template exists")
	}
}

func TestReadDataFromExcelToReturnErrorInstanceElementExists(t *testing.T) {
	file := model.DriverFile{
		Template:    packageMockFolder + "template.xml",
		InstanceTag: "WrongInstance",
	}
	_, err := ReadDataFromExcel(&file)
	if err == nil {
		t.Errorf("Error reading data from excel to XOG file. Debug: not validating if Instance element exists")
	}
}

func TestReadDataFromExcelToReturnErrorExcelFileExists(t *testing.T) {
	file := model.DriverFile{
		Template: packageMockFolder + "template.xml",
	}
	_, err := ReadDataFromExcel(&file)
	if err == nil {
		t.Errorf("Error reading data from excel to XOG file. Debug: not validating if ExcelFile exists")
	}
}

func TestReadDataFromExcelToReturnErrorMatchElementExists(t *testing.T) {
	file := model.DriverFile{
		Template:      packageMockFolder + "template.xml",
		ExcelFile:     packageMockFolder + "data.xlsx",
		InstanceTag:   "instance",
		ExcelStartRow: "1",
		MatchExcel: []model.MatchExcel{
			{
				Col:           1,
				XPath:         "invalid_xpath",
				AttributeName: "name",
			},
		},
	}
	_, err := ReadDataFromExcel(&file)
	if err == nil {
		t.Errorf("Error reading data from excel to XOG file. Debug: not validating if Match Tag element exists")
	}
}

func TestReadDataFromExcelToReturnXMLResult(t *testing.T) {
	file := model.DriverFile{
		Template:      packageMockFolder + "template.xml",
		ExcelFile:     packageMockFolder + "data.xlsx",
		InstanceTag:   "instance",
		ExcelStartRow: "1",
		MatchExcel: []model.MatchExcel{
			{
				Col:           1,
				AttributeName: "instanceCode",
			},
			{
				Col:   1,
				XPath: "//ColumnValue[@name='code']",
			},
			{
				Col:   2,
				XPath: "//ColumnValue[@name='name']",
			},
			{
				Col:   3,
				XPath: "//ColumnValue[@name='status_novo']",
			},
			{
				Col:         4,
				XPath:       "//ColumnValue[@name='multivalue_status']",
				MultiValued: true,
				Separator:   ";",
			},
			{
				Col:   5,
				XPath: "//ColumnValue[@name='analista']",
			},
		},
	}

	result, err := ReadDataFromExcel(&file)

	if err != nil {
		t.Fatalf("Error reading data from excel to XOG file. Debug: %s", err.Error())
	}

	expectedResult := etree.NewDocument()
	expectedResult.ReadFromFile(packageMockFolder + "result.xml")
	expectedResult.IndentTabs()
	expectedResultString, _ := expectedResult.WriteToString()

	if result != expectedResultString {
		t.Errorf("Error reading data from excel to XOG file. Debug: incorrect result")
	}
}

func TestExportInstancesToExcelToReturnErrorExcelPath(t *testing.T) {
	file := model.DriverFile{
		Type:          constant.CustomObjectInstances,
		ExportToExcel: true,
		InstanceTag:   "instance",
		MatchExcel: []model.MatchExcel{
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
	result.ReadFromFile(packageMockFolder + "result.xml")
	folder := ""
	err := ExportInstancesToExcel(result, &file, folder)
	if err == nil {
		t.Fatalf("Error exporting instances to excel file. Debug: not validating if ExcelFile and Folder exists")
	}
}

func TestExportInstancesToExcelToReturnErrorXPath(t *testing.T) {
	file := model.DriverFile{
		Type:          constant.CustomObjectInstances,
		ExportToExcel: true,
		InstanceTag:   "instance",
		MatchExcel: []model.MatchExcel{
			{
				XPath: "//WrongXPath",
			},
		},
	}

	result := etree.NewDocument()
	result.ReadFromFile(packageMockFolder + "result.xml")
	folder := ""
	err := ExportInstancesToExcel(result, &file, folder)
	if err == nil {
		t.Fatalf("Error exporting instances to excel file. Debug: not validating if XPath element")
	}
}

func TestExportInstancesToExcelToReturnExcelFile(t *testing.T) {
	file := model.DriverFile{
		Type:          constant.CustomObjectInstances,
		ExportToExcel: true,
		ExcelFile:     "instances.xlsx",
		InstanceTag:   "instance",
		MatchExcel: []model.MatchExcel{
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
				Separator:   ";",
			},
			{
				XPath: "//ColumnValue[@name='analista']",
			},
		},
	}

	result := etree.NewDocument()
	result.ReadFromFile(packageMockFolder + "result.xml")
	folder := "../" + constant.FolderRead + file.Type + "/"
	err := ExportInstancesToExcel(result, &file, folder)
	if err != nil {
		t.Fatalf("Error exporting instances to excel file. Debug: %s", err.Error())
	}

	f1, _ := xlsx.OpenFile(packageMockFolder + "data.xlsx")
	f2, _ := xlsx.OpenFile(folder + file.ExcelFile)

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
