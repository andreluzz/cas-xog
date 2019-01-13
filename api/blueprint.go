package api

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/andreluzz/cas-xog/constant"
	"github.com/andreluzz/cas-xog/model"
	"github.com/andreluzz/cas-xog/util"
)

type blueprint struct {
	ID   int    `json:"_internalId"`
	Name string `json:"name"`
	Code string `json:"code"`
	Type struct {
		ID string `json:"id"`
	} `json:"type"`
	Sections     []*blueprintSection     `json:"sections"`
	Visuals      []*blueprintVisual      `json:"visuals"`
	ExternalApps []*blueprintExternalApp `json:"externalApps"`
}

func (bp *blueprint) getNewBlueprintBody() []byte {
	body := `{"name": "` + bp.Name + `", "type": "` + bp.Type.ID + `"}`
	return []byte(body)
}

type blueprintExternalApp struct {
	ID           int    `json:"_internalId"`
	BlueprintID  int    `json:"blueprintId"`
	BaseURL      string `json:"baseUrl"`
	ConcreteURL  string `json:"concreteUrl"`
	ReferrerURLs string `json:"referrerUrls"`
	Name         string `json:"name"`
	Visual       struct {
		ID           string `json:"id"`
		Type         string `json:"_type"`
		DisplayValue string `json:"displayValue"`
	} `json:"visualId"`
}

func (extApp *blueprintExternalApp) getNewExternalApp(bpID int) []byte {
	body := `{
				"baseUrl": "` + extApp.BaseURL + `",
				"name": "` + extApp.Name + `",
				"concreteUrl": "` + extApp.ConcreteURL + `",
				"blueprintId": ` + strconv.Itoa(bpID) + `,
				"referrerUrls": "` + extApp.ReferrerURLs + `",
				"visualId": "` + extApp.Visual.ID + `"
			}`
	return []byte(body)
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
	ExtAppName    string `json:"extAppName"`
}

func (visual *blueprintVisual) getNewVisualBody() []byte {
	body := `{
				"resourceName": "` + visual.ResourceName + `",
				"label": "` + visual.Label + `",
				"type": "` + visual.Type + `",
				"visualId": ` + strconv.Itoa(visual.VisualID) + `,
				"sequence": ` + strconv.Itoa(visual.Sequence) + `,
				"extAppName": "` + visual.ExtAppName + `",
				"blueprintType": "` + visual.BlueprintType + `",
				"attributeName": "` + visual.AttributeName + `",
				"colorCode": "` + visual.ColorCode + `",
				"category": "` + visual.Category + `"
			}`
	return []byte(body)
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

func (section *blueprintSection) getNewSectionBody(bpID int) []byte {
	body := `{
				"name": "` + section.Name + `",
				"sequence": ` + strconv.Itoa(section.Sequence) + `, 
				"blueprintId": ` + strconv.Itoa(bpID) + `
			}`
	return []byte(body)
}

type blueprintField struct {
	ID          int    `json:"_internalId"`
	Name        string `json:"name"`
	MetadataURL string `json:"metadataURL"`
	Column      int    `json:"column"`
	Width       int    `json:"width"`
	Row         int    `json:"row"`
	Height      int    `json:"height"`
}

func (field *blueprintField) getNewFieldBody(sectionID int) []byte {
	body := `{
				"name": "` + field.Name + `",
				"metadataURL": "` + field.MetadataURL + `",
				"column": ` + strconv.Itoa(field.Column) + `,
				"width": ` + strconv.Itoa(field.Width) + `,
				"sectionId": ` + strconv.Itoa(sectionID) + `,
				"row": ` + strconv.Itoa(field.Row) + `,
				"height": ` + strconv.Itoa(field.Height) + `
			}`
	return []byte(body)
}

type blueprintResults struct {
	Results []struct {
		ID  int    `json:"_internalId"`
		URL string `json:"_self"`
	} `json:"_results"`
}

type blueprintResponse struct {
	ID   int    `json:"_internalId"`
	Code string `json:"code"`
}

func writeBlueprint(file *model.DriverFile, sourceFolder, outputFolder string, environments *model.Environments, restFunc util.Rest) error {
	bpPath := sourceFolder + file.Type + "/" + file.Path
	jsonFile, err := ioutil.ReadFile(bpPath)
	if err != nil {
		return err
	}

	endpoint := environments.Target.URL + "/ppm/rest/v1/"

	bp := &blueprint{}
	json.Unmarshal(jsonFile, bp)

	if file.TargetID != constant.Undefined {
		//Get target blueprint code
		response, status, err := restFunc(nil, endpoint+"blueprints/"+file.TargetID, http.MethodGet, environments.Target.AuthToken, nil)
		if err != nil {
			return err
		}
		if status != 200 {
			return errors.New("status code " + strconv.Itoa(status) + " - response: " + string(response))
		}
		resp := &blueprintResponse{}
		json.Unmarshal(response, resp)
		//Make blueprint editable
		body := `{
			"source": "` + resp.Code + `",
			"action": "edit"
		}`
		response, status, err = restFunc([]byte(body), endpoint+"private/copyBlueprint", http.MethodPost, environments.Target.AuthToken, nil)
		if err != nil {
			return err
		}
		if status != 200 {
			return errors.New("status code " + strconv.Itoa(status) + " - response: " + string(response))
		}
		resp = &blueprintResponse{}
		json.Unmarshal(response, resp)
		bp.ID = resp.ID
		//Update blueprint
		response, status, err = restFunc(bp.getNewBlueprintBody(), endpoint+"blueprints"+"/"+strconv.Itoa(bp.ID), http.MethodPatch, environments.Target.AuthToken, nil)
		if err != nil {
			return err
		}
		if status != 200 {
			return errors.New("status code " + strconv.Itoa(status) + " - response: " + string(response))
		}
		//Delete editable blueprint content
		err = deleteBlueprintContent(endpoint, strconv.Itoa(bp.ID), environments.Target.AuthToken, restFunc)
		if err != nil {
			return err
		}
	} else {
		response, status, err := restFunc(bp.getNewBlueprintBody(), endpoint+"blueprints", http.MethodPost, environments.Target.AuthToken, nil)
		if err != nil {
			return err
		}
		if status != 200 {
			return errors.New("status code " + strconv.Itoa(status) + " - response: " + string(response))
		}

		json.Unmarshal(response, bp)
	}

	//post sections
	url := endpoint + "blueprints/" + strconv.Itoa(bp.ID) + "/sections"
	for _, s := range bp.Sections {
		response, status, err := restFunc(s.getNewSectionBody(bp.ID), url, http.MethodPost, environments.Target.AuthToken, nil)
		if err != nil {
			return err
		}
		if status != 200 {
			return errors.New("status code " + strconv.Itoa(status) + " - response: " + string(response))
		}
		resp := &blueprintResponse{}
		json.Unmarshal(response, resp)

		for _, f := range s.Fields {
			response, status, err := restFunc(f.getNewFieldBody(resp.ID), url+"/"+strconv.Itoa(resp.ID)+"/fields", http.MethodPost, environments.Target.AuthToken, nil)
			if err != nil {
				return err
			}
			if status != 200 {
				return errors.New("status code " + strconv.Itoa(status) + " - response: " + string(response))
			}
		}
	}

	//post visuals
	url = endpoint + "blueprints/" + strconv.Itoa(bp.ID) + "/visuals"
	for _, v := range bp.Visuals {
		response, status, err := restFunc(v.getNewVisualBody(), url, http.MethodPost, environments.Target.AuthToken, nil)
		if err != nil {
			return err
		}
		if status != 200 {
			return errors.New("status code " + strconv.Itoa(status) + " - response: " + string(response))
		}
	}

	//post externalApps
	url = endpoint + "externalApps"
	for _, e := range bp.ExternalApps {
		response, status, err := restFunc(e.getNewExternalApp(bp.ID), url, http.MethodPost, environments.Target.AuthToken, nil)
		if err != nil {
			return err
		}
		if status != 200 {
			return errors.New("status code " + strconv.Itoa(status) + " - response: " + string(response))
		}
	}

	if file.TargetID != constant.Undefined {
		//publish edited blueprint
		body := `{"mode": "PUBLISHED"}`
		response, status, err := restFunc([]byte(body), endpoint+"blueprints/"+strconv.Itoa(bp.ID), http.MethodPut, environments.Target.AuthToken, nil)
		if err != nil {
			return err
		}
		if status != 200 {
			return errors.New("status code " + strconv.Itoa(status) + " - response: " + string(response))
		}
	}
	return nil
}

func readBlueprint(file *model.DriverFile, outputFolder string, environments *model.Environments, restFunc util.Rest) error {
	if file.ID == constant.Undefined {
		return errors.New("Required attribute id not found")
	}
	endpoint := environments.Source.URL + "/ppm/rest/v1/"
	response, status, err := restFunc(nil, endpoint+"blueprints/"+file.ID, http.MethodGet, environments.Source.AuthToken, nil)
	if err != nil {
		return err
	}
	if status != 200 {
		return errors.New("status code " + strconv.Itoa(status) + " - response: " + string(response))
	}

	bp := &blueprint{}
	json.Unmarshal(response, bp)

	//read bp sections
	response, status, err = restFunc(nil, endpoint+"blueprints/"+file.ID+"/sections", http.MethodGet, environments.Source.AuthToken, nil)
	if err != nil {
		return err
	}
	if status != 200 {
		return errors.New("status code " + strconv.Itoa(status) + " - response: " + string(response))
	}
	sections := &blueprintResults{}
	json.Unmarshal(response, sections)
	for sectionIndex, s := range sections.Results {
		response, status, err = restFunc(nil, s.URL, http.MethodGet, environments.Source.AuthToken, nil)
		if err != nil {
			return err
		}
		if status != 200 {
			return errors.New("status code " + strconv.Itoa(status) + " - response: " + string(response))
		}
		section := &blueprintSection{}
		json.Unmarshal(response, section)
		bp.Sections = append(bp.Sections, section)

		// read bp section fields
		response, status, err = restFunc(nil, section.FieldsAddr.URL, http.MethodGet, environments.Source.AuthToken, nil)
		if err != nil {
			return err
		}
		if status != 200 {
			return errors.New("status code " + strconv.Itoa(status) + " - response: " + string(response))
		}
		fields := &blueprintResults{}
		json.Unmarshal(response, fields)
		for _, f := range fields.Results {
			response, status, err = restFunc(nil, f.URL, http.MethodGet, environments.Source.AuthToken, nil)
			if err != nil {
				return err
			}
			if status != 200 {
				return errors.New("status code " + strconv.Itoa(status) + " - response: " + string(response))
			}
			field := &blueprintField{}
			json.Unmarshal(response, field)
			bp.Sections[sectionIndex].Fields = append(bp.Sections[sectionIndex].Fields, field)
		}
	}

	//read bp visuals
	response, status, err = restFunc(nil, endpoint+"blueprints/"+file.ID+"/visuals", http.MethodGet, environments.Source.AuthToken, nil)
	if err != nil {
		return err
	}
	if status != 200 {
		return errors.New("status code " + strconv.Itoa(status) + " - response: " + string(response))
	}
	visuals := &blueprintResults{}
	json.Unmarshal(response, visuals)
	for _, v := range visuals.Results {
		response, status, err = restFunc(nil, v.URL, http.MethodGet, environments.Source.AuthToken, nil)
		if err != nil {
			return err
		}
		if status != 200 {
			return errors.New("status code " + strconv.Itoa(status) + " - response: " + string(response))
		}
		visual := &blueprintVisual{}
		json.Unmarshal(response, visual)
		bp.Visuals = append(bp.Visuals, visual)
	}

	//read bp external apps
	param := make(map[string]string)
	param["filter"] = "(blueprintId = " + file.ID + ")"
	response, status, err = restFunc(nil, endpoint+"externalApps", http.MethodGet, environments.Source.AuthToken, param)
	if err != nil {
		return err
	}
	if status != 200 {
		return errors.New("status code " + strconv.Itoa(status) + " - response: " + string(response))
	}
	externalApps := &blueprintResults{}
	json.Unmarshal(response, externalApps)

	for _, e := range externalApps.Results {
		response, status, err = restFunc(nil, e.URL, http.MethodGet, environments.Source.AuthToken, nil)
		if err != nil {
			return err
		}
		if status != 200 {
			return errors.New("status code " + strconv.Itoa(status) + " - response: " + string(response))
		}
		externalApp := &blueprintExternalApp{}
		json.Unmarshal(response, externalApp)
		bp.ExternalApps = append(bp.ExternalApps, externalApp)
	}

	data, _ := json.Marshal(bp)
	bpPath := outputFolder + file.Type + "/" + file.Path
	ioutil.WriteFile(bpPath, util.JSONAvoidEscapeText(data), 0644)

	return nil
}

func deleteBlueprintContent(endpoint, bpID, token string, restFunc util.Rest) error {
	//delete sections
	response, status, err := restFunc(nil, endpoint+"blueprints/"+bpID+"/sections", http.MethodGet, token, nil)
	if err != nil {
		return err
	}
	if status != 200 {
		return errors.New("status code " + strconv.Itoa(status) + " - response: " + string(response))
	}
	sections := &blueprintResults{}
	json.Unmarshal(response, sections)
	for _, s := range sections.Results {
		response, status, err := restFunc(nil, s.URL, http.MethodDelete, token, nil)
		if err != nil {
			return err
		}
		if status != 200 {
			return errors.New("status code " + strconv.Itoa(status) + " - response: " + string(response))
		}
	}
	//delete visuals
	response, status, err = restFunc(nil, endpoint+"blueprints/"+bpID+"/visuals", http.MethodGet, token, nil)
	if err != nil {
		return err
	}
	if status != 200 {
		return errors.New("status code " + strconv.Itoa(status) + " - response: " + string(response))
	}
	visuals := &blueprintResults{}
	json.Unmarshal(response, visuals)
	for _, v := range visuals.Results {
		response, status, err := restFunc(nil, v.URL, http.MethodDelete, token, nil)
		if err != nil {
			return err
		}
		if status != 200 {
			return errors.New("status code " + strconv.Itoa(status) + " - response: " + string(response))
		}
	}
	//delete externalApps
	param := make(map[string]string)
	param["filter"] = "(blueprintId = " + bpID + ")"
	response, status, err = restFunc(nil, endpoint+"externalApps", http.MethodGet, token, param)
	if err != nil {
		return err
	}
	if status != 200 {
		return errors.New("status code " + strconv.Itoa(status) + " - response: " + string(response))
	}
	externalApps := &blueprintResults{}
	json.Unmarshal(response, externalApps)

	for _, e := range externalApps.Results {
		response, status, err := restFunc(nil, e.URL, http.MethodDelete, token, nil)
		if err != nil {
			return err
		}
		if status != 200 {
			return errors.New("status code " + strconv.Itoa(status) + " - response: " + string(response))
		}
	}
	return nil
}
