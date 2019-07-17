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

type teamResults struct {
	Results []result `json:"_results"`
}

type team struct {
	ID              int                   `json:"_internalId,omitempty"`
	Code            string                `json:"code"`
	Name            string                `json:"name"`
	Active          bool                  `json:"isActive"`
	TeamAllocations []*teamDefAllocations `json:"teamdefallocations"`
}

type teamDefAllocations struct {
	ID         int     `json:"_internalId,omitempty"`
	ResourceID int     `json:"resourceId"`
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

	for index, row := range xlFile.Sheets[0].Rows {
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
				resID, err := strconv.Atoi(val)
				if err != nil {
					return fmt.Errorf("migration - team code %s column resourceID is not a valid number. Debug: %s", rowMap["code"], err.Error())
				}
				alloc := teamDefAllocations{
					ResourceID: resID,
				}
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

			teams[rowMap["code"]] = t
		}
	}

	teamSlice := []team{}
	for _, value := range teams {
		teamSlice = append(teamSlice, value)
	}

	data, _ := json.Marshal(teamSlice)
	teamsPath := outputFolder + file.Type + "/" + file.Path
	ioutil.WriteFile(teamsPath, util.JSONAvoidEscapeText(data), 0644)
	return nil
}

func readTeam(file *model.DriverFile, outputFolder string, environments *model.Environments, restFunc util.Rest) error {
	if file.Code == constant.Undefined {
		return errors.New("Required attribute code not found")
	}
	endpoint := environments.Source.URL + constant.APIEndpoint

	filter := ""
	if file.Code != "*" {
		filter = fmt.Sprintf("?filter=(code =  '%s')", file.Code)
	}

	url := fmt.Sprintf("%steamdefinitions%s", endpoint, filter)
	response, status, err := restFunc(nil, url, http.MethodGet, environments.Source.AuthToken, environments.Source.Proxy, environments.Source.Cookie, nil)
	if err != nil {
		return err
	}
	if status != 200 {
		return fmt.Errorf("status code: %d | response: %s | url: %s", status, string(response), url)
	}

	tr := &teamResults{}
	json.Unmarshal(response, tr)

	teams := []team{}

	for _, t := range tr.Results {
		// GET Team details
		urlString, err := t.getURL(environments.Source.URL)
		if err != nil {
			return err
		}
		response, status, err = restFunc(nil, urlString, http.MethodGet, environments.Source.AuthToken, environments.Source.Proxy, environments.Source.Cookie, nil)
		if err != nil {
			return err
		}
		if status != 200 {
			return fmt.Errorf("status code: %d | response: %s | url: %s", status, string(response), urlString)
		}

		team := team{}
		json.Unmarshal(response, &team)

		// GET Allocations
		response, status, err = restFunc(nil, urlString+"/teamdefallocations", http.MethodGet, environments.Source.AuthToken, environments.Source.Proxy, environments.Source.Cookie, nil)
		ar := &teamDefAllocationsResults{}
		json.Unmarshal(response, ar)

		for _, a := range ar.Results {
			urlString, err := a.getURL(environments.Source.URL)
			if err != nil {
				return err
			}
			response, status, err = restFunc(nil, urlString, http.MethodGet, environments.Source.AuthToken, environments.Source.Proxy, environments.Source.Cookie, nil)
			teamAllocation := &teamDefAllocations{}
			json.Unmarshal(response, teamAllocation)
			team.TeamAllocations = append(team.TeamAllocations, teamAllocation)
		}

		teams = append(teams, team)
	}

	data, _ := json.Marshal(teams)
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

	endpoint := environments.Target.URL + constant.APIEndpoint

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
	url := fmt.Sprintf("%steamdefinitions", endpoint)
	body := fmt.Sprintf(`{
			"code": "%s",
			"name": "%s",
			"isActive": %t
		}`, t.Code, t.Name, t.Active)
	response, status, err := restFunc([]byte(body), url, http.MethodPost, environments.Target.AuthToken, environments.Target.Proxy, environments.Target.Cookie, nil)
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
	for _, a := range t.TeamAllocations {
		url := fmt.Sprintf("%steamdefinitions/%d/teamdefallocations", endpoint, newTeam.ID)
		body := fmt.Sprintf(`{
				"resourceId": %d,
				"totalAllocation": %f,
				"teamId": %d
			  }`, a.ResourceID, a.Allocation, newTeam.ID)

		response, status, err := restFunc([]byte(body), url, http.MethodPost, environments.Target.AuthToken, environments.Target.Proxy, environments.Target.Cookie, nil)
		if err != nil {
			return err
		}
		if status != 200 {
			return fmt.Errorf("status code: %d | response: %s | url: %s", status, string(response), url)
		}

		res := &result{}
		json.Unmarshal(response, res)
		body = fmt.Sprintf(`{"allocation": %f}`, a.Allocation)
		urlString, err := res.getURL(environments.Target.URL)
		if err != nil {
			return err
		}

		response, status, err = restFunc([]byte(body), urlString, http.MethodPut, environments.Target.AuthToken, environments.Target.Proxy, environments.Target.Cookie, nil)
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

	filter := fmt.Sprintf("?filter=(code =  '%s')", t.Code)

	url := fmt.Sprintf("%steamdefinitions%s", endpoint, filter)
	response, status, err := restFunc(nil, url, http.MethodGet, environments.Target.AuthToken, environments.Target.Proxy, environments.Target.Cookie, nil)
	if err != nil {
		return err
	}
	if status != 200 {
		return fmt.Errorf("status code: %d | response: %s | url: %s", status, string(response), url)
	}

	tr := &teamResults{}
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
	response, status, err = restFunc([]byte(body), url, http.MethodPut, environments.Target.AuthToken, environments.Target.Proxy, environments.Target.Cookie, nil)
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
	response, status, err = restFunc(nil, url+"/teamdefallocations", http.MethodGet, environments.Target.AuthToken, environments.Target.Proxy, environments.Target.Cookie, nil)
	ar := &teamDefAllocationsResults{}
	json.Unmarshal(response, ar)

	teamCurrentAllocations := []*teamDefAllocations{}
	for _, a := range ar.Results {
		urlString, err := a.getURL(environments.Target.URL)
		if err != nil {
			return err
		}
		response, status, err = restFunc(nil, urlString, http.MethodGet, environments.Target.AuthToken, environments.Target.Proxy, environments.Target.Cookie, nil)
		teamAllocation := &teamDefAllocations{}
		json.Unmarshal(response, teamAllocation)
		teamCurrentAllocations = append(teamCurrentAllocations, teamAllocation)
		index := getIndex(teamAllocation, t.TeamAllocations)
		if index != -1 {
			updateTeam := t.TeamAllocations[index]
			url := fmt.Sprintf("%steamdefinitions/%d/teamdefallocations/%d", endpoint, newTeam.ID, teamAllocation.ID)
			body = fmt.Sprintf(`{"allocation": %f}`, updateTeam.Allocation)

			response, status, err := restFunc([]byte(body), url, http.MethodPut, environments.Target.AuthToken, environments.Target.Proxy, environments.Target.Cookie, nil)
			if err != nil {
				return err
			}
			if status != 200 {
				return fmt.Errorf("status code: %d | response: %s | url: %s", status, string(response), url)
			}
		}
	}

	for _, team := range t.TeamAllocations {
		index := getIndex(team, teamCurrentAllocations)
		if index == -1 {
			url := fmt.Sprintf("%steamdefinitions/%d/teamdefallocations", endpoint, newTeam.ID)
			body := fmt.Sprintf(`{
				"resourceId": %d,
				"totalAllocation": %f,
				"teamId": %d
			  }`, team.ResourceID, team.Allocation, newTeam.ID)

			response, status, err := restFunc([]byte(body), url, http.MethodPost, environments.Target.AuthToken, environments.Target.Proxy, environments.Target.Cookie, nil)
			if err != nil {
				return err
			}
			if status != 200 {
				return fmt.Errorf("status code: %d | response: %s | url: %s", status, string(response), url)
			}

			res := &result{}
			json.Unmarshal(response, res)
			body = fmt.Sprintf(`{"allocation": %f}`, team.Allocation)
			urlString, err := res.getURL(environments.Target.URL)
			if err != nil {
				return err
			}

			response, status, err = restFunc([]byte(body), urlString, http.MethodPut, environments.Target.AuthToken, environments.Target.Proxy, environments.Target.Cookie, nil)
			if err != nil {
				return err
			}
			if status != 200 {
				return fmt.Errorf("status code: %d | response: %s | url: %s", status, string(response), urlString)
			}
		}
	}

	return nil
}

func getIndex(alloc *teamDefAllocations, slice []*teamDefAllocations) int {
	for i, a := range slice {
		if a.ResourceID == alloc.ResourceID {
			return i
		}
	}
	return -1
}
