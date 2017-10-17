package transform

import (
	"github.com/andreluzz/cas-xog/common"
	"github.com/beevik/etree"
)

func specificViewTransformations(xog, aux *etree.Document, file common.DriverFile) error {
	removeElementFromParent(xog, "//lookups")
	removeElementFromParent(xog, "//objects")

	if file.TargetPartition != "" {
		changePartition(xog, file)
	}

	if file.Code != "*" {
		validateCodeAndRemoveElementsFromParent(xog, "//views/property", file.Code)
		validateCodeAndRemoveElementsFromParent(xog, "//views/filter", file.Code)
		validateCodeAndRemoveElementsFromParent(xog, "//views/list", file.Code)
		//auxiliary xog file
		removeElementFromParent(aux, "//lookups")
		removeElementFromParent(aux, "//objects")

		mergeViews(xog, aux, file)
	}

	return nil
}

func mergeViews(xog, aux *etree.Document, file common.DriverFile) {
	sourcePropertySetView := xog.FindElement("//propertySet/update/view[@code='" + file.Code + "']")
	if sourcePropertySetView != nil {
		auxPropertySetView := aux.FindElement("//propertySet/update/view[@code='" + file.Code + "']")
		if auxPropertySetView != nil {
			auxPropertySetView.Parent().InsertChild(auxPropertySetView, sourcePropertySetView)
			auxPropertySetView.Parent().RemoveChild(auxPropertySetView)
		} else {
			propertySetUpdate := aux.FindElement("//propertySet/update")
			if propertySetUpdate != nil {
				propertySetUpdate.InsertChild(propertySetUpdate.FindElement("//nls[1]"), sourcePropertySetView)
			}
		}
	}

	auxPropertySet := aux.FindElement("//propertySet")
	xogPropertySet := xog.FindElement("//propertySet")
	xogPropertySet.Parent().AddChild(auxPropertySet)
	xogPropertySet.Parent().RemoveChild(xogPropertySet)
}

func changePartition(xog *etree.Document, file common.DriverFile) {
	var elems []*etree.Element
	if file.SourcePartition == "" {
		elems = xog.FindElements("//*[@partitionCode]")
	} else {
		elems = xog.FindElements("//*[@partitionCode='" + file.SourcePartition + "']")
	}

	for _, e := range elems {
		e.CreateAttr("partitionCode", file.TargetPartition)
	}
	if file.SourcePartition == "" {
		for _, e := range xog.FindElements("//*[@dataProviderPartitionId]") {
			e.CreateAttr("dataProviderPartitionId", file.TargetPartition)
		}
	} else {
		for _, e := range xog.FindElements("//*[@dataProviderPartitionId='" + file.SourcePartition + "']") {
			e.CreateAttr("dataProviderPartitionId", file.TargetPartition)
		}
	}
}
