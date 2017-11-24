package xog

import (
	"encoding/xml"
	"errors"
	"fmt"
	"github.com/andreluzz/cas-xog/constant"
	"github.com/andreluzz/cas-xog/migration"
	"github.com/andreluzz/cas-xog/model"
	"github.com/andreluzz/cas-xog/transform"
	"github.com/andreluzz/cas-xog/util"
	"github.com/andreluzz/cas-xog/validate"
	"github.com/beevik/etree"
	"io/ioutil"
	"reflect"
	"sort"
	"strconv"
	"strings"
)

var driverXOG *model.Driver

func LoadDriver(path string) (int, error) {
	driverXOG = &model.Driver{}
	driverXOG.Clear()
	xmlFile, err := ioutil.ReadFile(path)
	if err != nil {
		return 0, errors.New("Error loading driver file - " + err.Error())
	}

	driverXOGTypePattern := model.DriverTypesPattern{}
	xml.Unmarshal(xmlFile, &driverXOGTypePattern)

	v, err := strconv.ParseFloat(driverXOGTypePattern.Version, 64)
	if err != nil || v < constant.VERSION {
		return 0, errors.New(fmt.Sprintf("invalid driver(%s) version, expected version %.1f or greater", path, constant.VERSION))
	}

	if len(driverXOGTypePattern.Files) > 0 {
		return 0, errors.New(fmt.Sprintf("invalid driver(%s) tag <file> is no longer supported", path))
	}

	types := reflect.ValueOf(&driverXOGTypePattern).Elem()
	typeOfT := types.Type()
	for i := 0; i < types.NumField(); i++ {
		t := types.Field(i)
		if t.Kind() == reflect.Slice {
			for _, f := range t.Interface().([]model.DriverFile) {
				f.Type = typeOfT.Field(i).Name
				f.ExecutionOrder = -1
				driverXOG.Files = append(driverXOG.Files, f)
			}
		}
	}

	doc := etree.NewDocument()
	doc.ReadFromBytes(xmlFile)

	for i, e := range doc.FindElements("//xogdriver/*") {
		tag := e.Tag
		path := e.SelectAttrValue("path", constant.UNDEFINED)
		code := e.SelectAttrValue("code", constant.UNDEFINED)
		for y, f := range driverXOG.Files {
			if f.ExecutionOrder == -1 && (strings.ToLower(f.GetXMLType()) == strings.ToLower(tag)) && (f.Path == path) && (f.Code == code) {
				driverXOG.Files[y].ExecutionOrder = i
				break
			}
		}
	}

	sort.Sort(model.ByExecutionOrder(driverXOG.Files))

	driverXOG.Version = driverXOGTypePattern.Version
	driverXOG.FilePath = path

	return len(driverXOG.Files), nil
}

func GetLoadedDriver() *model.Driver {
	return driverXOG
}

func ValidateLoadedDriver() bool {
	if driverXOG == nil {
		return false
	}
	return len(driverXOG.Files) > 0
}

func GetDriversList(folder string) ([]model.Driver, error) {
	driversFileList, err := ioutil.ReadDir(folder)

	if err != nil || len(driversFileList) == 0 {
		return nil, errors.New("driver folder not found or empty")
	}

	var driversList []model.Driver
	for _, f := range driversFileList {
		driver := new(model.Driver)
		driver.Info = f
		driver.FilePath = folder + f.Name()
		driversList = append(driversList, *driver)
	}

	return append(driversList, GetPackagesDriversFileInfoList()...), nil
}

func ProcessDriverFile(file *model.DriverFile, action, sourceFolder, outputFolder string, environments *model.Environments, soapFunc util.Soap) model.Output {
	output := model.Output{Code: constant.OUTPUT_SUCCESS, Debug: ""}
	transformedString := ""

	if action == constant.MIGRATE && file.Type != constant.MIGRATION {
		output.Code = constant.OUTPUT_WARNING
		output.Debug = "Use action 'r' to this type(" + file.Type + ") of file"
		return output
	} else if action == constant.READ && file.Type == constant.MIGRATION {
		output.Code = constant.OUTPUT_WARNING
		output.Debug = "Use action 'm' to this type(" + file.Type + ") of file"
		return output
	}

	if action == constant.MIGRATE {
		resp, err := migration.ReadDataFromExcel(file)
		if err != nil {
			output.Code = constant.OUTPUT_ERROR
			output.Debug = err.Error()
			return output
		}
		transformedString, _ = resp.WriteToString()
		file.SetXML(transformedString)
		file.Write(outputFolder)
		return output
	}

	err := file.InitXML(action, sourceFolder)
	if err != nil {
		output.Code = constant.OUTPUT_ERROR
		output.Debug = err.Error()
		return output
	}
	if action == constant.WRITE {
		iniTagRegexpStr, endTagRegexpStr := file.TagCDATA()
		if iniTagRegexpStr != "" && endTagRegexpStr != "" {
			transformedString := transform.IncludeCDATA(file.GetXML(), iniTagRegexpStr, endTagRegexpStr)
			file.SetXML(transformedString)
		}
	}
	err = file.RunXML(action, sourceFolder, environments, soapFunc)
	if err != nil {
		output.Code = constant.OUTPUT_ERROR
		output.Debug = err.Error()
		return output
	}
	xogResponse := etree.NewDocument()
	xogResponse.ReadFromString(file.GetXML())
	output, err = validate.Check(xogResponse)
	if err != nil {
		output.Code = constant.OUTPUT_ERROR
		output.Debug = err.Error()
		return output
	}
	if action == constant.READ {
		var auxResponse *etree.Document
		if file.NeedAuxXML() {
			auxResponse = etree.NewDocument()
			auxResponse.ReadFromString(file.GetAuxXML())
			output, err = validate.Check(auxResponse)
			if err != nil {
				output.Code = constant.OUTPUT_ERROR
				output.Debug = "[aux] " + err.Error()
				return output
			}
		}
		err = transform.Execute(xogResponse, auxResponse, file)
		str, _ := xogResponse.WriteToString()
		file.SetXML(str)
		if err != nil {
			output.Code = constant.OUTPUT_ERROR
			output.Debug = err.Error()
			return output
		}
		iniTagRegexpStr, endTagRegexpStr := file.TagCDATA()
		if iniTagRegexpStr != "" && endTagRegexpStr != "" {
			transformedString := transform.IncludeCDATA(file.GetXML(), iniTagRegexpStr, endTagRegexpStr)
			file.SetXML(transformedString)
		}
	}

	file.Write(outputFolder)
	return output
}

func CreateFileFolder(action, fileType string) (string, string) {
	sourceFolder := ""
	outputFolder := ""
	switch action {
	case constant.READ:
		sourceFolder = constant.FOLDER_READ
		outputFolder = constant.FOLDER_WRITE
		util.ValidateFolder(constant.FOLDER_READ + fileType)
	case constant.WRITE:
		sourceFolder = constant.FOLDER_WRITE
		outputFolder = constant.FOLDER_DEBUG
	case constant.MIGRATE:
		sourceFolder = constant.FOLDER_WRITE
		outputFolder = constant.FOLDER_MIGRATION
	}

	util.ValidateFolder(outputFolder + fileType)

	return sourceFolder, outputFolder
}
