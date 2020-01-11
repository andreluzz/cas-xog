package api

import (
	"encoding/json"
	"errors"
	"fmt"
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
	ID         int               `json:"_internalId"`
	Name       string            `json:"name"`
	Sequence   int               `json:"sequence"`
	FieldsAddr result            `json:"fields"`
	Fields     []*blueprintField `json:"fieldsData"`
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

type targetVisualsResults struct {
	Results []struct {
		ID   int    `json:"_internalId"`
		Name string `json:"attributeName"`
	} `json:"_results"`
}

type blueprintResults struct {
	Results []result `json:"_results"`
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

	endpoint := environments.Target.URL + environments.Target.API.Context + constant.APIEndpoint

	targetConfig := util.APIConfig{
		Client: environments.Target.API.Client,
		Cookie: environments.Target.Cookie,
		Proxy:  environments.Target.Proxy,
	}

	targetConfig.Token = environments.Target.API.Token
	if environments.Target.AuthToken != "" {
		targetConfig.Token = environments.Target.AuthToken
	}

	bp := &blueprint{}
	json.Unmarshal(jsonFile, bp)

	if file.TargetID != constant.Undefined {
		//Get target blueprint code
		targetConfig.Endpoint = endpoint + "private/blueprints/" + file.TargetID
		targetConfig.Method = http.MethodGet
		response, status, err := restFunc(nil, targetConfig, nil)
		if err != nil {
			return err
		}
		if status != 200 {
			return fmt.Errorf("status code: %d | response: %s | url: %s", status, string(response), targetConfig.Endpoint)
		}
		resp := &blueprintResponse{}
		json.Unmarshal(response, resp)
		//Make blueprint editable
		body := `{
			"source": "` + resp.Code + `",
			"action": "edit"
		}`
		targetConfig.Endpoint = endpoint + "private/copyBlueprint"
		targetConfig.Method = http.MethodPost
		response, status, err = restFunc([]byte(body), targetConfig, nil)
		if err != nil {
			return err
		}
		if status != 200 {
			return fmt.Errorf("status code: %d | response: %s | url: %s", status, string(response), targetConfig.Endpoint)
		}
		resp = &blueprintResponse{}
		json.Unmarshal(response, resp)
		bp.ID = resp.ID
		//Update blueprint
		targetConfig.Endpoint = endpoint + "private/blueprints/" + strconv.Itoa(bp.ID)
		targetConfig.Method = http.MethodPatch
		response, status, err = restFunc(bp.getNewBlueprintBody(), targetConfig, nil)
		if err != nil {
			return err
		}
		if status != 200 {
			return fmt.Errorf("status code: %d | response: %s | url: %s", status, string(response), targetConfig.Endpoint)
		}
		//Delete editable blueprint content
		err = deleteBlueprintContent(environments, strconv.Itoa(bp.ID), restFunc)
		if err != nil {
			return err
		}
	} else {
		targetConfig.Endpoint = endpoint + "private/blueprints"
		targetConfig.Method = http.MethodPost
		response, status, err := restFunc(bp.getNewBlueprintBody(), targetConfig, nil)
		if err != nil {
			return err
		}
		if status != 200 {
			return fmt.Errorf("status code: %d | response: %s | url: %s", status, string(response), targetConfig.Endpoint)
		}

		json.Unmarshal(response, bp)
	}

	//post sections
	url := endpoint + "private/blueprints/" + strconv.Itoa(bp.ID) + "/sections"
	for _, s := range bp.Sections {
		targetConfig.Endpoint = url
		targetConfig.Method = http.MethodPost
		response, status, err := restFunc(s.getNewSectionBody(bp.ID), targetConfig, nil)
		if err != nil {
			return err
		}
		if status != 200 {
			return fmt.Errorf("status code: %d | response: %s | url: %s", status, string(response), url)
		}
		resp := &blueprintResponse{}
		json.Unmarshal(response, resp)

		for _, f := range s.Fields {
			targetConfig.Endpoint = url + "/" + strconv.Itoa(resp.ID) + "/fields"
			targetConfig.Method = http.MethodPost
			response, status, err := restFunc(f.getNewFieldBody(resp.ID), targetConfig, nil)
			if err != nil {
				return err
			}
			if status != 200 {
				return fmt.Errorf("status code: %d | response: %s | url: %s", status, string(response), url+"/"+strconv.Itoa(resp.ID)+"/fields")
			}
		}
	}

	//post visuals
	for _, v := range bp.Visuals {
		targetConfig.Endpoint = endpoint + "private/blueprints/" + strconv.Itoa(bp.ID) + "/visuals"
		targetConfig.Method = http.MethodPost
		response, status, err := restFunc(v.getNewVisualBody(), targetConfig, nil)
		if err != nil {
			return err
		}
		if status != 200 {
			return fmt.Errorf("status code: %d | response: %s | url: %s", status, string(response), targetConfig.Endpoint)
		}
	}

	//post externalApps
	for _, e := range bp.ExternalApps {
		targetConfig.Endpoint = endpoint + "private/externalApps"
		targetConfig.Method = http.MethodPost
		response, status, err := restFunc(e.getNewExternalApp(bp.ID), targetConfig, nil)
		if err != nil {
			return err
		}
		if status != 200 {
			return fmt.Errorf("status code: %d | response: %s | url: %s", status, string(response), targetConfig.Endpoint)
		}
	}

	//publish edited blueprint
	body := `{"mode": "PUBLISHED"}`
	targetConfig.Endpoint = endpoint + "private/blueprints/" + strconv.Itoa(bp.ID)
	targetConfig.Method = http.MethodPut
	response, status, err := restFunc([]byte(body), targetConfig, nil)
	if err != nil {
		return err
	}
	if status != 200 {
		return fmt.Errorf("status code: %d | response: %s | url: %s", status, string(response), endpoint+"private/blueprints/"+strconv.Itoa(bp.ID))
	}

	return nil
}

func readBlueprint(file *model.DriverFile, outputFolder string, environments *model.Environments, restFunc util.Rest) error {
	if file.ID == constant.Undefined {
		return errors.New("Required attribute id not found")
	}

	endpoint := environments.Source.URL + environments.Source.API.Context + constant.APIEndpoint

	sourceConfig := util.APIConfig{
		Client: environments.Source.API.Client,
		Cookie: environments.Source.Cookie,
		Proxy:  environments.Source.Proxy,
	}
	sourceConfig.Token = environments.Source.API.Token
	if environments.Source.AuthToken != "" {
		sourceConfig.Token = environments.Source.AuthToken
	}

	targetConfig := util.APIConfig{
		Client: environments.Target.API.Client,
		Cookie: environments.Target.Cookie,
		Proxy:  environments.Target.Proxy,
	}
	targetConfig.Token = environments.Target.API.Token
	if environments.Target.AuthToken != "" {
		targetConfig.Token = environments.Target.AuthToken
	}

	sourceConfig.Endpoint = endpoint + "private/blueprints/" + file.ID
	sourceConfig.Method = http.MethodGet
	response, status, err := restFunc(nil, sourceConfig, nil)
	if err != nil {
		return err
	}
	if status != 200 {
		return fmt.Errorf("status code: %d | response: %s | url: %s", status, string(response), endpoint+"private/blueprints/"+file.ID)
	}

	bp := &blueprint{}
	json.Unmarshal(response, bp)

	//read bp sections
	sourceConfig.Endpoint = endpoint + "private/blueprints/" + file.ID + "/sections"
	sourceConfig.Method = http.MethodGet
	response, status, err = restFunc(nil, sourceConfig, nil)
	if err != nil {
		return err
	}
	if status != 200 {
		return fmt.Errorf("status code: %d | response: %s | url: %s", status, string(response), endpoint+"private/blueprints/"+file.ID+"/sections")
	}
	sections := &blueprintResults{}
	json.Unmarshal(response, sections)
	for sectionIndex, s := range sections.Results {
		urlString, err := s.getURL(environments.Source.URL, environments.Source.API.Context)
		if err != nil {
			return err
		}
		sourceConfig.Endpoint = urlString
		sourceConfig.Method = http.MethodGet
		response, status, err = restFunc(nil, sourceConfig, nil)
		if err != nil {
			return err
		}
		if status != 200 {
			return fmt.Errorf("status code: %d | response: %s | url: %s", status, string(response), urlString)
		}
		section := &blueprintSection{}
		json.Unmarshal(response, section)
		bp.Sections = append(bp.Sections, section)

		// read bp section fields
		urlString, err = section.FieldsAddr.getURL(environments.Source.URL, environments.Source.API.Context)
		if err != nil {
			return err
		}
		sourceConfig.Endpoint = urlString
		sourceConfig.Method = http.MethodGet
		response, status, err = restFunc(nil, sourceConfig, nil)
		if err != nil {
			return err
		}
		if status != 200 {
			return fmt.Errorf("status code: %d | response: %s | url: %s", status, string(response), urlString)
		}
		fields := &blueprintResults{}
		json.Unmarshal(response, fields)
		for _, f := range fields.Results {
			urlString, err = f.getURL(environments.Source.URL, environments.Source.API.Context)
			if err != nil {
				return err
			}
			sourceConfig.Endpoint = urlString
			sourceConfig.Method = http.MethodGet
			response, status, err = restFunc(nil, sourceConfig, nil)
			if err != nil {
				return err
			}
			if status != 200 {
				return fmt.Errorf("status code: %d | response: %s | url: %s", status, string(response), urlString)
			}
			field := &blueprintField{}
			json.Unmarshal(response, field)
			bp.Sections[sectionIndex].Fields = append(bp.Sections[sectionIndex].Fields, field)
		}
	}

	//read bp visuals
	sourceConfig.Endpoint = endpoint + "private/blueprints/" + file.ID + "/visuals"
	sourceConfig.Method = http.MethodGet
	response, status, err = restFunc(nil, sourceConfig, nil)
	if err != nil {
		return err
	}
	if status != 200 {
		return fmt.Errorf("status code: %d | response: %s | url: %s", status, string(response), endpoint+"private/blueprints/"+file.ID+"/visuals")
	}
	visuals := &blueprintResults{}
	json.Unmarshal(response, visuals)
	targetVisuals := make(map[string]int)
	if len(visuals.Results) > 0 {
		targetEndpoint := environments.Target.URL + environments.Target.API.Context + constant.APIEndpoint
		// read target environment available modules
		targetConfig.Endpoint = targetEndpoint + "private/availableVisuals?filter=((blueprintType = '" + bp.Type.ID + "') and ( category = 'MODULE'))"
		targetConfig.Method = http.MethodGet
		response, status, err = restFunc(nil, targetConfig, nil)
		if err != nil {
			return err
		}
		if status != 200 {
			return fmt.Errorf("status code: %d | response: %s | url: %s", status, string(response), targetConfig.Endpoint)
		}
		modules := &targetVisualsResults{}
		json.Unmarshal(response, modules)
		for _, m := range modules.Results {
			targetVisuals[m.Name] = m.ID
		}
		// read target environment available visuals
		targetConfig.Endpoint = targetEndpoint + "private/availableVisuals?filter=((blueprintType = '" + bp.Type.ID + "') and ( category = 'VISUAL'))"
		targetConfig.Method = http.MethodGet
		response, status, err = restFunc(nil, targetConfig, nil)
		if err != nil {
			return err
		}
		if status != 200 {
			return fmt.Errorf("status code: %d | response: %s | url: %s", status, string(response), targetConfig.Endpoint)
		}
		visuals := &targetVisualsResults{}
		json.Unmarshal(response, visuals)
		for _, v := range visuals.Results {
			targetVisuals[v.Name] = v.ID
		}
	}

	for _, v := range visuals.Results {
		urlString, err := v.getURL(environments.Source.URL, environments.Source.API.Context)
		if err != nil {
			return err
		}
		sourceConfig.Endpoint = urlString
		sourceConfig.Method = http.MethodGet
		response, status, err = restFunc(nil, sourceConfig, nil)
		if err != nil {
			return err
		}
		if status != 200 {
			return fmt.Errorf("status code: %d | response: %s | url: %s", status, string(response), urlString)
		}
		visual := &blueprintVisual{}
		json.Unmarshal(response, visual)
		//validate visual exists in target environment
		targetVisualID, ok := targetVisuals[visual.AttributeName]
		if ok {
			visual.VisualID = targetVisualID
			bp.Visuals = append(bp.Visuals, visual)
		}
	}

	//read bp external apps
	param := make(map[string]string)
	param["filter"] = "(blueprintId = " + file.ID + ")"
	sourceConfig.Endpoint = endpoint + "private/externalApps"
	sourceConfig.Method = http.MethodGet
	response, status, err = restFunc(nil, sourceConfig, param)
	if err != nil {
		return err
	}
	if status != 200 {
		return fmt.Errorf("status code: %d | response: %s | url: %s", status, string(response), endpoint+"private/externalApps")
	}
	externalApps := &blueprintResults{}
	json.Unmarshal(response, externalApps)

	for _, e := range externalApps.Results {
		urlString, err := e.getURL(environments.Source.URL, environments.Source.API.Context)
		if err != nil {
			return err
		}
		sourceConfig.Endpoint = urlString
		sourceConfig.Method = http.MethodGet
		response, status, err = restFunc(nil, sourceConfig, nil)
		if err != nil {
			return err
		}
		if status != 200 {
			return fmt.Errorf("status code: %d | response: %s | url: %s", status, string(response), urlString)
		}
		externalApp := &blueprintExternalApp{}
		json.Unmarshal(response, externalApp)
		bp.ExternalApps = append(bp.ExternalApps, externalApp)
	}

	data, _ := json.MarshalIndent(bp, "", "    ")
	bpPath := outputFolder + file.Type + "/" + file.Path
	ioutil.WriteFile(bpPath, util.JSONAvoidEscapeText(data), 0644)

	return nil
}

func deleteBlueprintContent(envs *model.Environments, bpID string, restFunc util.Rest) error {
	env := envs.Target
	endpoint := env.URL + env.API.Context + constant.APIEndpoint

	config := util.APIConfig{
		Client: env.API.Client,
		Cookie: env.Cookie,
		Proxy:  env.Proxy,
	}

	config.Token = env.API.Token
	if env.AuthToken != "" {
		config.Token = env.AuthToken
	}
	//delete sections
	config.Endpoint = endpoint + "private/blueprints/" + bpID + "/sections"
	config.Method = http.MethodGet
	response, status, err := restFunc(nil, config, nil)
	if err != nil {
		return err
	}
	if status != 200 {
		return errors.New("status code " + strconv.Itoa(status) + " - response: " + string(response))
	}
	sections := &blueprintResults{}
	json.Unmarshal(response, sections)
	for _, s := range sections.Results {
		urlString, err := s.getURL(env.URL, env.API.Context)
		if err != nil {
			return err
		}
		config.Endpoint = urlString
		config.Method = http.MethodDelete
		response, status, err = restFunc(nil, config, nil)
		if err != nil {
			return err
		}
		if status != 200 {
			return errors.New("status code " + strconv.Itoa(status) + " - response: " + string(response))
		}
	}
	//delete visuals
	config.Endpoint = endpoint + "private/blueprints/" + bpID + "/visuals"
	config.Method = http.MethodGet
	response, status, err = restFunc(nil, config, nil)
	if err != nil {
		return err
	}
	if status != 200 {
		return errors.New("status code " + strconv.Itoa(status) + " - response: " + string(response))
	}
	visuals := &blueprintResults{}
	json.Unmarshal(response, visuals)
	for _, v := range visuals.Results {
		urlString, err := v.getURL(env.URL, env.API.Context)
		if err != nil {
			return err
		}
		config.Endpoint = urlString
		config.Method = http.MethodDelete
		response, status, err := restFunc(nil, config, nil)
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
	config.Endpoint = endpoint + "private/externalApps"
	config.Method = http.MethodGet
	response, status, err = restFunc(nil, config, param)
	if err != nil {
		return err
	}
	if status != 200 {
		return errors.New("status code " + strconv.Itoa(status) + " - response: " + string(response))
	}

	externalApps := &blueprintResults{}
	json.Unmarshal(response, externalApps)
	for _, e := range externalApps.Results {
		urlString, err := e.getURL(env.URL, env.API.Context)
		if err != nil {
			return err
		}
		config.Endpoint = urlString
		config.Method = http.MethodGet
		response, status, err = restFunc(nil, config, nil)
		if err != nil {
			return err
		}
		if status != 200 {
			return errors.New("status code " + strconv.Itoa(status) + " - response: " + string(response))
		}
	}
	return nil
}
