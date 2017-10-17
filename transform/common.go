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
