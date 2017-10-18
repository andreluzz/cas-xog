package transform

import (
	"github.com/andreluzz/cas-xog/common"
	"github.com/beevik/etree"
)

func specificObjectTransformations(xog *etree.Document, file common.DriverFile) error {

	object := xog.FindElement("//objects/object")
	for _, e := range object.FindElements("//object") {
		object.RemoveChild(e)
	}

	if file.SourcePartition != "" {
		removeOtherPartitionsAttributes(xog, file)
	}

	if len(file.Includes) > 0 {
		removeUndefinedIncludes(xog, file.Includes)
		processObjectIncludes(xog, file.Includes)
	}

	if file.TargetPartition != "" {
		changePartition(xog, file)
	}

	if file.PartitionModel != "" {
		xog.FindElement("//object[@code='" + file.Code + "']").CreateAttr("partitionModelCode", file.PartitionModel)
	}

	return nil
}

func removeOtherPartitionsAttributes(xog *etree.Document, file common.DriverFile) {
	partitionElements := xog.FindElements("//[@partitionCode='" + file.SourcePartition + "']")
	var includes []common.Include
	for _, e := range partitionElements {
		var include common.Include
		switch e.Tag{
		case "customAttribute":
			include.Type = "attribute"
		case "link":
			include.Type = "link"
		case "action":
			include.Type = "action"
		}
		include.Code = e.SelectAttrValue("code", "")
		includes = append(includes, include)
	}
	if len(includes) > 0 {
		removeUndefinedIncludes(xog, includes)
		processObjectIncludes(xog, includes)
	}
}

func removeUndefinedIncludes(xog *etree.Document, includes []common.Include) {
	removeActions := true
	removeLinks := true
	removeAttributes := true
	for _, include := range includes {
		if include.Type == "action" {
			removeActions = false
		}
		if include.Type == "link" {
			removeLinks = false
		}
		if include.Type == "attribute" {
			removeAttributes = false
		}
	}

	if removeAttributes {
		removeElementsFromParent(xog, "//customAttribute")
		removeElementsFromParent(xog, "//attributeDefault")
		removeElementsFromParent(xog, "//attributeAutonumbering")
		removeElementsFromParent(xog, "//displayMapping")
		removeElementsFromParent(xog, "//audit/attribute")
	}
	if removeLinks {
		removeElementFromParent(xog, "//links")
	}
	if removeActions {
		removeElementFromParent(xog, "actions")
	}
}

func processObjectIncludes(xog *etree.Document, includes []common.Include) {
	validateAttributesToRemove(xog, includes, "//customAttribute", "code", "attribute")
	validateAttributesToRemove(xog, includes, "//attributeDefault", "code", "attribute")
	validateAttributesToRemove(xog, includes, "//attributeAutonumbering", "code", "attribute")
	validateAttributesToRemove(xog, includes, "//displayMapping", "attributeCode", "attribute")
	validateAttributesToRemove(xog, includes, "//audit/attribute", "code", "attribute")
	validateAttributesToRemove(xog, includes, "//link", "code", "link")
	validateAttributesToRemove(xog, includes, "//action", "code", "action")
}

func validateInclude(includeType, code string, includes []common.Include) bool {
	for _, include := range includes {
		if include.Type == includeType {
			if include.Code == code {
				return false
			}
		}
	}
	return true
}

func validateAttributesToRemove(xog *etree.Document, includes []common.Include, path, attributeKey, includeType string) {
	for _, e := range xog.FindElements(path) {
		code := e.SelectAttrValue(attributeKey, "")
		if validateInclude(includeType, code, includes) {
			e.Parent().RemoveChild(e)
		}
	}
}
