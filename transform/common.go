package transform

import (
	"errors"
	"github.com/andreluzz/cas-xog/common"
	"github.com/beevik/etree"
	"strings"
)

func Execute(xog, aux *etree.Document, file common.DriverFile) error {
	err := errors.New("")
	err = nil

	headerElement := xog.FindElement("//NikuDataBus/Header")
	if headerElement == nil {
		return errors.New("[transform error] no header element")
	}
	headerElement.CreateAttr("version", "8.0")

	switch file.Type {
	case common.LOOKUP:
		specificLookupTransformations(xog, file)
	case common.PROCESS:
		err = specificProcessTransformations(xog, aux, file)
		if err != nil {
			return errors.New("[transform error] " + err.Error())
		}
	case common.OBJECT:
		specificObjectTransformations(xog, file)
	case common.VIEW:
		err = specificViewTransformations(xog, aux, file)
		if err != nil {
			return errors.New("[transform error] " + err.Error())
		}
	case common.PORTLET, common.QUERY:
		removeElementFromParent(xog, "//lookups")
		removeElementFromParent(xog, "//objects")
	case common.MENU:
		err = specificMenuTransformations(xog, aux, file)
		if err != nil {
			return errors.New("[transform error] " + err.Error())
		}
	case common.RESOURCE_CLASS_INSTANCE, common.WIP_CLASS_INSTANCE, common.TRANSACTION_CLASS_INSTANCE:
		headerElement.CreateAttr("version", "12.0")
	case common.INVESTMENT_CLASS_INSTANCE:
		headerElement.CreateAttr("version", "14.1")
	}

	if len(file.Elements) > 0 {
		for _,e := range file.Elements {
			if e.Action == common.ACTION_REMOVE && e.XPath != "" && e.Type == "" && e.Code == "" {
				if strings.HasPrefix(e.XPath, "/") {
					e.XPath = "." + e.XPath
				}
				removeElementsFromParent(xog, e.XPath)
			}
		}
	}

	removeElementFromParent(xog, "//partitionModels")
	removeElementFromParent(xog, "//XOGOutput")

	if len(file.Replace) > 0 {
		findAndReplace(xog, file.Replace)
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

func findAndReplace(xog *etree.Document, replace []common.FileReplace) {
	xogString, _ := xog.WriteToString()
	for _, r := range replace {
		xogString = strings.Replace(xogString, r.From, r.To, -1)
	}
	xmlResult := etree.NewDocument()
	xmlResult.ReadFromString(xogString)
	xog.SetRoot(xmlResult.Root())
}

func changePartition(xog *etree.Document, sourcePartition, targetPartition string) {
	var elems []*etree.Element
	if sourcePartition == "" {
		elems = xog.FindElements("//*[@partitionCode]")
	} else {
		elems = xog.FindElements("//*[@partitionCode='" + sourcePartition + "']")
	}

	for _, e := range elems {
		e.CreateAttr("partitionCode", targetPartition)
	}
	if sourcePartition == "" {
		for _, e := range xog.FindElements("//*[@dataProviderPartitionId]") {
			e.CreateAttr("dataProviderPartitionId", targetPartition)
		}
	} else {
		for _, e := range xog.FindElements("//*[@dataProviderPartitionId='" + sourcePartition + "']") {
			e.CreateAttr("dataProviderPartitionId", targetPartition)
		}
	}
}
