package transform

import (
	"github.com/andreluzz/cas-xog/constant"
	"github.com/andreluzz/cas-xog/model"
	"github.com/andreluzz/cas-xog/util"
	"github.com/beevik/etree"
	"strings"
)

func ProcessPackageFile(file model.DriverFile, packageFolder, writeFolder string, definitions []model.Definition) error {

	xog := etree.NewDocument()
	err := xog.ReadFromFile(packageFolder + file.Path)
	if err != nil {
		return err
	}

	for _, def := range definitions {
		if def.Value == def.Default {
			continue
		}
		switch def.Action {
		case constant.PACKAGE_ACTION_CHANGE_PARTITION_MODEL:
			if file.Type == constant.OBJECT {
				e := xog.FindElement("//object[@partitionModelCode]")
				if e != nil {
					e.CreateAttr("partitionModelCode", def.Value)
				}
			}
		case constant.PACKAGE_ACTION_CHANGE_PARTITION:
			if file.Type == constant.OBJECT || file.Type == constant.VIEW {
				changePartition(xog, "", def.Value)
			}
		case constant.PACKAGE_ACTION_REPLACE_STRING:
			if def.Value == "" {
				continue
			}
			replaced := strings.Replace(def.To, "##DEFINITION_VALUE##", def.Value, 1)
			if replaced == def.From {
				continue
			}
			if def.TransformTypes == "" || strings.Contains(def.TransformTypes, file.Type) {
				findAndReplace(xog, []model.FileReplace{{From: def.From, To: replaced}})
			}
		}
	}

	xog.IndentTabs()
	util.ValidateFolder(writeFolder)
	xog.WriteToFile(writeFolder + "/" + file.Path)

	return nil
}
