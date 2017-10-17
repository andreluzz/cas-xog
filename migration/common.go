package migration

import (
	"errors"
	"github.com/andreluzz/cas-xog/common"
	"github.com/beevik/etree"
	"github.com/tealeg/xlsx"
	"strconv"
	"strings"
)

func ReadDataFromExcel(file common.DriverFile) (*etree.Document, error) {
	err := errors.New("")
	err = nil

	excelStartRowIndex := 0
	if file.ExcelStartRow != "" {
		excelStartRowIndex, err = strconv.Atoi(file.ExcelStartRow)
		if err != nil {
			return nil, errors.New("[migration error] tag 'startRow' not a number. Error:  " + err.Error())
		}
		excelStartRowIndex -= 1
	}

	xog := etree.NewDocument()
	err = xog.ReadFromFile(file.Template)
	if err != nil {
		return nil, errors.New("[migration error] invalid template file. Error: " + err.Error())
	}

	instance := "instance"
	if file.InstanceTag != "" {
		instance = file.InstanceTag
	}
	templateInstanceElement := xog.FindElement("//" + instance)
	if templateInstanceElement == nil {
		return nil, errors.New("[migration error] no instance element found")
	}

	parent := templateInstanceElement.Parent()
	instanceCopy := templateInstanceElement.Copy()
	templateInstanceElement.Parent().RemoveChild(templateInstanceElement)

	xlFile, err := xlsx.OpenFile(file.ExcelFile)
	if err != nil {
		return nil, errors.New("[migration error] " + err.Error())
	}

	for index, row := range xlFile.Sheets[0].Rows {
		if index >= excelStartRowIndex {
			element := instanceCopy.Copy()
			for _, match := range file.MatchExcel {
				var e *etree.Element
				if match.AttributeValue == "" {
					e = element.FindElement("//" + match.Tag + "[@" + match.AttributeName + "']")
				} else {
					e = element.FindElement("//" + match.Tag + "[@" + match.AttributeName + "='" + match.AttributeValue + "']")
				}
				if element.Tag == match.Tag {
					e = element
				}

				if e == nil {
					return nil, errors.New("[migration error] invalid attribute name(" + match.AttributeName + ") or tag(" + match.Tag + ")")
				}

				value := row.Cells[match.Col-1].String()

				if match.IsAttribute {
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

func ExportInstancesToExcel(xog *etree.Document, file common.DriverFile) error {
	xlsxFile := xlsx.NewFile()
	sheet, err := xlsxFile.AddSheet("Instances")
	if err != nil {
		return errors.New("[migration error] ExportInstancesToExcel: " + err.Error())
	}

	for _, instance := range xog.FindElements("//" + file.InstanceTag) {
		row := sheet.AddRow()

		for _, match := range file.MatchExcel {
			var e *etree.Element
			if match.AttributeValue == "" {
				e = instance.FindElement("//" + match.Tag + "[@" + match.AttributeName + "']")
			} else {
				e = instance.FindElement("//" + match.Tag + "[@" + match.AttributeName + "='" + match.AttributeValue + "']")
			}
			if instance.Tag == match.Tag {
				e = instance
			}

			cell := row.AddCell()

			if e == nil {
				cell.Value = ""
				continue
			}

			value := ""
			if match.IsAttribute {
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

	err = xlsxFile.Save(common.FOLDER_READ + file.Type + "/" + file.ObjCode + "_instances.xlsx")
	if err != nil {
		return errors.New("[migration error] ExportInstancesToExcel: " + err.Error())
	}

	return nil
}
