package transform

import (
	"github.com/andreluzz/cas-xog/constant"
	"github.com/andreluzz/cas-xog/model"
	"github.com/beevik/etree"
)

func specificLookupTransformations(xog *etree.Document, file *model.DriverFile) {
	if file.TargetPartition != "" {
		elems := xog.FindElements("//lookupValue")
		for _, e := range elems {
			currentPartitionCode := e.SelectAttrValue("partitionCode", constant.UNDEFINED)
			if file.SourcePartition == constant.UNDEFINED {
				e.CreateAttr("partitionCode", file.TargetPartition)
			} else if file.SourcePartition == currentPartitionCode || (file.SourcePartition == "NIKU.ROOT" && currentPartitionCode == constant.UNDEFINED) {
				e.CreateAttr("partitionCode", file.TargetPartition)
			}

			currentPartitionModeCode := e.SelectAttrValue("partitionModeCode", constant.UNDEFINED)
			if currentPartitionModeCode == constant.UNDEFINED {
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
