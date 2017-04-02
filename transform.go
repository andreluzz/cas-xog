package main

import (
	"github.com/beevik/etree"
	"strings"
)

var doc *etree.Document

func initDoc(path string) bool {
	doc = etree.NewDocument()
	if err := doc.ReadFromFile(path); err != nil {
		//panic(err)
		return false
	}
	return true
}

func removeTag(parent *etree.Element, elems []*etree.Element, attrCode string, codes string) {
	for _, e := range elems {
		code := e.SelectAttrValue(attrCode, "")
		if !strings.Contains(codes, code) {
			parent.RemoveChild(e)
		}
	}
}

func Transform(xogfile XogDriverFile, path string) bool {
	if initStatus := initDoc(path); initStatus == false {
		return false
	}

	root := doc.SelectElement("NikuDataBus")
	//Replace version for compatibility reasons
	header := root.SelectElement("Header")
	if header != nil {
		header.CreateAttr("version", "8.0")
	}

	var attributes []string
	for _, att := range xogfile.Attributes {
		attributes = append(attributes, att.Code)
	}

	objectAttributesFiltered := false
	tagsRemoved := RemoveUnnecessaryTags(xogfile.Type)
	partitionReplaced := ReplacePartition(xogfile.SourcePartition, xogfile.TargetPartition)
	if xogfile.SingleView && xogfile.Type == "views" {
		SingleView(xogfile.Code, xogfile.CopyToView)
		FilterViewsAtributes(attributes)
	}
	if len(xogfile.Attributes) > 0 && xogfile.Type == "objects" {
		objectAttributesFiltered = FilterObjectAtributes(attributes)
	}

	doc.Indent(4)
	if err := doc.WriteToFile(path); err != nil {
		panic(err)
	}

	return tagsRemoved || partitionReplaced || xogfile.SingleView || objectAttributesFiltered
}

func FilterViewsAtributes(attributes []string) {

}

func FilterObjectAtributes(attributes []string) bool {
	object := doc.FindElement("//object")
	codes := strings.Join(attributes, ",")

	if object != nil {
		//remove unnecessary customAttribute
		removeTag(object, doc.FindElements("//customAttribute"), "code", codes)

		//remove unnecessary attributeDefault
		removeTag(object, doc.FindElements("//attributeDefault"), "code", codes)

		//remove links
		links := doc.FindElement("//links")
		object.RemoveChild(links)

		//remove unnecessary displayMappings
		displayMappings := doc.FindElement("//displayMappings")
		if displayMappings != nil {
			removeTag(displayMappings, doc.FindElements("//displayMapping"), "attributeCode", codes)
		}

		//remove links
		actions := doc.FindElement("//actions")
		object.RemoveChild(actions)

		//remove links autonumbering
		autonumbering := doc.FindElement("//autonumbering")
		if autonumbering != nil {
			removeTag(autonumbering, doc.FindElements("//attributeAutonumbering"), "code", codes)
		}
	}

	return true
}

func SingleView(viewCode string, copyToView string) {
	root := doc.SelectElement("NikuDataBus")
	content := root.SelectElement("contentPack")
	views := content.SelectElement("views")

	for _, e := range doc.FindElements("//property") {
		code := e.SelectAttrValue("code", "")
		if viewCode != code {
			views.RemoveChild(e)
		} else {
			if copyToView != "" {
				e.CreateAttr("code", copyToView)
				if strings.Contains(copyToView, "Create") {
					e.CreateAttr("type", "create")
				} else {
					e.CreateAttr("type", "update")
				}
			}
		}
	}
}

func ReplacePartition(source string, target string) bool {
	if target == "" {
		return false
	}

	var elems []*etree.Element
	if source == "" {
		elems = doc.FindElements("//*[@partitionCode]")
	} else {
		elems = doc.FindElements("//*[@partitionCode='" + source + "']")
	}

	for _, e := range elems {
		e.CreateAttr("partitionCode", target)
	}

	return true
}

func RemoveUnnecessaryTags(action string) bool {
	transf := false
	var removeTags []string
	removeTags = append(removeTags, "partitionModels")

	switch action {
	case "views":
		removeTags = append(removeTags, "objects")
		removeTags = append(removeTags, "lookups")
		transf = true
	case "processes", "portlets":
		removeTags = append(removeTags, "lookups")
		transf = true
	}

	root := doc.SelectElement("NikuDataBus")
	content := root.SelectElement("contentPack")

	//Remove unecessary removeTags
	if content != nil {
		for i := range removeTags {
			e := content.SelectElement(removeTags[i])
			if e != nil {
				content.RemoveChild(e)
			}
		}
	}

	return transf
}

func Validate(path string) (bool, string) {
	if initStatus := initDoc(path); initStatus == false {
		//ERROR-0: Reading file does not exist
		return false, "\033[91mERROR-0\033[0m"
	}

	elem_status := doc.FindElement("//XOGOutput/Status")
	message := ""
	status := false

	if elem_status != nil {
		s := elem_status.SelectAttrValue("state", "unknown")
		elem_statistics := doc.FindElement("//XOGOutput/Statistics")
		totalRecords := "0"
		if elem_statistics != nil {
			totalRecords = elem_statistics.SelectAttrValue("totalNumberOfRecords", "unknown")
		}
		if s == "SUCCESS" && elem_statistics != nil && totalRecords != "0" {
			message = "\033[92m" + s + "\033[0m"
			status = true
		} else {
			message = "\033[91m" + s + "\033[0m"
			status = false
		}
	} else {
		//ERROR-1: Output file does not have the XOGOutput Status tag
		message = "\033[91mERROR-1\033[0m"
		status = false
	}

	return status, message
}
