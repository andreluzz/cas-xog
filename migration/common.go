package migration

import (
	"errors"
	"strconv"
	"strings"
	"github.com/tealeg/xlsx"
	"github.com/beevik/etree"
	"github.com/andreluzz/cas-xog/common"
)

func ReadDataFromExcel(file common.DriverFile) (*etree.Document, error) {

	err := errors.New("")
	err = nil

	excelStartRowIndex := 0
	if file.ExcelStartRow != "" {
		excelStartRowIndex, err = strconv.Atoi(file.ExcelStartRow)
		if err != nil {
			return nil, errors.New("migration - tag 'startRow' not a number. Debug:  " + err.Error())
		}
		excelStartRowIndex -= 1
	}

	xog := etree.NewDocument()
	err = xog.ReadFromFile(file.Template)
	if err != nil {
		return nil, errors.New("migration - invalid template file. Debug: " + err.Error())
	}

	instance := "instance"
	if file.InstanceTag != "" {
		instance = file.InstanceTag
	}
	templateInstanceElement := xog.FindElement("//" + instance)
	if templateInstanceElement == nil {
		return nil, errors.New("migration - template invalid no instance element found")
	}

	parent := templateInstanceElement.Parent()
	instanceCopy := templateInstanceElement.Copy()
	templateInstanceElement.Parent().RemoveChild(templateInstanceElement)

	xlFile, err := xlsx.OpenFile(file.ExcelFile)
	if err != nil {
		return nil, errors.New("migration - error opening excel. Debug: " + err.Error())
	}

	for index, row := range xlFile.Sheets[0].Rows {
		if index >= excelStartRowIndex {
			element := instanceCopy.Copy()
			for _, match := range file.MatchExcel {
				var e *etree.Element
				if match.XPath == "" {
					e = element
				} else {
					e = element.FindElement(match.XPath)
				}

				if e == nil {
					return nil, errors.New("migration - invalid xpath element not found in template file")
				}

				value := row.Cells[match.Col-1].String()

				if match.AttributeName != "" {
					e.CreateAttr(match.AttributeName, value)
				} else {
					if match.MultiValued && value != "" {
						separator := ";"
						if match.Separator != "" {
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
			parent.AddChild(element)
		}
	}

	return xog, nil
}

func ExportInstancesToExcel(xog *etree.Document, file common.DriverFile, folder string) error {
	xlsxFile := xlsx.NewFile()
	sheet, _ := xlsxFile.AddSheet("Instances")

	for _, instance := range xog.FindElements("//" + file.InstanceTag) {
		row := sheet.AddRow()

		for _, match := range file.MatchExcel {
			var e *etree.Element
			if match.XPath == "" {
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

	common.ValidateFolder(folder)
	err := xlsxFile.Save(folder + file.ExcelFile)
	if err != nil {
		return errors.New("migration - ExportInstancesToExcel saving excel error. Debug: " + err.Error())
	}

	return nil
}
