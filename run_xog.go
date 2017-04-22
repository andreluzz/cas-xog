package main

import (
	"encoding/xml"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"time"
)

type XogRead struct {
	Types []struct {
		Type  string `xml:"type,attr"`
		Value string `xml:"value"`
	} `xml:"xogtype"`
	Examples []struct {
		Type  string `xml:"type,attr"`
		Value string `xml:"value"`
	} `xml:"xogExample"`
}

type XogMenu struct {
	Code           string `xml:"code,attr"`
	Action         string `xml:"action,attr"`
	TargetPosition int    `xml:"targetPosition,attr"`
	Links          []struct {
		Code           string `xml:"code,attr"`
		TargetPosition int    `xml:"targetPosition,attr"`
	} `xml:"link"`
}

type XogUnit struct {
	Name             string `xml:"name,attr"`
	ParentName       string `xml:"parentName,attr"`
	RemoveUnitChilds bool   `xml:"removeUnitChilds,attr"`
	Remove           bool   `xml:"remove,attr"`
}

type XogViewSection struct {
	SourceSectionPosition string `xml:"sourceSectionPosition,attr"`
	TargetSectionPosition string `xml:"targetSectionPosition,attr"`
	Action                string `xml:"action,attr"`
	Attributes            []struct {
		Code         string `xml:"code,attr"`
		Column       string `xml:"column,attr"`
		Remove       bool   `xml:"remove,attr"`
		InsertBefore string `xml:"insertBefore,attr"`
	} `xml:"attribute"`
}

type XogViewAction struct {
	Code         string `xml:"code,attr"`
	Remove       bool   `xml:"remove,attr"`
	GroupCode    string `xml:"groupCode,attr"`
	InsertBefore string `xml:"insertBefore,attr"`
}

type XogDriverFile struct {
	Code            string           `xml:"code,attr"`
	Path            string           `xml:"path,attr"`
	Type            string           `xml:"type,attr"`
	ObjCode         string           `xml:"objectCode,attr"`
	SingleView      bool             `xml:"singleView,attr"`
	CopyToView      string           `xml:"copyToView,attr"`
	EnvTarget       string           `xml:"envTarget,attr"`
	IgnoreReading   bool             `xml:"ignoreReading,attr"`
	SourcePartition string           `xml:"sourcePartition,attr"`
	TargetPartition string           `xml:"targetPartition,attr"`
	OnlyStructure   bool             `xml:"onlyStructure,attr"`
	InsertBefore    string           `xml:"insertBefore,attr"`
	Sections        []XogViewSection `xml:"section"`
	Actions         []XogViewAction  `xml:"action"`
	Menus           []XogMenu        `xml:"menu"`
	Units           []XogUnit        `xml:"unit"`
	Includes        []struct {
		Type string `xml:"type,attr"`
		Code string `xml:"code,attr"`
	} `xml:"include"`
}

type XogDriver struct {
	Files []XogDriverFile `xml:"file"`
}

type XogEnv struct {
	GlobalVars []struct {
		Name  string `xml:"name,attr"`
		Value string `xml:"value,attr"`
	} `xml:"global>var"`
	Environments []struct {
		Name   string `xml:"name,attr"`
		Params []struct {
			Name  string `xml:"name,attr"`
			Value string `xml:"value,attr"`
		} `xml:"param"`
	} `xml:"environments>env"`
}

var global_env_version = "8.0"
var xog *XogDriver
var env *XogEnv
var readDefault *XogRead
var inputAction string

func main() {

	fmt.Println("")
	fmt.Println("--------------------------------------------")
	fmt.Println("#### Processing XOG Files - Version 1.0 ####")
	fmt.Println("--------------------------------------------")

	env = new(XogEnv)
	xml.Unmarshal(loadFile("xogEnv.xml"), env)

	readDefault = new(XogRead)
	xml.Unmarshal(loadFile("xogRead.xml"), readDefault)

	//read command-line parameters to use in silent mode option
	var driverPath string
	flag.StringVar(&driverPath, "xogdriver", "", "The path to the xog driver file, required for silent mode")
	var create bool
	flag.BoolVar(&create, "create", false, "Create the xog read files")
	var readEnv int
	flag.IntVar(&readEnv, "read", -1, "Read xog from environment")
	var writeEnv int
	flag.IntVar(&writeEnv, "write", -1, "Write xog to environment")
	flag.Parse()

	if driverPath != "" {
		Debug("[XOG] Silent Mode: \033[92mON\033[0m\n")
	}

	loadXogDriverFile(driverPath)

	fmt.Println("")

	if driverPath != "" {
		if create || readEnv != -1 || writeEnv != -1 {
			if create {
				fmt.Println("[XOG] Silent Mode - Create XOGs Read files ")
				scanActions("c", -1)
			}
			if readEnv != -1 {
				fmt.Println("[XOG] Silent Mode - Reading XOGs")
				scanActions("r", readEnv)
			}
			if writeEnv != -1 {
				fmt.Println("[XOG] Silent Mode - Writing XOGs")
				scanActions("w", writeEnv)
			}
		}
	} else {
		var exit = false
		for {
			exit = scanActions("", -1)
			if exit {
				break
			}
		}
	}
}

func loadXogDriverFile(silentXogDriverPathFile string) {
	//Define xog driver path
	var driverIndex = 0
	xogDriverPath := "drivers/"
	xogDriverFileName := silentXogDriverPathFile

	if silentXogDriverPathFile == "" {
		xogDriverFileList, _ := ioutil.ReadDir(xogDriverPath)

		if len(xogDriverFileList) == 0 {
			Debug("\n[XOG]\033[91mERROR\033[0m - XogDriver folders or file not found! Press any key to exit...\n")
			scanexit := ""
			fmt.Scanln(&scanexit)
			os.Exit(0)
		}

		fmt.Println("")
		fmt.Println("Available drivers:")
		for k, f := range xogDriverFileList {
			Debug("%d - %s\n", k, f.Name())
		}
		fmt.Print("Choose driver [0]: ")
		var input string = "0"
		fmt.Scanln(&input)

		var err error
		driverIndex, err = strconv.Atoi(input)

		if err != nil || driverIndex < 0 || driverIndex+1 > len(xogDriverFileList) {
			Debug("\n[XOG]\033[91mERROR\033[0m - Invalid XOG driver! Press any key to exit...\n")
			scanexit := ""
			fmt.Scanln(&scanexit)
			os.Exit(0)
		}

		xogDriverFileName = xogDriverFileList[driverIndex].Name()
	}
	xogDriverPathFile := xogDriverPath + xogDriverFileName

	xog = new(XogDriver)
	xml.Unmarshal(loadFile(xogDriverPathFile), xog)

	Debug("\n[XOG]\033[92mLoaded XOG Driver file\033[0m: %s\n", xogDriverPathFile)
}

func scanActions(silentAction string, silentEnv int) bool {
	inputAction = silentAction
	var envIndex = silentEnv
	if silentAction == "" {
		//Define action: Write, Read ou Create
		fmt.Print("Choose action (l = Load new XOG Driver, c = Create XOGs Read files, r = Read XOGs, w = Write XOGs or x = eXit): ")
		fmt.Scanln(&inputAction)

		envIndex = 0
		if inputAction == "w" || inputAction == "r" {
			//Define environment
			fmt.Println("")
			fmt.Println("Available environments:")
			for k, e := range env.Environments {
				Debug("%d - %s\n", k, e.Name)
			}
			fmt.Print("Choose environment [0]: ")
			var input string = "0"
			fmt.Scanln(&input)

			var err error
			envIndex, err = strconv.Atoi(input)

			if err != nil || envIndex < 0 || envIndex+1 > len(env.Environments) {
				Debug("\n[XOG]\033[91mERROR\033[0m - Invalid environment!\n\n")
				return false
			}
		}
	}

	start := time.Now()

	switch strings.ToLower(inputAction) {
	case "w":
		ExecuteXOG(xog, env, envIndex, "write")
	case "r":
		ExecuteXOG(xog, env, envIndex, "read")
	case "c":
		createReadFilesXOG(xog)
	case "l":
		loadXogDriverFile("")
	case "x":
		return true
	default:
		Debug("\n[XOG]\033[91mERROR\033[0m - Action not implemented!\n\n")
		return false
	}

	elapsed := time.Since(start)

	Debug("\n------------------------------")
	Debug("\n#### Processing concluded ####")
	Debug("\nExecuted in:  %s", elapsed)
	Debug("\n------------------------------\n\n")

	return false
}

func loadFile(path string) []byte {
	xmlFile, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err)
	}
	return xmlFile
}
