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

type teamResults struct {
	Results []result `json:"_results"`
}

type team struct {
	ID              int                   `json:"_internalId"`
	Code            string                `json:"code"`
	Name            string                `json:"name"`
	Active          bool                  `json:"isActive"`
	TeamAllocations []*teamDefAllocations `json:"teamdefallocations"`
}

type teamDefAllocations struct {
	TeamID     int     `json:"teamId"`
	ResourceID int     `json:"resourceId"`
	Allocation float64 `json:"allocation"`
}

type teamDefAllocationsResults struct {
	Results []result `json:"_results"`
}

type result struct {
	ID  int    `json:"_internalId"`
	URL string `json:"_self"`
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
	response, status, err := restFunc(nil, url, http.MethodGet, environments.Source.AuthToken, environments.Source.Proxy, nil)
	if err != nil {
		return err
	}
	if status != 200 {
		return errors.New("status code " + strconv.Itoa(status) + " - response: " + string(response))
	}

	tr := &teamResults{}
	json.Unmarshal(response, tr)

	teams := []team{}

	for _, t := range tr.Results {
		// GET Team details
		response, status, err = restFunc(nil, t.URL, http.MethodGet, environments.Source.AuthToken, environments.Source.Proxy, nil)
		if err != nil {
			return err
		}
		if status != 200 {
			return errors.New("status code " + strconv.Itoa(status) + " - response: " + string(response))
		}

		team := team{}
		json.Unmarshal(response, &team)

		// GET Allocations
		response, status, err = restFunc(nil, t.URL+"/teamdefallocations", http.MethodGet, environments.Source.AuthToken, environments.Source.Proxy, nil)
		ar := &teamDefAllocationsResults{}
		json.Unmarshal(response, ar)

		for _, a := range ar.Results {
			response, status, err = restFunc(nil, a.URL, http.MethodGet, environments.Source.AuthToken, environments.Source.Proxy, nil)
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
		url := fmt.Sprintf("%steamdefinitions", endpoint)
		body := fmt.Sprintf(`{
			"code": "%s",
			"name": "%s",
			"isActive": %t
		}`, t.Code, t.Name, t.Active)
		response, status, err := restFunc([]byte(body), url, http.MethodPost, environments.Target.AuthToken, environments.Target.Proxy, nil)
		if err != nil {
			return err
		}
		if status != 200 {
			return errors.New("status code " + strconv.Itoa(status) + " - response: " + string(response))
		}

		newTeam := &team{}
		json.Unmarshal(response, newTeam)

		for _, a := range t.TeamAllocations {
			url := fmt.Sprintf("%steamdefinitions/%d/teamdefallocations", endpoint, newTeam.ID)
			body := fmt.Sprintf(`{
				"resourceId": %d,
				"totalAllocation": %f,
				"teamId": %d
			  }`, a.ResourceID, a.Allocation, newTeam.ID)

			response, status, err := restFunc([]byte(body), url, http.MethodPost, environments.Target.AuthToken, environments.Target.Proxy, nil)
			if err != nil {
				return err
			}
			if status != 200 {
				return errors.New("status code " + strconv.Itoa(status) + " - response: " + string(response))
			}

			res := &result{}
			json.Unmarshal(response, res)
			body = fmt.Sprintf(`{"allocation": %f}`, a.Allocation)
			response, status, err = restFunc([]byte(body), res.URL, http.MethodPut, environments.Target.AuthToken, environments.Target.Proxy, nil)
			if err != nil {
				return err
			}
			if status != 200 {
				return errors.New("status code " + strconv.Itoa(status) + " - response: " + string(response))
			}
		}
	}

	return nil
}
