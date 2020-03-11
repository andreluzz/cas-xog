package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/andreluzz/cas-xog/view"
)

var version = "Development Build"

func main() {
	if len(os.Args) > 1 {
		arg := strings.ToLower(os.Args[1])
		if strings.Contains(arg, "version") {
			fmt.Printf("CAS-XOG version: %s\n", version)
			return
		}
	}

	view.Home(version)
	var exit = false
	for {
		exit = view.Interface()
		if exit {
			break
		}
	}
}
