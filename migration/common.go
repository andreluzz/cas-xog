package migration

import (
	"errors"
	"github.com/andreluzz/cas-xog/constant"
	"github.com/andreluzz/cas-xog/model"
	"github.com/andreluzz/cas-xog/util"
	"github.com/beevik/etree"
	"github.com/tealeg/xlsx"
	"strconv"
	"strings"
)

//ReadDataFromExcel used to create xog file from data in excel format. Only accept .xlsx extension
func ReadDataFromExcel(file *model.DriverFile) (string, error) {

	excelStartRowIndex, xog, templateInstanceElement, err := validateReadDataFromExcelDriverAttributes(file)
	if err != nil {
		return constant.Undefined, err
	}

	parent := templateInstanceElement.Parent()
	instanceCopy := templateInstanceElement.Copy()
	templateInstanceElement.Parent().RemoveChild(templateInstanceElement)

	xlFile, err := xlsx.OpenFile(file.ExcelFile)
	if err != nil {
		return constant.Undefined, errors.New("migration - error opening excel. Debug: " + err.Error())
	}

	for index, row := range xlFile.Sheets[0].Rows {
		if index >= excelStartRowIndex {
			element := instanceCopy.Copy()
			for _, match := range file.MatchExcel {
				var e *etree.Element
				if match.XPath == constant.Undefined {
					e = element
				} else {
					e = element.FindElement(match.XPath)
				}

				if e == nil {
					return constant.Undefined, errors.New("migration - invalid xpath element not found in template file")
				}

				value := constant.Undefined
				if match.Col-1 < len(row.Cells) {
					value = row.Cells[match.Col-1].String()
				}

				if match.RemoveIfNull && match.XPath != constant.Undefined && value == constant.Undefined {
					e.Parent().RemoveChild(e)
				} else {
					if match.AttributeName != constant.Undefined {
						e.CreateAttr(match.AttributeName, value)
					} else {
						if match.MultiValued && value != constant.Undefined {
							separator := ";"
							if match.Separator != constant.Undefined {
								separator = match.Separator
							}
							for _, val := range strings.Split(value, separator) {
								v := e.CreateElement("Value")
								v.SetText(strings.TrimSpace(val))
							}
						} else {
							e.SetText(value)
						}
					}
				}
			}
			parent.AddChild(element)
		}
	}
	xog.IndentTabs()
	return xog.WriteToString()
}

func validateReadDataFromExcelDriverAttributes(file *model.DriverFile) (int, *etree.Document, *etree.Element, error) {

	xog := etree.NewDocument()
	err := xog.ReadFromFile(file.Template)
	if err != nil {
		return 0, nil, nil, errors.New("migration - invalid template file. Debug: " + err.Error())
	}

	excelStartRowIndex := 0
	if file.ExcelStartRow != constant.Undefined {
		excelStartRowIndex, err = strconv.Atoi(file.ExcelStartRow)
		if err != nil {
			return 0, nil, nil, errors.New("migration - tag 'startRow' not a number. Debug:  " + err.Error())
		}
		excelStartRowIndex--
	}

	instance := constant.DefaultInstanceTag
	if file.InstanceTag != constant.Undefined {
		instance = file.InstanceTag
	}
	templateInstanceElement := xog.FindElement("//" + instance)
	if templateInstanceElement == nil {
		return 0, nil, nil, errors.New("migration - template invalid no instance element found")
	}

	return excelStartRowIndex, xog, templateInstanceElement, nil
}

//ExportInstancesToExcel used to create excel file with the data from xog file
func ExportInstancesToExcel(xog *etree.Document, file *model.DriverFile, folder string) error {
	xlsxFile := xlsx.NewFile()
	sheet, _ := xlsxFile.AddSheet("Instances")

	for _, instance := range xog.FindElements("//" + file.InstanceTag) {
		row := sheet.AddRow()

		for _, match := range file.MatchExcel {
			var e *etree.Element
			if match.XPath == constant.Undefined {
				e = instance
			} else {
				e = instance.FindElement(match.XPath)
			}

			cell := row.AddCell()

			if e == nil {
				cell.Value = ""
				continue
			}

			value := ""
			if match.AttributeName != "" {
				value = e.SelectAttrValue(match.AttributeName, "")
			} else {
				if match.MultiValued {
					separator := ";"
					if match.Separator != "" {
						separator = match.Separator
					}
					for _, v := range e.FindElements("//Value") {
						value += v.Text() + separator
					}
					value = value[:len(value)-1]
				} else {
					value = e.Text()
				}
			}
			cell.Value = value
		}
	}

	util.ValidateFolder(folder + util.GetPathFolder(file.ExcelFile))
	err := xlsxFile.Save(folder + file.ExcelFile)
	if err != nil {
		return errors.New("migration - ExportInstancesToExcel saving excel error. Debug: " + err.Error())
	}

	return nil
}
