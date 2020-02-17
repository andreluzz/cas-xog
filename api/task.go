package api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"

	"github.com/andreluzz/cas-xog/constant"
	"github.com/andreluzz/cas-xog/model"
	"github.com/andreluzz/cas-xog/util"
	"github.com/tealeg/xlsx"
)

type project struct {
	ID    int              `json:"_internalId,omitempty"`
	Code  string           `json:"code"`
	Tasks map[string]*task `json:"tasks,omitempty"`
}

type task struct {
	ID               int                  `json:"_internalId,omitempty"`
	Name             string               `json:"name"`
	Code             string               `json:"code"`
	Start            string               `json:"startDate"`
	Finish           string               `json:"finishDate"`
	Status           string               `json:"status"`
	PercentComplete  float64              `json:"percentComplete"`
	CustomAttributes map[string]string    `json:"customAttributes,omitempty"`
	Resources        map[string]*resource `json:"resource,omitempty"`
}

type resource struct {
	ResourceID    string `json:"resource"`
	Start         string `json:"startDate"`
	Finish        string `json:"finishDate"`
	EstimateCurve struct {
		SegmentList struct {
			Segments []*segment `json:"segments,omitempty"`
		} `json:"segmentList,omitempty"`
	} `json:"estimateCurve,omitempty"`
}

type segment struct {
	Start  string  `json:"start"`
	Finish string  `json:"finish"`
	Value  float64 `json:"value"`
}

func migrateTask(file *model.DriverFile, outputFolder string, environments *model.Environments, restFunc util.Rest) error {
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

	projects := make(map[string]project, len(xlFile.Sheets[0].Rows))

	for index, row := range xlFile.Sheets[0].Rows {
		if index >= excelStartRowIndex {
			rowMap := make(map[string]string, len(file.MatchExcel))
			for _, match := range file.MatchExcel {
				index := match.Col - 1
				if index >= len(row.Cells) {
					break
				}
				cell := row.Cells[index]
				if cell.IsTime() {
					cellDate, _ := cell.GetTime(false)
					fmt.Println("Date: ", cellDate.String())
					rowMap[match.AttributeName] = cellDate.Format("YYYY-MM-DDTHH:MM:SS")
				}
				rowMap[match.AttributeName] = cell.String()
			}

			p, ok := projects[rowMap["project"]]
			if !ok {
				p = project{
					Code:  rowMap["project"],
					Tasks: make(map[string]*task),
				}
			}

			t, ok := p.Tasks[rowMap["code"]]
			if !ok {
				perc, _ := strconv.ParseFloat(rowMap["percentComplete"], 64)
				t = &task{
					Name:             rowMap["name"],
					Code:             rowMap["code"],
					Start:            rowMap["start"],
					Finish:           rowMap["finish"],
					PercentComplete:  perc,
					Status:           "0",
					Resources:        make(map[string]*resource),
					CustomAttributes: make(map[string]string),
				}
				if perc > 0 {
					t.Status = "1"
				}
				if perc == 1 {
					t.Status = "2"
				}
				p.Tasks[rowMap["code"]] = t
			}

			for attr, val := range rowMap {
				if isCustomAttribute(attr) {
					t.CustomAttributes[attr] = val
				}
			}

			if len(t.CustomAttributes) == 0 {
				t.CustomAttributes = nil
			}

			if rowMap["resourceId"] == "" {
				t.Resources = nil
				break
			}

			r, ok := t.Resources[rowMap["resourceId"]]
			if !ok {
				r = &resource{
					ResourceID: rowMap["resourceId"],
					Start:      rowMap["start"],
					Finish:     rowMap["finish"],
				}
				t.Resources[rowMap["resourceId"]] = r
			}

			val, _ := strconv.ParseFloat(rowMap["segmentValue"], 64)
			s := segment{
				Start:  rowMap["segmentStart"],
				Finish: rowMap["segmentFinish"],
				Value:  val,
			}

			r.EstimateCurve.SegmentList.Segments = append(r.EstimateCurve.SegmentList.Segments, &s)

			projects[rowMap["project"]] = p
		}
	}

	projectslice := []project{}
	for _, value := range projects {
		projectslice = append(projectslice, value)
	}

	data, _ := json.MarshalIndent(projectslice, "", "    ")
	taskPath := outputFolder + file.Type + "/" + file.Path
	ioutil.WriteFile(taskPath, util.JSONAvoidEscapeText(data), 0644)
	return nil
}

func writeTask(file *model.DriverFile, sourceFolder, outputFolder string, environments *model.Environments, restFunc util.Rest) error {
	tmPath := sourceFolder + file.Type + "/" + file.Path
	jsonFile, err := ioutil.ReadFile(tmPath)
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

	projects := []project{}
	json.Unmarshal(jsonFile, &projects)

	errors := []string{}

	for _, p := range projects {
		fmt.Println("endpoint: " + endpoint)
		url := endpoint + "projects?filter=(code = '" + p.Code + "')"
		targetConfig.Endpoint = url
		targetConfig.Method = http.MethodGet
		// get project id from project code
		response, status, err := restFunc(nil, targetConfig, nil)
		if err != nil {
			errors = append(errors, err.Error())
			break
		}
		if status != 200 {
			errors = append(errors, fmt.Sprintf("status code: %d | response: %s | url: %s", status, string(response), url))
			break
		}
		pr := &results{}
		json.Unmarshal(response, pr)
		if len(pr.Results) == 0 {
			errors = append(errors, fmt.Sprintf("invalid project id %s", p.Code))
			break
		}
		p.ID = pr.Results[0].ID

		for _, t := range p.Tasks {
			// check if task exists
			url := endpoint + "projects/" + strconv.Itoa(p.ID) + "/tasks?filter=(code = '" + t.Code + "')"
			targetConfig.Endpoint = url
			targetConfig.Method = http.MethodGet
			response, status, err := restFunc(nil, targetConfig, nil)
			newTask := false
			if err != nil || status != 200 {
				newTask = true
			}

			url = endpoint + "projects/" + strconv.Itoa(p.ID) + "/tasks"
			targetConfig.Endpoint = url
			targetConfig.Method = http.MethodPost

			existingTask := &results{}
			json.Unmarshal(response, existingTask)
			if !newTask && len(existingTask.Results) > 0 {
				t.ID = existingTask.Results[0].ID
				url = endpoint + "projects/" + strconv.Itoa(p.ID) + "/tasks/" + strconv.Itoa(t.ID)
				targetConfig.Endpoint = url
				targetConfig.Method = http.MethodPut
			}

			body := fmt.Sprintf(`{
				"name": "%s",
				"code": "%s",
				"startDate": "%s",
				"finishDate": "%s",
				"status": "%s",
				"percentComplete": %f`, t.Name, t.Code, t.Start, t.Finish, t.Status, t.PercentComplete)

			if len(t.CustomAttributes) > 0 {
				customAttr := ""
				for attr, val := range t.CustomAttributes {
					customAttr = fmt.Sprintf(`%s, "%s": "%s"`, customAttr, attr, val)
				}
				body = body + customAttr
			}

			body = body + "}"

			// create new task and get task id
			response, status, err = restFunc([]byte(body), targetConfig, nil)
			if err != nil {
				errors = append(errors, err.Error())
				break
			}
			if status != 200 {
				errors = append(errors, fmt.Sprintf("status code: %d | response: %s | url: %s", status, string(response), url))
				break
			}
			tr := &result{}
			json.Unmarshal(response, tr)
			t.ID = tr.ID

			url = endpoint + "projects/" + strconv.Itoa(p.ID) + "/tasks/" + strconv.Itoa(t.ID) + "/assignments"

			// assign resources to task
			for _, r := range t.Resources {
				// check resource already assigned to the task
				targetConfig.Endpoint = url + "?filter=(resource = '" + r.ResourceID + "')"
				targetConfig.Method = http.MethodGet
				response, status, err := restFunc(nil, targetConfig, nil)
				newAssignment := false
				if err != nil || status != 200 {
					newAssignment = true
				}

				targetConfig.Endpoint = url
				targetConfig.Method = http.MethodPost

				existingAssignment := &results{}
				json.Unmarshal(response, existingAssignment)
				if !newAssignment && len(existingAssignment.Results) > 0 {
					targetConfig.Endpoint = url + "/" + strconv.Itoa(existingAssignment.Results[0].ID)
					targetConfig.Method = http.MethodPut
				}

				body, _ := json.Marshal(r)
				response, status, err = restFunc(body, targetConfig, nil)
				if err != nil {
					errors = append(errors, err.Error())
				}
				if status != 200 {
					errors = append(errors, fmt.Sprintf("status code: %d | response: %s | url: %s", status, string(response), url))
				}
			}
		}
	}

	if len(errors) > 0 {
		return fmt.Errorf("%s", strings.Join(errors, "\n"))
	}

	return nil
}

func isCustomAttribute(attributeLabel string) bool {
	switch attributeLabel {
	case "project", "name", "code", "start", "finish", "percentComplete", "resourceId", "segmentStart", "segmentFinish", "segmentValue":
		return false
	default:
		return true
	}
}
