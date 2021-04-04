package ynaboosh

import (
	"log"

	"github.com/thebenwaters/ynaboosh-desktop/pkg/internal/ynaboosh/forms"
	"github.com/thebenwaters/ynaboosh-desktop/pkg/internal/ynaboosh/models"
	"github.com/thebenwaters/ynaboosh-desktop/pkg/internal/ynaboosh/tables"
	"github.com/thebenwaters/ynaboosh-desktop/pkg/internal/ynaboosh/ynab"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"

	_ "embed"

	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/widget"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	externalYnab "go.bmvs.io/ynab"
)

var (
	externalYNABClient externalYnab.ClientServicer
)

func CreateWindow() fyne.Window {
	myApp := app.NewWithID("com.github.thebenwaters.ynaboosh-desktop")
	myApp.Settings().SetTheme(&hebrewFontTheme{})
	rootStorage := myApp.Storage().RootURI()
	log.Println(rootStorage)
	dbPath, err := storage.Child(rootStorage, "ynaboosh.db")
	if err != nil {
		log.Panicln(err)
	}
	db, _ := gorm.Open(sqlite.Open(dbPath.Path()), &gorm.Config{})
	manager := &models.DBManager{db}
	err = InitializeDB(db)
	log.Println(err)
	appWindow := myApp.NewWindow("YNABoosh")
	transactionTable := tables.NewTransactionTable()
	transactionTable.WrapTableWidth()
	syncForm := forms.NewSyncTransactionsUploadForm(transactionTable, manager, appWindow)

	topContainer := container.NewVBox(syncForm, container.NewHBox(widget.NewButton("Approve All", func() {
		transactionTable.ApprovedSetAll(true)
	}), widget.NewButton("Unapprove All", func() {
		transactionTable.ApprovedSetAll(false)
	})))
	editRulesForm := forms.NewRuleEditForm(manager)

	mappings := container.NewGridWithRows(2, editRulesForm, tables.NewRulesTable(editRulesForm, manager))

	syncContainer := container.NewBorder(topContainer, widget.NewButton("Submit to YNAB", func() {}), nil, nil, transactionTable)

	settingsContainer := container.NewVBox(
		widget.NewButton("Login to YNAB", func() {
			externalYNABClient = ynab.Login(*manager)
		}),
		widget.NewButton("Logout of YNAB", func() {
			log.Println("tapped")
		}),
		widget.NewButton("Refresh YNAB Data", func() {
			log.Println("refreshed")
		}),
	)

	tabs := container.NewAppTabs(
		container.NewTabItem("Sync Transactions", syncContainer),
		container.NewTabItem("Historical Transactions", widget.NewLabel("foo")),
		container.NewTabItem("Rules", mappings),
		container.NewTabItem("Settings", settingsContainer),
	)

	tabs.SetTabLocation(container.TabLocationLeading)

	appWindow.SetContent(tabs)
	return appWindow
}
