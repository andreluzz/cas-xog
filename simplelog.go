package main

import (
	"fmt"
	"github.com/mattn/go-colorable"
)

func Debug(format string, args ...interface{}) {
	fmt.Fprintf(colorable.NewColorableStdout(), format, args...)
}
