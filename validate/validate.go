package validate

import (
	"errors"
	"github.com/andreluzz/cas-xog/constant"
	"github.com/andreluzz/cas-xog/model"
	"github.com/beevik/etree"
)

//Check verify the xog response output looking for errors and warnings
func Check(xog *etree.Document) (model.Output, error) {
	errorOutput := model.Output{Code: constant.OutputError, Debug: ""}
	warningOutput := model.Output{Code: constant.OutputWarning, Debug: ""}

	if xog == nil {
		errorOutput.Debug = "invalid xog"
		return errorOutput, errors.New(errorOutput.Debug)
	}

	output := xog.FindElement("//XOGOutput")

	if output == nil {
		errorOutput.Debug = "no output tag defined"
		return errorOutput, errors.New(errorOutput.Debug)
	}

	statusElement := output.FindElement("Status")

	if statusElement == nil {
		errorOutput.Debug = "no status tag defined"
		return errorOutput, errors.New(errorOutput.Debug)
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
			}
			errorOutput.Debug = desc
			return errorOutput, errors.New(errorOutput.Debug)
		}
	}

	statisticsElement := output.FindElement("Statistics")

	if statisticsElement != nil {
		statTotalNumberOfRecords := statisticsElement.SelectAttrValue("totalNumberOfRecords", "0")
		if statTotalNumberOfRecords == "0" {
			errorOutput.Debug = "output statistics totalNumberOfRecords = 0"
			return errorOutput, errors.New(errorOutput.Debug)
		}
		statFailureRecords := statisticsElement.SelectAttrValue("failureRecords", "0")
		if statFailureRecords != "0" {
			errorOutput.Debug = "output statistics failure on " + statFailureRecords + " records out of " + statTotalNumberOfRecords
			return errorOutput, errors.New(errorOutput.Debug)
		}
	}

	elapsedTime := statusElement.SelectAttrValue("elapsedTime", constant.Undefined)

	debug := ""
	if elapsedTime != constant.Undefined {
		debug = "| Elapsed time: " + elapsedTime
	}

	return model.Output{Code: constant.OutputSuccess, Debug: debug}, nil
}
