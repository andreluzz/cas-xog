package main

import (
	"github.com/beevik/etree"
	"io/ioutil"
)

type Helper struct {
}

func (h *Helper) Validate(path string) (bool, string) {
	status := false
	message := ""

	doc := etree.NewDocument()
	if err := doc.ReadFromFile(path); err != nil {
		//ERROR-0: Reading file does not exist
		return false, "\033[91mERROR-0\033[0m"
	}

	elem_status := doc.FindElement("//XOGOutput/Status")

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

func (h *Helper) SingleView(path string, viewCode string) {
	doc := etree.NewDocument()
	if err := doc.ReadFromFile(path); err != nil {
		panic(err)
	}
	root := doc.SelectElement("NikuDataBus")
	content := root.SelectElement("contentPack")
	views := content.SelectElement("views")

	for _, e := range doc.FindElements("//property") {
		code := e.SelectAttrValue("code", "")
		if viewCode != code {
			views.RemoveChild(e)
		}
	}

	doc.Indent(4)
	if err := doc.WriteToFile(path); err != nil {
		panic(err)
	}
}

func (h *Helper) ReplacePartition(path string, source string, target string) bool {
	if target == "" {
		return false
	}

	doc := etree.NewDocument()
	if err := doc.ReadFromFile(path); err != nil {
		panic(err)
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

	doc.Indent(4)
	if err := doc.WriteToFile(path); err != nil {
		panic(err)
	}

	return true
}

func (h *Helper) RemoveUnnecessaryTags(path string, action string) bool {
	transform := true
	doc := etree.NewDocument()
	if err := doc.ReadFromFile(path); err != nil {
		panic(err)
	}
	root := doc.SelectElement("NikuDataBus")
	content := root.SelectElement("contentPack")

	//Replace version for compatibility reasons
	header := root.SelectElement("Header")
	if header != nil {
		header.CreateAttr("version", "8.0")
	}

	var removeTags []string
	removeTags = append(removeTags, "partitionModels")

	switch action {
	case "views":
		removeTags = append(removeTags, "objects")
		removeTags = append(removeTags, "lookups")
	case "processes", "portlets":
		removeTags = append(removeTags, "lookups")
	default:
		transform = false
	}

	//Remove unecessary removeTags
	if content != nil {
		for i := range removeTags {
			e := content.SelectElement(removeTags[i])
			if e != nil {
				content.RemoveChild(e)
			}
		}
	}

	doc.Indent(4)
	if err := doc.WriteToFile(path); err != nil {
		panic(err)
	}

	return transform
}

func (h *Helper) loadFile(path string) []byte {
	xmlFile, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err)
	}
	return xmlFile
}
