package xog

import (
	"github.com/andreluzz/cas-xog/transform"
	"github.com/andreluzz/cas-xog/common"
	"github.com/beevik/etree"
)

func debug(index, total int, action, status, path, err string) {

	actionLabel := "Write"
	if action == "r" {
		actionLabel = "Read"
	} else if action == "m" {
		actionLabel = "Create"
	}

	color := "green"
	statusLabel := "success"
	if status == transform.OUTPUT_WARNING {
		statusLabel = "warning"
		color = "yellow"
	} else if status == transform.OUTPUT_ERROR {
		statusLabel = "error  "
		color = "red"
	}

	if err != "" {
		err = "| Debug: " + err
	}

	output[status] += 1

	common.Debug("\r[CAS-XOG][%s[%s %s]] %03d/%03d | file: %s %s", color, actionLabel, statusLabel, index, total, path, err)
}

func loadAndValidate(action, folder string, file *common.DriverFile, env *EnvType) (*etree.Document, common.XOGOutput, error) {

	if action != "w" && file.Type == common.MIGRATION {
		return nil, common.XOGOutput{Code: transform.OUTPUT_ERROR, Debug: ""}, nil
	}

	body, err := GetXMLFile(action, file, env)
	errorOutput := common.XOGOutput{Code: transform.OUTPUT_ERROR, Debug: ""}
	if err != nil {
		return nil, errorOutput, err
	}

	resp, err := common.SoapCall(body, env.URL)

	if err != nil {
		return nil, errorOutput, err
	}

	resp.IndentTabs()
	resp.WriteToFile(folder + file.Type + "/" + file.Path)

	validateOutput, err := transform.Validate(resp)

	if err != nil {
		return nil, validateOutput, err
	}

	return resp, validateOutput, nil
}