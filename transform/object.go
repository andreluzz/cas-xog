package transform

import (
	"strings"

	"github.com/andreluzz/cas-xog/constant"
	"github.com/andreluzz/cas-xog/model"
	"github.com/beevik/etree"
)

func getWildcardAttributesElements(xog *etree.Document, file *model.DriverFile) {
	for index, f := range file.Elements {
		if f.Type == constant.ElementTypeAttribute && strings.Contains(f.Code, "*") {
			for _, e := range xog.FindElements("//customAttribute") {
				code := e.SelectAttrValue("code", constant.Undefined)
				if strings.HasPrefix(code, f.Code[:len(f.Code)-1]) {
					file.Elements = append(file.Elements, model.Element{
						Code: code,
						Type: constant.ElementTypeAttribute,
					})
				}
			}
			file.Elements = append(file.Elements[:index], file.Elements[index+1:]...)
		}
	}
}

func specificObjectTransformations(xog, aux *etree.Document, file *model.DriverFile) {

	removeChildObjects(xog)

	if hasElementsToProcess(file) {
		getWildcardAttributesElements(xog, file)
		objectProcessElements(xog, aux, file)
	}

	if file.SourcePartition != constant.Undefined {
		var codesToRemove []string

		for _, e := range xog.FindElements("//customAttribute") {
			partition := e.SelectAttrValue("partitionCode", constant.Undefined)
			if partition != file.SourcePartition {
				codesToRemove = append(codesToRemove, e.SelectAttrValue("code", constant.Undefined))

			}
		}

		for _, code := range codesToRemove {
			for _, e := range xog.FindElements("//*[@code='" + code + "']") {
				e.Parent().RemoveChild(e)
			}
			for _, e := range xog.FindElements("//*[@attributeCode='" + code + "']") {
				e.Parent().RemoveChild(e)
			}
		}
	}

	if file.TargetPartition != constant.Undefined {
		changePartition(xog, file.SourcePartition, file.TargetPartition)
	}

	if file.PartitionModel != constant.Undefined {
		element := xog.FindElement("//object[@code='" + file.Code + "']")
		element.CreateAttr("partitionModelCode", file.PartitionModel)
	}
}

func hasElementsToProcess(file *model.DriverFile) bool {
	for _, f := range file.Elements {
		if f.Code != constant.Undefined && (f.Type == constant.ElementTypeAction || f.Type == constant.ElementTypeLink || f.Type == constant.ElementTypeAttribute) {
			return true
		}
	}
	return false
}

func objectProcessElements(xog, aux *etree.Document, file *model.DriverFile) {
	removeChildObjects(aux)

	if file.OnlyElements {
		removeElementsFromParent(aux, "//customAttribute")
		removeElementsFromParent(aux, "//action")
		removeElementsFromParent(aux, "//attributeAutonumbering")
		removeElementsFromParent(aux, "//attributeDefault")
		removeElementsFromParent(aux, "//link")
		removeElementsFromParent(aux, "//displayMapping")
		removeElementsFromParent(aux, "//scoreContributions")
		removeElementsFromParent(aux, "//capabilities")
	}

	for _, f := range file.Elements {
		if f.Code != constant.Undefined && (f.Type == constant.ElementTypeAction || f.Type == constant.ElementTypeLink || f.Type == constant.ElementTypeAttribute) {
			for _, e := range xog.FindElements("//[@code='" + f.Code + "']") {
				removeElementFromParent(aux, "//"+e.Tag+"[@code='"+f.Code+"']")
				parentTag := e.Parent().Tag
				if parentTag == "object" {
					targetElement := aux.FindElement("//customAttribute")
					if e.Tag == "attributeDefault" || targetElement == nil {
						targetElement = aux.FindElement("//links")
					}
					targetElement.Parent().InsertChild(targetElement, e)
				} else {
					aux.FindElement("//" + parentTag).AddChild(e)
				}
			}
			if f.Type == constant.ElementTypeAttribute {
				for _, e := range xog.FindElements("//*[@attributeCode='" + f.Code + "']") {
					removeElementFromParent(aux, "//"+e.Tag+"[@attributeCode='"+f.Code+"']")
					aux.FindElement("//" + e.Parent().Tag).AddChild(e)
				}
			}
		}
	}
	xog.SetRoot(aux.Root())
}

func removeChildObjects(doc *etree.Document) {
	object := doc.FindElement("//objects/object")
	for _, e := range object.FindElements("//object") {
		object.RemoveChild(e)
	}
}
