package main

import (
	"log"

	internal "github.com/thebenwaters/ynab-desktop-importer/internal/pkg/app"
)

func main() {

	log.Println("pre")
	log.Println("hello")

	globalState := internal.NewApp()
	/*
		wrappedState := flows.InternalState(*globalState)
		parent := app.NewApp()
		parent.Object = con
		uploadFlow := wrappedState.NewUploadFlow(w)
		homeScreen := screens.NewHomeScreen(screens.NewUploadScreen(*globalState, uploadFlow), container.NewVBox(), container.NewVBox())
		globalState.Homescreen = homeScreen
	*/
	globalState.Run()

}
