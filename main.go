package main

import (
	"log"

	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"github.com/thebenwaters/ynab-desktop-importer/internal/pkg/flows"
	"github.com/thebenwaters/ynab-desktop-importer/internal/pkg/screens"
	"github.com/thebenwaters/ynab-desktop-importer/internal/pkg/ynabimporter"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func main() {
	a := app.NewWithID("dev.benwaters.ynabimporter")
	w := a.NewWindow("YNAB Desktop Importer")

	globalState := ynabimporter.NewGlobalState()

	db, err := gorm.Open(sqlite.Open("gorm.db"), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}
	err = ynabimporter.InitializeDB(db)
	if err != nil {
		log.Fatalln(err)
	}
	globalState.DB = db
	wrappedState := flows.InternalState(*globalState)
	uploadFlow := wrappedState.NewUploadFlow(w)
	homeScreen := screens.NewHomeScreen(screens.NewUploadScreen(w, uploadFlow), container.NewVBox(), container.NewVBox())
	w.SetMainMenu(screens.NewMainMenu())
	w.SetContent(homeScreen)
	w.SetFullScreen(true)

	w.ShowAndRun()

}
