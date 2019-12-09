package transform

import (
	"errors"
	"strconv"
	"strings"
	"github.com/beevik/etree"
	"github.com/tealeg/xlsx"	
	
	"github.com/andreluzz/cas-xog/util"
	"github.com/andreluzz/cas-xog/model"
	"github.com/andreluzz/cas-xog/constant"
)

// Node represents a obs unit
type Node struct {
	name string
	id string
	xpath string
}

func specificObsTransformations(xog *etree.Document, file *model.DriverFile) error {
	if  file.ExcelFile == constant.Undefined {
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
					i := strings.LastIndex(node.xpath, "/unit")
					node.xpath = node.xpath[0:i]
					break
				}
				if index%2 == 0 {
					node.xpath += "/unit[@code='" + cell.Value + "']"
					node.id = cell.Value
				} else {
					node.name = cell.Value
				}
			}
			nodes = append(nodes, node)
		}
	}

	removeElementsFromParent(xog, "//unit")
	removeElementsFromParent(xog, "//associatedObject")
	
	obs := xog.FindElement("//obs")
	
	for _, n := range nodes {
		unitElement := etree.NewElement("unit")
		unitElement.CreateAttr("code", n.id)
		unitElement.CreateAttr("name", n.name)
		if n.xpath == "" {
			obs.AddChild(unitElement)
		} else {
			obs.FindElement("/"+n.xpath).AddChild(unitElement)
		}
	}

	return nil
}