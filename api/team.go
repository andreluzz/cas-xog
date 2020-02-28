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
	"github.com/tealeg/xlsx"
)

type team struct {
	ID              int                   `json:"_internalId,omitempty"`
	Code            string                `json:"code"`
	Name            string                `json:"name"`
	Active          bool                  `json:"isActive"`
	TeamAllocations []*teamDefAllocations `json:"teamdefallocations,omitempty"`
}

type teamDefAllocations struct {
	ID         int `json:"_internalId,omitempty"`
	ResourceID struct {
		ID string `json:"id"`
	} `json:"resourceId"`
	Allocation float64 `json:"allocation"`
}

type teamDefAllocationsResults struct {
	Results []result `json:"_results"`
}

func migrateTeam(file *model.DriverFile, outputFolder string, environments *model.Environments, restFunc util.Rest) error {
	xlFile, err := xlsx.OpenFile(util.ReplacePathSeparatorByOS(file.ExcelFile))
	if err != nil {
		return fmt.Errorf("migration - error opening excel. Debug: %s", err.Error())
	}

	excelStartRowIndex := 0
	if file.ExcelStartRow != constant.Undefined {
		excelStartRowIndex, err = strconv.Atoi(file.ExcelStartRow)
		if err != nil {
			return fmt.Errorf("migration - tag 'startRow' not a number. Debug: %s", err.Error())
		}
		excelStartRowIndex--
	}

	teams := make(map[string]team, len(xlFile.Sheets[0].Rows))

	excelEndRowIndex := 0
	if file.ExcelStartRow != constant.Undefined {
		excelEndRowIndex, err = strconv.Atoi(file.ExcelEndRow)
		if err != nil {
			return fmt.Errorf("migration - tag 'endRow' not a number. Debug: %s", err.Error())
		}
	}

	for index, row := range xlFile.Sheets[0].Rows {
		if excelEndRowIndex != 0 && excelEndRowIndex == index {
			break
		}
		if index >= excelStartRowIndex {
			rowMap := make(map[string]string, len(file.MatchExcel))
			for _, match := range file.MatchExcel {
				rowMap[match.AttributeName] = row.Cells[match.Col-1].String()
			}

			t, ok := teams[rowMap["code"]]
			if !ok {
				t = team{
					Name:   rowMap["name"],
					Code:   rowMap["code"],
					Active: true,
				}
			}

			if val, ok := rowMap["resourceId"]; ok {
				if val != "" {
					if err != nil {
						return fmt.Errorf("migration - team code %s column resourceID (%s) is not a valid number. Debug: %s", rowMap["code"], val, err.Error())
					}
					alloc := teamDefAllocations{}
					alloc.ResourceID.ID = val
					if val, ok := rowMap["allocation"]; ok {
						if allocationValue, err := strconv.ParseFloat(val, 64); err == nil {
							alloc.Allocation = allocationValue
						} else {
							return fmt.Errorf("migration - team code %s column Allocation is not a valid number. Debug: %s", rowMap["code"], err.Error())
						}
					} else {
						alloc.Allocation = 1
					}
					t.TeamAllocations = append(t.TeamAllocations, &alloc)
				}
			}

			teams[rowMap["code"]] = t
		}
	}

	teamSlice := []team{}
	for _, value := range teams {
		teamSlice = append(teamSlice, value)
	}

	data, _ := json.MarshalIndent(teamSlice, "", "    ")
	teamsPath := outputFolder + file.Type + "/" + file.Path
	ioutil.WriteFile(teamsPath, util.JSONAvoidEscapeText(data), 0644)
	return nil
}

func readTeam(file *model.DriverFile, outputFolder string, environments *model.Environments, restFunc util.Rest) error {
	if file.Code == constant.Undefined {
		return errors.New("Required attribute code not found")
	}
	endpoint := environments.Target.URL + environments.Target.API.Context + constant.APIEndpoint

	sourceConfig := util.APIConfig{
		Client: environments.Source.API.Client,
		Cookie: environments.Source.Cookie,
		Proxy:  environments.Source.Proxy,
	}

	sourceConfig.Token = environments.Source.API.Token
	if environments.Source.AuthToken != "" {
		sourceConfig.Token = environments.Source.AuthToken
	}

	filter := ""
	if file.Code != "*" {
		filter = fmt.Sprintf("?filter=(code =  '%s')", file.Code)
	}

	url := fmt.Sprintf("%steamdefinitions%s", endpoint, filter)
	sourceConfig.Endpoint = url
	sourceConfig.Method = http.MethodGet
	response, status, err := restFunc(nil, sourceConfig, nil)
	if err != nil {
		return err
	}
	if status != 200 {
		return fmt.Errorf("status code: %d | response: %s | url: %s", status, string(response), url)
	}

	tr := &results{}
	json.Unmarshal(response, tr)

	teams := []team{}

	for _, t := range tr.Results {
		// GET Team details
		urlString, err := t.getURL(environments.Source.URL, environments.Source.API.Context)
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

		team := team{}
		json.Unmarshal(response, &team)

		// GET Allocations
		sourceConfig.Endpoint = urlString + "/teamdefallocations"
		sourceConfig.Method = http.MethodGet
		response, status, err = restFunc(nil, sourceConfig, nil)
		ar := &teamDefAllocationsResults{}
		json.Unmarshal(response, ar)

		for _, a := range ar.Results {
			urlString, err := a.getURL(environments.Source.URL, environments.Source.API.Context)
			if err != nil {
				return err
			}
			sourceConfig.Endpoint = urlString
			sourceConfig.Method = http.MethodGet
			response, status, err = restFunc(nil, sourceConfig, nil)
			teamAllocation := &teamDefAllocations{}
			json.Unmarshal(response, teamAllocation)
			team.TeamAllocations = append(team.TeamAllocations, teamAllocation)
		}

		teams = append(teams, team)
	}

	data, _ := json.MarshalIndent(teams, "", "    ")
	teamsPath := outputFolder + file.Type + "/" + file.Path
	ioutil.WriteFile(teamsPath, util.JSONAvoidEscapeText(data), 0644)
	return nil
}

func writeTeam(file *model.DriverFile, sourceFolder, outputFolder string, environments *model.Environments, restFunc util.Rest) error {
	tmPath := sourceFolder + file.Type + "/" + file.Path
	jsonFile, err := ioutil.ReadFile(tmPath)
	if err != nil {
		return err
	}

	endpoint := environments.Target.URL + environments.Target.API.Context + constant.APIEndpoint

	tm := []team{}
	json.Unmarshal(jsonFile, &tm)

	for _, t := range tm {
		if file.Action == "update" {
			if err := updateTeam(t, endpoint, environments, restFunc); err != nil {
				return err
			}
		} else {
			if err := createTeam(t, endpoint, environments, restFunc); err != nil {
				return err
			}
		}
	}

	return nil
}

func createTeam(t team, endpoint string, environments *model.Environments, restFunc util.Rest) error {

	targetConfig := util.APIConfig{
		Client: environments.Target.API.Client,
		Cookie: environments.Target.Cookie,
		Proxy:  environments.Target.Proxy,
	}

	targetConfig.Token = environments.Target.API.Token
	if environments.Target.AuthToken != "" {
		targetConfig.Token = environments.Target.AuthToken
	}

	url := fmt.Sprintf("%steamdefinitions", endpoint)
	body := fmt.Sprintf(`{
			"code": "%s",
			"name": "%s",
			"isActive": %t
		}`, t.Code, t.Name, t.Active)

	targetConfig.Endpoint = url
	targetConfig.Method = http.MethodPost
	response, status, err := restFunc([]byte(body), targetConfig, nil)
	if err != nil {
		return err
	}
	if status == 400 {
		fmt.Printf("\nstatus code: %d | Code: %s | response: %s | url: %s\n", status, t.Code, string(response), url)
		return nil
	}
	if status != 200 {
		return fmt.Errorf("status code: %d | Code: %s | response: %s | url: %s", status, t.Code, string(response), url)
	}

	newTeam := &team{}
	err = json.Unmarshal(response, newTeam)
	if err != nil {
		return fmt.Errorf("status code: %d | response: %s | url: %s | error: %s", status, string(response), url, err.Error())
	}
	for _, a := range t.TeamAllocations {
		url := fmt.Sprintf("%steamdefinitions/%d/teamdefallocations", endpoint, newTeam.ID)
		body := fmt.Sprintf(`{
				"resourceId": %s,
				"totalAllocation": %f,
				"teamId": %d
			  }`, a.ResourceID.ID, a.Allocation, newTeam.ID)

		targetConfig.Endpoint = url
		targetConfig.Method = http.MethodPost
		response, status, err := restFunc([]byte(body), targetConfig, nil)
		if err != nil {
			return err
		}
		if status == 400 {
			fmt.Printf("\nstatus code: %d | resourceId: %s | response: %s | url: %s", status, a.ResourceID.ID, string(response), url)
			continue
		}
		if status != 200 {
			return fmt.Errorf("status code: %d | resourceId: %s | response: %s | url: %s", status, a.ResourceID.ID, string(response), url)
		}

		res := &result{}
		json.Unmarshal(response, res)
		body = fmt.Sprintf(`{"allocation": %f}`, a.Allocation)
		urlString, err := res.getURL(environments.Target.URL, environments.Target.API.Context)
		if err != nil {
			return err
		}

		targetConfig.Endpoint = urlString
		targetConfig.Method = http.MethodPut
		response, status, err = restFunc([]byte(body), targetConfig, nil)
		if err != nil {
			return err
		}
		if status != 200 {
			return fmt.Errorf("status code: %d | response: %s | url: %s", status, string(response), urlString)
		}
	}
	return nil
}

func updateTeam(t team, endpoint string, environments *model.Environments, restFunc util.Rest) error {

	targetConfig := util.APIConfig{
		Client: environments.Target.API.Client,
		Cookie: environments.Target.Cookie,
		Proxy:  environments.Target.Proxy,
	}

	targetConfig.Token = environments.Target.API.Token
	if environments.Target.AuthToken != "" {
		targetConfig.Token = environments.Target.AuthToken
	}

	filter := fmt.Sprintf("?filter=(code =  '%s')", t.Code)

	url := fmt.Sprintf("%steamdefinitions%s", endpoint, filter)
	targetConfig.Endpoint = url
	targetConfig.Method = http.MethodGet
	response, status, err := restFunc(nil, targetConfig, nil)
	if err != nil {
		return err
	}
	if status != 200 {
		return fmt.Errorf("status code: %d | response: %s | url: %s", status, string(response), url)
	}

	tr := &results{}
	err = json.Unmarshal(response, tr)
	if err != nil {
		return fmt.Errorf("status code: %d | response: %s | url: %s | error: %s", status, string(response), url, err.Error())
	}
	if len(tr.Results) <= 0 {
		return fmt.Errorf("invalid team code - %s", t.Code)
	}

	teamdefinitionsInternalID := tr.Results[0].ID

	url = fmt.Sprintf("%steamdefinitions/%d", endpoint, teamdefinitionsInternalID)
	body := fmt.Sprintf(`{
			"name": "%s",
			"isActive": %t
		}`, t.Name, t.Active)

	targetConfig.Endpoint = url
	targetConfig.Method = http.MethodPut
	response, status, err = restFunc([]byte(body), targetConfig, nil)
	if err != nil {
		return err
	}
	if status != 200 {
		return fmt.Errorf("status code: %d | response: %s | url: %s", status, string(response), url)
	}

	newTeam := &team{}
	err = json.Unmarshal(response, newTeam)
	if err != nil {
		return fmt.Errorf("status code: %d | response: %s | url: %s | error: %s", status, string(response), url, err.Error())
	}

	// GET Allocations
	targetConfig.Endpoint = url + "/teamdefallocations"
	targetConfig.Method = http.MethodGet
	response, status, err = restFunc(nil, targetConfig, nil)
	ar := &teamDefAllocationsResults{}
	json.Unmarshal(response, ar)

	teamCurrentAllocations := []*teamDefAllocations{}
	for _, a := range ar.Results {
		urlString, err := a.getURL(environments.Target.URL, environments.Target.API.Context)
		if err != nil {
			return err
		}

		targetConfig.Endpoint = urlString
		targetConfig.Method = http.MethodGet
		response, status, err = restFunc(nil, targetConfig, nil)
		teamAllocation := &teamDefAllocations{}
		json.Unmarshal(response, teamAllocation)
		teamCurrentAllocations = append(teamCurrentAllocations, teamAllocation)
		index := getIndex(teamAllocation, t.TeamAllocations)
		if index != -1 {
			updateTeam := t.TeamAllocations[index]
			url := fmt.Sprintf("%steamdefinitions/%d/teamdefallocations/%d", endpoint, newTeam.ID, teamAllocation.ID)
			body = fmt.Sprintf(`{"allocation": %f}`, updateTeam.Allocation)

			targetConfig.Endpoint = url
			targetConfig.Method = http.MethodGet
			response, status, err = restFunc([]byte(body), targetConfig, nil)
			if err != nil {
				return err
			}
			if status != 200 {
				return fmt.Errorf("Get teamdefallocations[%s] status code: %d | response: %s | url: %s", t.Code, status, string(response), url)
			}
		}
	}

	for _, team := range t.TeamAllocations {
		index := getIndex(team, teamCurrentAllocations)
		if index == -1 {
			url := fmt.Sprintf("%steamdefinitions/%d/teamdefallocations", endpoint, newTeam.ID)
			body := fmt.Sprintf(`{
				"resourceId": %s,
				"totalAllocation": %f,
				"teamId": %d
			  }`, team.ResourceID.ID, team.Allocation, newTeam.ID)

			targetConfig.Endpoint = url
			targetConfig.Method = http.MethodPost
			response, status, err = restFunc([]byte(body), targetConfig, nil)
			if err != nil {
				return err
			}
			if status == 400 {
				fmt.Printf("\n[%s] status code: %d | response: %s | url: %s", t.Code, status, string(response), url)
				continue
			} else if status != 200 {
				return fmt.Errorf("Post teamdefallocations [%s] status code: %d | response: %s | url: %s", t.Code, status, string(response), url)
			}

			res := &result{}
			json.Unmarshal(response, res)
			body = fmt.Sprintf(`{"allocation": %f}`, team.Allocation)
			urlString, err := res.getURL(environments.Target.URL, environments.Target.API.Context)
			if err != nil {
				return err
			}

			targetConfig.Endpoint = urlString
			targetConfig.Method = http.MethodPut
			response, status, err = restFunc([]byte(body), targetConfig, nil)
			if err != nil {
				return err
			}
			if status != 200 {
				return fmt.Errorf("[%s] status code: %d | response: %s | url: %s", t.Code, status, string(response), url)
			}
		}
	}

	return nil
}

func getIndex(alloc *teamDefAllocations, slice []*teamDefAllocations) int {
	for i, a := range slice {
		if a.ResourceID.ID == alloc.ResourceID.ID {
			return i
		}
	}
	return -1
}
