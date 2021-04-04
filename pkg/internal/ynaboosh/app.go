package ynaboosh

import (
	"log"

	"github.com/thebenwaters/ynaboosh-desktop/pkg/internal/ynaboosh/extensions"
	"github.com/thebenwaters/ynaboosh-desktop/pkg/internal/ynaboosh/forms"
	"github.com/thebenwaters/ynaboosh-desktop/pkg/internal/ynaboosh/models"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"

	_ "embed"

	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/widget"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var (
	boundFileSelected binding.String
	fileType          string = "NONE"
	transactions      []models.Transaction
)

func CreateWindow() fyne.Window {
	boundFileSelected = binding.NewString()
	boundFileSelected.Set("")
	myApp := app.NewWithID("com.github.thebenwaters.ynaboosh-desktop")
	myApp.Settings().SetTheme(&hebrewFontTheme{})
	rootStorage := myApp.Storage().RootURI()
	log.Println(rootStorage)
	dbPath, err := storage.Child(rootStorage, "ynaboosh.db")
	if err != nil {
		log.Panicln(err)
	}
	db, err := gorm.Open(sqlite.Open(dbPath.Path()), &gorm.Config{})
	err = InitializeDB(db)
	log.Println(err)
	appWindow := myApp.NewWindow("YNABoosh")
	transactionTable := extensions.NewTransactionTable()
	syncForm := forms.NewSyncTransactionsUploadForm(transactionTable, appWindow)

	topContainer := container.NewVBox(syncForm, container.NewHBox(widget.NewButton("Approve All", func() {
		transactionTable.ApprovedSetAll(true)
	}), widget.NewButton("Unapprove All", func() {
		transactionTable.ApprovedSetAll(false)
	})))
	editRulesForm := forms.NewRuleEditForm()

	mappings := container.NewGridWithRows(2, editRulesForm, extensions.NewRulesList(editRulesForm))

	syncContainer := container.NewBorder(topContainer, widget.NewButton("Submit to YNAB", func() {}), nil, nil, transactionTable)

	settingsContainer := container.NewVBox(
		widget.NewButton("Login to YNAB", func() {
			log.Println("login")
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
		container.NewTabItem("Rules", mappings),
		container.NewTabItem("Settings", settingsContainer),
	)

	tabs.SetTabLocation(container.TabLocationLeading)

	appWindow.SetContent(tabs)
	return appWindow
}
