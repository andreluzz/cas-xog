package transform

import (
	"github.com/andreluzz/cas-xog/common"
	"github.com/beevik/etree"
)

var packageMockFolder string

func init() {
	packageMockFolder = "../" + common.FOLDER_MOCK + "transform/"
}

func readMockResultAndCompare(xog *etree.Document, compareXml string) bool {
	xog.IndentTabs()
	xogString, _ := xog.WriteToString()

	xogProcessedToCompare := etree.NewDocument()
	xogProcessedToCompare.ReadFromFile(packageMockFolder + compareXml)
	xogProcessedToCompare.IndentTabs()

	xogProcessedToCompareString, _ := xogProcessedToCompare.WriteToString()
	if xogString != xogProcessedToCompareString {
		return false
	}
	return true
}