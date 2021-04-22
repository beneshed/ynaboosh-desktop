package ynaboosh

import (
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

	"go.uber.org/zap"

	externalYnab "go.bmvs.io/ynab"
)

var (
	externalYNABClient externalYnab.ClientServicer
)

func CreateWindow() (fyne.Window, error) {
	logger, _ := zap.NewProduction()
	defer logger.Sync() // flushes buffer, if any
	sugar := logger.Sugar()
	myApp := app.NewWithID("com.github.thebenwaters.ynaboosh-desktop")
	myApp.Settings().SetTheme(&hebrewFontTheme{})
	rootStorage := myApp.Storage().RootURI()
	sugar.Infof("Storage root: %s", rootStorage)
	dbPath, err := storage.Child(rootStorage, "ynaboosh.db")
	if err != nil {
		sugar.Error(err)
		return nil, err
	}
	db, _ := gorm.Open(sqlite.Open(dbPath.Path()), &gorm.Config{})
	manager := &models.DBManager{db}
	err = InitializeDB(db)
	if err != nil {
		sugar.Error(err)
		return nil, err
	}
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

	editRulesFormContainer := container.NewVBox(editRulesForm, widget.NewButton("Fix RTL Lanuages", func() {}))

	rulesTable := tables.NewRulesTable(editRulesForm, manager)

	mappings := container.NewGridWithRows(2, editRulesFormContainer, rulesTable)

	submitButton := widget.NewButton("Submit to YNAB", func() {})
	submitButton.Importance = widget.HighImportance

	syncContainer := container.NewBorder(topContainer, submitButton, nil, nil, transactionTable)

	settingsContainer := container.NewVBox(
		widget.NewButton("Login to YNAB", func() {
			externalYNABClient = ynab.Login(*manager)
		}),
		widget.NewButton("Logout of YNAB", func() {

		}),
		widget.NewButton("Refresh YNAB Data", func() {

		}),
	)

	tabs := container.NewAppTabs(
		container.NewTabItem("Sync Transactions", syncContainer),
		container.NewTabItem("Historical Transactions", widget.NewLabel("foo")),
		container.NewTabItem("Rules", mappings),
		container.NewTabItem("Settings", settingsContainer),
	)

	tabs.SetTabLocation(container.TabLocationLeading)

	// menu
	mainMenu := fyne.NewMainMenu(fyne.NewMenu("File", fyne.NewMenuItem("Import Rules", rulesTable.ImportRules)))

	appWindow.SetContent(tabs)
	appWindow.SetMainMenu(mainMenu)
	return appWindow, nil
}
