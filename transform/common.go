package transform

import (
	"errors"
	"github.com/andreluzz/cas-xog/common"
	"github.com/beevik/etree"
	"strings"
)

func Process(xog, aux *etree.Document, file common.DriverFile) error {
	err := errors.New("")
	err = nil

	headerElement := xog.FindElement("//NikuDataBus/Header")
	if headerElement == nil {
		return errors.New("[transform error] no header element")
	}
	headerElement.CreateAttr("version", "8.0")

	switch file.Type {
	case common.PROCESS:
		err = specificProcessTransformations(xog, aux, file)
		if err != nil {
			return errors.New("[transform error] " + err.Error())
		}
	case common.OBJECT:
		err = specificObjectTransformations(xog, file)
		if err != nil {
			return errors.New("[transform error] " + err.Error())
		}
	case common.VIEW:
		err = specificViewTransformations(xog, aux, file)
		if err != nil {
			return errors.New("[transform error] " + err.Error())
		}
	case common.PORTLET:
		removeElementFromParent(xog, "//lookups")
	case common.MENU:
		removeElementFromParent(xog, "//objects")
		removeElementFromParent(xog, "//pages")
	case common.OBS:
		if file.RemoveObjAssoc {
			removeElementsFromParent(xog, "//associatedObject")
		}
		if file.RemoveSecurity {
			removeElementsFromParent(xog, "//Security")
			removeElementsFromParent(xog, "//rights")
		}
	case common.RESOURCE_CLASS_INSTANCE, common.WIP_CLASS_INSTANCE, common.TRANSACTION_CLASS_INSTANCE:
		headerElement.CreateAttr("version", "12.0")
	case common.INVESTMENT_CLASS_INSTANCE:
		headerElement.CreateAttr("version", "14.1")
	}

	removeElementFromParent(xog, "//partitionModels")
	removeElementFromParent(xog, "//XOGOutput")

	if len(file.Replace) > 0 {
		findAndReplace(xog, file)
	}

	return err
}

func removeElementFromParent(xog *etree.Document, path string) {
	element := xog.FindElement(path)
	if element != nil {
		element.Parent().RemoveChild(element)
	}
}

func removeElementsFromParent(xog *etree.Document, path string) {
	for _, e := range xog.FindElements(path) {
		e.Parent().RemoveChild(e)
	}
}

func validateCodeAndRemoveElementsFromParent(xog *etree.Document, path, code string) {
	for _, e := range xog.FindElements(path) {
		elementCode := e.SelectAttrValue("code", "")
		if elementCode != code {
			e.Parent().RemoveChild(e)
		}
	}
}

func findAndReplace(xog *etree.Document, file common.DriverFile) {
	xogString, _ := xog.WriteToString()
	for _, r := range file.Replace {
		xogString = strings.Replace(xogString, r.From, r.To, -1)
	}
	xmlResult := etree.NewDocument()
	xmlResult.ReadFromString(xogString)
	xog.SetRoot(xmlResult.Root())
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
