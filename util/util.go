package util

import (
	"github.com/andreluzz/cas-xog/constant"
	"os"
	"reflect"
	"unsafe"
)

func BytesToString(b []byte) string {
	bh := (*reflect.SliceHeader)(unsafe.Pointer(&b))
	sh := reflect.StringHeader{Data: bh.Data, Len: bh.Len}
	return *(*string)(unsafe.Pointer(&sh))
}

func ValidateFolder(folder string) error {
	_, dirErr := os.Stat(folder)
	if os.IsNotExist(dirErr) {
		err := os.MkdirAll(folder, os.ModePerm)
		if err != nil {
			return err
		}
	}
	return nil
}

func GetOutputDebug(debug string) string {
	if debug == "" {
		return ""
	}
	return "| Debug: " + debug
}

func GetStatusColorFromOutput(code string) (string, string) {
	switch code {
	case constant.OUTPUT_SUCCESS:
		return "success", "green"
	case constant.OUTPUT_WARNING:
		return "warning", "yellow"
	case constant.OUTPUT_ERROR:
		return "error  ", "red"
	}
	return "", ""
}

func GetActionLabel(action string) string {
	switch action {
	case constant.READ:
		return "Read"
	case constant.WRITE:
		return "Write"
	case constant.MIGRATE:
		return "Create"
	}
	return ""
}