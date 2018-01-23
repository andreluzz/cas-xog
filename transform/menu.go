package transform

import (
	"errors"
	"github.com/andreluzz/cas-xog/constant"
	"github.com/andreluzz/cas-xog/model"
	"github.com/beevik/etree"
	"strconv"
)

func specificMenuTransformations(xog, aux *etree.Document, file *model.DriverFile) error {
	removeElementFromParent(xog, "//objects")
	removeElementFromParent(xog, "//pages")

	if len(file.Sections) > 0 {
		removeElementFromParent(aux, "//objects")
		removeElementFromParent(aux, "//pages")

		for _, s := range file.Sections {
			sourceSectionElement := xog.FindElement("//section[@code='" + s.Code + "']")
			if sourceSectionElement == nil {
				return errors.New("invalid source menu section code(" + s.Code + ")")
			}

			targetSectionElement := aux.FindElement("//section[@code='" + s.Code + "']")
			switch s.Action {
			case constant.ActionUpdate:
				err := executeActionUpdate(targetSectionElement, sourceSectionElement, s)
				if err != nil {
					return err
				}
			case constant.ActionInsert:
				err := executeActionInsert(targetSectionElement, sourceSectionElement, aux, s)
				if err != nil {
					return err
				}
			}
		}
		xog.SetRoot(aux.Root())
	}

	return nil
}

func executeActionInsert(targetSectionElement, sourceSectionElement *etree.Element, aux *etree.Document, section model.Section) error {
	if targetSectionElement != nil {
		return errors.New("cannot insert section code(" + section.Code + ") because it already exists in target")
	}
	position := "-1"
	if section.TargetPosition != constant.Undefined {
		position = section.TargetPosition
	}
	if len(section.Links) > 0 {
		for _, e := range sourceSectionElement.FindElements("//link") {
			removeLink := true
			for _, l := range section.Links {
				if l.Code == e.SelectAttrValue("pageCode", "") {
					removeLink = false
				}
			}
			if removeLink {
				e.Parent().RemoveChild(e)
			}
		}
	}
	targetElementAtPosition := aux.FindElement("//section[" + position + "]")
	if targetElementAtPosition == nil {
		return errors.New("invalid target section position(" + position + ")")
	}
	targetElementAtPosition.Parent().InsertChild(targetElementAtPosition, sourceSectionElement)
	for i, e := range aux.FindElements("//section") {
		e.CreateAttr("position", strconv.Itoa(i))
	}

	return nil
}

func executeActionUpdate(targetSectionElement, sourceSectionElement *etree.Element, section model.Section) error {
	if targetSectionElement == nil {
		return errors.New("invalid target menu section code(" + section.Code + ")")
	}
	if len(section.Links) <= 0 {
		return errors.New("can't update menu section code(" + section.Code + ") without tag link")
	}
	for _, l := range section.Links {
		sourceLinkElement := sourceSectionElement.FindElement("//link[@pageCode='" + l.Code + "']")
		if sourceLinkElement == nil {
			return errors.New("invalid source menu section link code(" + l.Code + ")")
		}
		targetSectionElement.AddChild(sourceLinkElement)
	}
	for i, e := range targetSectionElement.FindElements("//link") {
		e.CreateAttr("position", strconv.Itoa(i))
	}

	return nil
}
