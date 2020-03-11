package transform

import (
	"errors"
	"strconv"
	"strings"

	"github.com/andreluzz/cas-xog/constant"
	"github.com/andreluzz/cas-xog/model"
	"github.com/andreluzz/cas-xog/util"
	"github.com/beevik/etree"
	"github.com/tealeg/xlsx"
)

func specificDepartmentTransformations(xog *etree.Document, file *model.DriverFile) error {

	if file.ExcelFile == constant.Undefined {
		return nil
	}

	xlFile, err := xlsx.OpenFile(util.ReplacePathSeparatorByOS(file.ExcelFile))
	if err != nil {
		return errors.New("OBS excel import - error opening excel. Debug: " + err.Error())
	}

	excelStartRowIndex := 0
	if file.ExcelStartRow != constant.Undefined {
		excelStartRowIndex, err = strconv.Atoi(file.ExcelStartRow)
		if err != nil {
			return errors.New("OBS excel import - tag 'startRow' not a number. Debug:  " + err.Error())
		}
		excelStartRowIndex--
	}

	nodes := []Node{}

	for rowIndex, row := range xlFile.Sheets[0].Rows {
		if rowIndex >= excelStartRowIndex {
			node := Node{}
			for index, cell := range row.Cells {
				if cell.Value == constant.Undefined || index == len(row.Cells)-1 {
					i := strings.LastIndex(node.xpath, "/Department")
					if i < 0 {
						break
					}
					node.xpath = node.xpath[0:i]
					if cell.Value != constant.Undefined {
						node.name = cell.Value
					}
					break
				}
				if index%2 == 0 {
					node.xpath += "/Department[@department_code='" + cell.Value + "']"
					node.id = cell.Value
				} else {
					node.name = cell.Value
				}
			}
			nodes = append(nodes, node)
		}
	}

	removeElementsFromParent(xog, "//Department")

	departments := xog.FindElement("//Departments")

	errs := []string{}

	for _, n := range nodes {
		departmentElement := etree.NewElement("Department")
		departmentElement.CreateAttr("department_code", n.id)
		departmentElement.CreateAttr("entity", file.Entity)
		departmentElement.CreateAttr("short_description", n.name)
		descriptionElement := etree.NewElement("Description")
		descriptionElement.CreateText(n.name)
		departmentElement.AddChild(descriptionElement)
		if n.xpath == "" {
			departments.AddChild(departmentElement)
		} else {
			parent := departments.FindElement("/" + n.xpath)
			if parent != nil {
				parent.AddChild(departmentElement)
			} else {
				errString := "[" + n.id + " - " + n.name + "] \n " + n.xpath + "\n"
				errs = append(errs, errString)
			}
		}
	}

	if len(errs) > 0 {
		return errors.New("ids:  " + strings.Join(errs, ","))
	}

	return nil
}
