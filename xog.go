package main

import (
	"bytes"
	"fmt"
	"github.com/beevik/etree"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

func createReadFilesXOG(xog *XogDriver) {

	os.RemoveAll("_read/")
	os.MkdirAll("_read/", os.ModePerm)

	xogTypes := make(map[string]string)

	for _, xog := range readDefault.Types {
		xogTypes[xog.Type] = xog.Value
	}

	var i int = 1
	total := len(xog.Files)

	for _, xogfile := range xog.Files {

		status, path := createReadFile(xogfile, xogTypes, false)
		if xogfile.Type == "views" && xogfile.TargetPartition != "" && xogfile.SingleView {
			status, _ = createReadFile(xogfile, xogTypes, true)
		}

		Debug("\n[CAS-XOG]Created read file %03d/%03d - %s | to file: %s", i, total, status, path)

		i += 1
	}
	fmt.Println("")
}

func createReadFile(xogfile XogDriverFile, xogTypes map[string]string, createTargetViewReadXOG bool) (string, string) {
	status := ""
	path := "_read/" + xogfile.Type + "/" + xogfile.Path
	if createTargetViewReadXOG {
		path = "_read/" + xogfile.Type + "/target_" + xogfile.Path
	}

	if xogfile.IgnoreReading {
		status = "\033[93mIGNORED\033[0m"
	} else {
		//check if dir exists
		_, dir_err := os.Stat("_read/" + xogfile.Type)
		if os.IsNotExist(dir_err) {
			_ = os.Mkdir("_read/"+xogfile.Type, os.ModePerm)
		}

		doc := etree.NewDocument()
		if err := doc.ReadFromString(xogTypes[xogfile.Type]); err != nil {
			panic(err)
		}

		code := xogfile.Code
		if xogfile.Type == "lookups" {
			code = strings.ToUpper(code)
		}

		filterCode := doc.FindElement("//Filter[@name='code']")
		if filterCode != nil {
			filterCode.SetText(code)
		}
		filterObjectCode := doc.FindElement("//Filter[@name='object_code']")
		if filterObjectCode != nil {
			switch xogfile.Type {
			case "objects":
				filterObjectCode.SetText(xogfile.Code)
			default:
				filterObjectCode.SetText(xogfile.ObjCode)
			}
		}
		filterInstanceCode := doc.FindElement("//Filter[@name='instanceCode']")
		if filterInstanceCode != nil {
			filterInstanceCode.SetText(code)
		}
		filterInstanceObjectCode := doc.FindElement("//Filter[@name='objectCode']")
		if filterInstanceObjectCode != nil {
			filterInstanceObjectCode.SetText(xogfile.ObjCode)
		}
		filterPartition := doc.FindElement("//Filter[@name='partition_code']")
		if filterPartition != nil {
			if xogfile.SourcePartition == "" {
				e := filterPartition.Parent()
				e.RemoveChild(filterPartition)
			} else {
				if createTargetViewReadXOG {
					filterPartition.SetText(xogfile.TargetPartition)
				} else {
					filterPartition.SetText(xogfile.SourcePartition)
				}
			}
		}

		status = "\033[92mSUCCESS\033[0m"

		doc.Indent(4)
		err := doc.WriteToFile(path)
		if err != nil {
			status = "\033[91mFAILURE\033[0m"
		}

		if xogfile.Type == "views" || xogfile.Type == "customObjectInstances" {
			if xogfile.ObjCode == "" {
				status = "\033[91mFAILURE\033[0m"
			}
		}
	}

	return status, path
}

func ExecuteXOG(xog *XogDriver, env *XogEnv, envIndex int, action string) {
	var inputDir, outputDir string

	if action == "write" {
		inputDir = "_write/"
		outputDir = "_debug/"
	} else {
		inputDir = "_read/"
		outputDir = "_write/"
	}

	var i int = 1
	total := len(xog.Files)

	os.RemoveAll(outputDir)
	os.MkdirAll(outputDir, os.ModePerm)

	for _, xogfile := range xog.Files {

		xogfileCompletePath := xogfile.Type + "/" + xogfile.Path

		inputPath := inputDir + xogfileCompletePath
		outputPath := outputDir + xogfileCompletePath

		//When IgnoreReading is true write from folder 'extra'
		if xogfile.IgnoreReading && action == "write" {
			inputPath = "extra/" + xogfileCompletePath
		}

		//When OnlyStructure is true write and read the file with 'pre_' in the name
		if xogfile.Type == "lookups" && xogfile.OnlyStructure {
			if action == "write" {
				inputPath = inputDir + xogfile.Type + "/pre_" + xogfile.Path
			} else {
				outputPath = outputDir + xogfile.Type + "/pre_" + xogfile.Path
			}
		}

		if xogfile.IgnoreReading && action != "write" {
			Debug("\n[CAS-XOG]Readed %03d/%03d - \033[93mIGNORED\033[0m | transform: NONE | to file: %s", i, total, outputPath)
		} else {
			//check if dir exists
			_, err1 := os.Stat(outputDir + xogfile.Type)
			if os.IsNotExist(err1) {
				_ = os.Mkdir(outputDir+xogfile.Type, os.ModePerm)
			}

			validatedWriteEnvironment := true
			//validate if envTarget is defined for the same environment select to write the XOGs
			if action == "write" && xogfile.EnvTarget != "" && xogfile.EnvTarget != strconv.Itoa(envIndex) {
				validatedWriteEnvironment = false
			}

			tempOutputPath := ""

			if validatedWriteEnvironment {
				execCommand(envIndex, inputPath, outputPath)
			}

			status, statusMessage := Validate(outputPath)
			transform := "NONE"

			if status && action != "write" {
				if Transform(xogfile, outputPath) {
					transform = "\033[96mTRUE\033[0m"
				}
				if xogfile.EnvTarget != "" {
					//read file from target environment
					tempOutputPath = outputDir + xogfile.Type + "/temp_" + xogfile.Path
					if xogfile.TargetPartition != "" {
						inputPath = inputDir + xogfile.Type + "/target_" + xogfile.Path
					}

					targetEnvironment, _ := strconv.Atoi(xogfile.EnvTarget)
					execCommand(targetEnvironment, inputPath, tempOutputPath)
					//Transform view to include the new attributes
					Transform(xogfile, tempOutputPath)

					//Merge menus
					if xogfile.Type == "menus" && len(xogfile.Menus) > 0 {
						_, statusMessage = MergeMenus(xogfile, outputPath, tempOutputPath)
					}

					//Merge views
					if xogfile.Type == "views" && xogfile.SingleView {
						_, statusMessage = MergeViews(xogfile, outputPath, tempOutputPath)
					}

					//Merge obs
					if xogfile.Type == "obs" {
						_, statusMessage = MergeOBS(xogfile, outputPath, tempOutputPath)
						transform = "\033[96mTRUE\033[0m"
					}

					os.Remove(tempOutputPath)
				}
			}

			if !validatedWriteEnvironment {
				//Transform views - general - trying to write attributes readed from a different target environment (envTarget)
				statusMessage = "\033[91mER-TVG1\033[0m"
			}

			if action != "write" {
				Debug("\n[CAS-XOG]Readed %03d/%03d - %s | transform: %s | to file: %s", i, total, statusMessage, transform, outputPath)
			} else {
				Debug("\n[CAS-XOG]Writed %03d/%03d - %s | to file: %s", i, total, statusMessage, outputPath)
			}
		}
		i += 1
	}
	fmt.Println("")
}

func execCommand(envIndex int, inputPath string, outputPath string) {
	var xog_path = env.GlobalVars[0].Value

	var args [5]string

	for k, param := range env.Environments[envIndex].Params {
		args[k] = param.Value
	}

	var out, stderr bytes.Buffer

	cmd := exec.Command(xog_path, "-username", args[0], "-password", args[1], "-servername", args[2], "-portnumber", args[3], "-sslenabled", args[4], "-input", inputPath, "-output", outputPath)
	cmd.Stdin = strings.NewReader("some input")
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		fmt.Println(fmt.Sprint(err) + ": " + stderr.String())
		panic(err)
	}
}
