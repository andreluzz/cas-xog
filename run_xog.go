package main

import (
	"encoding/xml"
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
	Sections        []XogViewSection `xml:"section"`
	Actions         []XogViewAction  `xml:"action"`
	Menus           []XogMenu        `xml:"menu"`
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

	loadXogDriverFile()

	fmt.Println("")

	var exit = false
	for {
		exit = scanActions()
		if exit {
			break
		}
	}
}

func loadXogDriverFile() {
	//Define xog driver path
	var driverIndex = 0
	xogDriverPath := "drivers/"
	xogDriverFileList, _ := ioutil.ReadDir(xogDriverPath)

	if len(xogDriverFileList) == 0 {
		fmt.Printf("\n[XOG]\033[91mERROR\033[0m - Drivers not found!\n")
		os.Exit(0)
	}

	fmt.Println("")
	fmt.Println("Available drivers:")
	for k, f := range xogDriverFileList {
		fmt.Printf("%d - %s\n", k, f.Name())
	}
	fmt.Print("Choose driver [0]: ")
	var input string = "0"
	fmt.Scanln(&input)

	driverIndex, _ = strconv.Atoi(input)

	xogDriverFileName := xogDriverFileList[driverIndex].Name()
	xogDriverPathFile := xogDriverPath + xogDriverFileName

	xog = new(XogDriver)
	xml.Unmarshal(loadFile(xogDriverPathFile), xog)

	fmt.Printf("\n[XOG]\033[92mLoaded XOG Driver file\033[0m: %s\n", xogDriverPathFile)
}

func scanActions() bool {
	inputAction = ""
	//Define action: Write, Read ou Create
	fmt.Print("Choose action (l = Load new XOG Driver, c = Create XOGs Read files, r = Read XOGs, w = Write XOGs or x = eXit): ")
	fmt.Scanln(&inputAction)

	var envIndex = 0
	if inputAction == "w" || inputAction == "r" {
		//Define environment
		fmt.Println("")
		fmt.Println("Available environments:")
		for k, e := range env.Environments {
			fmt.Printf("%d - %s\n", k, e.Name)
		}
		fmt.Print("Choose environment [0]: ")
		var input string = "0"
		fmt.Scanln(&input)

		envIndex, _ = strconv.Atoi(input)
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
		loadXogDriverFile()
	case "x":
		return true
	default:
		fmt.Printf("\n[XOG]\033[91mERROR\033[0m - Action not implemented!\n")
	}

	elapsed := time.Since(start)

	fmt.Printf("\n------------------------------")
	fmt.Printf("\n#### Processing concluded ####")
	fmt.Printf("\nExecuted in:  %s", elapsed)
	fmt.Printf("\n------------------------------\n\n")

	return false
}

func loadFile(path string) []byte {
	xmlFile, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err)
	}
	return xmlFile
}
