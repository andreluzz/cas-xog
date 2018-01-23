package transform

import (
	"github.com/andreluzz/cas-xog/constant"
	"github.com/andreluzz/cas-xog/model"
	"github.com/andreluzz/cas-xog/util"
	"github.com/andreluzz/cas-xog/validate"
	"github.com/beevik/etree"
	"strings"
)

//ProcessPackageFile transform package drivers when needed using auxiliary xog from the installation environment
func ProcessPackageFile(file *model.DriverFile, packageFolder, writeFolder string, definitions []model.Definition) model.Output {
	output := model.Output{Code: constant.OutputSuccess, Debug: constant.Undefined}

	if file == nil {
		output.Code = constant.OutputError
		output.Debug = "trying to process a nil DriverFile"
		return output
	}

	xog := etree.NewDocument()
	err := xog.ReadFromFile(packageFolder + file.Path)
	if err != nil {
		output.Code = constant.OutputError
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
			output.Code = constant.OutputError
			output.Debug = err.Error()
			return output
		}
	} else if file.PackageTransform && !file.NeedPackageTransform() {
		output.Code = constant.OutputWarning
		output.Debug = "only single views and menu can be transformed in packages processing"
	}

	processPackageDefinitions(definitions, file, xog)

	xog.IndentTabs()
	util.ValidateFolder(writeFolder + util.GetPathFolder(file.Path))
	xog.WriteToFile(writeFolder + "/" + file.Path)

	return output
}

func processPackageDefinitions(definitions []model.Definition, file *model.DriverFile, xog *etree.Document) {
	for _, def := range definitions {
		if def.Value == def.Default {
			continue
		}
		switch def.Action {
		case constant.PackageActionChangePartitionModel:
			if file.Type == constant.TypeObject {
				e := xog.FindElement("//object[@partitionModelCode]")
				if e != nil {
					e.CreateAttr("partitionModelCode", def.Value)
				}
			}
		case constant.PackageActionChangePartition:
			if file.Type == constant.TypeObject || file.Type == constant.TypeView {
				changePartition(xog, constant.Undefined, def.Value)
			}
		case constant.PackageActionReplaceString:
			if def.Value == constant.Undefined {
				continue
			}
			replaced := strings.Replace(def.To, "##DEFINITION_VALUE##", def.Value, 1)
			if replaced == def.From {
				continue
			}
			if def.TransformTypes == constant.Undefined || strings.Contains(def.TransformTypes, file.Type) {
				findAndReplace(xog, []model.FileReplace{{From: def.From, To: replaced}})
			}
		}
	}
}
