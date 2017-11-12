package transform

import (
	"errors"
	"regexp"
	"strings"

	"github.com/andreluzz/cas-xog/common"
	"github.com/beevik/etree"
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
		for _, e := range file.Elements {
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

	for _, e := range xog.FindElements("//*[@dataProviderPartitionId='" + sourcePartition + "']") {
		e.CreateAttr("dataProviderPartitionId", targetPartition)
	}
}

func IncludeCDATA(xogString string, iniTagRegexpStr string, endTagRegexpStr string) (string, error) {
	iniTagRegexp, _ := regexp.Compile(iniTagRegexpStr)
	endTagRegexp, _ := regexp.Compile(endTagRegexpStr)

	iniIndex := iniTagRegexp.FindAllStringIndex(xogString, -1)
	endIndex := endTagRegexp.FindAllStringIndex(xogString, -1)

	shiftIndex := 0

	for i := 0; i < len(iniIndex); i++ {
		index := iniIndex[i][1] + shiftIndex
		xogString = xogString[:index] + "<![CDATA[" + xogString[index:]

		sqlString := xogString[index:endIndex[i][1]]

		paramRegexp, _ := regexp.Compile(`<(.*):param(.*)/>`)
		paramIndex := paramRegexp.FindStringIndex(sqlString)

		shiftIndex += 9

		eIndex := endIndex[i][0] + shiftIndex
		if len(paramIndex) > 0 {
			eIndex = endIndex[i][0] + 12 - (len(sqlString) - paramIndex[0])
		}

		xogString = xogString[:eIndex] + "]]>" + xogString[eIndex:]

		shiftIndex += 3
	}

	replacer := strings.NewReplacer("&gt;", ">", "&lt;", "<", "&apos;", "'", "&quot;", "\"")
	xogString = replacer.Replace(xogString)

	return xogString, nil
}
