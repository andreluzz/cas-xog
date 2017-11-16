package main

import (
	"github.com/andreluzz/cas-xog/render"
)

func main() {
	render.Home()
	var exit = false
	for {
		exit = render.Interface()
		if exit {
			break
		}
	}
}
