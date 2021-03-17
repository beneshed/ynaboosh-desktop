package main

import (
	"log"

	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/thebenwaters/ynab-desktop-importer/internal/pkg/flows"
	"github.com/thebenwaters/ynab-desktop-importer/internal/pkg/screens"
	"github.com/thebenwaters/ynab-desktop-importer/internal/pkg/ynabimporter"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func main() {
	a := app.NewWithID("dev.benwaters.ynabimporter")
	w := a.NewWindow("YNAB CSV Importer")

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

	log.Println(db, err)

	table := NewMappedTable()
	uploadFlow := globalState.NewUploadFlow(w)
	uploadCreditCardTranastionButton := widget.NewButton("Upload Credit Card Transactions", func() {
		w.SetContent(uploadFlow)
	})
	uploadBankTransactions := widget.NewButton("Upload Bank Transactions", func() {
		w.SetContent(uploadFlow)
	})

	testContainer := container.NewMax()
	testContainer.Add(flows.NewMappedTable())

	uploadContainer := container.NewVBox()
	uploadContainer.Add(uploadCreditCardTranastionButton)
	uploadContainer.Add(widget.NewSeparator())
	uploadContainer.Add(uploadBankTransactions)
	tabs := container.NewAppTabs(
		container.NewTabItem("File Uploads", uploadContainer),
		container.NewTabItem("History", widget.NewLabel("Historical Uploaded Transactions")),
		container.NewTabItem("Suggested Categories", table),
		container.NewTabItem("Settings", globalState.NewSettingsScreen()),
		container.NewTabItem("About", screens.NewAbout().NewAboutScreen()),
		container.NewTabItem("Test", testContainer),
	)

	w.SetContent(tabs)
	w.SetFullScreen(true)

	w.ShowAndRun()

}
