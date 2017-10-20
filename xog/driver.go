package xog

import (
	"encoding/xml"
	"errors"
	"github.com/andreluzz/cas-xog/common"
	"github.com/andreluzz/cas-xog/migration"
	"github.com/andreluzz/cas-xog/transform"
	"github.com/beevik/etree"
	"io/ioutil"
	"os"
	"time"
)

var driverXOG *common.Driver
var driverPath string

func LoadDriver(path string) error {
	xmlFile, err := ioutil.ReadFile(path)
	if err != nil {
		return errors.New("Error loading driver file: " + err.Error())
	}
	driverPath = path
	driverXOG = new(common.Driver)
	xml.Unmarshal(xmlFile, driverXOG)
	return nil
}

var output map[string]int

func ProcessDriverFiles(action string) {
	start := time.Now()

	output = map[string]int{transform.OUTPUT_SUCCESS: 0, transform.OUTPUT_WARNING: 0, transform.OUTPUT_ERROR: 0}

	common.Debug("\n------------------------------------------------------------------")
	common.Debug("\n[blue[Initiated at]]: %s", start.Format("Mon _2 Jan 2006 - 15:04:05"))
	common.Debug("\nProcessing driver: %s", driverPath)
	common.Debug("\n------------------------------------------------------------------\n")

	var env *EnvType

	if action == "r" {
		os.RemoveAll(common.FOLDER_READ)
		os.MkdirAll(common.FOLDER_READ, os.ModePerm)
		os.RemoveAll(common.FOLDER_WRITE)
		os.MkdirAll(common.FOLDER_WRITE, os.ModePerm)
		env = SourceEnv.copyEnv()
	} else if action == "w" {
		os.RemoveAll(common.FOLDER_DEBUG)
		os.MkdirAll(common.FOLDER_DEBUG, os.ModePerm)
		env = TargetEnv.copyEnv()
	} else if action == "m" {
		os.RemoveAll(common.FOLDER_MIGRATION)
		os.MkdirAll(common.FOLDER_MIGRATION, os.ModePerm)
	}

	for i, f := range driverXOG.Files {
		folder := common.FOLDER_WRITE
		actionLabel := "Reading"
		if action == "r" {
			_, dirErr := os.Stat(common.FOLDER_READ + f.Type)
			if os.IsNotExist(dirErr) {
				_ = os.Mkdir(common.FOLDER_READ+f.Type, os.ModePerm)
			}
		} else if action == "w" {
			folder = common.FOLDER_DEBUG
			actionLabel = "Writing"
		} else if action == "m" {
			folder = common.FOLDER_MIGRATION
			actionLabel = "Creating"
		}

		common.Debug("\n[CAS-XOG][blue[%s]] %03d/%03d | file: %s", actionLabel, i+1, len(driverXOG.Files), f.Path)

		if f.IgnoreReading && action == "r" {
			debug(i+1, len(driverXOG.Files), action, transform.OUTPUT_WARNING, f.Path, "File reading ignored")
			continue
		}

		if action == "m" && f.Type != common.MIGRATION {
			debug(i+1, len(driverXOG.Files), action, transform.OUTPUT_WARNING, f.Path, "Use action 'r' to this type("+f.Type+") of file")
			continue
		} else if action == "r" && f.Type == common.MIGRATION {
			debug(i+1, len(driverXOG.Files), action, transform.OUTPUT_WARNING, f.Path, "Use action 'm' to this type("+f.Type+") of file")
			continue
		}

		//check if target folder type dir exists
		_, dirErr := os.Stat(folder + f.Type)
		if os.IsNotExist(dirErr) {
			_ = os.Mkdir(folder+f.Type, os.ModePerm)
		}

		resp, validateOutput, err := loadAndValidate(action, folder, &f, env)

		if err != nil {
			debug(i+1, len(driverXOG.Files), action, validateOutput.Code, f.Path, err.Error())
			continue
		}

		if action == "r" {
			var aux *etree.Document
			var auxFile common.DriverFile
			var auxEnv *EnvType
			loadAuxFile := false

			switch f.Type {
			case common.PROCESS:
				if f.CopyPermissions != "" {
					loadAuxFile = true
					auxEnv = env.copyEnv()
					auxFile = common.DriverFile{Code: f.CopyPermissions, Path: "aux_" + f.CopyPermissions + ".xml", Type: common.PROCESS}
				}
			case common.VIEW:
				if f.Code != "*" {
					loadAuxFile = true
					auxEnv = TargetEnv.copyEnv()
					partition := f.SourcePartition
					if f.TargetPartition != "" {
						partition = f.TargetPartition
					}
					auxFile = common.DriverFile{Code: f.Code, ObjCode: f.ObjCode, Path: "aux_" + f.Path + ".xml", SourcePartition: partition, Type: common.VIEW}
				}
			}

			if loadAuxFile {
				aux, _, err = loadAndValidate(action, folder, &auxFile, auxEnv)
				if err != nil {
					debug(i+1, len(driverXOG.Files), action, transform.OUTPUT_ERROR, f.Path, "[Auxiliary XOG] "+err.Error())
					continue
				}
			}

			err := transform.Process(resp, aux, f)
			if err != nil {
				debug(i+1, len(driverXOG.Files), action, transform.OUTPUT_ERROR, f.Path, err.Error())
				continue
			}

			if f.ExportToExcel {
				err := migration.ExportInstancesToExcel(resp, f)
				if err != nil {
					debug(i+1, len(driverXOG.Files), action, transform.OUTPUT_ERROR, f.Path, err.Error())
					continue
				}
			}
		}

		if action == "m" {
			resp, err = migration.ReadDataFromExcel(f)
			if err != nil {
				debug(i+1, len(driverXOG.Files), action, transform.OUTPUT_ERROR, f.Path, err.Error())
				continue
			}
		}

		nikuDataBusElement := resp.FindElement("//NikuDataBus")
		if nikuDataBusElement != nil {
			resp.SetRoot(nikuDataBusElement)
		}

		resp.IndentTabs()
		resp.WriteToFile(folder + f.Type + "/" + f.Path)

		actionLabel = "Read"
		if action == "w" {
			actionLabel = "Write"
		}
		if action == "m" {
			actionLabel = "Create"
		}

		debug(i+1, len(driverXOG.Files), action, validateOutput.Code, f.Path, validateOutput.Debug)
	}

	elapsed := time.Since(start)

	TargetEnv.logout()
	SourceEnv.logout()

	common.Debug("\n\n------------------------------------------------------------------")
	common.Debug("\nStats: total = %d | failure = %d | success = %d | warning = %d", len(driverXOG.Files), output[transform.OUTPUT_ERROR], output[transform.OUTPUT_SUCCESS], output[transform.OUTPUT_WARNING])
	common.Debug("\n[blue[Concluded in]]: %.3f seconds", elapsed.Seconds())
	common.Debug("\n------------------------------------------------------------------\n")
}

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
