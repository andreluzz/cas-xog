package api

import (
	"github.com/andreluzz/cas-xog/constant"
	"github.com/andreluzz/cas-xog/model"
	"github.com/andreluzz/cas-xog/util"
)

//ProcessDriverFile execute an api resquest return the response
func ProcessDriverFile(file *model.DriverFile, action, sourceFolder, outputFolder string, environments *model.Environments, restFunc util.Rest) model.Output {
	var err error
	if action == "r" {
		switch file.Type {
		case constant.APITypeBlueprint:
			err = readBlueprint(file, outputFolder, environments, restFunc)
		}
	} else {
		switch file.Type {
		case constant.APITypeBlueprint:
			//err = readBlueprint(file, sourceFolder, outputFolder, environments, restFunc)
		}
	}

	if err != nil {
		return model.Output{Code: constant.OutputError, Debug: err.Error()}
	}

	return model.Output{Code: constant.OutputSuccess, Debug: constant.Undefined}
}
