package util

import (
	"bytes"
	"os"
	"reflect"
	"regexp"
	"runtime"
	"strings"
	"unsafe"

	"github.com/andreluzz/cas-xog/constant"
)

//BytesToString convert an array of bytes to a string
func BytesToString(b []byte) string {
	bh := (*reflect.SliceHeader)(unsafe.Pointer(&b))
	sh := reflect.StringHeader{Data: bh.Data, Len: bh.Len}
	return *(*string)(unsafe.Pointer(&sh))
}

//ValidateFolder creates the folder structure if it do not exists
func ValidateFolder(folder string) error {
	if folder != "" && folder[len(folder)-1:] == "/" {
		folder = folder[0 : len(folder)-1]
	}
	_, dirErr := os.Stat(folder)
	if os.IsNotExist(dirErr) {
		err := os.MkdirAll(folder, os.ModePerm)
		if err != nil {
			return err
		}
	}
	return nil
}

//GetOutputDebug formats the string from the debug when it has errors or warnings
func GetOutputDebug(code, debug string) string {
	if code != constant.OutputSuccess {
		return "| Debug: " + debug
	}
	return debug
}

//GetStatusColorFromOutput returns the status and color from the output struct
func GetStatusColorFromOutput(code string) (string, string) {
	switch code {
	case constant.OutputSuccess:
		return "success", "green"
	case constant.OutputWarning:
		return "warning", "yellow"
	case constant.OutputError:
		return "error  ", "red"
	}
	return "", ""
}

//GetActionLabel returns the properly formatted string according to the constant action
func GetActionLabel(action string) string {
	switch action {
	case constant.Read:
		return "Read"
	case constant.Write:
		return "Write"
	case constant.Migrate:
		return "Create"
	}
	return ""
}

//RightPad insert a defined number of characters on the right of the string
func RightPad(s, padStr string, length int) string {
	var padCountInt int
	padCountInt = 1 + ((length - len(padStr)) / len(padStr))
	var retStr = s + strings.Repeat(padStr, padCountInt)
	return retStr[:length]
}

//GetPathFolder returns only the folders without filename and extension
func GetPathFolder(path string) string {
	folder := ""

	re := regexp.MustCompile(`.*[/\\]`)
	match := re.FindStringSubmatch(path)

	if len(match) > 0 {
		folder = match[0]
		matchInit, _ := regexp.MatchString(`^[/\\]`, path)

		if !matchInit {
			folder = "/" + folder
		}
	}

	return folder
}

//GetPathWithoutExtension returns the folders and filename without the file extension
func GetPathWithoutExtension(path string) string {
	extIndex := strings.LastIndex(path, ".")
	return path[:extIndex]
}

//GetExtension returns the file extension
func GetExtension(path string) string {
	extIndex := strings.LastIndex(path, ".")
	return path[extIndex:]
}

//GetDirectFolder returns the file closest folder
func GetDirectFolder(path string) string {
	extIndex := strings.LastIndex(path, GetPathSeparator())
	folder := path[:extIndex]
	extIndex = strings.LastIndex(folder, GetPathSeparator())
	return folder[extIndex+1:]
}

//GetPathSeparator returns the folder separator by OS
func GetPathSeparator() string {
	return string(os.PathSeparator)
}

//ReplacePathSeparatorByOS returns path with right folder separator by OS
func ReplacePathSeparatorByOS(path string) string {
	var result string
	if runtime.GOOS == "windows" {
		result = strings.Replace(path, "/", GetPathSeparator(), -1)
	} else {
		result = strings.Replace(path, "\\", GetPathSeparator(), -1)
	}
	return result
}

//JSONEscapeText avoid escaping caracters like <, > and & from json byte array
func JSONAvoidEscapeText(data []byte) []byte {
	data = bytes.Replace(data, []byte("\\u003c"), []byte("<"), -1)
	data = bytes.Replace(data, []byte("\\u003e"), []byte(">"), -1)
	data = bytes.Replace(data, []byte("\\u0026"), []byte("&"), -1)
	return data
}
