package transform

import (
	"github.com/andreluzz/cas-xog/constant"
	"github.com/andreluzz/cas-xog/model"
	"github.com/beevik/etree"
)

func specificObjectTransformations(xog, aux *etree.Document, file *model.DriverFile) {

	removeChildObjects(xog)

	if hasElementsToProcess(file) {
		removeChildObjects(aux)

		for _, f := range file.Elements {
			if f.Code != constant.Undefined && (f.Type == constant.ElementTypeAction || f.Type == constant.ElementTypeLink || f.Type == constant.ElementTypeAttribute) {
				for _, e := range xog.FindElements("//[@code='" + f.Code + "']") {
					removeElementFromParent(aux, "//"+e.Tag+"[@code='"+f.Code+"']")
					parentTag := e.Parent().Tag
					if parentTag == "object" {
						targetElement := aux.FindElement("//customAttribute")
						if e.Tag == "attributeDefault" {
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

func removeChildObjects(doc *etree.Document) {
	object := doc.FindElement("//objects/object")
	for _, e := range object.FindElements("//object") {
		object.RemoveChild(e)
	}
}
