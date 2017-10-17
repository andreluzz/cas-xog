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

	if len(file.Includes) > 0 {
		removeUndefinedIncludes(xog, file)
		processObjectIncludes(xog, file)
	}

	return nil
}

func removeUndefinedIncludes(xog *etree.Document, file common.DriverFile) {
	removeActions := true
	removeLinks := true
	removeAttributes := true
	for _, include := range file.Includes {
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

func processObjectIncludes(xog *etree.Document, file common.DriverFile) {
	validateAttributesToRemove(xog, file, "//customAttribute", "code", "attribute")
	validateAttributesToRemove(xog, file, "//attributeDefault", "code", "attribute")
	validateAttributesToRemove(xog, file, "//attributeAutonumbering", "code", "attribute")
	validateAttributesToRemove(xog, file, "//displayMapping", "attributeCode", "attribute")
	validateAttributesToRemove(xog, file, "//audit/attribute", "code", "attribute")
	validateAttributesToRemove(xog, file, "//link", "code", "link")
	validateAttributesToRemove(xog, file, "//action", "code", "action")
}

func validateInclude(includeType, code string, file common.DriverFile) bool {
	for _, include := range file.Includes {
		if include.Type == includeType {
			if include.Code == code {
				return false
			}
		}
	}
	return true
}

func validateAttributesToRemove(xog *etree.Document, file common.DriverFile, path, attributeKey, includeType string) {
	for _, e := range xog.FindElements(path) {
		code := e.SelectAttrValue(attributeKey, "")
		if validateInclude(includeType, code, file) {
			e.Parent().RemoveChild(e)
		}
	}
}
