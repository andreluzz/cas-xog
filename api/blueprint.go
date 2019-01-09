package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/andreluzz/cas-xog/constant"
	"github.com/andreluzz/cas-xog/model"
	"github.com/andreluzz/cas-xog/util"
)

type blueprintRead struct {
	Name       string        `json:"name"`
	ID         int           `json:"_internalId"`
	TypeObject blueprintType `json:"type"`
}

type blueprintWrite struct {
	Name string `json:"name"`
	ID   int    `json:"_internalId"`
	Type string `json:"type"`
}

type blueprintType struct {
	ID string `json:"id"`
}

func writeBlueprint(file *model.DriverFile, sourceFolder, outputFolder string, environments *model.Environments, restFunc util.Rest) error {
	path := sourceFolder + file.Type + "/" + file.Path
	jsonBytes, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}

	bp := &blueprintWrite{}

	if file.APIAction == "new" {
		endpoint := environments.Target.URL + "/ppm/rest/v1/blueprints"

		fmt.Println()
		fmt.Println(endpoint)
		fmt.Println()

		response, err := restFunc(jsonBytes, endpoint, http.MethodPost, environments.Source.AuthToken)
		if err != nil {
			return err
		}
		json.Unmarshal(response, bp)

		fmt.Println()
		fmt.Println(bp.ID)
		fmt.Println()

		path = outputFolder + file.Type + "/" + file.Path
		ioutil.WriteFile(path, response, 0644)
	}

	return nil
}

func readBlueprint(file *model.DriverFile, outputFolder string, environments *model.Environments, restFunc util.Rest) error {
	if file.ID == constant.Undefined {
		return errors.New("Required attribute id not found")
	}
	endpoint := environments.Source.URL + "/ppm/rest/v1/blueprints/" + file.ID
	response, err := restFunc(nil, endpoint, http.MethodGet, environments.Source.AuthToken)
	if err != nil {
		return err
	}

	bpRead := &blueprintRead{}
	json.Unmarshal(response, bpRead)
	bpWrite := &blueprintWrite{
		Name: bpRead.Name,
		Type: bpRead.TypeObject.ID,
	}
	data, _ := json.Marshal(bpWrite)

	path := outputFolder + file.Type + "/" + file.Path

	ioutil.WriteFile(path, data, 0644)

	return nil
}
