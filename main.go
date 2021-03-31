package main

import (
	"log"

	"github.com/thebenwaters/ynab-utils-desktop/institutions"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/dialog"

	_ "embed"

	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

var (
	boundFileSelected binding.String
	fileType          string = "NONE"
	transactions      []institutions.Transaction
)

type MappedTableRow struct {
	Date     string `fyne:"Date of Transaction"`
	Business string `fyne:"Name of Business"`
	Amount   string `fyne:"Amount"`
	Category string `fune:"Category"`
	Verify   bool   `fyne:"Verify"`
}

type MappedTableState struct {
	headers []string
	data    map[int]MappedTableRow
}

func NewMappedTableState() *MappedTableState {
	return &MappedTableState{
		headers: []string{"Date", "Business", "Amount", "Category", "Verify"},
		data:    make(map[int]MappedTableRow),
	}
}

func (s MappedTableState) Length() (int, int) {
	return len(s.data) + 1, len(s.headers)
}

func (s MappedTableState) Init() *fyne.Container {
	hiddenCheck := widget.NewCheck("", func(b bool) {
		log.Println("test")
	})
	hiddenCheck.Hide()
	hiddenDropDown := widget.NewSelect([]string{"a", "b", "c"}, func(selected string) {
		log.Println(selected)
	})
	hiddenDropDown.PlaceHolder = " "
	hiddenDropDown.Hide()
	return container.NewMax(widget.NewLabel(""), hiddenDropDown, hiddenCheck)
}

func NewMappedTable() *widget.Table {
	state := NewMappedTableState()
	/*
		state.data[1] = MappedTableRow{
			Date:     "1/1/1",
			Business: "Fake Business",
			Amount:   "100",
			Category: "",
			Verify:   false,
		}
		state.data[2] = MappedTableRow{
			Date:     "1/1/1",
			Business: "Fake Business",
			Amount:   "100",
			Category: "",
			Verify:   false,
		}
		state.data[3] = MappedTableRow{
			Date:     "1/1/1",
			Business: "Fake Business",
			Amount:   "100",
			Category: "",
			Verify:   false,
		}
	*/
	table := widget.NewTable(func() (int, int) {
		return state.Length()
	}, func() fyne.CanvasObject { return state.Init() }, func(location widget.TableCellID, cell fyne.CanvasObject) {
		c := cell.(*fyne.Container)
		for _, obj := range c.Objects {
			switch typedObj := obj.(type) {
			case *widget.Label:
				if location.Row == 0 {
					switch location.Col {
					case 0:
						typedObj.SetText("Date")
					case 1:
						typedObj.SetText("Business")
					case 2:
						typedObj.SetText("Amount")
					case 3:
						typedObj.SetText("Category")
					case 4:
						typedObj.SetText("Verified")
					}
				} else {
					switch location.Col {
					case 0:
						typedObj.SetText(state.data[location.Row].Date)
					case 1:
						typedObj.SetText(state.data[location.Row].Business)
					case 2:
						typedObj.SetText(state.data[location.Row].Amount)
					case 3:
						typedObj.SetText(state.data[location.Row].Category)
					}
				}
			case *widget.Select:
				if location.Col == len(state.headers)-2 && location.Row != 0 {
					typedObj.Show()
				}
			case *widget.Check:
				if location.Col == len(state.headers)-1 && location.Row != 0 {
					typedObj.Show()
				}
			}
		}
	})
	table.OnSelected = func(id widget.TableCellID) {
		if id.Col == len(state.headers)-1 && id.Row == 0 {
			log.Println("select all")
		}
	}
	for i := 0; i < 5; i++ {
		table.SetColumnWidth(i, 250)
	}
	return table
}

func main() {
	boundFileSelected = binding.NewString()
	boundFileSelected.Set("")
	myApp := app.New()
	myWindow := myApp.NewWindow("YNAB Desktop Importer")

	filePicker := dialog.NewFileOpen(func(info fyne.URIReadCloser, err error) {
		log.Println(info, err)
		if err == nil && info != nil {
			log.Println("should work")
			err = boundFileSelected.Set(info.URI().Path())
			log.Println(err)
		}
	}, myWindow)

	transactionTable := NewTransactionTable()

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
	form.OnSubmit = func() {
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

	mappings := container.NewBorder(nil, container.NewMax(widget.NewButton("Next", func() { log.Println("test") })), nil, nil, NewMappedTable())

	syncContainer := container.NewBorder(form, widget.NewButton("Submit to YNAB", func() {}), nil, nil, transactionTable)

	settingsContainer := container.NewVBox(
		widget.NewButton("Login to YNAB", func() {
			log.Println("login")
		}),
		widget.NewButton("Logout of YNAB", func() {
			log.Println("tapped")
		}),
		widget.NewButton("Refresh YNAB Categories", func() {
			log.Println("refreshed")
		}),
	)

	tabs := container.NewAppTabs(
		container.NewTabItem("Sync Transactions", syncContainer),
		container.NewTabItem("Rules", mappings),
		container.NewTabItem("Settings", settingsContainer),
	)

	//tabs.Append(container.NewTabItemWithIcon("Home", theme.HomeIcon(), widget.NewLabel("Home tab")))

	tabs.SetTabLocation(container.TabLocationLeading)

	myWindow.SetContent(tabs)
	myWindow.ShowAndRun()
}
