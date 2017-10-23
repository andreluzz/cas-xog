package transform

import (
	"os"
	"github.com/beevik/etree"
	"github.com/andreluzz/cas-xog/common"
	"strings"
)

func ProcessPackage(file common.DriverFile, definitions []common.Definition) error {

	xog := etree.NewDocument()
	err := xog.ReadFromFile(file.PackageFolder + file.Type + "/" + file.Path)
	if err != nil {
		return err
	}

	for _, def := range definitions {
		if def.Value == def.Default {
			continue
		}
		switch def.Action {
		case "targetPartitionModel":
			if file.Type == common.OBJECT {
				e := xog.FindElement("//object[@partitionModelCode]")
				if e != nil {
					e.CreateAttr("partitionModelCode", def.Value)
				}
			}
		case "targetPartition":
			if file.Type == common.OBJECT || file.Type == common.VIEW {
				changePartition(xog, "", def.Value)
			}
		case "replaceString":
			if def.Value == "" {
				continue
			}
			replace := strings.Replace(def.Replace, "##DEFINITION_VALUE##", def.Value, 1)
			if replace == def.Match {
				continue
			}
			if def.TransformTypes == "" || strings.Contains(def.TransformTypes, file.Type) {
				findAndReplace(xog, []common.FileReplace{{From: def.Match, To: replace}})
			}
		}
	}

	//check if target folder type dir exists
	_, dirErr := os.Stat(common.FOLDER_WRITE + file.Type)
	if os.IsNotExist(dirErr) {
		_ = os.Mkdir(common.FOLDER_WRITE + file.Type, os.ModePerm)
	}

	xog.IndentTabs()
	err = xog.WriteToFile(common.FOLDER_WRITE + file.Type + "/" + file.Path)
	if err != nil {
		return err
	}

	return nil
}