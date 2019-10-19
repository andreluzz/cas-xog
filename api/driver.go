package api

import (
	"fmt"
	"net/url"

	"github.com/andreluzz/cas-xog/constant"
	"github.com/andreluzz/cas-xog/model"
	"github.com/andreluzz/cas-xog/util"
)

//ProcessDriverFile execute an api resquest return the response
func ProcessDriverFile(file *model.DriverFile, action, sourceFolder, outputFolder string, environments *model.Environments, restFunc util.Rest) model.Output {
	var err error
	switch action {
	case "r":
		switch file.APIType() {
		case constant.APITypeBlueprint:
			err = readBlueprint(file, outputFolder, environments, restFunc)
		case constant.APITypeTeam:
			err = readTeam(file, outputFolder, environments, restFunc)
		}
	case "w":
		switch file.APIType() {
		case constant.APITypeBlueprint:
			err = writeBlueprint(file, sourceFolder, outputFolder, environments, restFunc)
		case constant.APITypeTeam:
			err = writeTeam(file, sourceFolder, outputFolder, environments, restFunc)
		}
	case "m":
		switch file.APIType() {
		case constant.APITypeTeam:
			err = migrateTeam(file, outputFolder, environments, restFunc)
		default:
			err = fmt.Errorf("invalid action for %s", file.APIType())
		}
	}

	if err != nil {
		return model.Output{Code: constant.OutputError, Debug: err.Error()}
	}

	return model.Output{Code: constant.OutputSuccess, Debug: constant.Undefined}
}

type result struct {
	ID  int    `json:"_internalId"`
	URL string `json:"_self"`
}

func (r *result) getURL(env, context string) (string, error) {
	restURL, err := url.Parse(r.URL)
	if err != nil {
		return "", err
	}
	envURL, err := url.Parse(env)
	if err != nil {
		return "", err
	}
	if envURL.Host != restURL.Host {
		restURL.Host = envURL.Host
	}
	return restURL.String(), nil
}
