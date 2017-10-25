package validate

import (
	"testing"
	"github.com/andreluzz/cas-xog/common"
	"github.com/beevik/etree"
)

var packageMockFolder string

func init() {
	packageMockFolder = "../" + common.FOLDER_MOCK + "validate/"
}

func TestCheck(t *testing.T) {
	xog:= etree.NewDocument()
	xog.ReadFromFile(packageMockFolder + "mockSuccessOutput.xml")

	output, err:= Check(xog)

	if output.Code != common.OUTPUT_SUCCESS {
		t.Errorf("Expected output %s and recëived output %s", common.OUTPUT_SUCCESS, output.Code)
	}
	if err != nil {
		t.Errorf("Encountered an error while checking the output. %s", err.Error())
	}
}

func TestCheckRegularSuccessOutput(t *testing.T) {

}

func TestCheckToReturnErrorOutput(t *testing.T) {

}

func TestCheckToReturnWarningOutput(t *testing.T) {

}

func TestCheckToReturnInvalidOutput(t *testing.T) {

}

func TestCheckToReturnZeroResultsOutput(t *testing.T) {

}

func TestCheckToReturnNilOutput(t *testing.T) {

}

