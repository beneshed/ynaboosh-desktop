package main

import (
	"log"

	"github.com/thebenwaters/ynaboosh-desktop/pkg/internal/ynaboosh"
)

func main() {
	window, err := ynaboosh.CreateWindow()
	if err != nil {
		log.Panicln(err)
	}
	window.ShowAndRun()
}
