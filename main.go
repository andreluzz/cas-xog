package main

import (
	"github.com/andreluzz/cas-xog/view"
)

var version string

func main() {
	view.Home(version)
	var exit = false
	for {
		exit = view.Interface()
		if exit {
			break
		}
	}
}
