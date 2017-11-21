package main

import (
	"github.com/andreluzz/cas-xog/view"
)

func main() {
	view.Home()
	var exit = false
	for {
		exit = view.Interface()
		if exit {
			break
		}
	}
}
