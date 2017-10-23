package transform

import (
	"os"
	"github.com/beevik/etree"
	"github.com/andreluzz/cas-xog/common"
)

func ProcessPackage(file common.DriverFile, definitions []common.Definition) error {

	xog := etree.NewDocument()
	err := xog.ReadFromFile(file.PackageFolder + file.Type + "/" + file.Path)
	if err != nil {
		return err
	}

	for _, def := range definitions {
		switch def.Type {
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
		case "processXOGUser":
			if file.Type == common.PROCESS {

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