package ynaboosh

import (
	"log"

	"github.com/thebenwaters/ynaboosh-desktop/pkg/internal/ynaboosh/extensions"
	"github.com/thebenwaters/ynaboosh-desktop/pkg/internal/ynaboosh/institutions"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/dialog"

	_ "embed"

	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var (
	boundFileSelected binding.String
	fileType          string = "NONE"
	transactions      []institutions.Transaction
)

const (
	exampleWhen = `Transaction.Payee.Contains("לסרפוש")`
	exampleThen = `Transaction.Category = "Groceries";
Transaction.Payee = "Shufersal";`
	exampleRuleName        = "Shufersal Groceries"
	exampleRuleDescription = "Assign grocery category to shufersal and switch to english"
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

	filePicker := dialog.NewFileOpen(func(info fyne.URIReadCloser, err error) {
		log.Println(info, err)
		if err == nil && info != nil {
			log.Println("should work")
			err = boundFileSelected.Set(info.URI().Path())
			log.Println(err)
		}
	}, appWindow)

	transactionTable := extensions.NewTransactionTable()

	// sync transactions
	// 1) form
	// 2) adjust / fix / confirm
	form := widget.NewForm(
		widget.NewFormItem(
			"Transaction File", widget.NewButtonWithIcon("", theme.FileIcon(), func() {
				filePicker.Show()
			}),
		),
		widget.NewFormItem(
			"File Type", widget.NewSelect([]string{"Discount Bank - Hebrew", "Discount Bank - English", "Leumi", "Cal", "Max"}, func(value string) {
				fileType = value
			}),
		),
		widget.NewFormItem(
			"File Path", widget.NewLabelWithData(boundFileSelected),
		),
		widget.NewFormItem(
			"Account", widget.NewSelect([]string{"foo", "bar"}, func(value string) {
				log.Println(value)
			}),
		),
	)
	clearableForm := extensions.NewClearableForm(form)
	clearableForm.OnSubmit = func() {
		// run file detection
		fileName, err := boundFileSelected.Get()
		if err != nil {
			log.Panicln(err)
		}
		log.Println("About to parse")
		switch fileType {
		case "Cal":
			transactions := institutions.ParseCalTransations(fileName)
			transactionTable.AddTransactions(transactions)
		case "Leumi":
			transactions = append(transactions, institutions.ParseLeumiTransactions(fileName)...)
			transactionTable.AddTransactions(transactions)
		case "Discount Bank - Hebrew":
			transactions = append(transactions, institutions.ParseDiscountTransactions(fileName, true)...)
			transactionTable.AddTransactions(transactions)
		case "Max":
			transactions = append(transactions, institutions.ParseMaxTransations(fileName)...)
			transactionTable.AddTransactions(transactions)
		}
	}
	form.SubmitText = "Load File"
	topContainer := container.NewVBox(form, container.NewHBox(widget.NewButton("Approve All", func() {
		transactionTable.ApprovedSetAll(true)
	}), widget.NewButton("Unapprove All", func() {
		transactionTable.ApprovedSetAll(false)
	})))

	whenEntry := widget.NewMultiLineEntry()
	whenEntry.SetPlaceHolder(exampleWhen)
	thenEntry := widget.NewMultiLineEntry()
	thenEntry.SetPlaceHolder(exampleThen)
	ruleNameEntry := widget.NewEntry()
	ruleNameEntry.SetPlaceHolder(exampleRuleName)
	descriptionEntry := widget.NewEntry()
	descriptionEntry.SetPlaceHolder(exampleRuleDescription)

	createRuleForm := widget.NewForm(
		widget.NewFormItem("Rule Name", ruleNameEntry),
		widget.NewFormItem("Description", descriptionEntry),
		widget.NewFormItem("When", whenEntry),
		widget.NewFormItem("Then", thenEntry),
	)
	clearableCreateRuleForm := extensions.NewClearableForm(createRuleForm)
	clearableCreateRuleForm.OnSubmit = func() {
		clearableCreateRuleForm.Clear()
	}
	clearableCreateRuleForm.SubmitText = "Add"
	clearableCreateRuleForm.OnCancel = func() {
		clearableCreateRuleForm.SubmitText = "Add"
		clearableCreateRuleForm.Clear()

	}

	mappings := container.NewGridWithRows(2, createRuleForm, extensions.NewRulesList(whenEntry, clearableCreateRuleForm))

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
