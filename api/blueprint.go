package api

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"

	"github.com/andreluzz/cas-xog/constant"
	"github.com/andreluzz/cas-xog/model"
	"github.com/andreluzz/cas-xog/util"
)

type blueprint struct {
	ID   int    `json:"_internalId"`
	Name string `json:"name"`
	Type struct {
		ID string `json:"id"`
	} `json:"type"`
	Sections []*blueprintSection `json:"sections"`
	Visuals  []*blueprintVisual  `json:"visuals"`
}

type blueprintVisual struct {
	ID            int    `json:"_internalId"`
	VisualID      int    `json:"visualId"`
	Sequence      int    `json:"sequence"`
	Type          string `json:"type"`
	Label         string `json:"label"`
	Category      string `json:"category"`
	ColorCode     string `json:"colorCode"`
	ResourceName  string `json:"resourceName"`
	AttributeName string `json:"attributeName"`
	BlueprintType string `json:"blueprintType"`
}

type blueprintSection struct {
	ID         int    `json:"_internalId"`
	Name       string `json:"name"`
	Sequence   int    `json:"sequence"`
	FieldsAddr struct {
		URL string `json:"_self"`
	} `json:"fields"`
	Fields []*blueprintField `json:"fieldsData"`
}

type blueprintField struct {
	ID          int    `json:"_internalId"`
	Name        string `json:"name"`
	MetadataURL string `json:"metadataURL"`
	Column      int    `json:"column"`
	Width       int    `json:"width"`
	Height      int    `json:"height"`
}

type blueprintResults struct {
	Results []struct {
		ID  int    `json:"_internalId"`
		URL string `json:"_self"`
	} `json:"_results"`
}

func readBlueprint(file *model.DriverFile, outputFolder string, environments *model.Environments, restFunc util.Rest) error {
	if file.ID == constant.Undefined {
		return errors.New("Required attribute id not found")
	}
	endpoint := environments.Source.URL + "/ppm/rest/v1/"
	response, err := restFunc(nil, endpoint+"blueprints/"+file.ID, http.MethodGet, environments.Source.AuthToken)
	if err != nil {
		return err
	}

	bp := &blueprint{}
	json.Unmarshal(response, bp)

	//read bp sections
	response, err = restFunc(nil, endpoint+"blueprints/"+file.ID+"/sections", http.MethodGet, environments.Source.AuthToken)
	if err != nil {
		return err
	}
	sections := &blueprintResults{}
	json.Unmarshal(response, sections)
	for sectionIndex, s := range sections.Results {
		response, err = restFunc(nil, s.URL, http.MethodGet, environments.Source.AuthToken)
		if err != nil {
			return err
		}
		section := &blueprintSection{}
		json.Unmarshal(response, section)
		bp.Sections = append(bp.Sections, section)

		// read bp section fields
		response, err = restFunc(nil, section.FieldsAddr.URL, http.MethodGet, environments.Source.AuthToken)
		if err != nil {
			return err
		}
		fields := &blueprintResults{}
		json.Unmarshal(response, fields)
		for _, f := range fields.Results {
			response, err = restFunc(nil, f.URL, http.MethodGet, environments.Source.AuthToken)
			if err != nil {
				return err
			}
			field := &blueprintField{}
			json.Unmarshal(response, field)
			bp.Sections[sectionIndex].Fields = append(bp.Sections[sectionIndex].Fields, field)
		}
	}

	//read bp visuals
	response, err = restFunc(nil, endpoint+"blueprints/"+file.ID+"/visuals", http.MethodGet, environments.Source.AuthToken)
	if err != nil {
		return err
	}
	visuals := &blueprintResults{}
	json.Unmarshal(response, visuals)
	for _, v := range visuals.Results {
		response, err = restFunc(nil, v.URL, http.MethodGet, environments.Source.AuthToken)
		if err != nil {
			return err
		}
		visual := &blueprintVisual{}
		json.Unmarshal(response, visual)
		bp.Visuals = append(bp.Visuals, visual)
	}

	data, _ := json.Marshal(bp)
	bpPath := outputFolder + file.Type + "/" + file.Path
	ioutil.WriteFile(bpPath, data, 0644)

	return nil
}
