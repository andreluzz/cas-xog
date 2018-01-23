package transform

import (
	"github.com/andreluzz/cas-xog/constant"
	"github.com/andreluzz/cas-xog/model"
	"github.com/beevik/etree"
)

func specificLookupTransformations(xog *etree.Document, file *model.DriverFile) {
	if file.OnlyStructure {
		xog.SetRoot(file.GetDummyLookup())
		xog.FindElement("//dynamicLookup").CreateAttr("code", file.Code)
		return
	}

	if file.TargetPartition != "" {
		elems := xog.FindElements("//lookupValue")
		for _, e := range elems {
			currentPartitionCode := e.SelectAttrValue("partitionCode", constant.Undefined)
			if file.SourcePartition == constant.Undefined {
				e.CreateAttr("partitionCode", file.TargetPartition)
			} else if file.SourcePartition == currentPartitionCode || (file.SourcePartition == "NIKU.ROOT" && currentPartitionCode == constant.Undefined) {
				e.CreateAttr("partitionCode", file.TargetPartition)
			}

			currentPartitionModeCode := e.SelectAttrValue("partitionModeCode", constant.Undefined)
			if currentPartitionModeCode == constant.Undefined {
				e.CreateAttr("partitionModeCode", "PARTITION_AND_ANSTRS_DESDNTS")
			}
		}
	}

	if file.NSQL != "" {
		nsqlElement := xog.FindElement("//nsql")
		if nsqlElement != nil {
			nsqlElement.SetText(file.NSQL)
		}
	}
}
