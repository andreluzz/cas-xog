package validate

import (
	"errors"
	"github.com/andreluzz/cas-xog/common"
	"github.com/beevik/etree"
)

func Check(xog *etree.Document) (common.XOGOutput, error) {
	errorOutput := common.XOGOutput{Code: common.OUTPUT_ERROR, Debug: ""}
	warningOutput := common.XOGOutput{Code: common.OUTPUT_WARNING, Debug: ""}

	if xog == nil {
		return errorOutput, errors.New("invalid xog")
	}

	output := xog.FindElement("//XOGOutput")

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

	elapsedTime := statusElement.SelectAttrValue("elapsedTime", common.UNDEFINED)

	debug := ""
	if elapsedTime != common.UNDEFINED {
		debug = "| Elapsed time: " + elapsedTime
	}

	return common.XOGOutput{Code: common.OUTPUT_SUCCESS, Debug: debug}, nil
}
