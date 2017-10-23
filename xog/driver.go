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
	"fmt"
	"strconv"
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

func RenderDrivers() {
	var driverIndex = 1
	driverPath := "drivers/"
	driversFileList, _ := ioutil.ReadDir(driverPath)

	if len(driversFileList) == 0 {
		common.Debug("\n[CAS-XOG][red[ERROR]] - XogDriver folders or file not found! Press enter key to exit...\n")
		scanexit := ""
		fmt.Scanln(&scanexit)
		os.Exit(0)
	}

	var driversList []common.Driver
	for _, f := range driversFileList {
		driver := new(common.Driver)
		driver.Info = f
		driver.FilePath = driverPath + f.Name()
		driversList = append(driversList, *driver)
	}

	driversList = append(driversList, packagesDriversFileInfo...)

	fmt.Println("")
	fmt.Println("Available drivers:")
	for k, d := range driversList {
		if d.PackageDriver {
			common.Debug("%d - [blue[Package driver:]] %s\n", k+1, d.Info.Name())
		} else {
			common.Debug("%d - %s\n", k+1, d.Info.Name())
		}
	}
	if startInstallingPackage == 0 {
		fmt.Print("Choose driver [1] or p = Install Package: ")
	}else {
		fmt.Print("Choose driver [1]: ")
	}

	input := "1"
	fmt.Scanln(&input)

	if input == "p" && startInstallingPackage == 0 {
		startInstallingPackage = 1
		return
	}
	startInstallingPackage = -1

	var err error
	driverIndex, err = strconv.Atoi(input)

	if err != nil || driverIndex-1 < 0 || driverIndex > len(driversList) {
		common.Debug("\n[CAS-XOG][red[ERROR]] - Invalid XOG driver!\n")
		return
	}

	err = LoadDriver( driversList[driverIndex-1].FilePath)
	if err != nil {
		common.Debug("\n[CAS-XOG][red[ERROR]] - %s", err.Error())
		common.Debug("\n[CAS-XOG][red[FATAL]] - Check your driver file. Press enter key to exit...")
		scanExit := ""
		fmt.Scanln(&scanExit)
		os.Exit(0)
	}

	common.Debug("\n[CAS-XOG][blue[Loaded XOG Driver file]]: %s | Total files: [green[%d]]\n",  driversList[driverIndex-1].FilePath, len(driverXOG.Files))
}

