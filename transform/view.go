package transform

import (
	"errors"
	"github.com/andreluzz/cas-xog/constant"
	"github.com/andreluzz/cas-xog/model"
	"github.com/beevik/etree"
	"strconv"
	"strings"
)

func specificViewTransformations(xog, aux *etree.Document, file *model.DriverFile) error {
	if len(file.Sections) > 0 && file.Code == "*" {
		return errors.New("tag <section> is only available for single view")
	}

	removeElementFromParent(xog, "//lookups")
	removeElementFromParent(xog, "//objects")

	if file.TargetPartition != "" && file.SourcePartition == "" {
		return errors.New("can't change partition without attribute sourcePartition defined")
	} else if file.TargetPartition != "" && file.SourcePartition != "" {
		changePartition(xog, file.SourcePartition, file.TargetPartition)
	}

	if file.Code != "*" {
		validateCodeAndRemoveElementsFromParent(xog, "//views/property", file.Code)
		validateCodeAndRemoveElementsFromParent(xog, "//views/filter", file.Code)
		validateCodeAndRemoveElementsFromParent(xog, "//views/list", file.Code)

		//auxiliary xog file
		removeElementFromParent(aux, "//lookups")
		removeElementFromParent(aux, "//objects")

		//Only executes is file code is '*' (single view)
		if strings.Contains(file.Code, "*") == false {
			elementsTransform := false

			if len(file.Elements) > 0 {
				var err error
				elementsTransform, err = processElements(xog, aux, file)
				if err != nil {
					return err
				}
			}

			if len(file.Sections) > 0 {
				err := updateSections(xog, aux, file)
				if err != nil {
					return err
				}
				return nil
			}

			if elementsTransform {
				xog.SetRoot(aux.Root())
				return nil
			}
		}

		updateSourceWithTargetPropertySet(xog, aux, file)

		orderViewChildElements(xog, "//property")
		orderViewChildElements(xog, "//propertySet")
		orderViewChildElements(xog, "//filter")
		orderViewChildElements(xog, "//list")
	}

	return nil
}

func orderViewChildElements(xog *etree.Document, sequenceElementPath string) {
	for _, e := range xog.FindElements(sequenceElementPath) {
		e.Parent().AddChild(e.Copy())
		e.Parent().RemoveChild(e)
	}
}

func processElementTypeActionGroup(element model.Element, xog, aux *etree.Document) error {
	if element.Action == constant.ActionInsert {
		sourceGroup := xog.FindElement("//actions/group[@code='" + element.Code + "']")
		if sourceGroup == nil {
			return errors.New("invalid source view action group code")
		}

		targetGroup := aux.FindElement("//actions/group[@code='" + element.Code + "']")

		if element.InsertBefore != constant.Undefined {
			insertBeforeGroup := aux.FindElement("//actions/group[@code='" + element.InsertBefore + "']")
			if insertBeforeGroup == nil {
				return errors.New("invalid insertBefore target view action group code")
			}
			insertBeforeGroup.Parent().InsertChild(insertBeforeGroup, sourceGroup)
		} else {
			if targetGroup == nil {
				actions := aux.FindElement("//actions")
				actions.AddChild(sourceGroup)
			} else {
				targetGroup.Parent().InsertChild(targetGroup, sourceGroup)
			}
		}

		if targetGroup != nil {
			targetGroup.Parent().RemoveChild(targetGroup)
		}

	} else if element.Action == constant.ActionRemove {
		targetGroup := aux.FindElement("//actions/group[@code='" + element.Code + "']")
		if targetGroup == nil {
			return errors.New("cannot remove target view action group - invalid code")
		}
		targetGroup.Parent().RemoveChild(targetGroup)
	}

	return nil
}

func processElementTypeAction(element model.Element, xog, aux *etree.Document) error {
	if element.Action == constant.ActionInsert {
		sourceAction := xog.FindElement("//actions/group/action[@code='" + element.Code + "']")
		if sourceAction == nil {
			return errors.New("invalid source view action code")
		}
		sourceGroupCode := sourceAction.Parent().SelectAttrValue("code", constant.Undefined)
		targetAction := aux.FindElement("//actions/group[@code='" + sourceGroupCode + "']/action[@code='" + element.Code + "']")

		if element.InsertBefore != constant.Undefined {
			insertBeforeAction := aux.FindElement("//actions/group/action[@code='" + element.InsertBefore + "']")
			if insertBeforeAction == nil {
				return errors.New("invalid insertBefore target view action code")
			}
			insertBeforeAction.Parent().InsertChild(insertBeforeAction, sourceAction)
		} else {
			if targetAction == nil {
				targetActionNLS := aux.FindElement("//actions/group[@code='" + sourceGroupCode + "']/nls[1]")
				targetActionNLS.Parent().InsertChild(targetActionNLS, sourceAction)
			} else {
				targetAction.Parent().InsertChild(targetAction, sourceAction)
			}
		}

		if targetAction != nil {
			targetAction.Parent().RemoveChild(targetAction)
		}
	} else if element.Action == constant.ActionRemove {
		targetAction := aux.FindElement("//actions/group/action[@code='" + element.Code + "']")
		if targetAction == nil {
			return errors.New("cannot remove target view action - invalid code")
		}
		targetAction.Parent().RemoveChild(targetAction)
	}

	return nil
}

func processElements(xog, aux *etree.Document, file *model.DriverFile) (bool, error) {

	validateCodeAndRemoveElementsFromParent(aux, "//views/property", file.Code)
	validateCodeAndRemoveElementsFromParent(aux, "//views/filter", file.Code)
	validateCodeAndRemoveElementsFromParent(aux, "//views/list", file.Code)

	validElements := false
	for _, e := range file.Elements {
		if e.XPath == constant.Undefined && (e.Type == constant.ElementTypeActionGroup || e.Type == constant.ElementTypeAction) {
			validElements = true
			switch e.Type {
			case constant.ElementTypeAction:
				err := processElementTypeAction(e, xog, aux)
				if err != nil {
					return false, err
				}
			case constant.ElementTypeActionGroup:
				err := processElementTypeActionGroup(e, xog, aux)
				if err != nil {
					return false, err
				}
			}
		}
	}

	if validElements {
		for i, g := range aux.FindElements("//actions/group") {
			g.CreateAttr("groupOrder", strconv.Itoa(i))
		}
	}

	return validElements, nil
}

func updateSections(xog, aux *etree.Document, file *model.DriverFile) error {

	validateCodeAndRemoveElementsFromParent(aux, "//views/property", file.Code)
	validateCodeAndRemoveElementsFromParent(aux, "//views/filter", file.Code)
	validateCodeAndRemoveElementsFromParent(aux, "//views/list", file.Code)

	nikuDataBus := aux.FindElement("//NikuDataBus")
	aux.SetRoot(nikuDataBus)
	aux.IndentTabs()

	targetView := aux.FindElement("//views/*[@code='" + file.Code + "']")
	if targetView == nil {
		return errors.New("can't process section because the view (" + file.Code + ") does not exist in target environment")
	}

	sourceView := xog.FindElement("//views/*[@code='" + file.Code + "']")
	if sourceView == nil {
		return errors.New("can't process section because the view (" + file.Code + ") does not exist in source environment")
	}

	err := processSections(file, sourceView, targetView, constant.ActionReplace)
	if err != nil {
		return err
	}

	err = processSections(file, sourceView, targetView, constant.ActionUpdate)
	if err != nil {
		return err
	}

	err = processSections(file, sourceView, targetView, constant.ActionRemove)
	if err != nil {
		return err
	}

	err = processSections(file, sourceView, targetView, constant.ActionInsert)
	if err != nil {
		return err
	}

	for _, section := range file.Sections {
		if section.Action != constant.ActionRemove && section.Action != constant.ActionReplace && section.Action != constant.ActionInsert && section.Action != constant.ActionUpdate {
			return errors.New("invalid action attribute (" + section.Action + ") on tag <section>")
		}
	}

	xog.SetRoot(aux.Root())

	for i, s := range xog.FindElements("//section") {
		s.CreateAttr("sequence", strconv.Itoa(i+1))
	}

	return nil
}

func processSections(file *model.DriverFile, sourceView, targetView *etree.Element, action string) error {
	for _, section := range file.Sections {
		if section.Action == action {
			err := processSectionByType(section, sourceView, targetView)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func validateSection(section model.Section, sourceView, targetView *etree.Element) (*etree.Element, *etree.Element, error) {
	var sourceSection *etree.Element
	if section.Action != constant.ActionRemove {
		if section.SourcePosition == "" {
			return nil, nil, errors.New("attribute sourcePosition from tag <section> is not defined")
		}

		sourceSection = sourceView.FindElement("//section[" + section.SourcePosition + "]")
		if sourceSection == nil {
			return nil, nil, errors.New("source section for position " + section.SourcePosition + " does not exist")
		}
	}

	targetSection := targetView.FindElement("//nls[1]")
	if section.TargetPosition != "" {
		targetSection = targetView.FindElement("//section[" + section.TargetPosition + "]")
		if targetSection == nil {
			return nil, nil, errors.New("target position " + section.TargetPosition + " does not exist")
		}
	}

	return sourceSection, targetSection, nil
}

func processViewSectionActionUpdate(section model.Section, sourceSection, targetSection *etree.Element) error {
	if len(section.Fields) == 0 {
		return errors.New("cannot update section because there is no tag <field> defined")
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
	for _, f := range section.Fields {
		if f.Remove {
			removeElement := targetSection.FindElement("//viewFieldDescriptor[@attributeCode='" + f.Code + "']")
			if removeElement == nil {
				return errors.New("cannot remove field because code does not exist in target environment section")
			}
			removeElement.Parent().RemoveChild(removeElement)
			continue
		}

		attributeElement := sourceSection.FindElement("//viewFieldDescriptor[@attributeCode='" + f.Code + "']")
		if attributeElement == nil {
			return errors.New("field attribute code does not exist in source environment view")
		}
		var targetAttribute *etree.Element
		if f.InsertBefore != "" {
			targetAttribute = targetSection.FindElement("//viewFieldDescriptor[@attributeCode='" + f.InsertBefore + "']")
			if targetAttribute == nil {
				//Transform views - section - trying to insert before an attribute that does not exists in target environment
				return errors.New("trying to insert before an field that does not exists in target environment")
			}
		}
		switch f.Column {
		case constant.ColumnLeft:
			if f.InsertBefore == "" {
				columnLeft.AddChild(attributeElement)
			} else {
				columnLeft.InsertChild(targetAttribute, attributeElement)
			}
		case constant.ColumnRight:
			if f.InsertBefore == "" {
				columnRight.AddChild(attributeElement)
			} else {
				columnRight.InsertChild(targetAttribute, attributeElement)
			}
		default:
			return errors.New("cannot update section, column value invalid, only 'right' or 'left' are available")
		}
	}
	return nil
}

func processSectionByType(section model.Section, sourceView, targetView *etree.Element) error {
	sourceSection, targetSection, err := validateSection(section, sourceView, targetView)
	if err != nil {
		return err
	}

	switch section.Action {
	case constant.ActionReplace:
		if section.TargetPosition == "" {
			return errors.New("cannot replace section because attribute targetPosition from tag <section> is not defined")
		}
		targetView.InsertChild(targetSection, sourceSection)
		targetView.RemoveChild(targetSection)
	case constant.ActionRemove:
		if section.TargetPosition == "" {
			return errors.New("cannot remove section because attribute targetPosition from tag <section> is not defined")
		}
		targetView.RemoveChild(targetSection)
	case constant.ActionInsert:
		targetView.InsertChild(targetSection, sourceSection)
	case constant.ActionUpdate:
		err := processViewSectionActionUpdate(section, sourceSection, targetSection)
		if err != nil {
			return err
		}
	}

	return nil
}

func validateCodeAndRemoveElementsFromParent(xog *etree.Document, path, code string) {
	for _, e := range xog.FindElements(path) {
		elementCode := e.SelectAttrValue("code", "")
		if strings.Contains(code, "*") {
			if strings.Contains(elementCode, code[1:]) == false {
				e.Parent().RemoveChild(e)
			}
		} else {
			if elementCode != code {
				e.Parent().RemoveChild(e)
			}
		}
	}
}

func updateSourceWithTargetPropertySet(xog, aux *etree.Document, file *model.DriverFile) {
	sourcePropertySetView := xog.FindElement("//propertySet/update/view[@code='" + file.Code + "']")
	if sourcePropertySetView != nil {
		auxPropertySetView := aux.FindElement("//propertySet/update/view[@code='" + file.Code + "']")
		if auxPropertySetView != nil {
			auxPropertySetView.Parent().InsertChild(auxPropertySetView, sourcePropertySetView)
			auxPropertySetView.Parent().RemoveChild(auxPropertySetView)
		} else {
			propertySetUpdate := aux.FindElement("//propertySet/update")
			if propertySetUpdate != nil {
				propertySetUpdate.InsertChild(propertySetUpdate.FindElement("//nls[1]"), sourcePropertySetView)
			}
		}
	}

	auxPropertySet := aux.FindElement("//propertySet")
	xogPropertySet := xog.FindElement("//propertySet")
	xogPropertySet.Parent().AddChild(auxPropertySet)
	xogPropertySet.Parent().RemoveChild(xogPropertySet)
}
