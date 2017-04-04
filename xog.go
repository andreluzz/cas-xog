package main

import (
	"bytes"
	"fmt"
	"github.com/beevik/etree"
	"os"
	"os/exec"
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

		status := ""
		path := "_read/" + xogfile.Type + "/" + xogfile.Path

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
					filterPartition.SetText(xogfile.SourcePartition)
				}
			}

			status = "\033[92mSUCCESS\033[0m"

			doc.Indent(4)
			err := doc.WriteToFile(path)
			if err != nil {
				status = "\033[91mFAILURE\033[0m"
			}

		}

		fmt.Printf("\n[XOG]Created read file %03d/%03d - %s | to file: %s", i, total, status, path)

		i += 1
	}
	fmt.Println("")
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

		//Para arquivos que estão sendo ignorados dos processos de leitura vamos escrever utilizando os arquivos da pasta "_extra"
		if xogfile.IgnoreReading && action == "write" {
			inputPath = "_extra/" + xogfileCompletePath
		}

		if xogfile.IgnoreReading && action != "write" {
			fmt.Printf("\n[XOG]Readed %03d/%03d - \033[93mIGNORED\033[0m | transform: NONE | to file: %s", i, total, outputPath)
		} else {
			//check if dir exists
			_, err1 := os.Stat(outputDir + xogfile.Type)
			if os.IsNotExist(err1) {
				_ = os.Mkdir(outputDir+xogfile.Type, os.ModePerm)
			}

			tempOutputPath := ""

			execCommand(envIndex, inputPath, outputPath)
			if xogfile.Type == "views" && len(xogfile.Includes) > 0 {
				tempOutputPath = outputDir + xogfile.Type + "/temp_" + xogfile.Path
				execCommand(xogfile.ViewEnvTarget, inputPath, tempOutputPath)
				Transform(xogfile, tempOutputPath)
			}

			status, statusMessage := Validate(outputPath)
			transform := "NONE"

			if status && action != "write" {
				if Transform(xogfile, outputPath) {
					transform = "\033[96mTRUE\033[0m"
				}
				if xogfile.Type == "views" && len(xogfile.Includes) > 0 {
					_, statusMessage = MergeViews(xogfile, outputPath, tempOutputPath)
					os.Remove(tempOutputPath)
				}
			}

			if action != "write" {
				fmt.Printf("\n[XOG]Readed %03d/%03d - %s | transform: %s | to file: %s", i, total, statusMessage, transform, outputPath)
			} else {
				fmt.Printf("\n[XOG]Writed %03d/%03d - %s | to file: %s", i, total, statusMessage, outputPath)
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
