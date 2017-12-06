package transform

import (
	"github.com/andreluzz/cas-xog/constant"
	"github.com/andreluzz/cas-xog/model"
	"github.com/andreluzz/cas-xog/util"
	"github.com/andreluzz/cas-xog/validate"
	"github.com/beevik/etree"
	"strings"
)

func ProcessPackageFile(file *model.DriverFile, packageFolder, writeFolder string, definitions []model.Definition) model.Output {
	output := model.Output{Code: constant.OUTPUT_SUCCESS, Debug: constant.UNDEFINED}

	if file == nil {
		output.Code = constant.OUTPUT_ERROR
		output.Debug = "trying to process a nil DriverFile"
		return output
	}

	xog := etree.NewDocument()
	err := xog.ReadFromFile(packageFolder + file.Path)
	if err != nil {
		output.Code = constant.OUTPUT_ERROR
		output.Debug = err.Error()
		return output
	}

	if file.PackageTransform && file.NeedPackageTransform() {
		auxResponse := etree.NewDocument()
		auxResponse.ReadFromString(file.GetAuxXML())
		output, err = validate.Check(auxResponse)
		if err != nil {
			output.Debug = "aux validation - " + output.Debug
			return output
		}
		err = Execute(xog, auxResponse, file)
		if err != nil {
			output.Code = constant.OUTPUT_ERROR
			output.Debug = err.Error()
			return output
		}
	} else if file.PackageTransform && !file.NeedPackageTransform() {
		output.Code = constant.OUTPUT_WARNING
		output.Debug = "only single views and menu can be transformed in packages processing"
	}

	for _, def := range definitions {
		if def.Value == def.Default {
			continue
		}
		switch def.Action {
		case constant.PACKAGE_ACTION_CHANGE_PARTITION_MODEL:
			if file.Type == constant.OBJECT {
				e := xog.FindElement("//object[@partitionModelCode]")
				if e != nil {
					e.CreateAttr("partitionModelCode", def.Value)
				}
			}
		case constant.PACKAGE_ACTION_CHANGE_PARTITION:
			if file.Type == constant.OBJECT || file.Type == constant.VIEW {
				changePartition(xog, constant.UNDEFINED, def.Value)
			}
		case constant.PACKAGE_ACTION_REPLACE_STRING:
			if def.Value == constant.UNDEFINED {
				continue
			}
			replaced := strings.Replace(def.To, "##DEFINITION_VALUE##", def.Value, 1)
			if replaced == def.From {
				continue
			}
			if def.TransformTypes == constant.UNDEFINED || strings.Contains(def.TransformTypes, file.Type) {
				findAndReplace(xog, []model.FileReplace{{From: def.From, To: replaced}})
			}
		}
	}

	xog.IndentTabs()
	util.ValidateFolder(writeFolder)
	xog.WriteToFile(writeFolder + "/" + file.Path)

	return output
}
