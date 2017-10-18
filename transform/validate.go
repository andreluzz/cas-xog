package transform

import (
	"errors"
	"github.com/andreluzz/cas-xog/common"
	"github.com/beevik/etree"
)

const UNDEFINED string = ""
const OUTPUT_ERROR string = "error"
const OUTPUT_WARNING string = "warning"
const OUTPUT_SUCCESS string = "success"

func Validate(xog *etree.Document) (common.XOGOutput, error) {
	output := xog.FindElement("//XOGOutput")
	errorOutput := common.XOGOutput{Code: OUTPUT_ERROR, Debug: ""}
	warningOutput := common.XOGOutput{Code: OUTPUT_WARNING, Debug: ""}

	if output == nil {
		return errorOutput, errors.New("no output tag defined")
	}

	statusElement := output.FindElement("Status")

	if statusElement == nil {
		return errorOutput, errors.New("no status tag defined")
	}

	errorInformationElement := output.FindElement("//ErrorInformation")

	if errorInformationElement != nil {
		desc := ""
		severityElement := errorInformationElement.FindElement("Severity")
		descriptionElement := errorInformationElement.FindElement("Description")
		if descriptionElement != nil {
			desc = descriptionElement.Text()
		}
		if severityElement != nil {
			if severityElement.Text() == "WARNING" {
				warningOutput.Debug = desc
				return warningOutput, nil
			} else {
				return errorOutput, errors.New(desc)
			}
		}
	}

	statisticsElement := output.FindElement("Statistics")

	if statisticsElement != nil {
		statTotalNumberOfRecords := statisticsElement.SelectAttrValue("totalNumberOfRecords", "0")
		if statTotalNumberOfRecords == "0" {
			return errorOutput, errors.New("output statistics totalNumberOfRecords = 0")
		}
		statFailureRecords := statisticsElement.SelectAttrValue("failureRecords", "0")
		if statFailureRecords != "0" {
			return errorOutput, errors.New("output statistics failure on " + statFailureRecords + " records out of " + statTotalNumberOfRecords)
		}
	}

	elapsedTime := statusElement.SelectAttrValue("elapsedTime", UNDEFINED)

	debug := ""
	if elapsedTime != UNDEFINED {
		debug = "| Elapsed time: " + elapsedTime
	}

	return common.XOGOutput{Code: OUTPUT_SUCCESS, Debug: debug}, nil
}
