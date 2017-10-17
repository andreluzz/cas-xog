package main

import (
	"github.com/andreluzz/cas-xog/xog"
)

func main() {
	xog.RenderHome()
	var exit = false
	for {
		exit = xog.RenderInterface()
		if exit {
			break
		}
	}
}
