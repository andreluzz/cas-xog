package main

import (
	"github.com/beevik/etree"
	"strconv"
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

func removeTagFromParentEqual(parent *etree.Element, elems []*etree.Element, attrCode string, viewCode string) {
	for _, e := range elems {
		code := e.SelectAttrValue(attrCode, "")
		if code != viewCode {
			parent.RemoveChild(e)
		}
	}
}

func removeTagFromParent(parent *etree.Element, elems []*etree.Element, attrCode string, codes string) {
	for _, e := range elems {
		code := e.SelectAttrValue(attrCode, "")
		if !strings.Contains(codes, code) {
			parent.RemoveChild(e)
		}
	}
}

func removeTags(elems []*etree.Element, attrCode string, codes string) {
	for _, e := range elems {
		code := e.SelectAttrValue(attrCode, "")
		if !strings.Contains(codes, code) {
			parent := e.Parent()
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

	SimplifyLookupStructure := false
	objectFiltered := false
	tagsRemoved := RemoveUnnecessaryTags(xogfile.Type)
	partitionReplaced := ReplacePartition(xogfile.SourcePartition, xogfile.TargetPartition)
	if xogfile.Type == "lookups" && xogfile.OnlyStructure {
		LookupOnlyStructure()
		SimplifyLookupStructure = true
	}
	if xogfile.SingleView && xogfile.Type == "views" {
		SingleView(xogfile.Code, xogfile.CopyToView)
	}
	if len(xogfile.Includes) > 0 && xogfile.Type == "objects" {
		objectFiltered = FilterObjectAtributes(xogfile)
	}

	xogOutputElement := doc.FindElement("//XOGOutput")
	if xogOutputElement != nil {
		errorInformationElement := doc.FindElement("//ErrorInformation/Severity")
		if errorInformationElement == nil {
			root.RemoveChild(xogOutputElement)
		}
	}

	doc.Indent(4)
	if err := doc.WriteToFile(path); err != nil {
		panic(err)
	}

	return tagsRemoved || partitionReplaced || xogfile.SingleView || objectFiltered || SimplifyLookupStructure
}

func LookupOnlyStructure() {
	lookupElement := doc.FindElement("//dynamicLookup")

	lookupElement.CreateAttr("hiddenAttributeName", "id")
	lookupElement.CreateAttr("objectCode", "")
	lookupElement.CreateAttr("sortAttributeName", "id")
	lookupElement.CreateAttr("sortDirection", "asc")

	lookupElement.RemoveChild(lookupElement.FindElement("//nsql"))
	lookupElement.RemoveChild(lookupElement.FindElement("//displayedSuggestionAttributes"))
	lookupElement.RemoveChild(lookupElement.FindElement("//searchedSuggestionAttributes"))
	lookupElement.RemoveChild(lookupElement.FindElement("//browsePage"))

	lookupStructureExampleDoc := etree.NewDocument()
	if err := lookupStructureExampleDoc.ReadFromString(readDefault.Examples[0].Value); err != nil {
		panic(err)
	}

	lookupElement.AddChild(lookupStructureExampleDoc.FindElement("//nsql"))
	lookupElement.AddChild(lookupStructureExampleDoc.FindElement("//displayedSuggestionAttributes"))
	lookupElement.AddChild(lookupStructureExampleDoc.FindElement("//searchedSuggestionAttributes"))
	lookupElement.AddChild(lookupStructureExampleDoc.FindElement("//browsePage"))
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
		removeTagFromParent(object, doc.FindElements("//customAttribute"), "code", attibutesCodes)

		//remove attributeDefault
		removeTagFromParent(object, doc.FindElements("//attributeDefault"), "code", attibutesCodes)

		//remove displayMappings
		displayMappings := doc.FindElement("//displayMappings")
		if displayMappings != nil {
			removeTagFromParent(displayMappings, doc.FindElements("//displayMapping"), "attributeCode", attibutesCodes)
		}

		//remove autonumbering
		autonumbering := doc.FindElement("//autonumbering")
		if autonumbering != nil {
			removeTagFromParent(autonumbering, doc.FindElements("//attributeAutonumbering"), "code", attibutesCodes)
		}

		//remove audit
		audit := doc.FindElement("//audit")
		if audit != nil {
			removeTagFromParent(audit, doc.FindElements("//audit/attribute"), "code", attibutesCodes)

			auditElements := audit.ChildElements()
			if len(auditElements) == 0 {
				object.RemoveChild(audit)
			}
		}

		links := doc.FindElement("//links")
		if linksCodes == "" {
			//remove links
			object.RemoveChild(links)
		} else {
			removeTagFromParent(links, doc.FindElements("//links/link"), "code", linksCodes)
		}

		actions := doc.FindElement("//actions")
		if actionsCodes == "" {
			//remove actions
			object.RemoveChild(actions)
		} else {
			removeTagFromParent(actions, doc.FindElements("//actions/action"), "code", actionsCodes)
		}
	}

	return true
}

func SingleView(viewCode string, copyToView string) {
	root := doc.SelectElement("NikuDataBus")
	content := root.SelectElement("contentPack")
	views := content.SelectElement("views")

	removeTagFromParentEqual(views, doc.FindElements("//property"), "code", viewCode)
	removeTagFromParentEqual(views, doc.FindElements("//filter"), "code", viewCode)
	removeTagFromParentEqual(views, doc.FindElements("//list"), "code", viewCode)

	if copyToView != "" {
		for _, e := range views.ChildElements() {
			code := e.SelectAttrValue("code", "")
			if viewCode == code {
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

	//Remove unecessary tags
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

	if action == "groups" {
		//remove members
		members := doc.FindElement("//group/members")
		if members != nil {
			parent := members.Parent()
			parent.RemoveChild(members)
			transf = true
		}
	}

	return transf
}

func MergeMenus(xogfile XogDriverFile, sourcePath string, targetPath string) (bool, string) {
	sourceDoc := etree.NewDocument()
	if err := sourceDoc.ReadFromFile(sourcePath); err != nil {
		//trying read source file that does not exists
		return false, "\033[91mERRO-04\033[0m"
	}
	targetDoc := etree.NewDocument()
	if err := targetDoc.ReadFromFile(targetPath); err != nil {
		//trying read target file that does not exists
		return false, "\033[91mERRO-05\033[0m"
	}

	targetMenuElement := targetDoc.FindElement("//menu")

	//process menus
	for _, m := range xogfile.Menus {
		sourceSectionElement := sourceDoc.FindElement("//menu/section[@code='" + m.Code + "']")
		targetSectionElement := targetMenuElement.FindElement("//section[" + strconv.Itoa(m.TargetPosition) + "]")

		switch m.Action {
		case "insert":
			if m.TargetPosition == 0 || targetSectionElement == nil {
				targetMenuElement.AddChild(sourceSectionElement)
			} else {
				targetMenuElement.InsertChild(targetSectionElement, sourceSectionElement)
			}
		case "update":
			targetSectionElement = targetMenuElement.FindElement("//section[@code='" + m.Code + "']")
			if targetSectionElement == nil {
				//Transform menus - cannot update a section that does not exist in target
				return false, "\033[91mER-TM02\033[0m"
			}
			if len(m.Links) <= 0 {
				//Transform menus - lacking link tags to update menu
				return false, "\033[91mER-TM03\033[0m"
			}
			//insert the links inside the section
			for _, l := range m.Links {
				sourceLinkElement := sourceSectionElement.FindElement("//link[@code='" + l.Code + "']")
				targetLinkElement := targetSectionElement.FindElement("//link[" + strconv.Itoa(l.TargetPosition) + "]")
				if targetLinkElement == nil {
					targetSectionElement.AddChild(sourceLinkElement)
				} else {
					targetSectionElement.InsertChild(targetLinkElement, sourceLinkElement)
				}
			}
			// update section links position value
			i := 1
			for _, s := range targetSectionElement.FindElements("//link") {
				s.CreateAttr("position", strconv.Itoa(i))
				i += 1
			}
		case "replace":
			targetSectionElement = targetMenuElement.FindElement("//section[@code='" + m.Code + "']")
			if targetSectionElement == nil {
				//Transform menus - cannot replace a section that does not exist in target
				return false, "\033[91mER-TM04\033[0m"
			}
			if m.TargetPosition != 0 {
				//If attribute targetPosition exists change the position of the section that is being replaced
				targetPositionElement := targetMenuElement.FindElement("//section[" + strconv.Itoa(m.TargetPosition) + "]")
				targetMenuElement.InsertChild(targetPositionElement, sourceSectionElement)
			} else {
				targetMenuElement.InsertChild(targetSectionElement, sourceSectionElement)
			}
			targetMenuElement.RemoveChild(targetSectionElement)
		default:
			//Transform menus - invalid action at menu tag
			return false, "\033[91mER-TM01\033[0m"
		}
	}

	// update section links position value
	i := 1
	for _, s := range targetMenuElement.FindElements("//section") {
		s.CreateAttr("position", strconv.Itoa(i))
		i += 1
	}

	targetDoc.Indent(4)
	if err := targetDoc.WriteToFile(sourcePath); err != nil {
		panic(err)
	}
	return true, "\033[92mSUCCESS\033[0m"
}

func MergeOBS(xogfile XogDriverFile, sourcePath string, targetPath string) (bool, string) {
	sourceDoc := etree.NewDocument()
	if err := sourceDoc.ReadFromFile(sourcePath); err != nil {
		//trying to merge views and source view file does not exists
		return false, "\033[91mERRO-04\033[0m"
	}
	targetDoc := etree.NewDocument()
	if err := targetDoc.ReadFromFile(targetPath); err != nil {
		//trying to merge views and target view file does not exists
		return false, "\033[91mERRO-05\033[0m"
	}

	//clean source - remove associated objects
	for _, a := range sourceDoc.FindElements("//associatedObject") {
		a.Parent().RemoveChild(a)
	}
	for _, a := range sourceDoc.FindElements("//UserSecurity") {
		a.Parent().RemoveChild(a)
	}

	obs := targetDoc.FindElement("//obs")
	obs.CreateAttr("complete", "true")

	for _, u := range xogfile.Units {
		var targetParent *etree.Element

		if u.ParentName == "" {
			targetParent = targetDoc.FindElement("//obs")
		} else {
			targetParent = targetDoc.FindElement("//unit[@name='" + u.ParentName + "']")
		}

		if targetParent == nil {
			//Transform obs - wrong unit's parent name in target environment
			return false, "\033[91mER-TO01\033[0m"
		}

		targetUnit := targetParent.FindElement("//unit[@name='" + u.Name + "']")
		if u.Remove {
			if targetUnit == nil {
				//Transform obs - cannot remove unit, name does not exist in target environment
				return false, "\033[91mER-TO02\033[0m"
			}
			targetParent.RemoveChild(targetUnit)
		} else {
			if targetUnit != nil {
				targetParent.RemoveChild(targetUnit)
			}
			var sourceParent *etree.Element
			if u.ParentName == "" {
				sourceParent = sourceDoc.FindElement("//obs")
			} else {
				sourceParent = sourceDoc.FindElement("//unit[@name='" + u.ParentName + "']")
			}

			unit := sourceParent.FindElement("//unit[@name='" + u.Name + "']").Copy()
			if unit == nil {
				//Transform obs - wrong unit's name in source environment
				return false, "\033[91mER-TO03\033[0m"
			}

			if u.RemoveUnitChilds {
				for _, child := range unit.SelectElements("unit") {
					unit.RemoveChild(child)
				}
			}

			if u.ParentName == "" {
				targetParent.AddChild(unit)
			} else {
				targetAssociatedObject := targetParent.FindElement("./associatedObject[0]")
				if targetAssociatedObject != nil {
					targetParent.InsertChild(targetAssociatedObject, unit)
				} else {
					targetRights := targetParent.FindElement("./rights[0]")
					if targetRights != nil {
						targetParent.InsertChild(targetRights, unit)
					} else {
						targetParent.AddChild(unit)
					}
				}
			}
		}
	}

	targetDoc.Indent(4)
	if err := targetDoc.WriteToFile(sourcePath); err != nil {
		panic(err)
	}
	return true, "\033[92mSUCCESS\033[0m"
}

func MergeViews(xogfile XogDriverFile, sourcePath string, targetPath string) (bool, string) {
	sourceDoc := etree.NewDocument()
	if err := sourceDoc.ReadFromFile(sourcePath); err != nil {
		//trying to merge views and source view file does not exists
		return false, "\033[91mERRO-04\033[0m"
	}
	targetDoc := etree.NewDocument()
	if err := targetDoc.ReadFromFile(targetPath); err != nil {
		//trying to merge views and target view file does not exists
		return false, "\033[91mERRO-05\033[0m"
	}

	status := false
	message := "\033[93mWARNING\033[0m"

	//get view from source and insert in target if it does not exists in target
	targetView := targetDoc.FindElement("//views/*[@code='" + xogfile.Code + "']")
	if targetView == nil {
		sourceView := sourceDoc.FindElement("//views/*[@code='" + xogfile.Code + "']")
		if sourceView != nil {
			targetPropertySet := targetDoc.FindElement("//propertySet")
			if targetPropertySet != nil {
				targetParent := targetPropertySet.Parent()
				targetParent.InsertChild(targetPropertySet, sourceView)
			}
		}
	}

	//process replace
	for _, s := range xogfile.Sections {
		if s.Action == "replace" {
			status, message = processSection(s, targetDoc, sourceDoc)
			if !status {
				return status, message
			}
		}
	}

	//process update
	for _, s := range xogfile.Sections {
		if s.Action == "update" {
			status, message = processSection(s, targetDoc, sourceDoc)
			if !status {
				return status, message
			}
		}
	}

	//process remove
	for _, s := range xogfile.Sections {
		if s.Action == "remove" {
			status, message = processSection(s, targetDoc, sourceDoc)
			if !status {
				return status, message
			}
		}
	}

	//process insert
	for _, s := range xogfile.Sections {
		if s.Action == "insert" {
			status, message = processSection(s, targetDoc, sourceDoc)
			if !status {
				return status, message
			}
		}
	}

	//process tag action
	for _, a := range xogfile.Actions {
		status, message = processAction(a, targetDoc, sourceDoc)
		if !status {
			return status, message
		}
	}

	//update target sections sequence value
	i := 1
	for _, s := range targetDoc.FindElements("//section") {
		s.CreateAttr("sequence", strconv.Itoa(i))
		i += 1
	}

	//update target propertySet including or replacing the view
	sourcePropertySetElements := sourceDoc.FindElements("//propertySet")
	for _, sourcePropertySetElement := range sourcePropertySetElements {
		sourcePropertySetViewElement := sourcePropertySetElement.FindElement("//view[@code='" + xogfile.Code + "']")
		if sourcePropertySetViewElement != nil {
			targetPropertySetElements := targetDoc.FindElements("//propertySet")
			for _, targetPropertySetElement := range targetPropertySetElements {
				var targetInsertBeforeViewElement *etree.Element

				targetCurrentViewElement := targetPropertySetElement.FindElement("//view[@code='" + xogfile.Code + "']")

				if xogfile.InsertBefore != "" {
					targetInsertBeforeViewElement = targetPropertySetElement.FindElement("//view[@code='" + xogfile.InsertBefore + "']")
				}

				//if exists the element from insertBefore use it to define the position of the view
				if targetInsertBeforeViewElement != nil {
					parent := targetInsertBeforeViewElement.Parent()
					parent.InsertChild(targetInsertBeforeViewElement, sourcePropertySetViewElement)
				} else {
					if targetCurrentViewElement != nil {
						parent := targetCurrentViewElement.Parent()
						parent.InsertChild(targetCurrentViewElement, sourcePropertySetViewElement)
					} else {
						//if there is no insertBefore defined insert as the last one
						nlsElement := targetPropertySetElement.FindElement("//nls[1]")
						targetPropertySetElement.InsertChild(nlsElement, sourcePropertySetViewElement)
					}
				}

				//if the target already have an element with the same view code then we need to remove it
				if targetCurrentViewElement != nil {
					parent := targetCurrentViewElement.Parent()
					parent.RemoveChild(targetCurrentViewElement)
				}

				status = true
				message = "\033[92mSUCCESS\033[0m"
			}
		}
	}

	targetDoc.Indent(4)
	if err := targetDoc.WriteToFile(sourcePath); err != nil {
		panic(err)
	}
	return status, message
}

func processAction(a XogViewAction, targetDoc *etree.Document, sourceDoc *etree.Document) (bool, string) {
	sourceGroup := sourceDoc.FindElement("//actions/group[@code='" + a.GroupCode + "']")

	if sourceGroup == nil {
		//Transform views - action - group code does not exist in source environment view
		return false, "\033[91mER-TVA1\033[0m"
	}

	targetGroup := targetDoc.FindElement("//actions/group[@code='" + a.GroupCode + "']")

	if sourceGroup == nil {
		//Transform views - action - group code does not exist in target environment view
		return false, "\033[91mER-TVA2\033[0m"
	}

	if a.Remove {
		action := targetGroup.FindElement("//action[@code='" + a.Code + "']")
		if action == nil {
			//Transform views - action - cannot remove action because there is no match code in target environment
			return false, "\033[91mER-TVA3\033[0m"
		}
		targetGroup.RemoveChild(action)
	} else {
		var targetAttribute *etree.Element
		if a.InsertBefore != "" {
			targetAttribute = targetGroup.FindElement("//action[@code='" + a.InsertBefore + "']")
		}
		if a.InsertBefore == "" || targetAttribute == nil {
			targetAttribute = targetGroup.FindElement("//nls[1]")
		}
		attributeElement := sourceGroup.FindElement("//action[@code='" + a.Code + "']")
		targetGroup.InsertChild(targetAttribute, attributeElement)
	}

	return true, "\033[92mSUCCESS\033[0m"
}

func processSection(s XogViewSection, targetDoc *etree.Document, sourceDoc *etree.Document) (bool, string) {
	var sourceSection *etree.Element
	if s.Action != "remove" {
		if s.SourceSectionPosition == "" {
			//Transform views - section - invalid value for sourceSectionPosition
			return false, "\033[91mER-TVS3\033[0m"
		}
		sourceSection = sourceDoc.FindElement("//section[" + s.SourceSectionPosition + "]")
	}

	if sourceSection == nil {
		if s.Action != "remove" {
			//Transform views - section - invalid value for sourceSectionPosition
			return false, "\033[91mER-TVS3\033[0m"
		}
	} else {
		//get all attributes codes from source section
		var sourceSectionAttributesCodes []string
		if len(s.Attributes) > 0 {
			for _, a := range s.Attributes {
				sourceSectionAttributesCodes = append(sourceSectionAttributesCodes, a.Code)
			}

			//remove unnecessary attributes from source section
			removeTags(sourceSection.FindElements("//viewFieldDescriptor"), "attributeCode", strings.Join(sourceSectionAttributesCodes, ";"))
		} else {
			elems := sourceSection.FindElements("//viewFieldDescriptor")
			if elems != nil {
				for _, e := range elems {
					sourceSectionAttributesCodes = append(sourceSectionAttributesCodes, e.SelectAttrValue("attributeCode", ""))
				}
			}
		}

		//remove attributes in target that will be included from source
		for i := range sourceSectionAttributesCodes {
			element := targetDoc.FindElement("//viewFieldDescriptor[@attributeCode='" + sourceSectionAttributesCodes[i] + "']")
			if element != nil {
				parent := element.Parent()
				parent.RemoveChild(element)
			}
		}

	}

	targetSection := targetDoc.FindElement("//section[" + s.TargetSectionPosition + "]")
	if targetSection == nil {
		if s.Action == "replace" || s.Action == "update" || s.Action == "remove" {
			//Transform views - section - invalid value for targetSectionPosition
			return false, "\033[91mER-TVS2\033[0m"
		}
	}

	switch s.Action {
	case "remove":
		parent := targetSection.Parent()
		parent.RemoveChild(targetSection)
	case "insert":
		if targetSection == nil {
			//If there is no section for TargetSectionPosition insert section as last one
			targetSection = targetDoc.FindElement("//nls[1]")
		}
		parent := targetSection.Parent()
		parent.InsertChild(targetSection, sourceSection)
	case "replace":
		parent := targetSection.Parent()
		parent.InsertChild(targetSection, sourceSection)
		parent.RemoveChild(targetSection)
	case "update":
		if len(s.Attributes) <= 0 {
			//Transform views - section - action update without attributes
			return false, "\033[91mER-TVS4\033[0m"
		}

		columnRight := targetSection.FindElement("//column[@sequence='2']")
		if columnRight == nil {
			//Create column if it does not exists
			columnRight = etree.NewElement("column")
			columnRight.CreateAttr("sequence", "2")
			nlsElement := targetSection.FindElement("//nls[1]")
			targetSection.InsertChild(nlsElement, columnRight)
		}
		columnLeft := targetSection.FindElement("//column[@sequence='1']")
		if columnLeft == nil {
			//Create column if it does not exists
			columnLeft = etree.NewElement("column")
			columnLeft.CreateAttr("sequence", "1")
			targetSection.InsertChild(columnRight, columnLeft)
		}

		for _, a := range s.Attributes {
			if !a.Remove {
				attributeElement := sourceSection.FindElement("//viewFieldDescriptor[@attributeCode='" + a.Code + "']")
				if attributeElement == nil {
					//Transform views - general - attribute code does not exist in source environment view
					return false, "\033[91mER-TVG2\033[0m"
				}
				var targetAttribute *etree.Element
				if a.InsertBefore != "" {
					targetAttribute = targetSection.FindElement("//viewFieldDescriptor[@attributeCode='" + a.InsertBefore + "']")
					if targetAttribute == nil {
						//Transform views - section - trying to insert before an attribute that does not exists in target environment
						return false, "\033[91mER-TVS6\033[0m"
					}
				}
				switch a.Column {
				case "left":
					if a.InsertBefore == "" {
						columnLeft.AddChild(attributeElement)
					} else {
						columnLeft.InsertChild(targetAttribute, attributeElement)
					}
				case "right":
					if a.InsertBefore == "" {
						columnRight.AddChild(attributeElement)
					} else {
						columnRight.InsertChild(targetAttribute, attributeElement)
					}
				default:
					//Transform views - section - column value invalid, only 'right' or 'left' are available
					return false, "\033[91mER-TVS5\033[0m"
				}
			}
		}
	default:
		//Transform views - section - invalid action at section tag
		return false, "\033[91mER-TVS1\033[0m"
	}

	return true, "\033[92mSUCCESS\033[0m"
}

func Validate(path string) (bool, string) {
	if initStatus := initDoc(path); initStatus == false {
		//General - Trying to validate a write file that does not exist
		return false, "\033[91mER-GN01\033[0m"
	}

	statusElement := doc.FindElement("//XOGOutput/Status")
	message := ""
	status := false

	if statusElement != nil {
		s := statusElement.SelectAttrValue("state", "UNKNOWN")
		statisticsElement := doc.FindElement("//XOGOutput/Statistics")
		totalRecords := "0"
		if statisticsElement != nil {
			totalRecords = statisticsElement.SelectAttrValue("totalNumberOfRecords", "UNKNOWN")
		}
		if s == "SUCCESS" && statisticsElement != nil && totalRecords != "0" {
			//validate warning
			errorInformationElement := doc.FindElement("//ErrorInformation/Severity")
			if errorInformationElement != nil {
				if errorInformationElement.Text() == "WARNING" {
					status = true
					message = "\033[93mWARNING\033[0m"
				} else {
					status = false
					message = "\033[91m" + errorInformationElement.Text() + "\033[0m"
				}
			} else {
				message = "\033[92m" + s + "\033[0m"
				status = true
			}
		} else {
			message = "\033[91m" + s + "\033[0m"
			status = false
		}
	} else {
		//General - Output file does not have the status tag inside XOGOutput block
		message = "\033[91mER-GN02\033[0m"
		status = false
	}

	return status, message
}
