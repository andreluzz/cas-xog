package transform

import (
	"errors"
	"github.com/beevik/etree"
	"github.com/andreluzz/cas-xog/common"
	"strconv"
)

func specificViewTransformations(xog, aux *etree.Document, file common.DriverFile) error {

	if len(file.Sections) > 0 && file.Code == "*" {
		return errors.New("tag <section> is only available for single view")
	}

	removeElementFromParent(xog, "//lookups")
	removeElementFromParent(xog, "//objects")

	if file.TargetPartition != "" {
		changePartition(xog, file)
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

func updateSections(xog, aux *etree.Document, file common.DriverFile) error {

	validateCodeAndRemoveElementsFromParent(aux, "//views/property", file.Code)
	validateCodeAndRemoveElementsFromParent(aux, "//views/filter", file.Code)
	validateCodeAndRemoveElementsFromParent(aux, "//views/list", file.Code)

	nikuDataBus := aux.FindElement("//NikuDataBus")
	aux.SetRoot(nikuDataBus)
	aux.IndentTabs()

	targetView := aux.FindElement("//property[@code='" + file.Code + "']")
	if targetView == nil {
		return errors.New("can't process section because the view does not exist in target environment")
	}

	sourceView := xog.FindElement("//property[@code='" + file.Code + "']")
	if sourceView == nil {
		return errors.New("can't process section because the view does not exist in source environment")
	}

	for _, section := range file.Sections {
		if section.Action == "replace" {
			err := processSectionByType(section, sourceView, targetView)
			if err != nil {
				return err
			}
		}
	}

	for _, section := range file.Sections {
		if section.Action == "update" {
			err := processSectionByType(section, sourceView, targetView)
			if err != nil {
				return err
			}
		}
	}

	for _, section := range file.Sections {
		if section.Action == "remove" {
			err := processSectionByType(section, sourceView, targetView)
			if err != nil {
				return err
			}
		}
	}

	for _, section := range file.Sections {
		if section.Action == "insert" {
			err := processSectionByType(section, sourceView, targetView)
			if err != nil {
				return err
			}
		}
	}

	xog.SetRoot(aux.Root())

	for i, s := range xog.FindElements("//section") {
		s.CreateAttr("sequence", strconv.Itoa(i+1))
	}

	return nil
}

func processSectionByType(section common.ViewSection, sourceView, targetView *etree.Element) error {
	var sourceSection *etree.Element
	if section.Action != "remove" {
		if section.SourcePosition == "" {
			return errors.New("attribute sourcePosition from tag <section> is not defined")
		}

		sourceSection = sourceView.FindElement("//section["+ section.SourcePosition +"]")
		if sourceSection == nil {
			return errors.New("source section for position " + section.SourcePosition + " does not exist")
		}
	}

	targetSection := targetView.FindElement("//nls[1]")
	if section.TargetPosition != "" {
		targetSection = targetView.FindElement("//section["+ section.TargetPosition +"]")
		if targetSection == nil {
			return errors.New("target position " + section.SourcePosition + " out of bounds")
		}
	}

	switch section.Action {
	case "replace":
		if section.TargetPosition == "" {
			return errors.New("cannot replace section because attribute targetPosition from tag <section> is not defined")
		}
		targetView.InsertChild(targetSection, sourceSection)
		targetView.RemoveChild(targetSection)
	case "remove":
		if section.TargetPosition == "" {
			return errors.New("cannot remove section because attribute targetPosition from tag <section> is not defined")
		}
		targetView.RemoveChild(targetSection)
	case "insert":
		targetView.InsertChild(targetSection, sourceSection)
	case "update":
		if len(section.Fields) == 0 {
			return errors.New("cannot update section because there is no tag <filed> defined")
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
			case "left":
				if f.InsertBefore == "" {
					columnLeft.AddChild(attributeElement)
				} else {
					columnLeft.InsertChild(targetAttribute, attributeElement)
				}
			case "right":
				if f.InsertBefore == "" {
					columnRight.AddChild(attributeElement)
				} else {
					columnRight.InsertChild(targetAttribute, attributeElement)
				}
			default:
				return errors.New("cannot update section, column value invalid, only 'right' or 'left' are available")
			}

		}
	default:
		return errors.New("invalid action attribute (" + section.Action + ") on tag <section>")
	}

	return nil
}

func updatePropertySet(xog, aux *etree.Document, file common.DriverFile) {
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
