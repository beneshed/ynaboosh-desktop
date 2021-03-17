package flows

import (
	"log"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
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
	return container.NewMax(widget.NewLabel(""), hiddenCheck)
}

func NewMappedTable() *widget.Table {
	state := NewMappedTableState()
	state.data[1] = MappedTableRow{
		Date:     "adsjfaldksjfadkl;sjfkadfjasdkljfadsl",
		Business: "adflkajdslfkjadsklf",
		Amount:   "adlfkjasdfkljas",
		Category: "adfadsfadsfadsfas",
		Verify:   false,
	}
	table := widget.NewTable(func() (int, int) {
		return state.Length()
	}, func() fyne.CanvasObject { return state.Init() }, func(location widget.TableCellID, cell fyne.CanvasObject) {
		c := cell.(*fyne.Container)
		log.Println("c", c)
		for _, obj := range c.Objects {
			log.Println("OBJECT FUCKER")
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
			case *widget.Check:
				if location.Col == len(state.headers)-1 && location.Row != 0 {
					typedObj.Show()
				}
			}
		}
	})
	for i := 0; i < 5; i++ {
		table.SetColumnWidth(i, 250)
	}
	return table
}
