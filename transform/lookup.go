package transform

import (
	"github.com/andreluzz/cas-xog/common"
	"github.com/beevik/etree"
)

func specificLookupTransformations(xog *etree.Document, file common.DriverFile) {
	if file.TargetPartition != "" {
		elems := xog.FindElements("//lookupValue")
		for _, e := range elems {
			currentPartitionCode := e.SelectAttrValue("partitionCode", common.UNDEFINED)
			if file.SourcePartition == common.UNDEFINED {
				e.CreateAttr("partitionCode", file.TargetPartition)
			} else if file.SourcePartition == currentPartitionCode || (file.SourcePartition == "NIKU.ROOT" && currentPartitionCode == common.UNDEFINED) {
				e.CreateAttr("partitionCode", file.TargetPartition)
			}

			currentPartitionModeCode := e.SelectAttrValue("partitionModeCode", common.UNDEFINED)
			if currentPartitionModeCode == common.UNDEFINED {
				e.CreateAttr("partitionModeCode", "PARTITION_AND_ANSTRS_DESDNTS")
			}
		}
	}
}
