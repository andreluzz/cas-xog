package transform

import (
	"github.com/andreluzz/cas-xog/common"
	"github.com/beevik/etree"
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
		case "targetPartition":
		case "xogUser":
		}
	}

	xog.IndentTabs()
	err = xog.WriteToFile(common.FOLDER_WRITE + file.Type + "/" + file.Path)
	if err != nil {
		return err
	}

	return nil
}