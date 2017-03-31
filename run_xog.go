package main

import (
	"encoding/xml"
	"fmt"
	"strconv"
	"strings"
	"time"
)

type XogRead struct {
	Types []struct {
		Type  string `xml:"type,attr"`
		Value string `xml:"value"`
	} `xml:"xogtype"`
}

type XogDriver struct {
	Files []struct {
		Code            string `xml:"code,attr"`
		Path            string `xml:"path,attr"`
		Type            string `xml:"type,attr"`
		IgnoreReading   bool   `xml:"ignoreReading,attr"`
		ObjCode         string `xml:"objectCode,attr"`
		SourcePartition string `xml:"sourcePartition,attr"`
		TargetPartition string `xml:"targetPartition,attr"`
		SingleView      bool   `xml:"singleView,attr"`
	} `xml:"files>file"`
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
var inputAction, xogDriverPath string
var helper Helper

func main() {

	fmt.Println("")
	fmt.Println("--------------------------------------------")
	fmt.Println("#### Processing XOG Files - Version 1.0 ####")
	fmt.Println("--------------------------------------------")

	env = new(XogEnv)
	xml.Unmarshal(helper.loadFile("xogEnv.xml"), env)

	readDefault = new(XogRead)
	xml.Unmarshal(helper.loadFile("xogRead.xml"), readDefault)

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
	xogDriverPath = "xogDriver.xml"
	fmt.Println("")
	fmt.Print("Enter XOG Driver path [xogDriver.xml]: ")
	fmt.Scanln(&xogDriverPath)

	xog = new(XogDriver)
	xml.Unmarshal(helper.loadFile(xogDriverPath), xog)

	fmt.Printf("\n[XOG]\033[92mLoaded XOG Driver file\033[0m: %s\n", xogDriverPath)
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
		fmt.Printf("\n[\033[91mError\033[0m] Action not implemented!\n")
	}

	elapsed := time.Since(start)

	fmt.Printf("\n------------------------------")
	fmt.Printf("\n#### Processing concluded ####")
	fmt.Printf("\nExecuted in:  %s", elapsed)
	fmt.Printf("\n------------------------------\n\n")

	return false
}
