package transform

import (
	"errors"
	"github.com/andreluzz/cas-xog/constant"
	"github.com/andreluzz/cas-xog/model"
	"github.com/beevik/etree"
	"strconv"
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

		if len(file.Sections) > 0 {
			err := updateSections(xog, aux, file)
			if err != nil {
				return err
			}
		} else {
			updatePropertySet(xog, aux, file)
		}
	}

	return nil
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

	for _, section := range file.Sections {
		if section.Action == constant.ACTION_REPLACE {
			err := processSectionByType(section, sourceView, targetView)
			if err != nil {
				return err
			}
		}
	}

	for _, section := range file.Sections {
		if section.Action == constant.ACTION_UPDATE {
			err := processSectionByType(section, sourceView, targetView)
			if err != nil {
				return err
			}
		}
	}

	for _, section := range file.Sections {
		if section.Action == constant.ACTION_REMOVE {
			err := processSectionByType(section, sourceView, targetView)
			if err != nil {
				return err
			}
		}
	}

	for _, section := range file.Sections {
		if section.Action == constant.ACTION_INSERT {
			err := processSectionByType(section, sourceView, targetView)
			if err != nil {
				return err
			}
		}
	}

	for _, section := range file.Sections {
		if section.Action != constant.ACTION_REMOVE && section.Action != constant.ACTION_REPLACE && section.Action != constant.ACTION_INSERT && section.Action != constant.ACTION_UPDATE {
			return errors.New("invalid action attribute (" + section.Action + ") on tag <section>")
		}
	}

	xog.SetRoot(aux.Root())

	for i, s := range xog.FindElements("//section") {
		s.CreateAttr("sequence", strconv.Itoa(i+1))
	}

	return nil
}

func processSectionByType(section model.Section, sourceView, targetView *etree.Element) error {
	var sourceSection *etree.Element
	if section.Action != constant.ACTION_REMOVE {
		if section.SourcePosition == "" {
			return errors.New("attribute sourcePosition from tag <section> is not defined")
		}

		sourceSection = sourceView.FindElement("//section[" + section.SourcePosition + "]")
		if sourceSection == nil {
			return errors.New("source section for position " + section.SourcePosition + " does not exist")
		}
	}

	targetSection := targetView.FindElement("//nls[1]")
	if section.TargetPosition != "" {
		targetSection = targetView.FindElement("//section[" + section.TargetPosition + "]")
		if targetSection == nil {
			return errors.New("target position " + section.TargetPosition + " does not exist")
		}
	}

	switch section.Action {
	case constant.ACTION_REPLACE:
		if section.TargetPosition == "" {
			return errors.New("cannot replace section because attribute targetPosition from tag <section> is not defined")
		}
		targetView.InsertChild(targetSection, sourceSection)
		targetView.RemoveChild(targetSection)
	case constant.ACTION_REMOVE:
		if section.TargetPosition == "" {
			return errors.New("cannot remove section because attribute targetPosition from tag <section> is not defined")
		}
		targetView.RemoveChild(targetSection)
	case constant.ACTION_INSERT:
		targetView.InsertChild(targetSection, sourceSection)
	case constant.ACTION_UPDATE:
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
			case constant.COLUMN_LEFT:
				if f.InsertBefore == "" {
					columnLeft.AddChild(attributeElement)
				} else {
					columnLeft.InsertChild(targetAttribute, attributeElement)
				}
			case constant.COLUMN_RIGHT:
				if f.InsertBefore == "" {
					columnRight.AddChild(attributeElement)
				} else {
					columnRight.InsertChild(targetAttribute, attributeElement)
				}
			default:
				return errors.New("cannot update section, column value invalid, only 'right' or 'left' are available")
			}
		}
	}

	return nil
}

func validateCodeAndRemoveElementsFromParent(xog *etree.Document, path, code string) {
	for _, e := range xog.FindElements(path) {
		elementCode := e.SelectAttrValue("code", "")
		if elementCode != code {
			e.Parent().RemoveChild(e)
		}
	}
}

func updatePropertySet(xog, aux *etree.Document, file *model.DriverFile) {
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
