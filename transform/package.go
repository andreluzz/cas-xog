package transform

import (
	"strings"
	"github.com/beevik/etree"
	"github.com/andreluzz/cas-xog/common"
)

func ProcessPackageFile(file common.DriverFile, packageFolder string, definitions []common.Definition) error {

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
			replaced := strings.Replace(def.To, "##DEFINITION_VALUE##", def.Value, 1)
			if replaced == def.From {
				continue
			}
			if def.TransformTypes == "" || strings.Contains(def.TransformTypes, file.Type) {
				findAndReplace(xog, []common.FileReplace{{From: def.From, To: replaced}})
			}
		}
	}

	xog.IndentTabs()
	folder := common.FOLDER_WRITE + file.Type
	common.ValidateFolder(folder)
	xog.WriteToFile(folder + "/" + file.Path)

	return nil
}