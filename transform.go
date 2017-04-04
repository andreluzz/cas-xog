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

	objectFiltered := false
	menuFiltered := false
	tagsRemoved := RemoveUnnecessaryTags(xogfile.Type)
	partitionReplaced := ReplacePartition(xogfile.SourcePartition, xogfile.TargetPartition)
	if xogfile.SingleView && xogfile.Type == "views" {
		SingleView(xogfile.Code, xogfile.CopyToView)
	}
	if len(xogfile.Includes) > 0 && xogfile.Type == "objects" {
		objectFiltered = FilterObjectAtributes(xogfile)
	}
	if len(xogfile.Includes) > 0 && xogfile.Type == "menus" {
		menuFiltered = FilterMenuItems(xogfile)
	}

	xogOutputElement := doc.FindElement("//XOGOutput")
	if xogOutputElement != nil {
		root.RemoveChild(xogOutputElement)
	}

	doc.Indent(4)
	if err := doc.WriteToFile(path); err != nil {
		panic(err)
	}

	return tagsRemoved || partitionReplaced || xogfile.SingleView || objectFiltered || menuFiltered
}

func FilterMenuItems(xogfile XogDriverFile) bool {
	menu := doc.FindElement("//menu")
	sectionsCodes := ""
	linksCodes := ""
	var cleanSectionLinks []string

	for _, i := range xogfile.Includes {
		switch i.Type {
		case "menuSection":
			sectionsCodes += i.Code + ";"
		case "menuLink":
			linksCodes += i.Code + ";"
			sectionsCodes += i.SectionCode + ";"
			cleanSectionLinks = append(cleanSectionLinks, i.SectionCode)
		}
	}

	if menu != nil {
		//remove unnecessary sections
		removeTag(menu, doc.FindElements("//section"), "code", sectionsCodes)

		//remove unnecessary links
		for index := range cleanSectionLinks {
			section := doc.FindElement("//section[@code='" + cleanSectionLinks[index] + "']")
			removeTag(section, doc.FindElements("//section[@code='"+cleanSectionLinks[index]+"']/link"), "pageCode", linksCodes)
		}
	}

	return true
}

func FilterObjectAtributes(xogfile XogDriverFile) bool {
	object := doc.FindElement("//object")
	attibutesCodes := ""
	actionsCodes := ""
	linksCodes := ""

	for _, i := range xogfile.Includes {
		switch i.Type {
		case "attribute":
			attibutesCodes += i.Code + ";"
		case "action":
			actionsCodes += i.Code + ";"
		case "link":
			linksCodes += xogfile.Code + "." + i.Code + ";"
		}
	}

	if object != nil {
		//remove customAttribute
		removeTag(object, doc.FindElements("//customAttribute"), "code", attibutesCodes)

		//remove attributeDefault
		removeTag(object, doc.FindElements("//attributeDefault"), "code", attibutesCodes)

		//remove displayMappings
		displayMappings := doc.FindElement("//displayMappings")
		if displayMappings != nil {
			removeTag(displayMappings, doc.FindElements("//displayMapping"), "attributeCode", attibutesCodes)
		}

		//remove autonumbering
		autonumbering := doc.FindElement("//autonumbering")
		if autonumbering != nil {
			removeTag(autonumbering, doc.FindElements("//attributeAutonumbering"), "code", attibutesCodes)
		}

		//remove audit
		audit := doc.FindElement("//audit")
		if audit != nil {
			removeTag(audit, doc.FindElements("//audit/attribute"), "code", attibutesCodes)
		}

		auditElements := audit.ChildElements()
		if len(auditElements) == 0 {
			object.RemoveChild(audit)
		}

		links := doc.FindElement("//links")
		if linksCodes == "" {
			//remove links
			object.RemoveChild(links)
		} else {
			removeTag(links, doc.FindElements("//links/link"), "code", linksCodes)
		}

		actions := doc.FindElement("//actions")
		if actionsCodes == "" {
			//remove actions
			object.RemoveChild(actions)
		} else {
			removeTag(actions, doc.FindElements("//actions/action"), "code", actionsCodes)
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

	//remove unnecessary views from propertyset
	propertySet := views.SelectElement("propertySet")
	propertySetCreate := propertySet.SelectElement("create")
	propertySetUpdate := propertySet.SelectElement("update")

	if propertySetCreate.SelectAttrValue("code", "") == viewCode {
		propertySet.RemoveChild(propertySetUpdate)
	} else {
		propertySet.RemoveChild(propertySetCreate)
		removeTag(propertySetUpdate, doc.FindElements("//update/view"), "code", viewCode)
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
	case "menus":
		removeTags = append(removeTags, "objects")
		removeTags = append(removeTags, "pages")
		transf = true
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

	if action == "objects" {
		//remove subobjects
		object := content.FindElement("//objects/object")
		for _, e := range doc.FindElements("//objects/object/object") {
			object.RemoveChild(e)
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

func MergeViews(xogfile XogDriverFile, sourcePath string, targetPath string) (bool, string) {
	sourceDoc := etree.NewDocument()
	if err := sourceDoc.ReadFromFile(sourcePath); err != nil {
		return false, "\033[91mERROR-2\033[0m"
	}
	targetDoc := etree.NewDocument()
	if err := targetDoc.ReadFromFile(targetPath); err != nil {
		return false, "\033[91mERROR-2\033[0m"
	}

	viewFieldDescriptorCodes := ""
	for _, i := range xogfile.Includes {
		viewFieldDescriptorCodes += i.Code + ";"
	}

	//remove source unnecessary attributes
	for _, e := range sourceDoc.FindElements("//viewFieldDescriptor") {
		code := e.SelectAttrValue("attributeCode", "")
		if !strings.Contains(viewFieldDescriptorCodes, code) {
			parent := e.Parent()
			parent.RemoveChild(e)
		}
	}

	for _, i := range xogfile.Includes {
		//remove attribute from target if exists
		targetViewFieldDescriptorElement := targetDoc.FindElement("//viewFieldDescriptor[@attributeCode='" + i.Code + "']")
		if targetViewFieldDescriptorElement != nil {
			targetColumnElement := targetViewFieldDescriptorElement.Parent()
			targetColumnElement.RemoveChild(targetViewFieldDescriptorElement)
		}

		//Get attribute information from source
		sourceViewFieldDescriptorElement := sourceDoc.FindElement("//viewFieldDescriptor[@attributeCode='" + i.Code + "']")
		if sourceViewFieldDescriptorElement != nil {
			sourceColumnElement := sourceViewFieldDescriptorElement.Parent()
			sourceSectionElement := sourceColumnElement.Parent()

			sourceSectionSequenceAttrValue := sourceSectionElement.SelectAttrValue("sequence", "")

			//Include attribute in target
			targetSectionElement := targetDoc.FindElement("//section[@sequence='" + sourceSectionSequenceAttrValue + "']")
			if targetSectionElement != nil {
				sourceColumnSequenceAttrValue := sourceColumnElement.SelectAttrValue("sequence", "")
				targetColumnElement := targetSectionElement.FindElement("//column[@sequence='" + sourceColumnSequenceAttrValue + "']")
				if targetColumnElement != nil {
					if i.InsertAfter == "" and i.InsertBefore == "" {
						targetColumnElement.AddChild(sourceViewFieldDescriptorElement)
					} else {
						//get all target column elements
						targetColumnElements := targetColumnElement.ChildElements()
						//remove elements from column
						removeTag(targetColumnElement, targetColumnElement.ChildElements(), "attributeCode", "")
						//insert elements in order
						for _, e := range targetColumnElements {
							if i.InsertBefore == e.SelectAttrValue("attributeCode", "") {
								targetColumnElement.AddChild(sourceViewFieldDescriptorElement)
							}
							targetColumnElement.AddChild(e)
							if i.InsertAfter == e.SelectAttrValue("attributeCode", "") {
								targetColumnElement.AddChild(sourceViewFieldDescriptorElement)
							}
						}
					}
				} else {
					//insert column from source
					targetSectionElement.AddChild(sourceColumnElement)
				}
			} else {
				//insert section from source
				targetPropertyElement := targetDoc.FindElement("//property")
				targetPropertyElement.AddChild(sourceSectionElement)
			}
		}
	}

	targetDoc.Indent(4)
	if err := targetDoc.WriteToFile(sourcePath); err != nil {
		panic(err)
	}

	return true, "\033[92mSUCCESS\033[0m"
}
